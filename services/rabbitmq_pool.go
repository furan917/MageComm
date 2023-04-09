package services

import (
	"errors"
	"github.com/streadway/amqp"
	"magecomm/config_manager"
	"sync"
)

var (
	ErrConnectionPoolClosed = errors.New("connection pool is closed")
	ErrChannelPoolClosed    = errors.New("channel pool is closed")
	rmqOnce                 sync.Once
	RmqConnectionPool       *RabbitMQConnectionPool
	RmqChannelPool          *RabbitMQChannelPool
)

type RabbitMQConnectionPool struct {
	pool *sync.Pool
}

type RabbitMQChannelPool struct {
	pool *sync.Pool
}

func Close() {
	if RmqConnectionPool != nil {
		RmqConnectionPool.Close()
	}
}

func NewRabbitMQConnectionPool(initialSize int) *RabbitMQConnectionPool {
	p := &sync.Pool{
		New: func() interface{} {
			rmqConn := NewRabbitMQConnection()
			if err := rmqConn.Connect(""); err != nil {
				return err
			}
			return rmqConn
		},
	}

	cp := &RabbitMQConnectionPool{pool: p}

	for i := 0; i < initialSize; i++ {
		p.Put(p.New())
	}

	return cp
}

func NewRabbitMQChannelPool(connPool *RabbitMQConnectionPool, initialSize int) *RabbitMQChannelPool {
	p := &sync.Pool{
		New: func() interface{} {
			conn, err := connPool.Get()
			if err != nil {
				return err
			}
			channel, err := conn.CreateChannel()
			if err != nil {
				return err
			}
			return channel
		},
	}

	cp := &RabbitMQChannelPool{pool: p}

	for i := 0; i < initialSize; i++ {
		p.Put(p.New())
	}

	return cp
}

func (cp *RabbitMQConnectionPool) Get() (*RabbitMQConnection, error) {
	conn := cp.pool.Get()
	if conn == nil {
		return nil, ErrConnectionPoolClosed
	}
	return conn.(*RabbitMQConnection), nil
}

func (cp *RabbitMQConnectionPool) Put(conn *RabbitMQConnection) {
	cp.pool.Put(conn)
}

func (cp *RabbitMQConnectionPool) Close() {
	cp.pool = nil
}

func (cp *RabbitMQChannelPool) Get() (*amqp.Channel, error) {
	channel := cp.pool.Get()
	if channel == nil {
		return nil, ErrChannelPoolClosed
	}
	return channel.(*amqp.Channel), nil
}

func (cp *RabbitMQChannelPool) Put(channel *amqp.Channel) {
	cp.pool.Put(channel)
}

func (cp *RabbitMQChannelPool) Close() {
	cp.pool = nil
}

func init() {
	rmqOnce.Do(func() {
		engine := config_manager.GetValue(config_manager.CommandConfigListenerEngine)
		if engine == EngineRabbitMQ {
			RmqConnectionPool = NewRabbitMQConnectionPool(1)
			RmqChannelPool = NewRabbitMQChannelPool(RmqConnectionPool, 5)
		}
	})
}

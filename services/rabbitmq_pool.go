package services

import (
	"errors"
	"fmt"

	"magecomm/config_manager"
	"magecomm/logger"
	"sync"

	"github.com/streadway/amqp"
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
	pool     *sync.Pool
	connPool *RabbitMQConnectionPool
}

func Close() {
	if RmqConnectionPool != nil {
		RmqConnectionPool.Close()
	}
}

func NewRabbitMQConnectionPool(initialSize int) *RabbitMQConnectionPool {
	p := &sync.Pool{
		New: func() interface{} {
			return NewRabbitMQConnection()
		},
	}

	cp := &RabbitMQConnectionPool{pool: p}

	for i := 0; i < initialSize; i++ {
		conn := NewRabbitMQConnection()
		if err := conn.Connect(""); err != nil {
			logger.Warnf("Failed to create initial RMQ connection: %v", err)
			continue
		}
		p.Put(conn)
	}

	return cp
}

func NewRabbitMQChannelPool(connPool *RabbitMQConnectionPool, initialSize int) *RabbitMQChannelPool {
	p := &sync.Pool{
		New: func() interface{} {
			return nil
		},
	}

	cp := &RabbitMQChannelPool{
		pool:     p,
		connPool: connPool,
	}

	for i := 0; i < initialSize; i++ {
		conn, err := connPool.Get()
		if err != nil {
			logger.Warnf("Failed to get connection for initial channel: %v", err)
			continue
		}
		channel, err := conn.CreateChannel()
		if err != nil {
			logger.Warnf("Failed to create initial RMQ channel: %v", err)
			connPool.Put(conn)
			continue
		}
		connPool.Put(conn)
		p.Put(channel)
	}

	return cp
}

func (cp *RabbitMQConnectionPool) Get() (*RabbitMQConnection, error) {
	obj := cp.pool.Get()
	if obj == nil {
		return nil, ErrConnectionPoolClosed
	}

	conn, ok := obj.(*RabbitMQConnection)
	if !ok || conn == nil {
		conn = NewRabbitMQConnection()
	}

	if conn.Connection == nil || conn.Connection.IsClosed() {
		if err := conn.Connect(""); err != nil {
			return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
		}
	}

	return conn, nil
}

func (cp *RabbitMQConnectionPool) Put(conn *RabbitMQConnection) {
	if conn != nil && conn.Connection != nil && !conn.Connection.IsClosed() {
		cp.pool.Put(conn)
	}
}

func (cp *RabbitMQConnectionPool) Close() {
	cp.pool = nil
}

func (cp *RabbitMQChannelPool) Get() (*amqp.Channel, error) {
	obj := cp.pool.Get()

	if obj != nil {
		if channel, ok := obj.(*amqp.Channel); ok && channel != nil {
			return channel, nil
		}
	}

	conn, err := cp.connPool.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	channel, err := conn.CreateChannel()
	if err != nil {
		cp.connPool.Put(conn)
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	cp.connPool.Put(conn)
	return channel, nil
}

func (cp *RabbitMQChannelPool) Put(channel *amqp.Channel) {
	if channel != nil {
		cp.pool.Put(channel)
	}
}

func (cp *RabbitMQChannelPool) Close() {
	cp.pool = nil
}

func InitializeRMQ() {
	rmqOnce.Do(func() {
		engine := config_manager.GetValue(config_manager.CommandConfigListenerEngine)
		if engine == EngineRabbitMQ {
			RmqConnectionPool = NewRabbitMQConnectionPool(1)
			RmqChannelPool = NewRabbitMQChannelPool(RmqConnectionPool, 5)
		}
	})
}

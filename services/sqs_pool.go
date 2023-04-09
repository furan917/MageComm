package services

import (
	"errors"
	"magecomm/config_manager"
	"sync"
)

const (
	maxSQSConnections = 5
)

type SQSConnectionPool struct {
	pool *sync.Pool
}

var sqsConnectionPool *SQSConnectionPool
var sqsOnce sync.Once

func GetSQSConnection() (*SQSConnection, error) {
	if sqsConnectionPool == nil {
		return nil, errors.New("SQS connection pool is not initialized")
	}

	conn := sqsConnectionPool.pool.Get()
	if conn == nil {
		return nil, ErrConnectionPoolClosed
	}

	return conn.(*SQSConnection), nil
}

func ReleaseSQSConnection(conn *SQSConnection) {
	if sqsConnectionPool != nil {
		sqsConnectionPool.pool.Put(conn)
	}
}

func init() {
	sqsOnce.Do(func() {
		engine := config_manager.GetValue(config_manager.CommandConfigListenerEngine)
		if engine == EngineSQS {
			sqsConnectionPool = &SQSConnectionPool{
				pool: &sync.Pool{
					New: func() interface{} {
						sqsConn := NewSQSConnection()
						if err := sqsConn.Connect(); err != nil {
							return err
						}
						return sqsConn
					},
				},
			}

			for i := 0; i < maxSQSConnections; i++ {
				sqsConnectionPool.pool.Put(sqsConnectionPool.pool.New())
			}
		}
	})
}

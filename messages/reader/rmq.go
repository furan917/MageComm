package reader

import (
    "fmt"
    "strings"

    "magecomm/logger"
    "magecomm/messages/queues"
    "magecomm/services"

    "github.com/streadway/amqp"
)

type RmqReader struct{}

func (r *RmqReader) DrainOutputQueue(queueName string) (int, error) {
    outputQueueName := queues.MapQueueToOutputQueue(queueName)

    channel, err := services.RmqChannelPool.Get()
    if err != nil {
        return 0, fmt.Errorf("error getting RMQ channel: %w", err)
    }
    defer services.RmqChannelPool.Put(channel)

    queueNameWithPrefix, err := services.CreateRmqQueue(channel, outputQueueName)
    if err != nil {
        return 0, fmt.Errorf("error creating queue: %w", err)
    }

    count := 0
    for {
        msg, ok, err := channel.Get(queueNameWithPrefix, false)
        if err != nil {
            return count, fmt.Errorf("error receiving message: %w", err)
        }

        if !ok {
            break
        }

        count++
        r.displayMessage(count, &msg)

        if err := msg.Ack(false); err != nil {
            logger.Warnf("Failed to ack message: %v", err)
        }
    }

    return count, nil
}

func (r *RmqReader) displayMessage(index int, msg *amqp.Delivery) {
    fmt.Printf("\n%s\n", strings.Repeat("=", 60))
    fmt.Printf("Output %d\n", index)
    fmt.Printf("Received: %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"))
    fmt.Printf("%s\n", strings.Repeat("=", 60))
    fmt.Println(string(msg.Body))
}

package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

const maxMessageRetries = 3

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewClient(conn *amqp.Connection, queueName string) (*Client, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	deadLetterQueueName := fmt.Sprintf("%s.dlq", queueName)

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": deadLetterQueueName,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	_, err = ch.QueueDeclare(
		deadLetterQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare DLQ: %w", err)
	}

	return &Client{
		conn:    conn,
		channel: ch,
		queue:   queueName,
	}, nil
}

func (c *Client) Publish(ctx context.Context, message model.Message) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}
	err = c.channel.PublishWithContext(
		ctx,
		"",      // exchange
		c.queue, // routing key (queue name)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (c *Client) Consume(ctx context.Context, cb func(ctx context.Context, n model.Message) error) error {
	msgs, err := c.channel.Consume(
		c.queue,
		"",
		false, // manual ack
		false, // not exclusive
		false, // no local
		false, // no wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			logger.GetLogger().Info("stopping RabbitMQ consumer")
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed")
			}

			var m model.Message
			if err := json.Unmarshal(msg.Body, &m); err != nil {
				logger.GetLogger().Errorf("invalid message: %v", err)
				// Malformed messages are discarded to DLQ.
				if err := msg.Nack(false, false); err != nil {
					logger.GetLogger().Errorf("failed to ack message: %v", err)
				}
				continue
			}

			retries := 0
			for retries < maxMessageRetries {
				if err := cb(ctx, m); err != nil {
					retries++
					logger.GetLogger().Errorf("cb failed (attempt %d/%d): %v", retries, maxMessageRetries, err)
					if retries == maxMessageRetries {
						logger.GetLogger().Errorf("max retries reached, discarding message")
						// Discarded to DLQ.
						if err := msg.Nack(false, false); err != nil {
							logger.GetLogger().Errorf("failed to nack message: %v", err)
						}
						break
					}

					delay := time.Duration(retries) * 2 * time.Second
					time.Sleep(delay)
					continue
				}

				if err := msg.Ack(false); err != nil {
					logger.GetLogger().Errorf("failed to ack message: %v", err)
				}
				break
			}
		}
	}
}

func (r *Client) Close() error {
	var closeErr error

	defer func() {
		if err := r.conn.Close(); err != nil && closeErr == nil {
			closeErr = fmt.Errorf("failed to close connection: %w", err)
		}
	}()

	if err := r.channel.Close(); err != nil {
		closeErr = fmt.Errorf("failed to close channel: %w", err)
	}

	return closeErr
}

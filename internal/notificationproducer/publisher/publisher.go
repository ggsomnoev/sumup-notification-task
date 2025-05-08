package publisher

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQPublisher(conn *amqp.Connection, queueName string) (*RabbitMQPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &RabbitMQPublisher{
		conn:    conn,
		channel: ch,
		queue:   queueName,
	}, nil
}

func (r *RabbitMQPublisher) Publish(ctx context.Context, message []byte) error {
	err := r.channel.PublishWithContext(
		ctx,
		"",      // exchange
		r.queue, // routing key (queue name)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (r *RabbitMQPublisher) Close() error {
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

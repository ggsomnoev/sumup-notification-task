package notificationproducer

import (
	"context"

	"github.com/ggsomnoev/sumup-notification-task/internal/notificationproducer/process"
	"github.com/ggsomnoev/sumup-notification-task/internal/notificationproducer/publisher"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Process(
	ctx context.Context,
	srv *echo.Echo,
	rabbitMQConn *amqp.Connection,
	queueName string,
) error {
	publisher, err := publisher.NewRabbitMQPublisher(rabbitMQConn, queueName)
	if err != nil {
		return err
	}
	process.RegisterHandlers(ctx, srv, publisher)
	return nil
}

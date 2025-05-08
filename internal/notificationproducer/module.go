package notificationproducer

import (
	"context"
	"fmt"

	"github.com/ggsomnoev/sumup-notification-task/internal/lifecycle"
	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"github.com/ggsomnoev/sumup-notification-task/internal/notificationproducer/process"
	"github.com/ggsomnoev/sumup-notification-task/internal/notificationproducer/publisher"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Process(
	procSpawnFn lifecycle.ProcessSpawnFunc,
	ctx context.Context,
	srv *echo.Echo,
	rabbitMQConn *amqp.Connection,
	queueName string,
) {
	procSpawnFn(func(ctx context.Context) error {
		pub, err := publisher.NewRabbitMQPublisher(rabbitMQConn, queueName)
		if err != nil {
			return err
		}

		process.RegisterHandlers(ctx, srv, pub)

		<-ctx.Done()
		logger.GetLogger().Info("closing the RabbitMQ connection due to app exit")

		if pub != nil {
			err := pub.Close()
			if err != nil {
				return fmt.Errorf("failed to close RabbitMQ connection: %w", err)
			}
		}

		return nil
	}, "Publisher Starter")
}

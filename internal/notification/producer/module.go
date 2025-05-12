package producer

import (
	"context"
	"fmt"

	"github.com/ggsomnoev/sumup-notification-task/internal/lifecycle"
	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/producer/handler"
	"github.com/labstack/echo/v4"
)

func Process(
	procSpawnFn lifecycle.ProcessSpawnFunc,
	ctx context.Context,
	srv *echo.Echo,
	publisher handler.Publisher,
) {
	procSpawnFn(func(ctx context.Context) error {
		handler.RegisterHandlers(ctx, srv, publisher)

		<-ctx.Done()
		logger.GetLogger().Info("closing the RabbitMQ connection due to app exit")

		if publisher != nil {
			err := publisher.Close()
			if err != nil {
				return fmt.Errorf("failed to close RabbitMQ connection: %w", err)
			}
		}

		return nil
	}, "Publisher")
}

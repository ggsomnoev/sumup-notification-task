package notificationproducer

import (
	"context"

	"github.com/ggsomnoev/sumup-notification-task/internal/notificationproducer/process"
	"github.com/ggsomnoev/sumup-notification-task/internal/notificationproducer/publisher"
	"github.com/labstack/echo/v4"
)

func Process(
	ctx context.Context,
	srv *echo.Echo,
) {
	publisher := publisher.NewPublisher()
	process.RegisterHandlers(ctx, srv, publisher)
}

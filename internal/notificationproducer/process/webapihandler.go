package process

import (
	"context"
	"net/http"

	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"github.com/ggsomnoev/sumup-notification-task/internal/model"
	"github.com/ggsomnoev/sumup-notification-task/internal/validator"
	"github.com/labstack/echo/v4"
)

const successfullyAddedNotification = "successfully added notification"

//counterfeiter:generate . Publisher
type Publisher interface {
	Publish(context.Context, model.Notification) error
}

func RegisterHandlers(ctx context.Context, srv *echo.Echo, publisher Publisher) {
	if srv != nil {
		srv.POST("/notifications", handleNotification(ctx, publisher))
	} else {
		logger.GetLogger().Warn("Running routes without a webapi server, did NOT register routes.")
	}
}

func handleNotification(ctx context.Context, publisher Publisher) echo.HandlerFunc {
	return func(c echo.Context) error {
		var notification model.Notification
		if err := c.Bind(&notification); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		if err := validator.ValidateNotification(notification); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		if err := publisher.Publish(ctx, notification); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": successfullyAddedNotification,
			"channel": notification.Channel,
		})
	}
}

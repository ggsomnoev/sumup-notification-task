package webapi

import (
	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func JSONLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.GetLogger()

		err := next(c)

		log.WithFields(logrus.Fields{
			"method":    c.Request().Method,
			"uri":       c.Request().RequestURI,
			"status":    c.Response().Status,
			"timestamp": time.Now().Format(time.RFC3339),
		}).Info("request")

		return err
	}
}

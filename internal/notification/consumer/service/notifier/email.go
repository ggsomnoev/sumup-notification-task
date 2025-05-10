package notifier

import (
	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
)

type EmailNotifier struct{}

func NewEmailNotifier() *EmailNotifier {
	return &EmailNotifier{}
}

func (es *EmailNotifier) Send(n model.Notification) error {
	logger.GetLogger().Infof("Trying to send email notification: %v", n)
	return nil
}

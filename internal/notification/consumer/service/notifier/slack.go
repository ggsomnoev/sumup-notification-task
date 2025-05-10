package notifier

import (
	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
)

type SlackNotifier struct{}

func NewSlackNotifier() *SlackNotifier {
	return &SlackNotifier{}
}

func (es *SlackNotifier) Send(n model.Notification) error {
	logger.GetLogger().Infof("Trying to send slack notification: %v", n)
	return nil
}

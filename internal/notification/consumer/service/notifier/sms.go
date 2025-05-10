package notifier

import (
	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
)

type SmsNotifier struct{}

func NewSmsNotifier() *SmsNotifier {
	return &SmsNotifier{}
}

func (es *SmsNotifier) Send(n model.Notification) error {
	logger.GetLogger().Infof("Trying to send sms notification: %v", n)
	return nil
}

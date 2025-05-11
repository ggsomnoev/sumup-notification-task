package notifier

import (
	"errors"
	"fmt"

	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
	"github.com/lateralusd/textbelt"
)

var ErrFailedToSendSMS = errors.New("failed to send sms")

//counterfeiter:generate . TextbeltClient
type TextbeltClient interface {
	Quota() (int, error)
	Send(string, string) (string, error)
	Status(string) (textbelt.MessageStatus, error)
}

type SmsNotifier struct {
	client TextbeltClient
}

func NewSmsNotifier(client TextbeltClient) *SmsNotifier {
	return &SmsNotifier{
		client: client,
	}
}

func (sn *SmsNotifier) Send(n model.Notification) error {
	rem, err := sn.client.Quota()
	if err != nil {
		return fmt.Errorf("failed to check message quota: %w", err)
	}
	if rem == 0 {
		logger.GetLogger().Infof("daily quota reacher, remaining messages: %d", rem)
		return nil
	}

	msg, err := sn.client.Send(n.From, n.Message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	status, err := sn.client.Status(msg)
	if err != nil {
		return fmt.Errorf("failed to check message status: %w", err)
	}
	if status == textbelt.StatusFailed {
		return fmt.Errorf("%w (status - %s): %w", ErrFailedToSendSMS, status, err)
	}

	return nil
}

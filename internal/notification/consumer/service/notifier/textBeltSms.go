package notifier

import (
	"fmt"
	"time"

	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
	"github.com/lateralusd/textbelt"
)

var texterTimeout = 2 * time.Second

type TextBeltSmsNotifier struct {
	texter *textbelt.Textbelt
}

func NewTextBeltSmsNotifier() *TextBeltSmsNotifier {
	texter := textbelt.New(
		textbelt.WithKey("textbelt"),
		textbelt.WithTimeout(texterTimeout),
	)
	return &TextBeltSmsNotifier{
		texter: texter,
	}
}

func (tbsn *TextBeltSmsNotifier) Send(n model.Notification) error {
	rem, err := tbsn.texter.Quota()
	if err != nil {
		return fmt.Errorf("failed to check message quota: %v", err)
	}
	if rem == 0 {
		logger.GetLogger().Infof("daily quota reacher, remaining messages: %d", rem)
		return nil
	}

	msg, err := tbsn.texter.Send(n.From, n.Message)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	status, err := tbsn.texter.Status(msg)
	if err != nil {
		return fmt.Errorf("failed to check message status: %v", err)
	}

	logger.GetLogger().WithField(
		"notification", n,
	).Infof("Message %s status is %s", msg, status)

	return nil
}

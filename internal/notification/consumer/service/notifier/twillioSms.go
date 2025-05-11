package notifier

import (
	"fmt"

	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioSmsNotifier struct {
	twilioClient *twilio.RestClient
}

func NewTwilioSmsNotifier(
	twilioClient *twilio.RestClient,
) *TwilioSmsNotifier {
	return &TwilioSmsNotifier{
		twilioClient: twilioClient,
	}
}

func (es *TwilioSmsNotifier) Send(n model.Notification) error {
	logger.GetLogger().Infof("Trying to send sms notification using Twilio...: %v", n)

	params := &openapi.CreateMessageParams{}
	params.SetTo(n.Recipient)
	params.SetFrom(n.From)
	params.SetBody(n.Message)

	// TODO: check the returned resp.
	// Does not work with PH phone numbers.
	_, err := es.twilioClient.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS to %s: %v", n.Recipient, err)
	}

	return nil
}

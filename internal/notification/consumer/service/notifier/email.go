package notifier

import (
	"errors"
	"fmt"

	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var ErrFailedToSendEmail = errors.New("email not sent successfully")

//counterfeiter:generate . SendGridClient
type SendGridClient interface {
	Send(*mail.SGMailV3) (*rest.Response, error)
}

type EmailNotifier struct {
	client         SendGridClient
	senderIdentity string
}

func NewEmailNotifier(client SendGridClient, senderIdenitity string) *EmailNotifier {
	return &EmailNotifier{
		client:         client,
		senderIdentity: senderIdenitity,
	}
}

func (es *EmailNotifier) Send(n model.Notification) error {
	from := mail.NewEmail("Notifier", es.senderIdentity)
	to := mail.NewEmail("Recipient", n.Recipient)
	message := mail.NewSingleEmail(from, n.Subject, to, n.Message, n.Message)

	response, err := es.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email via SendGrid: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("%w, status code: %d, body: %s", ErrFailedToSendEmail, response.StatusCode, response.Body)
	}
	return nil
}

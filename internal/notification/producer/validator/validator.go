package validator

import (
	"errors"
	"regexp"

	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
)

const (
	phoneRegEx = `^\+\d{8,15}$`
	emailRegEx = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
)

var (
	ErrInvalidEmailFormat    = errors.New("invalid email format")
	ErrInvalidPhoneNumFormat = errors.New("invalid phone number format")
	ErrMissingFields         = errors.New("missing required fields")
	ErrUnsupportedChannel    = errors.New("unsupported notification channel")
)

func ValidateNotification(n model.Notification) error {
	switch n.Channel {
	case model.ChannelEmail:
		if n.Subject == "" || n.Message == "" {
			return ErrMissingFields
		}
		if !validateField(n.Recipient, emailRegEx) {
			return ErrInvalidEmailFormat
		}
	case model.ChannelSMS:
		if n.Message == "" {
			return ErrMissingFields
		}
		if !validateField(n.Recipient, phoneRegEx) {
			return ErrInvalidPhoneNumFormat
		}
	case model.ChannelSlack:
		if n.Subject == "" || n.Message == "" {
			return ErrMissingFields
		}
	default:
		return ErrUnsupportedChannel
	}
	return nil
}

func validateField(recipient string, regexps string) bool {
	regex := regexp.MustCompile(regexps)
	return regex.MatchString(recipient)
}

package validator

import (
	"errors"
	"regexp"

	"github.com/ggsomnoev/sumup-notification-task/internal/model"
	"github.com/go-playground/validator/v10"
)

const (
	phoneRegEx = `^\+\d{8,15}$`
	slackRegEx = `^[@#][a-zA-Z0-9._-]+$`
	emailRegEx = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateNotification(n model.Notification) error {
	if err := validate.Struct(n); err != nil {
		return err
	}

	switch n.Channel {
	case "email":
		if !validateField(n.Recipient, emailRegEx) {
			return errors.New("invalid email format")
		}
	case "sms":
		if !validateField(n.Recipient, phoneRegEx) {
			return errors.New("invalid phone number")
		}
	case "slack":
		if !validateField(n.Recipient, slackRegEx) {
			return errors.New("invalid slack recipient")
		}
	default:
		return errors.New("unsupported channel")
	}

	return nil
}

func validateField(recipient string, regexps string) bool {
	regex := regexp.MustCompile(regexps)
	return regex.MatchString(recipient)
}

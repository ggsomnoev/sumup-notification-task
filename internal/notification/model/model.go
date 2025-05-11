package model

import "github.com/google/uuid"

type Notification struct {
	From      string `json:"from" validate:"required"`
	Channel   string `json:"channel" validate:"required,oneof=email sms slack"`
	Recipient string `json:"recipient" validate:"required"`
	Subject   string `json:"subject,omitempty"`
	Message   string `json:"message" validate:"required"`
}

type Message struct {
	UUID uuid.UUID `json:"uuid"`
	Notification
}

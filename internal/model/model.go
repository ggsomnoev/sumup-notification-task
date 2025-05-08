package model

type Notification struct {
	Channel   string `json:"channel" validate:"required,oneof=email sms slack"`
	Recipient string `json:"recipient" validate:"required"`
	Subject   string `json:"subject,omitempty"`
	Message   string `json:"message" validate:"required"`
}

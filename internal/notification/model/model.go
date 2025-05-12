package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type ChannelType string

const (
	ChannelSlack ChannelType = "slack"
	ChannelEmail ChannelType = "email"
	ChannelSMS   ChannelType = "sms"
)

var validChannels = map[string]ChannelType{
	"slack": ChannelSlack,
	"email": ChannelEmail,
	"sms":   ChannelSMS,
}

func (c *ChannelType) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("invalid channel type: %w", err)
	}

	normalized := strings.ToLower(raw)
	val, ok := validChannels[normalized]
	if !ok {
		return fmt.Errorf("unsupported channel type: %s", raw)
	}

	*c = val
	return nil
}

type Notification struct {
	Channel   ChannelType `json:"channel"`
	Recipient string      `json:"recipient"`
	Subject   string      `json:"subject"`
	Message   string      `json:"message"`
}

type Message struct {
	UUID uuid.UUID `json:"uuid"`
	Notification
}

// TODO: add custom marshal/unmarshal funcs to store the message type as enum val.
// String values tend to change with time change.

package publisher

import (
	"context"

	"github.com/ggsomnoev/sumup-notification-task/internal/model"
)

type Publisher struct{}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Publish(_ context.Context, _ model.Notification) error {
	return nil
}

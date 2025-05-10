package service

import (
	"context"
	"fmt"

	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
	"github.com/google/uuid"
)

type Store interface {
	AddMessage(context.Context, model.Message) error
	MarkCompleted(context.Context, uuid.UUID) error
	MessageExists(context.Context, uuid.UUID) (bool, error)
	RunInAtomically(context.Context, func(context.Context) error) error
}

type Notifier interface {
	Send(n model.Notification) error
}

type Service struct {
	store     Store
	notifiers map[string]Notifier
}

func NewService(store Store, notifiers map[string]Notifier) *Service {
	return &Service{
		store:     store,
		notifiers: notifiers,
	}
}

func (s *Service) Send(ctx context.Context, message model.Message) error {
	notifier, ok := s.notifiers[message.Channel]
	if !ok {
		return fmt.Errorf("no notifier registered for channel: %s", message.Channel)
	}

	return s.store.RunInAtomically(ctx, func(ctx context.Context) error {
		exists, err := s.store.MessageExists(ctx, message.UUID)
		if err != nil {
			return fmt.Errorf("failed to check message: %w", err)
		}
		if exists {
			logger.GetLogger().Infof("Skipping notification already in progress for UUID: %s", message.UUID)
			return nil
		}

		logger.GetLogger().Infof("Sending notification via %s to %s", message.Channel, message.Recipient)
		if err := s.store.AddMessage(ctx, message); err != nil {
			return fmt.Errorf("failed to persist notification event: %w", err)
		}

		if err := notifier.Send(message.Notification); err != nil {
			return fmt.Errorf("sending notification failed: %w", err)
		}

		if err := s.store.MarkCompleted(ctx, message.UUID); err != nil {
			return fmt.Errorf("failed to mark event as completed: %w", err)
		}

		return nil
	})
}

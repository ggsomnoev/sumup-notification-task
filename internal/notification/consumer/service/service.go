package service

import (
	"context"
	"fmt"

	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
	"github.com/google/uuid"
)

//counterfeiter:generate . Store
type Store interface {
	AddMessage(context.Context, model.Message) error
	MarkCompleted(context.Context, uuid.UUID) error
	MessageExists(context.Context, uuid.UUID) (bool, error)
	RunInAtomically(context.Context, func(context.Context) error) error
}

//counterfeiter:generate . Notifier
type Notifier interface {
	Send(n model.Notification) error
}

type Service struct {
	store     Store
	notifiers map[model.ChannelType]Notifier
}

func NewService(
	store Store,
	notifiers map[model.ChannelType]Notifier,
) *Service {
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

		if err := s.store.AddMessage(ctx, message); err != nil {
			return fmt.Errorf("failed to persist notification event: %w", err)
		}

		logger.GetLogger().Infof("Trying to send notification via %s...", message.Channel)
		if err := notifier.Send(message.Notification); err != nil {
			return fmt.Errorf("sending notification failed: %w", err)
		}
		logger.GetLogger().Info("Message send successfully!")

		if err := s.store.MarkCompleted(ctx, message.UUID); err != nil {
			return fmt.Errorf("failed to mark event as completed: %w", err)
		}

		return nil
	})
}

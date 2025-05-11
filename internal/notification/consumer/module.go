package consumer

import (
	"context"
	"fmt"

	"github.com/ggsomnoev/sumup-notification-task/internal/lifecycle"
	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/consumer/service"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/consumer/service/notifier"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/consumer/store"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Consumer interface {
	Consume(context.Context, func(context.Context, model.Message) error) error
	Close() error
}

func Process(
	procSpawnFn lifecycle.ProcessSpawnFunc,
	ctx context.Context,
	pool *pgxpool.Pool,
	consumer Consumer,
	smsClient notifier.TextbeltClient,
	mailClient notifier.SendGridClient,
) {
	procSpawnFn(func(ctx context.Context) error {
		store := store.NewStore(pool)

		senders := map[string]service.Notifier{
			"email": notifier.NewEmailNotifier(mailClient),
			"sms":   notifier.NewSmsNotifier(smsClient),
			"slack": notifier.NewSlackNotifier(),
		}

		notificationSvc := service.NewService(store, senders)
		err := consumer.Consume(ctx, notificationSvc.Send)
		if err != nil {
			return fmt.Errorf("consume failed: %w", err)
		}

		<-ctx.Done()
		logger.GetLogger().Info("closing the RabbitMQ connection due to app exit")

		if err := consumer.Close(); err != nil {
			return fmt.Errorf("failed to close consumer: %w", err)
		}

		return nil
	}, "Consumer")
}

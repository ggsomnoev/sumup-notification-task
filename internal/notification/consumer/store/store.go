package store

import (
	"context"
	"fmt"
	"time"

	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const NotificationEventsTable = "notification_events"

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) AddEvent(ctx context.Context, m model.Message) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (uuid, payload, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (uuid) DO NOTHING
	`, NotificationEventsTable)
	_, err := s.pool.Exec(ctx, query, m.UUID, m, time.Now().UTC())
	return err
}

func (s *Store) MarkCompleted(ctx context.Context, uuid uuid.UUID) error {
	query := fmt.Sprintf(`UPDATE %s SET completed_at = $1 WHERE uuid = $2`, NotificationEventsTable)
	_, err := s.pool.Exec(ctx, query, time.Now().UTC(), uuid)
	return err
}

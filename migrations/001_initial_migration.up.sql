BEGIN;

CREATE TABLE IF NOT EXISTS notification_events (    
    uuid UUID PRIMARY KEY,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

COMMIT;
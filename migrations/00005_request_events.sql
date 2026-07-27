-- +goose Up
CREATE TABLE request_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id UUID NOT NULL REFERENCES requests(id) ON DELETE CASCADE,
    actor_id UUID REFERENCES users(id),
    event_type TEXT NOT NULL,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX request_events_request_id_created_at_idx
    ON request_events(request_id, created_at);

-- +goose Down
DROP TABLE request_events;

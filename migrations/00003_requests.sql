-- +goose Up
CREATE TABLE requests(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reference TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    client_name TEXT NOT NULL,
    client_email TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);
CREATE INDEX requests_status_idx
ON requests(status);
CREATE INDEX requests_created_at_idx
ON requests(created_at);
CREATE INDEX requests_client_email_idx
ON requests(client_email);
-- +goose Down
DROP TABLE requests;

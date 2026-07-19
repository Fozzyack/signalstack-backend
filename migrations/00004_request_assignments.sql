-- +goose Up
CREATE TABLE request_assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id UUID NOT NULL REFERENCES requests(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK (role IN ('lead', 'contributor')),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    unassigned_at TIMESTAMPTZ,
    personal_deadline TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE request_assignments;

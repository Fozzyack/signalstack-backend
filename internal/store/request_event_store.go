package store

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Fozzyack/signalstack-backend/internal/models"
)

type RequestEventStore interface {
	GetRequestEventsByRequestID(ctx context.Context, requestID string) ([]*models.RequestEvent, error)
	CreateRequestEvent(ctx context.Context, requestID string, actorID *string, eventType string, metadata json.RawMessage) (*models.RequestEvent, error)
}

func NewRequestEventStore(db *sql.DB) RequestEventStore {
	return NewPostgresStore(db)
}

const requestEventColumns = `
	id, request_id, actor_id, event_type, metadata, created_at
`

func scanRequestEvent(scanner interface{ Scan(...any) error }) (*models.RequestEvent, error) {
	event := &models.RequestEvent{}
	var actorID sql.NullString
	var metadata []byte

	err := scanner.Scan(
		&event.ID,
		&event.RequestID,
		&actorID,
		&event.EventType,
		&metadata,
		&event.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if actorID.Valid {
		event.ActorID = &actorID.String
	}
	if metadata != nil {
		event.Metadata = json.RawMessage(metadata)
	}

	return event, nil
}

func (ps *PostgresStore) GetRequestEventsByRequestID(ctx context.Context, requestID string) ([]*models.RequestEvent, error) {
	query := `
		SELECT ` + requestEventColumns + `
		FROM request_events
		WHERE request_id = $1
		ORDER BY created_at DESC
	`

	rows, err := ps.db.QueryContext(ctx, query, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*models.RequestEvent, 0)
	for rows.Next() {
		event, err := scanRequestEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (ps *PostgresStore) CreateRequestEvent(ctx context.Context, requestID string, actorID *string, eventType string, metadata json.RawMessage) (*models.RequestEvent, error) {
	var metadataValue any
	if len(metadata) > 0 {
		metadataValue = metadata
	}

	query := `
		INSERT INTO request_events (request_id, actor_id, event_type, metadata)
		VALUES ($1, $2, $3, $4)
		RETURNING ` + requestEventColumns + `
	`

	return scanRequestEvent(ps.db.QueryRowContext(
		ctx,
		query,
		requestID,
		nullableString(actorID),
		eventType,
		metadataValue,
	))
}

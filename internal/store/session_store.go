package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/Fozzyack/signalstack-backend/internal/models"
)

type SessionStore interface {
	CreateSession(ctx context.Context, userId string, token string) (*models.Session, error)
}

func NewSessionStore(db *sql.DB) SessionStore {
	return NewPostgresStore(db)
}

func (ps *PostgresStore) CreateSession(ctx context.Context, userId string, token string) (*models.Session, error) {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	expiresAt := time.Now().Add(24 * time.Hour)
	query := `
		INSERT INTO sessions (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, token, expires_at, created_at
	`

	createdSession := &models.Session{}
	err = tx.QueryRowContext(ctx, query, userId, token, expiresAt).Scan(
		&createdSession.ID,
		&createdSession.UserID,
		&createdSession.Token,
		&createdSession.ExpiresAt,
		&createdSession.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return createdSession, nil
}

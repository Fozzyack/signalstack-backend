package store

import (
	"context"
	"database/sql"

	"github.com/Fozzyack/signalstack-backend/internal/models"
)

type UserStore interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

func NewUserStore(db *sql.DB) UserStore {
	return NewPostgresStore(db)
}

func (ps *PostgresStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {

	query := `
		SELECT id, name, email, password_hash, created_at, updated_at FROM users
		WHERE email=$1
	`

	var foundUser models.User
	err := ps.db.QueryRowContext(ctx, query, email).Scan(
		&foundUser.ID,
		&foundUser.Name,
		&foundUser.Email,
		&foundUser.PasswordHash,
		&foundUser.CreatedAt,
		&foundUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &foundUser, nil

}

func (ps *PostgresStore) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at FROM users
		WHERE id = $1
	`

	foundUser := &models.User{}
	err := ps.db.QueryRowContext(ctx, query, id).Scan(
		&foundUser.ID,
		&foundUser.Name,
		&foundUser.Email,
		&foundUser.PasswordHash,
		&foundUser.CreatedAt,
		&foundUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

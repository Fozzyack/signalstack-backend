package store

import (
	"context"
	"database/sql"

	"github.com/Fozzyack/signalstack-backend/internal/models"
)

type UserStore interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, name, email, passwordHash string) (*models.User, error)
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

func (ps *PostgresStore) UpdateUser(ctx context.Context, id string, name, email, passwordHash string) (*models.User, error) {

	query := `
		UPDATE users
		SET name = $2, email = $3, password_hash = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, email, password_hash, created_at, updated_at
	`

	updatedUser := &models.User{}
	err := ps.db.QueryRowContext(ctx, query, id, name, email, passwordHash).Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.Email,
		&updatedUser.PasswordHash,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

package store

import (
	"context"
	"database/sql"

	"github.com/Fozzyack/signalstack-backend/internal/models"
)

type RequestStore interface {
	GetRequestById(ctx context.Context, id string) (*models.Request, error)
	GetRequests(ctx context.Context) ([]*models.Request, error)
	CreateRequest(ctx context.Context, title, description, clientName, clientEmail string) (*models.Request, error)
	UpdateRequest(ctx context.Context, id, reference, title, description, clientName, clientEmail, status, resolvedAt string) (*models.Request, error)
	DeleteRequest(ctx context.Context, id string) error
}

func NewRequestStore(db *sql.DB) RequestStore {
	return NewPostgresStore(db)
}

func (ps *PostgresStore) GetRequestById(ctx context.Context, id string) (*models.Request, error) {
	query := `
		SELECT id, reference, title, description, client_name, client_email,
			status, created_at, updated_at, resolved_at
		FROM requests
		WHERE id = $1
	`

	request := &models.Request{}
	err := ps.db.QueryRowContext(ctx, query, id).Scan(
		&request.ID,
		&request.Reference,
		&request.Title,
		&request.Description,
		&request.ClientName,
		&request.ClientEmail,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
		&request.ResolvedAt,
	)
	if err != nil {
		return nil, err
	}

	request.Assignments, err = ps.GetRequestAssignmentsByRequestID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (ps *PostgresStore) GetRequests(ctx context.Context) ([]*models.Request, error) {
	query := `
		SELECT id, reference, title, description, client_name, client_email,
			status, created_at, updated_at, resolved_at
		FROM requests
		ORDER BY created_at DESC
	`

	rows, err := ps.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := make([]*models.Request, 0)
	for rows.Next() {
		request := &models.Request{}
		err := rows.Scan(
			&request.ID,
			&request.Reference,
			&request.Title,
			&request.Description,
			&request.ClientName,
			&request.ClientEmail,
			&request.Status,
			&request.CreatedAt,
			&request.UpdatedAt,
			&request.ResolvedAt,
		)
		if err != nil {
			return nil, err
		}

		request.Assignments, err = ps.GetRequestAssignmentsByRequestID(ctx, request.ID)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (ps *PostgresStore) CreateRequest(ctx context.Context, title, description, clientName, clientEmail string) (*models.Request, error) {

	query := `
		INSERT INTO requests (
			title, description, client_name, client_email
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, reference, title, description, client_name, client_email,
			status, created_at, updated_at, resolved_at
	`

	request := &models.Request{}
	err := ps.db.QueryRowContext(
		ctx,
		query,
		title,
		description,
		clientName,
		clientEmail,
	).Scan(
		&request.ID,
		&request.Reference,
		&request.Title,
		&request.Description,
		&request.ClientName,
		&request.ClientEmail,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
		&request.ResolvedAt,
	)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (ps *PostgresStore) UpdateRequest(ctx context.Context, id, reference, title, description, clientName, clientEmail, status, resolvedAt string) (*models.Request, error) {
	var resolvedAtValue any
	if resolvedAt != "" {
		resolvedAtValue = resolvedAt
	}

	query := `
		UPDATE requests
		SET reference = $1,
			title = $2,
			description = $3,
			client_name = $4,
			client_email = $5,
			status = $6,
			updated_at = NOW(),
			resolved_at = $7
		WHERE id = $8
		RETURNING id, reference, title, description, client_name, client_email,
			status, created_at, updated_at, resolved_at
	`

	request := &models.Request{}
	err := ps.db.QueryRowContext(
		ctx,
		query,
		reference,
		title,
		description,
		clientName,
		clientEmail,
		status,
		resolvedAtValue,
		id,
	).Scan(
		&request.ID,
		&request.Reference,
		&request.Title,
		&request.Description,
		&request.ClientName,
		&request.ClientEmail,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
		&request.ResolvedAt,
	)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (ps *PostgresStore) DeleteRequest(ctx context.Context, id string) error {
	result, err := ps.db.ExecContext(ctx, `DELETE FROM requests WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

package store

import (
	"context"
	"database/sql"

	"github.com/Fozzyack/signalstack-backend/internal/models"
)

type RequestsAssignmentsStore interface {
	GetRequestAssignmentByID(ctx context.Context, id string) (*models.RequestAssignment, error)
	GetRequestAssignmentsByRequestID(ctx context.Context, requestID string) ([]*models.RequestAssignment, error)
	GetRequestAssignmentsByUserID(ctx context.Context, userID string) ([]*models.RequestAssignment, error)
	CreateRequestAssignment(ctx context.Context, requestID, userID, role string) (*models.RequestAssignment, error)
	UpdateRequestAssignment(ctx context.Context, id, role string, unassignedAt, personalDeadline, completedAt *string) (*models.RequestAssignment, error)
	UnassignRequestAssignment(ctx context.Context, id string) error
	DeleteRequestAssignment(ctx context.Context, id string) error
}

func NewRequestsAssignmentsStore(db *sql.DB) RequestsAssignmentsStore {
	return NewPostgresStore(db)
}

const requestAssignmentColumns = `
	id, request_id, user_id, role, assigned_at, unassigned_at,
	personal_deadline, completed_at
`

func scanRequestAssignment(scanner interface{ Scan(...any) error }) (*models.RequestAssignment, error) {
	assignment := &models.RequestAssignment{}
	err := scanner.Scan(
		&assignment.ID,
		&assignment.RequestID,
		&assignment.UserID,
		&assignment.Role,
		&assignment.AssignedAt,
		&assignment.UnassignedAt,
		&assignment.PersonalDeadline,
		&assignment.CompletedAt,
	)
	if err != nil {
		return nil, err
	}
	return assignment, nil
}

func (ps *PostgresStore) GetRequestAssignmentByID(ctx context.Context, id string) (*models.RequestAssignment, error) {
	query := `SELECT ` + requestAssignmentColumns + ` FROM request_assignments WHERE id = $1`
	return scanRequestAssignment(ps.db.QueryRowContext(ctx, query, id))
}

func (ps *PostgresStore) GetRequestAssignmentsByRequestID(ctx context.Context, requestID string) ([]*models.RequestAssignment, error) {
	query := `
		SELECT ` + requestAssignmentColumns + `
		FROM request_assignments
		WHERE request_id = $1
		ORDER BY assigned_at DESC
	`
	return ps.getRequestAssignments(ctx, query, requestID)
}

func (ps *PostgresStore) GetRequestAssignmentsByUserID(ctx context.Context, userID string) ([]*models.RequestAssignment, error) {
	query := `
		SELECT ` + requestAssignmentColumns + `
		FROM request_assignments
		WHERE user_id = $1 AND unassigned_at IS NULL
		ORDER BY assigned_at DESC
	`
	return ps.getRequestAssignments(ctx, query, userID)
}

func (ps *PostgresStore) getRequestAssignments(ctx context.Context, query, argument string) ([]*models.RequestAssignment, error) {
	rows, err := ps.db.QueryContext(ctx, query, argument)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := make([]*models.RequestAssignment, 0)
	for rows.Next() {
		assignment, err := scanRequestAssignment(rows)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, assignment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return assignments, nil
}

func (ps *PostgresStore) CreateRequestAssignment(ctx context.Context, requestID, userID, role string) (*models.RequestAssignment, error) {
	query := `
		INSERT INTO request_assignments (request_id, user_id, role)
		VALUES ($1, $2, $3)
		RETURNING ` + requestAssignmentColumns + `
	`
	return scanRequestAssignment(ps.db.QueryRowContext(ctx, query, requestID, userID, role))
}

func (ps *PostgresStore) UpdateRequestAssignment(ctx context.Context, id, role string, unassignedAt, personalDeadline, completedAt *string) (*models.RequestAssignment, error) {
	query := `
		UPDATE request_assignments
		SET role = $1, unassigned_at = $2, personal_deadline = $3, completed_at = $4
		WHERE id = $5
		RETURNING ` + requestAssignmentColumns + `
	`
	return scanRequestAssignment(ps.db.QueryRowContext(
		ctx,
		query,
		role,
		nullableString(unassignedAt),
		nullableString(personalDeadline),
		nullableString(completedAt),
		id,
	))
}

func (ps *PostgresStore) UnassignRequestAssignment(ctx context.Context, id string) error {
	result, err := ps.db.ExecContext(ctx, `
		UPDATE request_assignments
		SET unassigned_at = NOW()
		WHERE id = $1 AND unassigned_at IS NULL
	`, id)
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

func (ps *PostgresStore) DeleteRequestAssignment(ctx context.Context, id string) error {
	result, err := ps.db.ExecContext(ctx, `DELETE FROM request_assignments WHERE id = $1`, id)
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

func nullableString(value *string) any {
	if value == nil || *value == "" {
		return nil
	}
	return *value
}

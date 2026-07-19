package models

import "time"

type Request struct {
	ID          string               `json:"id"`
	Reference   string               `json:"reference"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	ClientName  string               `json:"client_name"`
	ClientEmail string               `json:"client_email"`
	Status      string               `json:"status"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	ResolvedAt  *time.Time           `json:"resolved_at,omitempty"`
	Assignments []*RequestAssignment `json:"assignments"`
}

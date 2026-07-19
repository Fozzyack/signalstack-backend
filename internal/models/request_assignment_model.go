package models

import "time"

type RequestAssignment struct {
	ID               string     `json:"id"`
	RequestID        string     `json:"request_id"`
	UserID           string     `json:"user_id"`
	Role             string     `json:"role"`
	AssignedAt       time.Time  `json:"assigned_at"`
	UnassignedAt     *time.Time `json:"unassigned_at,omitempty"`
	PersonalDeadline *time.Time `json:"personal_deadline,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
}

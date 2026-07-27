package models

import (
	"encoding/json"
	"time"
)

type RequestEvent struct {
	ID        string          `json:"id"`
	RequestID string          `json:"request_id"`
	ActorID   *string         `json:"actor_id,omitempty"`
	EventType string          `json:"event_type"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

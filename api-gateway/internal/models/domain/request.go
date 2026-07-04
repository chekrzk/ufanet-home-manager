package domain

import "time"

type Request struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	AssignedTo  string    `json:"assigned_to,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

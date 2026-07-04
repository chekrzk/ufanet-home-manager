package domain

import "time"

type Request struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Category      string    `json:"category"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	AssignedTo    string    `json:"assigned_to,omitempty"`
	PreferredDate string    `json:"preferred_date,omitempty"`
	Address       string    `json:"address,omitempty"`
	Apartment     string    `json:"apartment,omitempty"`
	Phone         string    `json:"phone,omitempty"`
	AcceptedAt    time.Time `json:"accepted_at,omitempty"`
	DeclinedAt    time.Time `json:"declined_at,omitempty"`
	CompletedAt   time.Time `json:"completed_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

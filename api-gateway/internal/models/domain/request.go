package domain

import "time"

type Request struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

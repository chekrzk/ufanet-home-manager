package domain

import "time"

type Device struct {
	UserID   string `json:"user_id"`
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id,omitempty"`
	HouseID   string    `json:"house_id,omitempty"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	EntityID  string    `json:"entity_id,omitempty"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

package domain

import "time"

type News struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	HouseID   string    `json:"house_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

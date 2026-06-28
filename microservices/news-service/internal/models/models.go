package models

import "time"

type UserContext struct {
	UserID string
	Role   string
}

type Pagination struct {
	Page  int
	Limit int
}

type News struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title     string `gorm:"size:255;not null"`
	Body      string `gorm:"not null"`
	HouseID   string `gorm:"size:64;index"`
	AuthorID  string `gorm:"type:uuid;index;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewsFilter struct {
	Actor      UserContext
	Pagination Pagination
	DateFrom   string
	DateTo     string
}

type NewsPage struct {
	Items []News
	Page  int
	Limit int
	Total int64
}

type CreateNewsCommand struct {
	Author  UserContext
	Title   string
	Body    string
	HouseID string
}

type NotificationEvent struct {
	UserID   string
	HouseID  string
	Type     string
	Title    string
	Body     string
	EntityID string
}

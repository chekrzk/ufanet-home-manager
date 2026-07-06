package models

import "time"

type UserContext struct {
	UserID string
	Role   string
}

type Device struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string `gorm:"type:uuid;index;not null"`
	Token     string `gorm:"size:512;uniqueIndex;not null"`
	Platform  string `gorm:"size:32;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Notification struct {
	ID        string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    *string `gorm:"type:uuid;index"`
	HouseID   string  `gorm:"size:64;index"`
	Type      string  `gorm:"size:64;not null;index"`
	Title     string  `gorm:"size:255;not null"`
	Body      string  `gorm:"not null"`
	EntityID  string  `gorm:"size:64;index"`
	Read      bool    `gorm:"not null;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RegisterDeviceCommand struct {
	User     UserContext
	Token    string
	Platform string
}

type UnregisterDeviceCommand struct {
	User  UserContext
	Token string
}

type PublishNotificationCommand struct {
	UserID   string
	HouseID  string
	Type     string
	Title    string
	Body     string
	EntityID string
}

type Pagination struct {
	Page  int
	Limit int
}

type ListNotificationsCommand struct {
	User       UserContext
	Pagination Pagination
}

type NotificationsPage struct {
	Items []Notification
	Page  int
	Limit int
	Total int64
}

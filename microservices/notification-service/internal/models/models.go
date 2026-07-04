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

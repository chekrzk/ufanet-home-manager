package models

import "time"

type UserContext struct {
	UserID string
	Role   string
}

type Profile struct {
	UserID    string `gorm:"type:uuid;primaryKey"`
	FullName  string `gorm:"size:255;not null"`
	HouseID   string `gorm:"size:64;index"`
	Apartment string `gorm:"size:32"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Worker struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         string `gorm:"type:uuid;uniqueIndex"`
	FullName       string `gorm:"size:255;not null"`
	Specialization string `gorm:"size:64;not null;index"`
	Phone          string `gorm:"size:32"`
	HouseID        string `gorm:"size:64;index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UpdateProfileCommand struct {
	Actor     UserContext
	FullName  string
	HouseID   string
	Apartment string
}

type AddWorkerCommand struct {
	Actor          UserContext
	UserID         string
	FullName       string
	Specialization string
	Phone          string
	HouseID        string
}

type ListWorkersFilter struct {
	Actor   UserContext
	HouseID string
}

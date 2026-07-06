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
	UserID         string `gorm:"size:64;index"`
	FullName       string `gorm:"size:255;not null"`
	Specialization string `gorm:"size:64;not null;index"`
	Phone          string `gorm:"size:32"`
	HouseID        string `gorm:"size:64;index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type WorkerAvailability struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WorkerID       string `gorm:"type:uuid;index"`
	UserID         string `gorm:"size:64;index;not null"`
	Specialization string `gorm:"size:64;not null;index"`
	HouseID        string `gorm:"size:64;index"`
	AvailableDate  string `gorm:"size:32;index;not null"`
	AvailableTime  string `gorm:"size:32;not null"`
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

type SetWorkerAvailabilityCommand struct {
	Worker          UserContext
	Specialization  string
	HouseID         string
	AvailableDate   string
	AvailableTime   string
}

type ListWorkerAvailabilityFilter struct {
	Actor          UserContext
	Specialization string
	HouseID        string
	AvailableDate  string
}

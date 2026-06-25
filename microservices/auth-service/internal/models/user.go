package models

import "time"

type Role string

const (
	RoleResident Role = "resident"
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
)

type User struct {
	ID           string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Phone        string `gorm:"uniqueIndex;size:32;not null"`
	PasswordHash string `gorm:"not null"`
	FullName     string `gorm:"size:255;not null"`
	Role         Role   `gorm:"type:varchar(32);not null;default:'resident'"`
	HouseID      string `gorm:"size:64"`
	Apartment    string `gorm:"size:32"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

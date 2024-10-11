// internal/model/user.go

package model

import (
	"time"
)

type UserType string

const (
	Employee          UserType = "employee"
	UserTypeCandidate UserType = "candidate" // Renamed to avoid conflict
)

type User struct {
	ID        uint     `gorm:"primaryKey"`
	UserType  UserType `gorm:"type:varchar(20)"`
	Username  string   `gorm:"unique;not null"`
	Password  string   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

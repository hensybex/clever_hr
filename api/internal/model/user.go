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
	TgID      string   `gorm:"uniqueIndex"`
	UserType  UserType `gorm:"type:varchar(20)"` // Change to varchar
	CreatedAt time.Time
	UpdatedAt time.Time
}

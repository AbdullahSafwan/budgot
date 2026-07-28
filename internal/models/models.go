package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username string `gorm:"size:100;not null"`
	Email    string `gorm:"size:255;uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

type Session struct {
	gorm.Model

	UserID        int       `gorm:"not null"`
	ExpiresAt     time.Time `gorm:"not null"`
	LastSeen      time.Time `gorm:"not null"`
	IPAddress     string    `gorm:"size:45;not null"`
	UserAgentHash string    `gorm:"size:64;not null"`
}

// type Account struct {
// 	gorm.Model

// 	Name string `gorm:"size:100;not null"`
// 	Balance float64 `gorm:"not null"`
// 	UserID int `gorm:"not null"`
// }

// type Category strucy {

// }

package models

import (
	"time"
)

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Email       string    `json:"email" gorm:"unique;not null"`
	Password    string    `json:"-" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	WalletAddr  string    `json:"wallet_address" gorm:"unique"`
	Role        string    `json:"role" gorm:"default:'researcher'"` // researcher, reviewer, admin
	Institution string    `json:"institution"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Papers  []Paper  `json:"papers,omitempty" gorm:"foreignKey:OwnerID"`
	Reviews []Review `json:"reviews,omitempty" gorm:"foreignKey:ReviewerID"`
}

package models

import (
	"time"

	"gorm.io/datatypes"
)

type Paper struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Title      string         `json:"title" gorm:"not null"`
	Abstract   string         `json:"abstract"`
	Authors    datatypes.JSON `json:"authors" gorm:"type:json"`
	Keywords   datatypes.JSON `json:"keywords" gorm:"type:json"`
	Category   string         `json:"category"`
	IPFSHash   string         `json:"ipfs_hash"`
	NFTTokenID *uint          `json:"nft_token_id"`
	OwnerID    uint           `json:"owner_id"`
	Status     string         `json:"status" gorm:"default:'draft'"` // draft, submitted, under_review, published
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`

	// Relationships
	Owner   User     `json:"owner" gorm:"foreignKey:OwnerID"`
	Reviews []Review `json:"reviews,omitempty" gorm:"foreignKey:PaperID"`
}

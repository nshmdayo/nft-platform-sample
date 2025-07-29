package models

import (
	"time"

	"gorm.io/datatypes"
)

type Review struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	PaperID        uint           `json:"paper_id"`
	ReviewerID     uint           `json:"reviewer_id"`
	Score          int            `json:"score" gorm:"check:score >= 1 AND score <= 10"`
	Comment        string         `json:"comment"`
	Recommendation string         `json:"recommendation"`                  // accept, reject, revision
	Status         string         `json:"status" gorm:"default:'pending'"` // pending, completed, rejected
	Metadata       datatypes.JSON `json:"metadata"`
	NFTTokenID     *uint          `json:"nft_token_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`

	// Relationships
	Paper    Paper `json:"paper" gorm:"foreignKey:PaperID"`
	Reviewer User  `json:"reviewer" gorm:"foreignKey:ReviewerID"`
}

type ReviewMetadata struct {
	ReviewCriteria map[string]interface{} `json:"review_criteria"`
	ReviewType     string                 `json:"review_type"`
	Version        int                    `json:"version"`
}

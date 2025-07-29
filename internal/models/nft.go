package models

import (
	"time"
)

type NFTMetadata struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TokenID     uint      `json:"token_id" gorm:"unique"`
	Type        string    `json:"type"`         // paper, review
	ReferenceID uint      `json:"reference_id"` // Paper ID or Review ID
	MetadataURI string    `json:"metadata_uri"` // IPFS URI
	TxHash      string    `json:"tx_hash"`      // Transaction hash
	CreatedAt   time.Time `json:"created_at"`
}

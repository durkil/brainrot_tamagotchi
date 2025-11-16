package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	WalletAddress string         `gorm:"uniqueIndex;not null" json:"wallet_address"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	NFTs []NFT `gorm:"foreignKey:OwnerAddress;references:WalletAddress" json:"nfts,omitempty"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}


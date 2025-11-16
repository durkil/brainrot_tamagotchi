package models

import (
	"time"

	"gorm.io/gorm"
)

// MarketListing represents an NFT listing on the marketplace
type MarketListing struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	TokenID       uint           `gorm:"uniqueIndex;not null" json:"token_id"`
	SellerAddress string         `gorm:"index;not null" json:"seller_address"`
	Price         float64        `json:"price"` // Price in BASE (ETH)
	IsActive      bool           `gorm:"default:true;index" json:"is_active"`
	ListedAt      time.Time      `json:"listed_at"`
	SoldAt        *time.Time     `json:"sold_at,omitempty"`
	BuyerAddress  *string        `json:"buyer_address,omitempty"`
	TxHash        string         `json:"tx_hash"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations (without constraint to avoid circular dependency during migration)
	// NFT *NFT `gorm:"foreignKey:TokenID;references:TokenID" json:"nft,omitempty"`
}

// TableName overrides the table name
func (MarketListing) TableName() string {
	return "market_listings"
}


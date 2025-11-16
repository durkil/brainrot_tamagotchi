package models

import (
	"time"

	"gorm.io/gorm"
)

// NFT represents a Brainrot NFT token
type NFT struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	TokenID      uint           `gorm:"uniqueIndex;not null" json:"token_id"`
	OwnerAddress string         `gorm:"index;not null" json:"owner_address"`
	MemeType     string         `json:"meme_type"`     // "pepe", "doge", etc.
	Rarity       string         `json:"rarity"`        // "common", "rare", "epic", "legendary"
	Level        int            `gorm:"default:1" json:"level"`
	ColorVariant int            `json:"color_variant"` // 0-4
	TokenURI     string         `json:"token_uri"`
	
	// Tamagotchi stats
	Hunger       int       `gorm:"default:100" json:"hunger"`       // 0-100
	Mood         int       `gorm:"default:100" json:"mood"`         // 0-100
	Energy       int       `gorm:"default:100" json:"energy"`       // 0-100
	LastFed      time.Time `json:"last_fed"`
	LastPlayed   time.Time `json:"last_played"`
	LastInteract time.Time `json:"last_interact"`
	
	// Blockchain data
	TxHash    string    `json:"tx_hash"`
	MintedAt  time.Time `json:"minted_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relations (without constraint to avoid circular dependency)
	// Listing *MarketListing `gorm:"foreignKey:TokenID;references:TokenID;constraint:OnDelete:CASCADE" json:"listing,omitempty"`
}

// TableName overrides the table name
func (NFT) TableName() string {
	return "nfts"
}

// IsAlive checks if the tamagotchi is alive
func (n *NFT) IsAlive() bool {
	return n.Hunger > 0 && n.Mood > 0 && n.Energy > 0
}

// NeedsFeeding checks if feeding is needed
func (n *NFT) NeedsFeeding() bool {
	return n.Hunger < 50
}

// CanFeedFree checks if user can feed for free (once per day)
func (n *NFT) CanFeedFree() bool {
	return time.Since(n.LastFed) >= 24*time.Hour
}

// CanPlayFree checks if user can play for free
func (n *NFT) CanPlayFree() bool {
	return time.Since(n.LastPlayed) >= 4*time.Hour
}


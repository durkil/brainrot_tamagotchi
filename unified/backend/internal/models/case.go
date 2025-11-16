package models

import (
	"time"

	"gorm.io/gorm"
)

// CaseOpening represents a case opening event
type CaseOpening struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserAddress string    `gorm:"index;not null" json:"user_address"`
	CaseType    string    `json:"case_type"` // "bronze", "silver", "gold"
	TokenID     uint      `json:"token_id"`
	Rarity      string    `json:"rarity"`
	MemeType    string    `json:"meme_type"`
	Price       float64   `json:"price"`
	TxHash      string    `json:"tx_hash"`
	OpenedAt    time.Time `json:"opened_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the table name
func (CaseOpening) TableName() string {
	return "case_openings"
}


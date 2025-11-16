package repository

import (
	"brainrot-tamagotchi/internal/models"
	"strings"

	"gorm.io/gorm"
)

type MarketListingRepository struct {
	db *gorm.DB
}

func NewMarketListingRepository(db *gorm.DB) *MarketListingRepository {
	return &MarketListingRepository{db: db}
}

// Create creates a new listing
func (r *MarketListingRepository) Create(listing *models.MarketListing) error {
	listing.SellerAddress = strings.ToLower(listing.SellerAddress)
	return r.db.Create(listing).Error
}

// GetByTokenID retrieves a listing by token ID
func (r *MarketListingRepository) GetByTokenID(tokenID uint) (*models.MarketListing, error) {
	var listing models.MarketListing
	err := r.db.Preload("NFT").Where("token_id = ? AND is_active = ?", tokenID, true).First(&listing).Error
	if err != nil {
		return nil, err
	}
	return &listing, nil
}

// GetActiveListings retrieves all active listings with pagination
func (r *MarketListingRepository) GetActiveListings(limit, offset int, filters map[string]interface{}) ([]models.MarketListing, error) {
	var listings []models.MarketListing
	query := r.db.Preload("NFT").Where("is_active = ?", true)

	// Apply filters
	if rarity, ok := filters["rarity"]; ok {
		query = query.Joins("JOIN nfts ON nfts.token_id = market_listings.token_id").
			Where("nfts.rarity = ?", rarity)
	}

	if minLevel, ok := filters["min_level"]; ok {
		query = query.Joins("JOIN nfts ON nfts.token_id = market_listings.token_id").
			Where("nfts.level >= ?", minLevel)
	}

	if maxPrice, ok := filters["max_price"]; ok {
		query = query.Where("price <= ?", maxPrice)
	}

	err := query.Limit(limit).Offset(offset).Order("listed_at DESC").Find(&listings).Error
	return listings, err
}

// GetBySeller retrieves all listings by seller
func (r *MarketListingRepository) GetBySeller(sellerAddress string, active bool) ([]models.MarketListing, error) {
	var listings []models.MarketListing
	query := r.db.Preload("NFT").Where("seller_address = ?", strings.ToLower(sellerAddress))
	if active {
		query = query.Where("is_active = ?", true)
	}
	err := query.Order("listed_at DESC").Find(&listings).Error
	return listings, err
}

// Update updates a listing
func (r *MarketListingRepository) Update(listing *models.MarketListing) error {
	return r.db.Save(listing).Error
}

// Deactivate deactivates a listing
func (r *MarketListingRepository) Deactivate(tokenID uint) error {
	return r.db.Model(&models.MarketListing{}).
		Where("token_id = ?", tokenID).
		Update("is_active", false).Error
}

// MarkAsSold marks a listing as sold
func (r *MarketListingRepository) MarkAsSold(tokenID uint, buyerAddress string) error {
	return r.db.Model(&models.MarketListing{}).
		Where("token_id = ?", tokenID).
		Updates(map[string]interface{}{
			"is_active":      false,
			"buyer_address":  strings.ToLower(buyerAddress),
			"sold_at":        gorm.Expr("NOW()"),
		}).Error
}


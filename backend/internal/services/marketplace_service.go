package services

import (
	"brainrot-tamagotchi/internal/blockchain"
	"brainrot-tamagotchi/internal/models"
	"brainrot-tamagotchi/internal/repository"
	"fmt"
	"time"
)

type MarketplaceService struct {
	listingRepo *repository.MarketListingRepository
	nftRepo     *repository.NFTRepository
	blockchain  *blockchain.Client
}

func NewMarketplaceService(
	listingRepo *repository.MarketListingRepository,
	nftRepo *repository.NFTRepository,
	blockchain *blockchain.Client,
) *MarketplaceService {
	return &MarketplaceService{
		listingRepo: listingRepo,
		nftRepo:     nftRepo,
		blockchain:  blockchain,
	}
}

// ListNFT creates a new marketplace listing
func (s *MarketplaceService) ListNFT(tokenID uint, sellerAddress string, price float64) error {
	// Verify ownership
	nft, err := s.nftRepo.GetByTokenID(tokenID)
	if err != nil {
		return err
	}

	if nft.OwnerAddress != sellerAddress {
		return fmt.Errorf("not the owner of this NFT")
	}

	// Check if already listed
	existing, _ := s.listingRepo.GetByTokenID(tokenID)
	if existing != nil && existing.IsActive {
		return fmt.Errorf("NFT already listed")
	}

	// Create listing
	listing := &models.MarketListing{
		TokenID:       tokenID,
		SellerAddress: sellerAddress,
		Price:         price,
		IsActive:      true,
		ListedAt:      time.Now(),
	}

	// In production: Call smart contract's listNFT function
	// TODO: Implement blockchain integration

	return s.listingRepo.Create(listing)
}

// BuyNFT processes an NFT purchase
func (s *MarketplaceService) BuyNFT(tokenID uint, buyerAddress string) error {
	// Get listing
	listing, err := s.listingRepo.GetByTokenID(tokenID)
	if err != nil {
		return err
	}

	if !listing.IsActive {
		return fmt.Errorf("listing not active")
	}

	if listing.SellerAddress == buyerAddress {
		return fmt.Errorf("cannot buy your own NFT")
	}

	// Verify NFT still owned by seller
	nft, err := s.nftRepo.GetByTokenID(tokenID)
	if err != nil {
		return err
	}

	if nft.OwnerAddress != listing.SellerAddress {
		return fmt.Errorf("seller no longer owns the NFT")
	}

	// In production: Call smart contract's buyNFT function
	// TODO: Implement blockchain integration

	// Update NFT owner
	nft.OwnerAddress = buyerAddress
	if err := s.nftRepo.Update(nft); err != nil {
		return err
	}

	// Mark listing as sold
	return s.listingRepo.MarkAsSold(tokenID, buyerAddress)
}

// CancelListing cancels an active listing
func (s *MarketplaceService) CancelListing(tokenID uint, sellerAddress string) error {
	listing, err := s.listingRepo.GetByTokenID(tokenID)
	if err != nil {
		return err
	}

	if listing.SellerAddress != sellerAddress {
		return fmt.Errorf("not the seller")
	}

	// In production: Call smart contract's cancelListing function
	// TODO: Implement blockchain integration

	return s.listingRepo.Deactivate(tokenID)
}

// GetActiveListings retrieves active listings with filters
func (s *MarketplaceService) GetActiveListings(
	limit, offset int,
	filters map[string]interface{},
) ([]models.MarketListing, error) {
	return s.listingRepo.GetActiveListings(limit, offset, filters)
}

// GetUserListings retrieves all listings by a user
func (s *MarketplaceService) GetUserListings(userAddress string, activeOnly bool) ([]models.MarketListing, error) {
	return s.listingRepo.GetBySeller(userAddress, activeOnly)
}

// UpdatePrice updates the price of a listing
func (s *MarketplaceService) UpdatePrice(tokenID uint, sellerAddress string, newPrice float64) error {
	listing, err := s.listingRepo.GetByTokenID(tokenID)
	if err != nil {
		return err
	}

	if listing.SellerAddress != sellerAddress {
		return fmt.Errorf("not the seller")
	}

	listing.Price = newPrice

	// In production: Call smart contract's updatePrice function
	// TODO: Implement blockchain integration

	return s.listingRepo.Update(listing)
}

// GetMarketplaceStats returns marketplace statistics
func (s *MarketplaceService) GetMarketplaceStats() (map[string]interface{}, error) {
	// TODO: Implement aggregation queries
	stats := map[string]interface{}{
		"total_listings":   0,
		"total_sales":      0,
		"total_volume":     0.0,
		"average_price":    0.0,
		"floor_price":      0.0,
	}

	return stats, nil
}


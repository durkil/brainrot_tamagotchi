package services

import (
	"brainrot-tamagotchi/internal/blockchain"
	"brainrot-tamagotchi/internal/models"
	"brainrot-tamagotchi/internal/repository"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type TamagotchiService struct {
	nftRepo    *repository.NFTRepository
	redis      *redis.Client
	blockchain *blockchain.Client
}

func NewTamagotchiService(
	nftRepo *repository.NFTRepository,
	redis *redis.Client,
	blockchain *blockchain.Client,
) *TamagotchiService {
	return &TamagotchiService{
		nftRepo:    nftRepo,
		redis:      redis,
		blockchain: blockchain,
	}
}

// GetPetState retrieves the current state of a pet
func (s *TamagotchiService) GetPetState(tokenID uint) (*models.NFT, error) {
	nft, err := s.nftRepo.GetByTokenID(tokenID)
	if err != nil {
		return nil, err
	}

	// Update stats based on time passed
	s.updateStatsBasedOnTime(nft)

	return nft, nil
}

// FeedPet feeds the pet (free once per day or paid)
func (s *TamagotchiService) FeedPet(tokenID uint, ownerAddress string, isPaid bool) error {
	nft, err := s.nftRepo.GetByTokenID(tokenID)
	if err != nil {
		return err
	}

	// Verify ownership
	if nft.OwnerAddress != ownerAddress {
		return fmt.Errorf("not the owner of this NFT")
	}

	// Check if can feed for free
	if !isPaid && !nft.CanFeedFree() {
		return fmt.Errorf("free feeding not available yet")
	}

	// Feed the pet
	nft.Hunger = min(100, nft.Hunger+50)
	nft.LastFed = time.Now()
	nft.LastInteract = time.Now()

	return s.nftRepo.UpdateStats(nft)
}

// PlayWithPet plays with the pet to improve mood
func (s *TamagotchiService) PlayWithPet(tokenID uint, ownerAddress string) error {
	nft, err := s.nftRepo.GetByTokenID(tokenID)
	if err != nil {
		return err
	}

	// Verify ownership
	if nft.OwnerAddress != ownerAddress {
		return fmt.Errorf("not the owner of this NFT")
	}

	// Check energy
	if nft.Energy < 10 {
		return fmt.Errorf("not enough energy to play")
	}

	// Play with pet
	nft.Mood = min(100, nft.Mood+30)
	nft.Energy = max(0, nft.Energy-10)
	nft.LastPlayed = time.Now()
	nft.LastInteract = time.Now()

	return s.nftRepo.UpdateStats(nft)
}

// RestorePet restores a dead pet (paid revival)
func (s *TamagotchiService) RestorePet(tokenID uint, ownerAddress string) error {
	nft, err := s.nftRepo.GetByTokenID(tokenID)
	if err != nil {
		return err
	}

	// Verify ownership
	if nft.OwnerAddress != ownerAddress {
		return fmt.Errorf("not the owner of this NFT")
	}

	// Restore to 50% stats
	nft.Hunger = 50
	nft.Mood = 50
	nft.Energy = 50
	nft.LastInteract = time.Now()

	return s.nftRepo.UpdateStats(nft)
}

// StartHungerDecayJob starts a background job to decay hunger/mood/energy
func (s *TamagotchiService) StartHungerDecayJob() {
	ticker := time.NewTicker(1 * time.Hour) // Check every hour
	defer ticker.Stop()

	log.Println("🔄 Hunger decay job started")

	for {
		<-ticker.C
		s.decreaseAllStats()
	}
}

// decreaseAllStats decreases stats for all alive pets
func (s *TamagotchiService) decreaseAllStats() {
	nfts, err := s.nftRepo.GetAliveNFTs()
	if err != nil {
		log.Printf("Error fetching alive NFTs: %v", err)
		return
	}

	for _, nft := range nfts {
		updated := false

		// Decrease hunger every 6 hours
		hoursSinceLastFed := time.Since(nft.LastFed).Hours()
		if hoursSinceLastFed >= 6 {
			periodsElapsed := int(hoursSinceLastFed / 6)
			nft.Hunger = max(0, nft.Hunger-(25*periodsElapsed))
			updated = true
		}

		// Decrease mood every 12 hours
		hoursSinceLastPlayed := time.Since(nft.LastPlayed).Hours()
		if hoursSinceLastPlayed >= 12 {
			periodsElapsed := int(hoursSinceLastPlayed / 12)
			nft.Mood = max(0, nft.Mood-(20*periodsElapsed))
			updated = true
		}

		// Regenerate energy slowly (10 per hour)
		if nft.Energy < 100 {
			nft.Energy = min(100, nft.Energy+10)
			updated = true
		}

		if updated {
			if err := s.nftRepo.UpdateStats(&nft); err != nil {
				log.Printf("Error updating stats for NFT %d: %v", nft.TokenID, err)
			}
		}
	}

	log.Printf("✅ Updated stats for %d pets", len(nfts))
}

// updateStatsBasedOnTime updates stats in real-time based on time passed
func (s *TamagotchiService) updateStatsBasedOnTime(nft *models.NFT) {
	now := time.Now()

	// Calculate hunger decay
	hoursSinceLastFed := now.Sub(nft.LastFed).Hours()
	if hoursSinceLastFed >= 6 {
		periodsElapsed := int(hoursSinceLastFed / 6)
		nft.Hunger = max(0, nft.Hunger-(25*periodsElapsed))
	}

	// Calculate mood decay
	hoursSinceLastPlayed := now.Sub(nft.LastPlayed).Hours()
	if hoursSinceLastPlayed >= 12 {
		periodsElapsed := int(hoursSinceLastPlayed / 12)
		nft.Mood = max(0, nft.Mood-(20*periodsElapsed))
	}

	// Energy regeneration
	hoursSinceLastInteract := now.Sub(nft.LastInteract).Hours()
	if nft.Energy < 100 && hoursSinceLastInteract > 0 {
		nft.Energy = min(100, nft.Energy+int(hoursSinceLastInteract*10))
	}
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}


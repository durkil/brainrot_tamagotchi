package repository

import (
	"brainrot-tamagotchi/internal/models"
	"strings"

	"gorm.io/gorm"
)

type NFTRepository struct {
	db *gorm.DB
}

func NewNFTRepository(db *gorm.DB) *NFTRepository {
	return &NFTRepository{db: db}
}

// Create creates a new NFT record
func (r *NFTRepository) Create(nft *models.NFT) error {
	nft.OwnerAddress = strings.ToLower(nft.OwnerAddress)
	return r.db.Create(nft).Error
}

// GetByTokenID retrieves an NFT by token ID
func (r *NFTRepository) GetByTokenID(tokenID uint) (*models.NFT, error) {
	var nft models.NFT
	err := r.db.Where("token_id = ?", tokenID).First(&nft).Error
	if err != nil {
		return nil, err
	}
	return &nft, nil
}

// GetByOwner retrieves all NFTs owned by an address
func (r *NFTRepository) GetByOwner(ownerAddress string) ([]models.NFT, error) {
	var nfts []models.NFT
	err := r.db.Where("owner_address = ?", strings.ToLower(ownerAddress)).Find(&nfts).Error
	return nfts, err
}

// Update updates an NFT
func (r *NFTRepository) Update(nft *models.NFT) error {
	return r.db.Save(nft).Error
}

// UpdateStats updates only the stats fields (hunger, mood, energy)
func (r *NFTRepository) UpdateStats(nft *models.NFT) error {
	return r.db.Model(nft).Updates(map[string]interface{}{
		"hunger":        nft.Hunger,
		"mood":          nft.Mood,
		"energy":        nft.Energy,
		"last_fed":      nft.LastFed,
		"last_played":   nft.LastPlayed,
		"last_interact": nft.LastInteract,
	}).Error
}

// GetAll retrieves all NFTs with pagination
func (r *NFTRepository) GetAll(limit, offset int) ([]models.NFT, error) {
	var nfts []models.NFT
	err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&nfts).Error
	return nfts, err
}

// GetAliveNFTs retrieves all NFTs that need stat updates
func (r *NFTRepository) GetAliveNFTs() ([]models.NFT, error) {
	var nfts []models.NFT
	err := r.db.Where("hunger > 0 AND mood > 0 AND energy > 0").Find(&nfts).Error
	return nfts, err
}

// Delete soft deletes an NFT (when burned)
func (r *NFTRepository) Delete(tokenID uint) error {
	return r.db.Where("token_id = ?", tokenID).Delete(&models.NFT{}).Error
}


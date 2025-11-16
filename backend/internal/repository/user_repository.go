package repository

import (
	"brainrot-tamagotchi/internal/models"
	"strings"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	user.WalletAddress = strings.ToLower(user.WalletAddress)
	return r.db.Create(user).Error
}

// GetByWalletAddress retrieves a user by wallet address
func (r *UserRepository) GetByWalletAddress(address string) (*models.User, error) {
	var user models.User
	err := r.db.Where("wallet_address = ?", strings.ToLower(address)).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetOrCreate gets or creates a user
func (r *UserRepository) GetOrCreate(address string) (*models.User, error) {
	user, err := r.GetByWalletAddress(address)
	if err == gorm.ErrRecordNotFound {
		user = &models.User{
			WalletAddress: strings.ToLower(address),
		}
		if err := r.Create(user); err != nil {
			return nil, err
		}
		return user, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}


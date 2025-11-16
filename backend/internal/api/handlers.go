package api

import (
	"brainrot-tamagotchi/internal/repository"
	"brainrot-tamagotchi/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	tamagotchiService  *services.TamagotchiService
	caseService        *services.CaseService
	marketplaceService *services.MarketplaceService
	userRepo           *repository.UserRepository
}

func NewHandler(
	tamagotchiService *services.TamagotchiService,
	caseService *services.CaseService,
	marketplaceService *services.MarketplaceService,
	userRepo *repository.UserRepository,
) *Handler {
	return &Handler{
		tamagotchiService:  tamagotchiService,
		caseService:        caseService,
		marketplaceService: marketplaceService,
		userRepo:           userRepo,
	}
}

// Health check
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Brainrot Tamagotchi API is running! 🎮",
		"version": "1.0.0",
	})
}

// ==================== Pet / Tamagotchi Endpoints ====================

// GetPet retrieves pet information
func (h *Handler) GetPet(c *gin.Context) {
	tokenID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token ID"})
		return
	}

	nft, err := h.tamagotchiService.GetPetState(uint(tokenID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
		return
	}

	c.JSON(http.StatusOK, nft)
}

// FeedPet feeds the pet
func (h *Handler) FeedPet(c *gin.Context) {
	tokenID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token ID"})
		return
	}

	var body struct {
		IsPaid bool `json:"is_paid"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		body.IsPaid = false
	}

	walletAddress := c.GetHeader("X-Wallet-Address")
	if walletAddress == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wallet address required"})
		return
	}

	err = h.tamagotchiService.FeedPet(uint(tokenID), walletAddress, body.IsPaid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pet fed successfully", "hunger": 100})
}

// PlayWithPet plays with the pet
func (h *Handler) PlayWithPet(c *gin.Context) {
	tokenID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token ID"})
		return
	}

	walletAddress := c.GetHeader("X-Wallet-Address")
	if walletAddress == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wallet address required"})
		return
	}

	err = h.tamagotchiService.PlayWithPet(uint(tokenID), walletAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Played with pet successfully"})
}

// ==================== Cases Endpoints ====================

// GetCasePrices returns prices for all case types
func (h *Handler) GetCasePrices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"prices": services.CasePrices,
	})
}

// BuyCase purchases a case
func (h *Handler) BuyCase(c *gin.Context) {
	var body struct {
		CaseType string `json:"case_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	walletAddress := c.GetHeader("X-Wallet-Address")
	if walletAddress == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wallet address required"})
		return
	}

	// TODO: Implement case purchase logic
	c.JSON(http.StatusOK, gin.H{
		"message":   "Case purchased (implementation pending)",
		"case_type": body.CaseType,
	})
}

// OpenCase opens a purchased case
func (h *Handler) OpenCase(c *gin.Context) {
	caseID := c.Param("id")

	// TODO: Implement case opening logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Case opened (implementation pending)",
		"case_id": caseID,
	})
}

// ==================== Marketplace Endpoints ====================

// GetMarketplace retrieves active marketplace listings
func (h *Handler) GetMarketplace(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	filters := make(map[string]interface{})
	if rarity := c.Query("rarity"); rarity != "" {
		filters["rarity"] = rarity
	}
	if minLevel := c.Query("min_level"); minLevel != "" {
		if level, err := strconv.Atoi(minLevel); err == nil {
			filters["min_level"] = level
		}
	}

	listings, err := h.marketplaceService.GetActiveListings(limit, offset, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch listings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"listings": listings,
		"count":    len(listings),
	})
}

// ListNFT creates a new marketplace listing
func (h *Handler) ListNFT(c *gin.Context) {
	var body struct {
		TokenID uint    `json:"token_id" binding:"required"`
		Price   float64 `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	walletAddress := c.GetHeader("X-Wallet-Address")
	if walletAddress == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wallet address required"})
		return
	}

	err := h.marketplaceService.ListNFT(body.TokenID, walletAddress, body.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "NFT listed successfully"})
}

// BuyNFT purchases an NFT from marketplace
func (h *Handler) BuyNFT(c *gin.Context) {
	tokenID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token ID"})
		return
	}

	walletAddress := c.GetHeader("X-Wallet-Address")
	if walletAddress == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wallet address required"})
		return
	}

	err = h.marketplaceService.BuyNFT(uint(tokenID), walletAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "NFT purchased successfully"})
}

// CancelListing cancels a marketplace listing
func (h *Handler) CancelListing(c *gin.Context) {
	tokenID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token ID"})
		return
	}

	walletAddress := c.GetHeader("X-Wallet-Address")
	if walletAddress == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wallet address required"})
		return
	}

	err = h.marketplaceService.CancelListing(uint(tokenID), walletAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Listing cancelled successfully"})
}

// ==================== User Endpoints ====================

// GetUser retrieves user information
func (h *Handler) GetUser(c *gin.Context) {
	address := c.Param("address")

	user, err := h.userRepo.GetOrCreate(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetInventory retrieves user's NFT inventory
func (h *Handler) GetInventory(c *gin.Context) {
	address := c.Param("address")

	// TODO: Implement inventory retrieval
	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"nfts":    []interface{}{},
	})
}


package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func (h *Handler) SetupRoutes(router *gin.Engine) {
	// API v1 group
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", h.HealthCheck)

		// Pet / Tamagotchi routes
		pets := api.Group("/pets")
		{
			pets.GET("/:id", h.GetPet)           // Get pet state
			pets.POST("/:id/feed", h.FeedPet)    // Feed pet
			pets.POST("/:id/play", h.PlayWithPet) // Play with pet
		}

		// Cases routes
		cases := api.Group("/cases")
		{
			cases.GET("/prices", h.GetCasePrices) // Get case prices
			cases.POST("/buy", h.BuyCase)          // Buy a case
			cases.POST("/:id/open", h.OpenCase)    // Open a case
		}

		// Marketplace routes
		marketplace := api.Group("/marketplace")
		{
			marketplace.GET("", h.GetMarketplace)           // Browse marketplace
			marketplace.POST("/list", h.ListNFT)            // List NFT for sale
			marketplace.POST("/:id/buy", h.BuyNFT)          // Buy NFT
			marketplace.DELETE("/:id", h.CancelListing)     // Cancel listing
		}

		// User routes
		users := api.Group("/users")
		{
			users.GET("/:address", h.GetUser)                // Get user info
			users.GET("/:address/inventory", h.GetInventory) // Get user's NFTs
		}
	}

	// Serve OpenAPI/Swagger docs (optional)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":    "Brainrot Tamagotchi API",
			"version": "1.0.0",
			"docs":    "/api/v1/health",
		})
	})
}


package main

import (
	"brainrot-tamagotchi/internal/api"
	"brainrot-tamagotchi/internal/blockchain"
	"brainrot-tamagotchi/internal/repository"
	"brainrot-tamagotchi/internal/services"
	"brainrot-tamagotchi/pkg/cache"
	"brainrot-tamagotchi/pkg/database"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database
	log.Println("📦 Connecting to database...")
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("✅ Database connected and migrated")

	// Initialize Redis
	log.Println("📦 Connecting to Redis...")
	redisClient := cache.NewRedisClient()
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("✅ Redis connected")

	// Initialize blockchain client
	log.Println("📦 Connecting to Base blockchain...")
	blockchainClient, err := blockchain.NewClient()
	if err != nil {
		log.Fatal("Failed to connect to blockchain:", err)
	}
	log.Println("✅ Blockchain connected")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	nftRepo := repository.NewNFTRepository(db)
	listingRepo := repository.NewMarketListingRepository(db)

	// Initialize services
	tamagotchiService := services.NewTamagotchiService(nftRepo, redisClient, blockchainClient)
	caseService := services.NewCaseService(blockchainClient, nftRepo)
	marketplaceService := services.NewMarketplaceService(listingRepo, nftRepo, blockchainClient)

	// Start background jobs
	log.Println("🔄 Starting background jobs...")
	go tamagotchiService.StartHungerDecayJob()
	log.Println("✅ Background jobs started")

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://brainrot-tamagotchi.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Wallet-Address"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize API handlers
	handler := api.NewHandler(
		tamagotchiService,
		caseService,
		marketplaceService,
		userRepo,
	)

	// Setup routes
	handler.SetupRoutes(router)

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s...\n", port)
	log.Println("🎮 Brainrot Tamagotchi API is ready!")
	log.Println("📖 Documentation: http://localhost:" + port + "/api/v1/health")

	// Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Server exited")
}


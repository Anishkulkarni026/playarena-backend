package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/cloudinary/cloudinary-go/v2"

	"github.com/JkD004/playarena-backend/api"
	"github.com/JkD004/playarena-backend/db"
	"github.com/JkD004/playarena-backend/venue"
	"github.com/JkD004/playarena-backend/user"
)

func main() {

	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Warning: .env file not found, continuing...")
	}

	// Initialize DB
	db.InitDB()

	// Initialize Cloudinary
	cld, err := cloudinary.New()
	if err != nil {
		log.Fatalf("❌ Failed to initialize Cloudinary: %v", err)
	}

	venue.SetCloudinary(cld)
	user.SetCloudinary(cld)

	// Setup Gin Router
	router := gin.Default()

	// CORS FIX
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true 
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// 🟢 ROOT ROUTE (Fix 404)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "PlayArena Backend is Running 🚀",
		})
	})

	// API Routes
	api.SetupRoutes(router)

	// Start Server
	log.Println("🚀 Backend running on port 8080...")
	router.Run(":8080")
}

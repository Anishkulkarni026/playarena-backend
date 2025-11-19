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

	// ✅ Load environment variables (.env)
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Warning: .env file not found, continuing...")
	}

	// ✅ Initialize Database
	db.InitDB()

	// ✅ Initialize Cloudinary
	cld, err := cloudinary.New()
	if err != nil {
		log.Fatalf("❌ Failed to initialize Cloudinary: %v", err)
	}

	// ✅ Pass Cloudinary into Venue module
	venue.SetCloudinary(cld)
	user.SetCloudinary(cld)

	// ✅ Setup Gin Router
	router := gin.Default()

	// ✅ CORS Configuration (THE FIX)
	config := cors.DefaultConfig()
	
	// Change 1: Allow ALL origins to stop the 403 errors
	config.AllowAllOrigins = true 
	
	// Change 2: Explicitly allow all the methods we are using
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	
	// Change 3: Allow headers for file uploads and auth
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	
	router.Use(cors.New(config))

	// ✅ Set up routes
	api.SetupRoutes(router)

	// ✅ Run Server
	router.Run(":8080")
}
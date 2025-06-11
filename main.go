package main

import (
	"dashlearn/modules/category"
	"dashlearn/modules/instructor"
	"dashlearn/modules/user"
	"dashlearn/utils"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Set Gin to release mode
	// gin.SetMode(gin.ReleaseMode)

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Warning: No .env file found")
	}

	// Initialize Gin
	router := gin.Default()

	// Enable CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // allow all origins
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		// ExposeHeaders:    []string{"Content-Length"},
		// AllowCredentials: true,
		// MaxAge: 12 * time.Hour,
	}))

	// Initialize API routes
	apiRoutesGroup := router.Group("/api")

	// Connect to database
	utils.ConnectDatabase()

	// craete superadmin
	utils.CreateSuperadminIfNotExists()

	// Register routes
	user.RegisterUserRoutes(apiRoutesGroup)
	instructor.RegisterInstructorRoutes(apiRoutesGroup)
	category.RegisterCategoryRoutes(apiRoutesGroup)

	// Run the server
	router.Run(":5000")
}

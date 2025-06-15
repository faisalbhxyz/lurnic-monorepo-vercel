package main

import (
	"dashlearn/models"
	"dashlearn/modules/category"
	"dashlearn/modules/course"
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

	// Load environment variables
	if gin.Mode() != gin.ReleaseMode {
		if err := godotenv.Load(); err != nil {
			log.Fatalln("Warning: No .env file found")
		}
	}

	// Initialize Gin
	router := gin.Default()

	// Enable CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Access-Control-Allow-Origin"},
		AllowCredentials: false,
	}))

	// Connect to database
	utils.ConnectDatabase()

	// Initialize API routes
	apiRoutesGroup := router.Group("/api")

	// craete superadmin
	CreateSuperadminIfNotExists()

	// Register routes
	user.RegisterUserRoutes(apiRoutesGroup)
	instructor.RegisterInstructorRoutes(apiRoutesGroup)
	category.RegisterCategoryRoutes(apiRoutesGroup)
	course.RegisterCourseRoutes(apiRoutesGroup)

	// Run the server
	router.Run(":5000")
}

func CreateSuperadminIfNotExists() {
	var r models.Role
	err := utils.DB.FirstOrCreate(&r, models.Role{Name: "superadmin", TenantID: nil}).Error
	if err != nil {
		panic("Failed to create or find superadmin role: " + err.Error())
	}
}

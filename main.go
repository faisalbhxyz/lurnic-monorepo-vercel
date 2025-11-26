package main

import (
	"dashlearn/internal/models"
	"dashlearn/internal/modules/banner"
	"dashlearn/internal/modules/category"
	"dashlearn/internal/modules/course"
	"dashlearn/internal/modules/enrollment"
	generalsettings "dashlearn/internal/modules/general_settings"
	"dashlearn/internal/modules/instructor"
	"dashlearn/internal/modules/order"
	paymentmethod "dashlearn/internal/modules/payment_method"
	"dashlearn/internal/modules/student"
	subcategory "dashlearn/internal/modules/sub_category"
	"dashlearn/internal/modules/user"
	"dashlearn/internal/utils"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

var Version = "v1.0.8"

func main() {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Load environment variables
	if gin.Mode() != gin.ReleaseMode {
		if err := godotenv.Load(); err != nil {
			log.Fatalln("Warning: No .env file found")
		}
	}

	c := cron.New()

	// Run the helper every minute
	c.AddFunc("@every 1m", func() {
		if err := course.CronJobForCoursesSchedule(utils.DB); err != nil {
			fmt.Println("⏰ Cron error:", err)
		}
	})

	c.Start()
	fmt.Println("⌛ Cron started for scheduled courses (Bangladesh time GMT+6)")

	// Initialize Gin
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 100<<20) // 100 MB
		c.Next()
	})

	// Enable CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "app-key"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Access-Control-Allow-Origin"},
		AllowCredentials: false,
	}))

	// Connect to database
	utils.ConnectDatabase()

	// Initialize API routes
	apiRoutesGroup := router.Group("/v1")

	// craete superadmin
	CreateSuperadminIfNotExists()

	// Register routes
	user.RegisterUserRoutes(apiRoutesGroup)
	instructor.RegisterInstructorRoutes(apiRoutesGroup)
	student.RegisterStudentRoutes(apiRoutesGroup)
	category.RegisterCategoryRoutes(apiRoutesGroup)
	subcategory.RegisterSubCategoryRoutes(apiRoutesGroup)
	course.RegisterCourseRoutes(apiRoutesGroup)
	enrollment.RegisterEnrollmentRoutes(apiRoutesGroup)
	banner.RegisterBannerRoutes(apiRoutesGroup)
	order.RegisterCourseRoutes(apiRoutesGroup)
	generalsettings.RegisterGeneralSettingsRoutes(apiRoutesGroup)
	paymentmethod.RegisterRoutes(apiRoutesGroup)

	// Run the server
	router.Run(":" + os.Getenv("APP_PORT"))
}

func CreateSuperadminIfNotExists() {
	var r models.Role
	err := utils.DB.FirstOrCreate(&r, models.Role{Name: "superadmin", TenantID: nil}).Error
	if err != nil {
		panic("Failed to create or find superadmin role: " + err.Error())
	}
}

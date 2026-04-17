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
	"dashlearn/internal/modules/role"
	"dashlearn/internal/modules/student"
	subcategory "dashlearn/internal/modules/sub_category"
	"dashlearn/internal/modules/user"
	"dashlearn/internal/observability"
	"dashlearn/internal/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/getsentry/sentry-go"
)

var Version = "v1.0.24"

func main() {
	fmt.Println("🚀 DashLearn Server Starting... Version:", Version)
	// Respect externally configured GIN_MODE; default to debug for local dev.
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// Load environment variables from .env when present (non-fatal).
	if err := godotenv.Load(); err != nil {
		log.Println("Info: No .env file found; relying on process environment")
	}

	debugRoutesEnabled := os.Getenv("ENABLE_DEBUG_ROUTES") == "true"

	flush, sentryEnabled := observability.InitSentry(observability.EnvSentryConfig(Version))
	if sentryEnabled {
		defer flush(2 * time.Second)
		defer func() {
			if r := recover(); r != nil {
				sentry.CurrentHub().Recover(r)
				flush(2 * time.Second)
				panic(r)
			}
		}()
	}

	// Initialize Gin
	router := gin.Default()
	if sentryEnabled {
		router.Use(sentrygin.New(sentrygin.Options{
			Repanic:         true,
			WaitForDelivery: false,
			Timeout:         2 * time.Second,
		}))
	}

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

	// Debug-only route to validate Sentry ingestion. Disabled by default.
	if debugRoutesEnabled {
		apiRoutesGroup.GET("/debug/sentry-test", func(c *gin.Context) {
			if c.ClientIP() != "127.0.0.1" && c.ClientIP() != "::1" {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			if requiredKey := os.Getenv("DEBUG_ROUTE_KEY"); requiredKey != "" && c.GetHeader("X-Debug-Key") != requiredKey {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			eventID := sentry.CaptureMessage("sentry-test: manual capture to validate Sentry ingestion")
			sentry.Flush(2 * time.Second)
			c.JSON(http.StatusOK, gin.H{
				"ok":      true,
				"eventId": eventID,
			})
		})
	}

	// Connect to database (required for all non-debug routes)
	if err := utils.ConnectDatabase(); err != nil {
		if debugRoutesEnabled {
			log.Printf("DB unavailable; starting in debug-only mode: %v", err)
		} else {
			log.Fatal("Failed to connect to database:", err)
		}
	} else {
		c := cron.New()

		// Run the helper every minute
		c.AddFunc("@every 1m", func() {
			if err := course.CronJobForCoursesSchedule(utils.DB); err != nil {
				fmt.Println("📚 Course Cron error:", err)
			}
			if err := course.CronJobForCourseLessonsSchedule(utils.DB); err != nil {
				fmt.Println("🧾 Lesson Cron error:", err)
			}
		})

		c.Start()
		fmt.Println("⌛ Cron started for scheduled courses and lessons (Bangladesh time GMT+6)")

		// craete superadmin
		CreateSuperadminIfNotExists()

		// Register routes
		user.RegisterUserRoutes(apiRoutesGroup)
		role.RegsiterRoleRoutes(apiRoutesGroup)
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
	}

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

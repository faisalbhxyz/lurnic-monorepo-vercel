package server

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
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

// NewEngine builds the Gin engine with all /v1 routes. Used by main (long-running)
// and by Vercel serverless (no cron when VERCEL is set).
// flush should be called on shutdown in long-running processes (main); serverless can ignore it.
func NewEngine(version string) (*gin.Engine, func(time.Duration) bool, error) {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	if err := godotenv.Load(); err != nil {
		log.Println("Info: No .env file found; relying on process environment")
	}

	debugRoutesEnabled := os.Getenv("ENABLE_DEBUG_ROUTES") == "true"
	onVercel := os.Getenv("VERCEL") != ""

	flush, sentryEnabled := observability.InitSentry(observability.EnvSentryConfig(version))

	router := gin.Default()
	// Avoid /path <-> /path/ redirect loops when behind Vercel rewrites (serverless).
	router.RedirectTrailingSlash = false

	if sentryEnabled {
		router.Use(sentrygin.New(sentrygin.Options{
			Repanic:         true,
			WaitForDelivery: false,
			Timeout:         2 * time.Second,
		}))
	}

	router.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 100<<20)
		c.Next()
	})

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "app-key"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Access-Control-Allow-Origin"},
		AllowCredentials: false,
	}))

	apiRoutesGroup := router.Group("/v1")

	if debugRoutesEnabled {
		apiRoutesGroup.GET("/debug/sentry-test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true, "note": "use main server for full sentry test"})
		})
	}

	if err := utils.ConnectDatabase(); err != nil {
		if debugRoutesEnabled {
			log.Printf("DB unavailable; starting without API routes: %v", err)
			return router, flush, nil
		}
		return nil, flush, err
	}

	if !onVercel {
		c := cron.New()
		c.AddFunc("@every 1m", func() {
			if err := course.CronJobForCoursesSchedule(utils.DB); err != nil {
				log.Println("Course Cron error:", err)
			}
			if err := course.CronJobForCourseLessonsSchedule(utils.DB); err != nil {
				log.Println("Lesson Cron error:", err)
			}
		})
		c.Start()
		log.Println("Cron started for scheduled courses and lessons")
	} else {
		log.Println("VERCEL=1: in-process cron disabled; use Vercel Cron → GET /v1/internal/cron/publish-scheduled with CRON_SECRET")
	}

	CreateSuperadminIfNotExists()

	user.RegisterUserRoutes(apiRoutesGroup)
	role.RegsiterRoleRoutes(apiRoutesGroup)
	instructor.RegisterInstructorRoutes(apiRoutesGroup)
	student.RegisterStudentRoutes(apiRoutesGroup)
	category.RegisterCategoryRoutes(apiRoutesGroup)
	subcategory.RegisterSubCategoryRoutes(apiRoutesGroup)
	course.RegisterCourseRoutes(apiRoutesGroup)

	// Vercel/serverless: in-process cron is disabled when VERCEL=1; schedule this path in vercel.json
	// and set CRON_SECRET (Vercel sends Authorization: Bearer <CRON_SECRET>).
	apiRoutesGroup.GET("/internal/cron/publish-scheduled", func(c *gin.Context) {
		secret := os.Getenv("CRON_SECRET")
		if secret == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "CRON_SECRET is not configured"})
			return
		}
		if c.GetHeader("Authorization") != "Bearer "+secret {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err := course.CronJobForCoursesSchedule(utils.DB); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := course.CronJobForCourseLessonsSchedule(utils.DB); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	enrollment.RegisterEnrollmentRoutes(apiRoutesGroup)
	banner.RegisterBannerRoutes(apiRoutesGroup)
	order.RegisterCourseRoutes(apiRoutesGroup)
	generalsettings.RegisterGeneralSettingsRoutes(apiRoutesGroup)
	paymentmethod.RegisterRoutes(apiRoutesGroup)

	return router, flush, nil
}

func CreateSuperadminIfNotExists() {
	var r models.Role
	err := utils.DB.FirstOrCreate(&r, models.Role{Name: "superadmin", TenantID: nil}).Error
	if err != nil {
		panic("Failed to create or find superadmin role: " + err.Error())
	}
}

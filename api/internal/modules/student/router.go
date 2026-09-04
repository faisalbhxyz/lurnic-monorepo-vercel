package student

import (
	"dashlearn/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterStudentRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/private/student", middleware.AuthMiddleware())
	{
		storefrontAdmin := NewStorefrontHandler()

		authGroup.GET("", GetStudents)
		authGroup.GET("/lite", GetStudentLite)
		authGroup.GET("/details/:id", GetStudentDetailsByID)
		authGroup.GET("/:studentId/learning-report", storefrontAdmin.GetAdminLearningReport)
		authGroup.POST("/register", CreateStudent)
		authGroup.PUT("/update/:id", UpdateStudent)
		authGroup.DELETE("/delete/:id", DeleteStudent)
	}

	adminStudents := rg.Group("/admin/students", middleware.AuthMiddleware())
	{
		storefrontAdmin := NewStorefrontHandler()
		adminStudents.GET("/:studentId/learning-report", storefrontAdmin.GetAdminLearningReport)
	}

	publicGroup := rg.Group("/student", middleware.GetTenantID())
	{
		storefront := NewStorefrontHandler()

		publicGroup.POST("/login", LoginStudent)
		publicGroup.POST("/logout", middleware.StudentAuthMiddleware(), LogoutStudent)
		publicGroup.POST("/register", CreateStudentPublic)
		publicGroup.POST("/forgot-password", ForgotPasswordStudent)
		publicGroup.POST("/reset-password", ResetPasswordStudent)
		publicGroup.GET("/details", middleware.StudentAuthMiddleware(), GetStudentDetails)
		publicGroup.PUT("/update", middleware.StudentAuthMiddleware(), UpdateStudentProfile)
		publicGroup.GET("/class-profile", middleware.StudentAuthMiddleware(), GetClassProfile)
		publicGroup.PUT("/class-profile", middleware.StudentAuthMiddleware(), PutClassProfile)

		publicGroup.GET("/learning-report", middleware.StudentAuthMiddleware(), storefront.GetLearningReport)
		publicGroup.POST("/watch-time", middleware.StudentAuthMiddleware(), storefront.PostWatchTime)
		publicGroup.POST("/watch-time/batch", middleware.StudentAuthMiddleware(), storefront.PostWatchTimeBatch)
		publicGroup.GET("/notifications", middleware.StudentAuthMiddleware(), storefront.ListNotifications)
		publicGroup.PATCH("/notifications/:id/read", middleware.StudentAuthMiddleware(), storefront.MarkNotificationRead)
		publicGroup.GET("/orders", middleware.StudentAuthMiddleware(), storefront.ListOrders)
	}
}

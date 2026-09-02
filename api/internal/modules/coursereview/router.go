package coursereview

import (
	"dashlearn/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	handler := NewHandler()

	courseGroup := rg.Group("/course", middleware.GetTenantID())
	{
		courseGroup.GET("/:slug/reviews", middleware.OptionalStudentAuthMiddleware(), handler.ListCourseReviews)
		courseGroup.POST("/:slug/review", middleware.StudentAuthMiddleware(), handler.CreateCourseReview)
	}
}

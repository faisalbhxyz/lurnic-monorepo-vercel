package courseprogress

import (
	"dashlearn/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCourseProgressRoutes(rg *gin.RouterGroup) {
	handler := NewHandler()

	studentCourseGroup := rg.Group("/course", middleware.GetTenantID(), middleware.StudentAuthMiddleware())
	{
		studentCourseGroup.GET("/:slug/progress", handler.GetCourseProgress)
		studentCourseGroup.POST("/:slug/lessons/:lessonId/complete", handler.MarkLessonComplete)
	}
}

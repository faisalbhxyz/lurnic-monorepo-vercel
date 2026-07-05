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
		studentCourseGroup.GET("/:slug/lessons/:lessonId/progress", handler.GetLessonVideoProgress)
		studentCourseGroup.PATCH("/:slug/lessons/:lessonId/progress", handler.UpdateLessonVideoProgress)
		studentCourseGroup.POST("/:slug/lessons/:lessonId/complete", handler.MarkLessonComplete)
	}
}

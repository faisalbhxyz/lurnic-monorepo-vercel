package student

import (
	"errors"
	"net/http"
	"strconv"

	"dashlearn/internal/apiresponse"
	"dashlearn/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StorefrontHandler struct {
	service *StorefrontService
}

func NewStorefrontHandler() *StorefrontHandler {
	return &StorefrontHandler{service: NewStorefrontService(utils.DB)}
}

func (h *StorefrontHandler) GetLearningReport(c *gin.Context) {
	period := c.DefaultQuery("period", "7d")
	data, err := h.service.GetLearningReport(c.GetUint("tenant_id"), c.GetUint("user_id"), period)
	if err != nil {
		if err.Error() == "invalid period" {
			apiresponse.Validation(c, "period must be 7d, 30d, or 90d")
			return
		}
		apiresponse.Error(c, http.StatusInternalServerError, "Failed to load learning report", "INTERNAL_ERROR")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *StorefrontHandler) ListNotifications(c *gin.Context) {
	data, err := h.service.ListNotifications(c.GetUint("tenant_id"), c.GetUint("user_id"))
	if err != nil {
		apiresponse.Error(c, http.StatusInternalServerError, "Failed to load notifications", "INTERNAL_ERROR")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *StorefrontHandler) MarkNotificationRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiresponse.Validation(c, "Invalid notification ID")
		return
	}

	if err := h.service.MarkNotificationRead(c.GetUint("tenant_id"), c.GetUint("user_id"), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			apiresponse.NotFound(c, "Notification not found")
			return
		}
		apiresponse.Error(c, http.StatusInternalServerError, "Failed to update notification", "INTERNAL_ERROR")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func (h *StorefrontHandler) ListOrders(c *gin.Context) {
	data, err := h.service.ListOrders(c.GetUint("tenant_id"), c.GetUint("user_id"))
	if err != nil {
		apiresponse.Error(c, http.StatusInternalServerError, "Failed to load orders", "INTERNAL_ERROR")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

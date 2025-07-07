package banner

import (
	"context"
	"dashlearn/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type BannerHandler struct {
	service BannerService
}

func NewBannerHandler(db *gorm.DB) *BannerHandler {
	return &BannerHandler{
		service: NewBannerService(db),
	}
}

func (h *BannerHandler) GetAll(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	banners, err := h.service.GetAll(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": banners})
}

func (h *BannerHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid banner ID"})
		return
	}
	tenantID := c.GetUint("tenant_id")

	banner, err := h.service.GetByID(tenantID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Banner not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": banner})
}

func (h *BannerHandler) Create(c *gin.Context) {
	var input CreateBannerInput
	if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {

		image, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
			return
		} else {
			imageURL, err := utils.UploadFile(context.Background(), image)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
				return
			}
			input.Image = imageURL
		}

		if err := h.service.Create(input, c.GetUint("tenant_id")); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Banner created successfully"})
	}
}

func (h *BannerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid banner ID"})
		return
	}

	var input UpdateBannerInput
	if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {

		log, _ := json.MarshalIndent(err, "", "  ")
		fmt.Println("Error binding input:", string(log))

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	image, err := c.FormFile("image")

	if err == nil {
		imageURL, err := utils.UploadFile(context.Background(), image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		}
		input.Image = &imageURL
	}

	if err := h.service.Update(id, input, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category updated successfully"})
}

func (h *BannerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid banner ID"})
		return
	}
	if err := h.service.Delete(id, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Banner deleted successfully"})
}

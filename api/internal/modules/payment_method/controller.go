package paymentmethod

import (
	"dashlearn/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type PaymentMethodHandler struct {
	service PaymentMethodService
}

func NewPaymentMethodHandler(db *gorm.DB) *PaymentMethodHandler {
	return &PaymentMethodHandler{
		service: NewPaymentMethodService(db),
	}
}

func (h *PaymentMethodHandler) GetAllPrivate(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	categories, err := h.service.GetAll(tenantID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *PaymentMethodHandler) GetAllPublic(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	status := true
	categories, err := h.service.GetAll(tenantID, &status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *PaymentMethodHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method ID"})
		return
	}
	tenantID := c.GetUint("tenant_id")

	category, err := h.service.GetByID(tenantID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment method not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": category})
}

func (h *PaymentMethodHandler) Create(c *gin.Context) {
	var input CreatePaymentMethodInput
	if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {
		handleValidationError(c, err)
		return
	}

	file_headers, err := c.FormFile("image")
	if err == nil {
		if file_headers.Size > 2*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Max image size is 2MB"})
			return
		}

		// ✅ 2. MIME type check
		file, err := file_headers.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}
		defer file.Close()

		// Detect content type
		buffer := make([]byte, 512)
		if _, err := file.Read(buffer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file content"})
			return
		}
		contentType := http.DetectContentType(buffer)

		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/jpg":  true,
		}
		if !allowedTypes[contentType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only PNG, JPG formats are supported"})
			return
		}

		// save file
		url, err := utils.UploadToBunny(file, file_headers)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		input.Image = &url
	} else {
		input.Image = nil
	}

	if err := h.service.Create(input, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Payment method created successfully"})
}

func (h *PaymentMethodHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method ID"})
		return
	}

	var input UpdatePaymentMethodInput
	if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {
		handleValidationError(c, err)
		return
	}

	file_headers, err := c.FormFile("image")
	if err == nil {
		if file_headers.Size > 2*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Max image size is 2MB"})
			return
		}

		// ✅ 2. MIME type check
		file, err := file_headers.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}
		defer file.Close()

		// Detect content type
		buffer := make([]byte, 512)
		if _, err := file.Read(buffer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file content"})
			return
		}
		contentType := http.DetectContentType(buffer)

		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/jpg":  true,
		}
		if !allowedTypes[contentType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only PNG, JPG formats are supported"})
			return
		}

		// save file
		url, err := utils.UploadToBunny(file, file_headers)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		input.Image = &url
	} else {
		input.Image = nil
	}

	if err := h.service.Update(uint(id), input, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Payment method created successfully"})
}

func (h *PaymentMethodHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method ID"})
		return
	}
	if err := h.service.Delete(uint(id), c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment method deleted successfully"})
}

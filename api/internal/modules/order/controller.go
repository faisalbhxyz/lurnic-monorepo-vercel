package order

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	service OrderService
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{
		service: NewOrderService(db),
	}
}

func (h *OrderHandler) GetAll(c *gin.Context) {
	orders, err := h.service.GetAll(c.GetUint("tenant_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func (h *OrderHandler) Create(c *gin.Context) {
	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	tenantID := c.GetUint("tenant_id")
	studentID := c.GetUint("user_id")

	res, err := h.service.Create(input, tenantID, studentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully", "order": res})
}

func (h *OrderHandler) MarkAsPaid(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
	}

	if err := h.service.MarkAsPaid(c.GetUint("tenant_id"), uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order marked as paid successfully"})
}

func (h *OrderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	if err := h.service.Delete(c.GetUint("tenant_id"), uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

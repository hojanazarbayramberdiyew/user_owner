package handler

import (
	"context"
	"net/http"
	"user_owner/internal/dto"
	"user_owner/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// CreateOrder creates a new order
// @Summary      Create a new order
// @Description  Create a new order with user information and location details
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization header string true "Bearer token"
// @Param        request body dto.CreateOrderRequest true "Order creation request"
// @Success      201  {object}  map[string]interface{}  "Order created successfully"
// @Failure      400  {object}  map[string]interface{}  "Invalid request"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /protected/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userPhone, exists := c.Get("user_phone")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user phone not found in token"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), userID, userPhone.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"data":    order,
	})
}

// GetAllOrders gets all orders
// @Summary Get all orders
// @Description Retrieve all orders from the database
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} map[string]interface{} "List of orders"
// @Failure 500 {object} map[string]interface{} "Failed to get orders"
// @Router /protected/orders [get]
func (h *OrderHandler) GetAllOrders(c *gin.Context) {

	ctx := context.Background()
	orders, err := h.orderService.GetAllOrders(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get all orders: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)

}

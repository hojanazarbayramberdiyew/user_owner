package handler

import (
	"context"
	"net/http"
	"user_owner/internal/dto"
	"user_owner/internal/service"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *userHandler {
	return &userHandler{service: service}
}

// CreateUser creates a new user
// @Summary      Create a new user
// @Description  Create a new user with name, password, phone number, email and logo
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateUserReq true "User creation request"
// @Success      201  {object}  map[string]interface{}  "User created successfully"
// @Failure      400  {object}  map[string]interface{}  "Invalid request body or missing required fields"
// @Failure      500  {object}  map[string]interface{}  "Failed to create user"
// @Router       /users [post]
func (h *userHandler) CreateUser(c *gin.Context) {

	var req dto.CreateUserReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if req.Name == "" || req.Password == "" || req.Email == "" || req.PhoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Name, password and email are required",
		})
		return
	}

	err := h.service.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    req})

}

// GetAllUsers gets all users
// @Summary      Get all users
// @Description  Retrieve all users from the database
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   map[string]interface{} "List of users"
// @Failure      500  {object}  map[string]interface{}  "Failed to get users"
// @Router       /users [get]
func (h *userHandler) GetAllUsers(c *gin.Context) {
	ctx := context.Background()
	users, err := h.service.GetUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get all users: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)

}

// Login user
// @Summary      User login
// @Description  Authenticate user and return token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginReq true "Login request"
// @Success      200  {object}  map[string]interface{}  "Login successful"
// @Failure      400  {object}  map[string]interface{}  "Invalid request"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /login [post]
func (h *userHandler) Login(c *gin.Context) {
	var loginReq dto.LoginReq
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if loginReq.Username == "" || loginReq.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username and password are required",
		})
		return
	}

	loginResp, err := h.service.Login(c.Request.Context(), &loginReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loginResp)

}

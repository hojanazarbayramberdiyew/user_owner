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

func (h *userHandler) GetAllUsers(c *gin.Context) {
	ctx := context.Background()
	users, err := h.service.GetUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get all users: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)

}

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

// Copyright 2026. Kimjibeom. All rights reserved.
package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.StaffLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	token, staff, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"staff": staff,
	})
}

// Register handles POST /api/auth/register (admin only)
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.StaffCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	staff, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		log.Printf("Register error: %v", err)
		if strings.Contains(err.Error(), "staffs_email_key") || strings.Contains(err.Error(), "duplicate key value") {
			c.JSON(http.StatusConflict, gin.H{"error": "이미 등록된 이메일입니다."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create staff member"})
		return
	}

	c.JSON(http.StatusCreated, staff)
}

// Me handles GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	staffID, _ := c.Get("staff_id")
	role, _ := c.Get("role")
	email, _ := c.Get("email")
	c.JSON(http.StatusOK, gin.H{
		"staff_id": staffID,
		"email":    email,
		"role":     role,
	})
}

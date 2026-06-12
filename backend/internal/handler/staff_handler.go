// Copyright 2026. Kimjibeom. All rights reserved.
package handler

import (
	"net/http"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type StaffHandler struct {
	staffService *service.StaffService
}

func NewStaffHandler(staffService *service.StaffService) *StaffHandler {
	return &StaffHandler{staffService: staffService}
}

func (h *StaffHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	staff, err := h.staffService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}
	c.JSON(http.StatusOK, staff)
}

func (h *StaffHandler) List(c *gin.Context) {
	staffs, err := h.staffService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch staff list"})
		return
	}
	c.JSON(http.StatusOK, staffs)
}

func (h *StaffHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req model.StaffUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.staffService.Update(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Staff updated"})
}

func (h *StaffHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.staffService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete staff"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Staff deactivated"})
}

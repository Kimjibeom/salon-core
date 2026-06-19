// Copyright 2026. Kimjibeom. All rights reserved.
package handler

import (
	"net/http"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ServiceMenuHandler struct {
	serviceService *service.ServiceMenuService
}

func NewServiceMenuHandler(serviceService *service.ServiceMenuService) *ServiceMenuHandler {
	return &ServiceMenuHandler{serviceService: serviceService}
}

func (h *ServiceMenuHandler) Create(c *gin.Context) {
	var req model.ServiceMenuCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	svc, err := h.serviceService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service menu"})
		return
	}

	c.JSON(http.StatusCreated, svc)
}

func (h *ServiceMenuHandler) GetAll(c *gin.Context) {
	services, err := h.serviceService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch service menus"})
		return
	}

	c.JSON(http.StatusOK, services)
}

func (h *ServiceMenuHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req model.ServiceMenuUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	svc, err := h.serviceService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service menu"})
		return
	}

	c.JSON(http.StatusOK, svc)
}

func (h *ServiceMenuHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.serviceService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service menu"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

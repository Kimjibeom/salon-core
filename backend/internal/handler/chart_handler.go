// Copyright 2026. Kimjibeom. All rights reserved.
package handler

import (
	"net/http"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ChartHandler struct {
	chartService *service.ChartService
}

func NewChartHandler(chartService *service.ChartService) *ChartHandler {
	return &ChartHandler{chartService: chartService}
}

func (h *ChartHandler) Create(c *gin.Context) {
	var req model.ChartCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	chart, err := h.chartService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chart"})
		return
	}
	c.JSON(http.StatusCreated, chart)
}

func (h *ChartHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	chart, err := h.chartService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chart not found"})
		return
	}
	c.JSON(http.StatusOK, chart)
}

func (h *ChartHandler) ListByCustomer(c *gin.Context) {
	customerID := c.Param("customerId")
	charts, err := h.chartService.ListByCustomer(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch charts"})
		return
	}
	c.JSON(http.StatusOK, charts)
}

func (h *ChartHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.chartService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chart"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Chart deleted"})
}

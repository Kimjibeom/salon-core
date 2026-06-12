// Copyright 2026. Kimjibeom. All rights reserved.
package handler

import (
	"net/http"

	"github.com/Kimjibeom/salon-core/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type MarketingHandler struct {
	marketingService *service.MarketingService
}

func NewMarketingHandler(marketingService *service.MarketingService) *MarketingHandler {
	return &MarketingHandler{marketingService: marketingService}
}

// ExtractTargets handles POST /api/marketing/targets
func (h *MarketingHandler) ExtractTargets(c *gin.Context) {
	var filters map[string]interface{}
	if err := c.ShouldBindJSON(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter criteria"})
		return
	}

	customers, err := h.marketingService.ExtractTargetCustomers(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract target customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":     len(customers),
		"customers": customers,
	})
}

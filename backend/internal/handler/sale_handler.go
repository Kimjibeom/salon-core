// Copyright 2026. Kimjibeom. All rights reserved.
package handler

import (
	"net/http"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type SaleHandler struct {
	saleService *service.SaleService
}

func NewSaleHandler(saleService *service.SaleService) *SaleHandler {
	return &SaleHandler{saleService: saleService}
}

func (h *SaleHandler) Create(c *gin.Context) {
	var req model.SaleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	sale, err := h.saleService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sale"})
		return
	}
	c.JSON(http.StatusCreated, sale)
}

func (h *SaleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	sale, err := h.saleService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sale not found"})
		return
	}
	c.JSON(http.StatusOK, sale)
}

func (h *SaleHandler) ListByDateRange(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date parameters required"})
		return
	}
	sales, err := h.saleService.ListByDateRange(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales"})
		return
	}
	c.JSON(http.StatusOK, sales)
}

func (h *SaleHandler) ListByStaff(c *gin.Context) {
	staffID := c.Param("staffId")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date parameters required"})
		return
	}
	sales, err := h.saleService.ListByStaff(c.Request.Context(), staffID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales"})
		return
	}
	c.JSON(http.StatusOK, sales)
}

func (h *SaleHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.saleService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sale"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sale deleted"})
}

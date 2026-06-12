// Copyright 2026. Kimjibeom. All rights reserved.
package handler

import (
	"net/http"
	"strconv"

	"github.com/Kimjibeom/salon-core/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) GetDailySummary(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date parameter required"})
		return
	}
	summary, err := h.analyticsService.GetDailySummary(c.Request.Context(), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch daily summary"})
		return
	}
	c.JSON(http.StatusOK, summary)
}

func (h *AnalyticsHandler) GetMonthlyTrend(c *gin.Context) {
	year, _ := strconv.Atoi(c.Query("year"))
	month, _ := strconv.Atoi(c.Query("month"))
	if year == 0 || month == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month parameters required"})
		return
	}
	trend, err := h.analyticsService.GetMonthlyTrend(c.Request.Context(), year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch monthly trend"})
		return
	}
	c.JSON(http.StatusOK, trend)
}

func (h *AnalyticsHandler) GetAnalyticsSummary(c *gin.Context) {
	months, _ := strconv.Atoi(c.DefaultQuery("months", "6"))
	if months <= 0 || months > 24 {
		months = 6
	}
	summary, err := h.analyticsService.GetAnalyticsSummary(c.Request.Context(), months)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch analytics summary"})
		return
	}
	c.JSON(http.StatusOK, summary)
}

func (h *AnalyticsHandler) GetPaymentBreakdown(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date parameters required"})
		return
	}
	breakdown, err := h.analyticsService.GetPaymentBreakdown(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment breakdown"})
		return
	}
	c.JSON(http.StatusOK, breakdown)
}

func (h *AnalyticsHandler) GetCustomerStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date parameters required"})
		return
	}
	newCount, returningCount, err := h.analyticsService.GetNewVsReturning(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customer stats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"new_customers":       newCount,
		"returning_customers": returningCount,
		"total":               newCount + returningCount,
	})
}

func (h *AnalyticsHandler) GetChurnAnalysis(c *gin.Context) {
	inactiveDays, _ := strconv.Atoi(c.DefaultQuery("inactive_days", "90"))
	if inactiveDays <= 0 {
		inactiveDays = 90
	}
	analysis, err := h.analyticsService.GetChurnAnalysis(c.Request.Context(), inactiveDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch churn analysis"})
		return
	}
	c.JSON(http.StatusOK, analysis)
}

func (h *AnalyticsHandler) GetStaffPerformance(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date parameters required"})
		return
	}
	performance, err := h.analyticsService.GetStaffPerformance(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch staff performance"})
		return
	}
	c.JSON(http.StatusOK, performance)
}

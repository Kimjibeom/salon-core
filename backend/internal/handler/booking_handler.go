// Copyright 2026. Kimjibeom. All rights reserved.
package handler

import (
	"net/http"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// BookingHandler handles public booking API endpoints.
type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

// GetAvailability returns available time slots for a given date, staff, and service.
// GET /api/booking/availability?date=YYYY-MM-DD&staff_id=UUID&service_id=UUID
func (h *BookingHandler) GetAvailability(c *gin.Context) {
	date := c.Query("date")
	staffID := c.Query("staff_id")
	serviceID := c.Query("service_id")

	if date == "" || staffID == "" || serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date, staff_id, and service_id are required"})
		return
	}

	availability, err := h.bookingService.GetAvailability(c.Request.Context(), date, staffID, serviceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, availability)
}

// CreatePublicReservation creates a reservation from the public booking site.
// POST /api/booking/reservations
func (h *BookingHandler) CreatePublicReservation(c *gin.Context) {
	var req model.PublicBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "입력값을 확인해주세요: " + err.Error()})
		return
	}

	reservation, err := h.bookingService.CreatePublicBooking(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reservation)
}

// GetStaffList returns active designers for the public booking page.
// GET /api/booking/staff
func (h *BookingHandler) GetStaffList(c *gin.Context) {
	staffs, err := h.bookingService.GetPublicStaffList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch staff list"})
		return
	}
	c.JSON(http.StatusOK, staffs)
}

// GetServiceList returns active services for the public booking page.
// GET /api/booking/services
func (h *BookingHandler) GetServiceList(c *gin.Context) {
	services, err := h.bookingService.GetPublicServiceList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch service list"})
		return
	}
	c.JSON(http.StatusOK, services)
}

// GetShopInfo returns basic shop info for the booking page header.
// GET /api/booking/shop
func (h *BookingHandler) GetShopInfo(c *gin.Context) {
	info, err := h.bookingService.GetShopInfo(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shop info"})
		return
	}
	c.JSON(http.StatusOK, info)
}

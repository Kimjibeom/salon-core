// Copyright 2026. Kimjibeom. All rights reserved.
package handler

import (
	"net/http"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ReservationHandler struct {
	reservationService *service.ReservationService
}

func NewReservationHandler(reservationService *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

func (h *ReservationHandler) Create(c *gin.Context) {
	var req model.ReservationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	res, err := h.reservationService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *ReservationHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.reservationService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *ReservationHandler) ListByDate(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date parameter required"})
		return
	}
	reservations, err := h.reservationService.ListByDate(c.Request.Context(), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
		return
	}
	c.JSON(http.StatusOK, reservations)
}

func (h *ReservationHandler) ListByDateRange(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date parameters required"})
		return
	}
	reservations, err := h.reservationService.ListByDateRange(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
		return
	}
	c.JSON(http.StatusOK, reservations)
}

func (h *ReservationHandler) GetWaitingQueue(c *gin.Context) {
	queue, err := h.reservationService.GetWaitingQueue(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch waiting queue"})
		return
	}
	c.JSON(http.StatusOK, queue)
}

func (h *ReservationHandler) AddToWaitingQueue(c *gin.Context) {
	var req model.ReservationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	res, err := h.reservationService.AddToWaitingQueue(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to waiting queue"})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *ReservationHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Status string `json:"status" binding:"required,oneof=reserved waiting in_progress completed canceled no_show"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}
	if err := h.reservationService.UpdateStatus(c.Request.Context(), id, body.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

func (h *ReservationHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req model.ReservationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.reservationService.Update(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reservation"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reservation updated"})
}

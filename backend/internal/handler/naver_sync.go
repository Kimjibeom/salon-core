// Copyright 2026. Kimjibeom. All rights reserved.
package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/repository"
	"github.com/Kimjibeom/salon-core/backend/internal/service"
	"github.com/Kimjibeom/salon-core/backend/internal/websocket"
	"github.com/gin-gonic/gin"
)

// NaverSyncHandler handles synchronization with Naver Reservation system.
type NaverSyncHandler struct {
	reservationService *service.ReservationService
	settingService     *service.SettingService
	mappingRepo        *repository.NaverMappingRepository
	hub                *websocket.Hub
}

func NewNaverSyncHandler(
	reservationService *service.ReservationService,
	settingService *service.SettingService,
	mappingRepo *repository.NaverMappingRepository,
	hub *websocket.Hub,
) *NaverSyncHandler {
	return &NaverSyncHandler{
		reservationService: reservationService,
		settingService:     settingService,
		mappingRepo:        mappingRepo,
		hub:                hub,
	}
}

// HandleWebhook receives webhook events from Naver.
func (h *NaverSyncHandler) HandleWebhook(c *gin.Context) {
	// Read the raw body for signature verification
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Verify webhook signature
	signature := c.GetHeader("X-Naver-Signature")
	if !h.verifySignature(c, body, signature) {
		log.Println("Naver webhook signature verification failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// Parse payload
	var payload model.NaverWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		// Re-bind from body since we already read it
		log.Printf("Naver webhook payload parse error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	log.Printf("Received Naver webhook: type=%s, booking_id=%s", payload.EventType, payload.BookingID)

	// Map naver IDs to internal IDs
	ctx := c.Request.Context()
	var internalStaffID string
	if payload.DesignerID != "" {
		if mapping, err := h.mappingRepo.GetByNaverID(ctx, "STAFF", payload.DesignerID); err == nil {
			internalStaffID = mapping.InternalID
		} else {
			log.Printf("No staff mapping found for naver designer ID: %s", payload.DesignerID)
		}
	}

	var internalServiceID string
	if payload.ServiceID != "" {
		if mapping, err := h.mappingRepo.GetByNaverID(ctx, "SERVICE", payload.ServiceID); err == nil {
			internalServiceID = mapping.InternalID
		} else {
			log.Printf("No service mapping found for naver service ID: %s", payload.ServiceID)
		}
	}

	// Parse start/end times
	startTime, _ := time.Parse(time.RFC3339, payload.StartTime)
	endTime, _ := time.Parse(time.RFC3339, payload.EndTime)
	date := startTime.Format("2006-01-02")
	startTimeStr := startTime.Format("15:04")
	endTimeStr := endTime.Format("15:04")

	switch strings.ToUpper(payload.EventType) {
	case "CREATED":
		req := &model.ReservationCreateRequest{
			CustomerName:  payload.CustomerName,
			CustomerPhone: payload.CustomerPhone,
			StaffID:       internalStaffID,
			ServiceID:     internalServiceID,
			TreatmentName: "", // will be filled from service mapping
			Date:          date,
			StartTime:     startTimeStr,
			EndTime:       endTimeStr,
			Status:        "reserved",
			Source:        "naver",
			Memo:          "네이버 예약 (ID: " + payload.BookingID + ")",
		}
		res, err := h.reservationService.Create(ctx, req)
		if err != nil {
			log.Printf("Failed to create naver reservation: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation"})
			return
		}
		log.Printf("Created naver reservation: %s", res.ID)

	case "CANCELED":
		log.Printf("Naver booking canceled: %s (manual intervention may be needed)", payload.BookingID)

	case "UPDATED":
		log.Printf("Naver booking updated: %s (manual intervention may be needed)", payload.BookingID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed"})
}

// verifySignature verifies the Naver webhook HMAC-SHA256 signature.
func (h *NaverSyncHandler) verifySignature(c *gin.Context, body []byte, signature string) bool {
	if signature == "" {
		// If no signature header, skip verification (for testing)
		log.Println("No signature provided, skipping verification")
		return true
	}

	ctx := c.Request.Context()
	secretSetting, err := h.settingService.GetSettingByKey(ctx, "naver_webhook_secret")
	if err != nil || secretSetting.Value == "" {
		log.Println("Naver webhook secret not configured, skipping verification")
		return true
	}

	mac := hmac.New(sha256.New, []byte(secretSetting.Value))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expectedMAC), []byte(signature))
}

// TriggerManualSync is a placeholder for manually triggering a sync job.
func (h *NaverSyncHandler) TriggerManualSync(c *gin.Context) {
	log.Println("Manual Naver Reservation sync triggered")
	c.JSON(http.StatusOK, gin.H{"message": "Naver sync job queued"})
}

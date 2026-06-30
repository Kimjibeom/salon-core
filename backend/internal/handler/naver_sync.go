package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NaverSyncHandler handles synchronization with Naver Reservation system.
// Since Naver Reservation does not provide an open public API for all users,
// this is a placeholder for where a webhook receiver or polling job would hook in,
// typically utilizing a partnered API integration or a crawling-based adapter.

type NaverSyncHandler struct{}

func NewNaverSyncHandler() *NaverSyncHandler {
	return &NaverSyncHandler{}
}

// HandleWebhook is a placeholder for receiving webhook events from Naver.
func (h *NaverSyncHandler) HandleWebhook(c *gin.Context) {
	// TODO: Parse Naver webhook payload
	// Verify signature
	// Map to internal model.Reservation
	// Save to database
	log.Println("Received Naver Reservation webhook sync request")

	c.JSON(http.StatusOK, gin.H{"message": "Naver sync placeholder handled"})
}

// TriggerManualSync is a placeholder for manually triggering a sync job.
func (h *NaverSyncHandler) TriggerManualSync(c *gin.Context) {
	log.Println("Manual Naver Reservation sync triggered")

	c.JSON(http.StatusOK, gin.H{"message": "Naver sync job queued"})
}

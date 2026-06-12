// Copyright 2026. Kimjibeom. All rights reserved.
package handler

import (
	"net/http"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type MembershipHandler struct {
	membershipService *service.MembershipService
}

func NewMembershipHandler(membershipService *service.MembershipService) *MembershipHandler {
	return &MembershipHandler{membershipService: membershipService}
}

func (h *MembershipHandler) Create(c *gin.Context) {
	var req model.MembershipCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	membership, err := h.membershipService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create membership"})
		return
	}
	c.JSON(http.StatusCreated, membership)
}

func (h *MembershipHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	membership, err := h.membershipService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Membership not found"})
		return
	}
	c.JSON(http.StatusOK, membership)
}

func (h *MembershipHandler) ListByCustomer(c *gin.Context) {
	customerID := c.Param("customerId")
	memberships, err := h.membershipService.ListByCustomer(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch memberships"})
		return
	}
	c.JSON(http.StatusOK, memberships)
}

func (h *MembershipHandler) GetTransactions(c *gin.Context) {
	id := c.Param("id")
	txns, err := h.membershipService.GetTransactions(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}
	c.JSON(http.StatusOK, txns)
}

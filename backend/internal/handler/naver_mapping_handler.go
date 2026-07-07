// Copyright 2026. Kimjibeom. All rights reserved.
package handler

import (
	"net/http"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

// NaverMappingHandler handles CRUD for naver_mappings.
type NaverMappingHandler struct {
	repo *repository.NaverMappingRepository
}

func NewNaverMappingHandler(repo *repository.NaverMappingRepository) *NaverMappingHandler {
	return &NaverMappingHandler{repo: repo}
}

// GetAll returns all naver mappings.
func (h *NaverMappingHandler) GetAll(c *gin.Context) {
	mappings, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch mappings"})
		return
	}
	if mappings == nil {
		mappings = []model.NaverMapping{}
	}
	c.JSON(http.StatusOK, mappings)
}

// Create creates a new naver mapping.
func (h *NaverMappingHandler) Create(c *gin.Context) {
	var req model.NaverMappingCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	mapping := &model.NaverMapping{
		InternalType: req.InternalType,
		InternalID:   req.InternalID,
		NaverID:      req.NaverID,
	}

	if err := h.repo.Create(c.Request.Context(), mapping); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create mapping: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapping)
}

// Delete removes a naver mapping.
func (h *NaverMappingHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete mapping"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Mapping deleted"})
}

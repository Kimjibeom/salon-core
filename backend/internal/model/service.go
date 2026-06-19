// Copyright 2026. Kimjibeom. All rights reserved.
package model

import (
	"time"
)

// ServiceMenu represents a salon service/product menu item.
type ServiceMenu struct {
	ID        string    `json:"id"`
	Category  string    `json:"category"` // 'cut', 'perm', 'color', 'clinic', 'product', etc.
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Duration  int       `json:"duration"` // in minutes
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ServiceMenuCreateRequest holds data for creating a new service menu item.
type ServiceMenuCreateRequest struct {
	Category string  `json:"category" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required,min=0"`
	Duration int     `json:"duration" binding:"required,min=0"`
	IsActive bool    `json:"is_active"`
}

// ServiceMenuUpdateRequest holds data for updating an existing service menu item.
type ServiceMenuUpdateRequest struct {
	Category *string  `json:"category"`
	Name     *string  `json:"name"`
	Price    *float64 `json:"price"`
	Duration *int     `json:"duration"`
	IsActive *bool    `json:"is_active"`
}

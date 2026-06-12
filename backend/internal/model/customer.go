// Copyright 2026. Kimjibeom. All rights reserved.
package model

import (
	"time"
)

// Customer represents a salon customer.
type Customer struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Phone         string     `json:"phone"`
	Email         string     `json:"email,omitempty"`
	BirthDate     *time.Time `json:"birth_date,omitempty"`
	Memo          string     `json:"memo"`
	Tags          []string   `json:"tags"`
	CreatedAt     time.Time  `json:"created_at"`
	LastVisitedAt *time.Time `json:"last_visited_at,omitempty"`
	VisitCount    int        `json:"visit_count"`
	IsDeleted     bool       `json:"-"`
}

// CustomerCreateRequest is the request body for creating a customer.
type CustomerCreateRequest struct {
	Name      string   `json:"name" binding:"required,max=100"`
	Phone     string   `json:"phone" binding:"required,max=20"`
	Email     string   `json:"email" binding:"omitempty,email,max=255"`
	BirthDate string   `json:"birth_date"` // YYYY-MM-DD
	Memo      string   `json:"memo" binding:"max=2000"`
	Tags      []string `json:"tags"`
}

// CustomerUpdateRequest is the request body for updating a customer.
type CustomerUpdateRequest struct {
	Name      *string  `json:"name" binding:"omitempty,max=100"`
	Phone     *string  `json:"phone" binding:"omitempty,max=20"`
	Email     *string  `json:"email" binding:"omitempty,email,max=255"`
	BirthDate *string  `json:"birth_date"`
	Memo      *string  `json:"memo" binding:"omitempty,max=2000"`
	Tags      []string `json:"tags"`
}

// CustomerSearchRequest holds the search criteria for customers.
type CustomerSearchRequest struct {
	Query  string `form:"q" binding:"max=100"`     // name or phone suffix
	Limit  int    `form:"limit" binding:"max=100"`
	Offset int    `form:"offset" binding:"gte=0"`
}

// CustomerWithHistory is used for detailed customer views with chart history.
type CustomerWithHistory struct {
	Customer
	Charts      []Chart      `json:"charts"`
	Memberships []Membership `json:"memberships"`
	Sales       []Sale       `json:"recent_sales"`
}

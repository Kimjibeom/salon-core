// Copyright 2026. Kimjibeom. All rights reserved.
package model

import (
	"time"
)

// Membership represents a prepaid or count-based membership.
type Membership struct {
	ID               string     `json:"id"`
	CustomerID       string     `json:"customer_id"`
	Type             string     `json:"type"` // money or count
	Name             string     `json:"name"`
	InitialBalance   float64    `json:"initial_balance"`
	Balance          float64    `json:"balance"`
	InitialCount     int        `json:"initial_count"`
	RemainingCount   int        `json:"remaining_count"`
	TargetTreatment  string     `json:"target_treatment,omitempty"`
	ExpiredAt        *time.Time `json:"expired_at,omitempty"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// MembershipCreateRequest is the request body for creating a membership.
type MembershipCreateRequest struct {
	CustomerID      string  `json:"customer_id" binding:"required,uuid"`
	Type            string  `json:"type" binding:"required,oneof=money count"`
	Name            string  `json:"name" binding:"required,max=200"`
	InitialBalance  float64 `json:"initial_balance" binding:"gte=0"`
	InitialCount    int     `json:"initial_count" binding:"gte=0"`
	TargetTreatment string  `json:"target_treatment" binding:"max=200"`
	ExpiredAt       string  `json:"expired_at"` // RFC3339
}

// MembershipTransaction represents a single deduction/addition to a membership.
type MembershipTransaction struct {
	ID            string    `json:"id"`
	MembershipID  string    `json:"membership_id"`
	SaleID        string    `json:"sale_id,omitempty"`
	Amount        float64   `json:"amount"`
	CountChange   int       `json:"count_change"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CountBefore   int       `json:"count_before"`
	CountAfter    int       `json:"count_after"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

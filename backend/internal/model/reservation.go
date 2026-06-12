// Copyright 2026. Kimjibeom. All rights reserved.
package model

import (
	"time"
)

// Reservation represents a booking or walk-in entry.
type Reservation struct {
	ID              string     `json:"id"`
	CustomerID      *string    `json:"customer_id,omitempty"`
	StaffID         *string    `json:"staff_id,omitempty"`
	CustomerName    string     `json:"customer_name"`
	CustomerPhone   string     `json:"customer_phone"`
	StaffName       string     `json:"staff_name,omitempty"`
	TreatmentName   string     `json:"treatment_name,omitempty"`
	Date            string     `json:"date"` // YYYY-MM-DD
	StartTime       string     `json:"start_time"` // HH:MM
	EndTime         string     `json:"end_time"` // HH:MM
	Status          string     `json:"status"` // reserved, waiting, in_progress, completed, canceled, no_show
	Source          string     `json:"source"` // online, offline, naver
	WaitingNumber   *int       `json:"waiting_number,omitempty"`
	WaitingStartedAt *time.Time `json:"waiting_started_at,omitempty"`
	Memo            string     `json:"memo"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ReservationCreateRequest is the request body for creating a reservation.
type ReservationCreateRequest struct {
	CustomerID    string `json:"customer_id" binding:"omitempty,uuid"`
	StaffID       string `json:"staff_id" binding:"omitempty,uuid"`
	CustomerName  string `json:"customer_name" binding:"required,max=100"`
	CustomerPhone string `json:"customer_phone" binding:"required,max=20"`
	TreatmentName string `json:"treatment_name" binding:"max=200"`
	Date          string `json:"date" binding:"required"` // YYYY-MM-DD
	StartTime     string `json:"start_time" binding:"required"` // HH:MM
	EndTime       string `json:"end_time" binding:"required"` // HH:MM
	Status        string `json:"status" binding:"omitempty,oneof=reserved waiting"`
	Source        string `json:"source" binding:"required,oneof=online offline naver"`
	Memo          string `json:"memo" binding:"max=2000"`
}

// ReservationUpdateRequest is the request body for updating a reservation.
type ReservationUpdateRequest struct {
	StaffID       *string `json:"staff_id" binding:"omitempty,uuid"`
	TreatmentName *string `json:"treatment_name" binding:"omitempty,max=200"`
	Date          *string `json:"date"`
	StartTime     *string `json:"start_time"`
	EndTime       *string `json:"end_time"`
	Status        *string `json:"status" binding:"omitempty,oneof=reserved waiting in_progress completed canceled no_show"`
	Memo          *string `json:"memo" binding:"omitempty,max=2000"`
}

// ReservationQuery holds filter criteria for fetching reservations.
type ReservationQuery struct {
	Date    string `form:"date"`    // YYYY-MM-DD
	StartDate string `form:"start_date"` // YYYY-MM-DD
	EndDate   string `form:"end_date"`   // YYYY-MM-DD
	StaffID string `form:"staff_id" binding:"omitempty,uuid"`
	Status  string `form:"status" binding:"omitempty,oneof=reserved waiting in_progress completed canceled no_show"`
}

// WaitingQueueEntry represents a single entry in the walk-in waiting queue.
type WaitingQueueEntry struct {
	Reservation
	WaitTimeMinutes int `json:"wait_time_minutes"`
	Position        int `json:"position"`
}

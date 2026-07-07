// Copyright 2026. Kimjibeom. All rights reserved.
package model

import "time"

// NaverMapping links internal staff/service IDs to Naver-side IDs.
type NaverMapping struct {
	ID           string    `json:"id"`
	InternalType string    `json:"internal_type"` // STAFF or SERVICE
	InternalID   string    `json:"internal_id"`
	NaverID      string    `json:"naver_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// NaverMappingCreateRequest is the request body for creating a mapping.
type NaverMappingCreateRequest struct {
	InternalType string `json:"internal_type" binding:"required,oneof=STAFF SERVICE"`
	InternalID   string `json:"internal_id" binding:"required,uuid"`
	NaverID      string `json:"naver_id" binding:"required"`
}

// TimeSlot represents a single available booking time slot.
type TimeSlot struct {
	StartTime string `json:"start_time"` // HH:MM
	EndTime   string `json:"end_time"`   // HH:MM
	Available bool   `json:"available"`
}

// AvailabilityResponse is the response for the availability API.
type AvailabilityResponse struct {
	Date      string     `json:"date"`
	StaffID   string     `json:"staff_id"`
	StaffName string     `json:"staff_name"`
	Slots     []TimeSlot `json:"slots"`
}

// PublicBookingRequest is the request body for a public booking.
type PublicBookingRequest struct {
	CustomerName  string `json:"customer_name" binding:"required,max=100"`
	CustomerPhone string `json:"customer_phone" binding:"required,max=20"`
	ServiceID     string `json:"service_id" binding:"required,uuid"`
	StaffID       string `json:"staff_id" binding:"required,uuid"`
	Date          string `json:"date" binding:"required"`
	StartTime     string `json:"start_time" binding:"required"`
	Memo          string `json:"memo" binding:"max=2000"`
}

// NaverWebhookPayload is the expected payload from Naver webhook.
type NaverWebhookPayload struct {
	EventType    string                `json:"event_type"`    // CREATED, UPDATED, CANCELED
	BookingID    string                `json:"booking_id"`
	PlaceID      string                `json:"place_id"`
	CustomerName string                `json:"customer_name"`
	CustomerPhone string               `json:"customer_phone"`
	DesignerID   string                `json:"designer_id"`   // Naver-side designer ID
	ServiceID    string                `json:"service_id"`    // Naver-side service ID
	StartTime    string                `json:"start_time"`    // ISO 8601
	EndTime      string                `json:"end_time"`      // ISO 8601
	Status       string                `json:"status"`
	Signature    string                `json:"signature"`
}

// PublicStaffResponse is a simplified staff response for public booking.
type PublicStaffResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	DayOff []int  `json:"day_off"`
}

// PublicServiceResponse is a simplified service response for public booking.
type PublicServiceResponse struct {
	ID       string  `json:"id"`
	Category string  `json:"category"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Duration int     `json:"duration"`
}

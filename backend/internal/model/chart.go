// Copyright 2026. Kimjibeom. All rights reserved.
package model

import (
	"time"
)

// Chart represents a treatment history record.
type Chart struct {
	ID                string    `json:"id"`
	CustomerID        string    `json:"customer_id"`
	StaffID           string    `json:"staff_id"`
	StaffName         string    `json:"staff_name,omitempty"`
	ServiceID         *string   `json:"service_id,omitempty"`
	Recipe            string    `json:"recipe"`
	TreatmentName     string    `json:"treatment_name,omitempty"`
	TreatmentDuration int       `json:"treatment_duration,omitempty"` // minutes
	Notes             string    `json:"notes"`
	BeforeImgURL      string    `json:"before_img_url,omitempty"`
	AfterImgURL       string    `json:"after_img_url,omitempty"`
	ConsentDocURL     string    `json:"consent_doc_url,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

// ChartCreateRequest is the request body for creating a chart entry.
type ChartCreateRequest struct {
	CustomerID        string `json:"customer_id" binding:"required,uuid"`
	StaffID           string `json:"staff_id" binding:"required,uuid"`
	ServiceID         string `json:"service_id" binding:"omitempty,uuid"`
	Recipe            string `json:"recipe" binding:"max=5000"`
	TreatmentName     string `json:"treatment_name" binding:"max=200"`
	TreatmentDuration int    `json:"treatment_duration" binding:"gte=0"`
	Notes             string `json:"notes" binding:"max=5000"`
	BeforeImgURL      string `json:"before_img_url" binding:"omitempty,url,max=2048"`
	AfterImgURL       string `json:"after_img_url" binding:"omitempty,url,max=2048"`
	ConsentDocURL     string `json:"consent_doc_url" binding:"omitempty,url,max=2048"`
}

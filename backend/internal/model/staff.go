// Copyright 2026. Kimjibeom. All rights reserved.
package model

import (
	"time"
)

// Staff represents a salon employee.
type Staff struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Email                string    `json:"email,omitempty"`
	PasswordHash         string    `json:"-"` // never exposed in JSON
	Role                 string    `json:"role"` // admin, designer, staff
	Phone                string    `json:"phone,omitempty"`
	ServiceIncentiveRate float64   `json:"service_incentive_rate"`
	ProductIncentiveRate float64   `json:"product_incentive_rate"`
	MonthlyTarget        float64   `json:"monthly_target"`
	DayOff               []int     `json:"day_off"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// StaffCreateRequest is the request body for creating a staff member.
type StaffCreateRequest struct {
	Name                 string  `json:"name" binding:"required,max=100"`
	Email                string  `json:"email" binding:"required,email,max=255"`
	Password             string  `json:"password" binding:"required,min=8,max=128"`
	Role                 string  `json:"role" binding:"required,oneof=admin designer staff"`
	Phone                string  `json:"phone" binding:"max=20"`
	ServiceIncentiveRate float64 `json:"service_incentive_rate" binding:"gte=0,lte=100"`
	ProductIncentiveRate float64 `json:"product_incentive_rate" binding:"gte=0,lte=100"`
	MonthlyTarget        float64 `json:"monthly_target" binding:"gte=0"`
}

// StaffUpdateRequest is the request body for updating a staff member.
type StaffUpdateRequest struct {
	Name                 *string  `json:"name" binding:"omitempty,max=100"`
	Role                 *string  `json:"role" binding:"omitempty,oneof=admin designer staff"`
	Phone                *string  `json:"phone" binding:"omitempty,max=20"`
	ServiceIncentiveRate *float64 `json:"service_incentive_rate" binding:"omitempty,gte=0,lte=100"`
	ProductIncentiveRate *float64 `json:"product_incentive_rate" binding:"omitempty,gte=0,lte=100"`
	MonthlyTarget        *float64 `json:"monthly_target" binding:"omitempty,gte=0"`
	DayOff               *[]int   `json:"day_off" binding:"omitempty,dive,gte=0,lte=6"`
	IsActive             *bool    `json:"is_active"`
}

// StaffLoginRequest is the request body for staff authentication.
type StaffLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// StaffPerformance holds computed performance metrics for a staff member.
type StaffPerformance struct {
	StaffID           string  `json:"staff_id"`
	StaffName         string  `json:"staff_name"`
	MonthlyTarget     float64 `json:"monthly_target"`
	ServiceRevenue    float64 `json:"service_revenue"`
	ProductRevenue    float64 `json:"product_revenue"`
	TotalRevenue      float64 `json:"total_revenue"`
	AchievementRate   float64 `json:"achievement_rate"`
	ServiceIncentive  float64 `json:"service_incentive"`
	ProductIncentive  float64 `json:"product_incentive"`
	TotalIncentive    float64 `json:"total_incentive"`
	NetServiceRevenue float64 `json:"net_service_revenue"`
	NetProductRevenue float64 `json:"net_product_revenue"`
	TotalPayroll      float64 `json:"total_payroll"`
}

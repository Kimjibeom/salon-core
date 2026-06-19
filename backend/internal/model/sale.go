// Copyright 2026. Kimjibeom. All rights reserved.
package model

import (
	"time"
)

// Sale represents a sales transaction.
type Sale struct {
	ID                string    `json:"id"`
	ReservationID     *string   `json:"reservation_id,omitempty"`
	CustomerID        *string   `json:"customer_id,omitempty"`
	StaffID           string    `json:"staff_id"`
	StaffName         string    `json:"staff_name,omitempty"`
	CustomerName      string    `json:"customer_name,omitempty"`
	ServiceID         *string   `json:"service_id,omitempty"`
	ItemName          string    `json:"item_name"`
	TotalAmount       float64   `json:"total_amount"`
	Category          string    `json:"category"` // service, product
	PaymentMethod     string    `json:"payment_method"` // card, cash, membership, mixed
	CardAmount        float64   `json:"card_amount"`
	CashAmount        float64   `json:"cash_amount"`
	MembershipAmount  float64   `json:"membership_amount"`
	CardCommissionRate float64  `json:"card_commission_rate"`
	MembershipID      *string   `json:"membership_id,omitempty"`
	Memo              string    `json:"memo"`
	CreatedAt         time.Time `json:"created_at"`
}

// SaleCreateRequest is the request body for creating a sale.
type SaleCreateRequest struct {
	ReservationID      string  `json:"reservation_id" binding:"omitempty,uuid"`
	CustomerID         string  `json:"customer_id" binding:"omitempty,uuid"`
	ServiceID          string  `json:"service_id" binding:"omitempty,uuid"`
	StaffID            string  `json:"staff_id" binding:"required,uuid"`
	ItemName           string  `json:"item_name" binding:"max=200"`
	TotalAmount        float64 `json:"total_amount" binding:"required,gt=0"`
	Category           string  `json:"category" binding:"required,oneof=service product"`
	PaymentMethod      string  `json:"payment_method" binding:"required,oneof=card cash membership mixed"`
	CardAmount         float64 `json:"card_amount" binding:"gte=0"`
	CashAmount         float64 `json:"cash_amount" binding:"gte=0"`
	MembershipAmount   float64 `json:"membership_amount" binding:"gte=0"`
	CardCommissionRate float64 `json:"card_commission_rate" binding:"gte=0,lte=100"`
	MembershipID       string  `json:"membership_id" binding:"omitempty,uuid"`
	Memo               string  `json:"memo" binding:"max=2000"`
}

// DailySummary holds aggregated sales data for a single day.
type DailySummary struct {
	Date              string  `json:"date"`
	TotalRevenue      float64 `json:"total_revenue"`
	ServiceRevenue    float64 `json:"service_revenue"`
	ProductRevenue    float64 `json:"product_revenue"`
	CardRevenue       float64 `json:"card_revenue"`
	CashRevenue       float64 `json:"cash_revenue"`
	MembershipRevenue float64 `json:"membership_revenue"`
	TransactionCount  int     `json:"transaction_count"`
	CustomerCount     int     `json:"customer_count"`
}

// AnalyticsSummary holds aggregated analytics data.
type AnalyticsSummary struct {
	Period            string  `json:"period"`
	TotalRevenue      float64 `json:"total_revenue"`
	ServiceRevenue    float64 `json:"service_revenue"`
	ProductRevenue    float64 `json:"product_revenue"`
	AvgServicePrice   float64 `json:"avg_service_price"`
	AvgProductPrice   float64 `json:"avg_product_price"`
	CardRevenue       float64 `json:"card_revenue"`
	CashRevenue       float64 `json:"cash_revenue"`
	MembershipRevenue float64 `json:"membership_revenue"`
	NewCustomers      int     `json:"new_customers"`
	ReturningCustomers int    `json:"returning_customers"`
	TotalCustomers    int     `json:"total_customers"`
}

// ChurnAnalysis holds customer churn statistics.
type ChurnAnalysis struct {
	VisitNumber       int     `json:"visit_number"`
	TotalCustomers    int     `json:"total_customers"`
	ChurnedCustomers  int     `json:"churned_customers"`
	ChurnRate         float64 `json:"churn_rate"`
	LastStaffName     string  `json:"last_staff_name,omitempty"`
}

// PaymentBreakdown is used for payment method pie charts.
type PaymentBreakdown struct {
	Method string  `json:"method"`
	Amount float64 `json:"amount"`
	Ratio  float64 `json:"ratio"`
}

// Copyright 2026. Kimjibeom. All rights reserved.
package service

import (
	"context"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/repository"
)

// AnalyticsService handles business logic for analytics and reporting.
type AnalyticsService struct {
	saleRepo *repository.SaleRepository
}

func NewAnalyticsService(saleRepo *repository.SaleRepository) *AnalyticsService {
	return &AnalyticsService{saleRepo: saleRepo}
}

// GetDailySummary returns aggregated daily sales data.
func (s *AnalyticsService) GetDailySummary(ctx context.Context, date string) (*model.DailySummary, error) {
	return s.saleRepo.GetDailySummary(ctx, date)
}

// GetMonthlyTrend returns daily aggregated data for a given month.
func (s *AnalyticsService) GetMonthlyTrend(ctx context.Context, year, month int) ([]model.DailySummary, error) {
	return s.saleRepo.GetMonthlyTrend(ctx, year, month)
}

// GetAnalyticsSummary returns monthly analytics for a number of months.
func (s *AnalyticsService) GetAnalyticsSummary(ctx context.Context, months int) ([]model.AnalyticsSummary, error) {
	return s.saleRepo.GetAnalyticsSummary(ctx, months)
}

// GetPaymentBreakdown returns payment method distribution.
func (s *AnalyticsService) GetPaymentBreakdown(ctx context.Context, startDate, endDate string) ([]model.PaymentBreakdown, error) {
	return s.saleRepo.GetPaymentBreakdown(ctx, startDate, endDate)
}

// GetNewVsReturning returns customer visit type counts.
func (s *AnalyticsService) GetNewVsReturning(ctx context.Context, startDate, endDate string) (int, int, error) {
	return s.saleRepo.GetNewVsReturning(ctx, startDate, endDate)
}

// GetChurnAnalysis returns customer churn statistics.
func (s *AnalyticsService) GetChurnAnalysis(ctx context.Context, inactiveDays int) ([]model.ChurnAnalysis, error) {
	return s.saleRepo.GetChurnAnalysis(ctx, inactiveDays)
}

// GetStaffPerformance returns performance metrics for all staff.
func (s *AnalyticsService) GetStaffPerformance(ctx context.Context, startDate, endDate string) ([]model.StaffPerformance, error) {
	return s.saleRepo.GetStaffPerformance(ctx, startDate, endDate)
}

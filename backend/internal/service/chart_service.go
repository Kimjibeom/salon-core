// Copyright 2026. Kimjibeom. All rights reserved.
package service

import (
	"context"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/repository"
)

// ChartService handles business logic for treatment charts.
type ChartService struct {
	chartRepo *repository.ChartRepository
}

func NewChartService(chartRepo *repository.ChartRepository) *ChartService {
	return &ChartService{chartRepo: chartRepo}
}

func (s *ChartService) Create(ctx context.Context, req *model.ChartCreateRequest) (*model.Chart, error) {
	c := &model.Chart{
		CustomerID:        req.CustomerID,
		StaffID:           req.StaffID,
		Recipe:            req.Recipe,
		TreatmentName:     req.TreatmentName,
		TreatmentDuration: req.TreatmentDuration,
		Notes:             req.Notes,
		BeforeImgURL:      req.BeforeImgURL,
		AfterImgURL:       req.AfterImgURL,
		ConsentDocURL:     req.ConsentDocURL,
	}
	if err := s.chartRepo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *ChartService) GetByID(ctx context.Context, id string) (*model.Chart, error) {
	return s.chartRepo.GetByID(ctx, id)
}

func (s *ChartService) ListByCustomer(ctx context.Context, customerID string) ([]model.Chart, error) {
	return s.chartRepo.ListByCustomer(ctx, customerID)
}

func (s *ChartService) Delete(ctx context.Context, id string) error {
	return s.chartRepo.Delete(ctx, id)
}

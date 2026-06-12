// Copyright 2026. Kimjibeom. All rights reserved.
package service

import (
	"context"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/repository"
)

// MarketingService handles targeted customer extraction.
type MarketingService struct {
	saleRepo *repository.SaleRepository
}

func NewMarketingService(saleRepo *repository.SaleRepository) *MarketingService {
	return &MarketingService{saleRepo: saleRepo}
}

// ExtractTargetCustomers applies multi-condition filters to extract customer lists.
func (s *MarketingService) ExtractTargetCustomers(ctx context.Context, filters map[string]interface{}) ([]model.Customer, error) {
	return s.saleRepo.TargetMarketingQuery(ctx, filters)
}

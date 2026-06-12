// Copyright 2026. Kimjibeom. All rights reserved.
package service

import (
	"context"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/repository"
)

// StaffService handles business logic for staff management.
type StaffService struct {
	staffRepo *repository.StaffRepository
}

func NewStaffService(staffRepo *repository.StaffRepository) *StaffService {
	return &StaffService{staffRepo: staffRepo}
}

func (s *StaffService) GetByID(ctx context.Context, id string) (*model.Staff, error) {
	return s.staffRepo.GetByID(ctx, id)
}

func (s *StaffService) List(ctx context.Context) ([]model.Staff, error) {
	return s.staffRepo.List(ctx)
}

func (s *StaffService) Update(ctx context.Context, id string, req *model.StaffUpdateRequest) error {
	return s.staffRepo.Update(ctx, id, req)
}

func (s *StaffService) Delete(ctx context.Context, id string) error {
	return s.staffRepo.Delete(ctx, id)
}

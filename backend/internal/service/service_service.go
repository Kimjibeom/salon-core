// Copyright 2026. Kimjibeom. All rights reserved.
package service

import (
	"context"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/repository"
)

type ServiceMenuService struct {
	repo *repository.ServiceMenuRepository
}

func NewServiceMenuService(repo *repository.ServiceMenuRepository) *ServiceMenuService {
	return &ServiceMenuService{repo: repo}
}

func (s *ServiceMenuService) Create(ctx context.Context, req *model.ServiceMenuCreateRequest) (*model.ServiceMenu, error) {
	return s.repo.Create(ctx, req)
}

func (s *ServiceMenuService) GetAll(ctx context.Context) ([]model.ServiceMenu, error) {
	return s.repo.FindAll(ctx)
}

func (s *ServiceMenuService) Update(ctx context.Context, id string, req *model.ServiceMenuUpdateRequest) (*model.ServiceMenu, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *ServiceMenuService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

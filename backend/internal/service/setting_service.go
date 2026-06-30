package service

import (
	"context"
	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/repository"
)

type SettingService struct {
	repo *repository.SettingRepository
}

func NewSettingService(repo *repository.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

func (s *SettingService) GetAllSettings(ctx context.Context) ([]model.Setting, error) {
	return s.repo.GetAll(ctx)
}

func (s *SettingService) GetSettingByKey(ctx context.Context, key string) (*model.Setting, error) {
	return s.repo.GetByKey(ctx, key)
}

func (s *SettingService) UpdateSetting(ctx context.Context, setting *model.Setting) error {
	return s.repo.Update(ctx, setting)
}

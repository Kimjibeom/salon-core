// Copyright 2026. Kimjibeom. All rights reserved.
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kimjibeom/salon-core/backend/internal/middleware"
	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication logic.
type AuthService struct {
	staffRepo *repository.StaffRepository
}

// NewAuthService creates a new AuthService.
func NewAuthService(staffRepo *repository.StaffRepository) *AuthService {
	return &AuthService{staffRepo: staffRepo}
}

// Login authenticates a staff member and returns a JWT token.
func (s *AuthService) Login(ctx context.Context, req *model.StaffLoginRequest) (string, *model.Staff, error) {
	staff, err := s.staffRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Generic error message to prevent user enumeration
		return "", nil, errors.New("invalid email or password")
	}

	if !staff.IsActive {
		return "", nil, errors.New("account is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password)); err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	token, err := middleware.GenerateToken(staff.ID, staff.Email, staff.Role)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, staff, nil
}

// Register creates a new staff member with hashed password.
func (s *AuthService) Register(ctx context.Context, req *model.StaffCreateRequest) (*model.Staff, error) {
	// Hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	staff := &model.Staff{
		Name:                 req.Name,
		Email:                req.Email,
		PasswordHash:         string(hashedPassword),
		Role:                 req.Role,
		Phone:                req.Phone,
		ServiceIncentiveRate: req.ServiceIncentiveRate,
		ProductIncentiveRate: req.ProductIncentiveRate,
		BaseSalary:           req.BaseSalary,
		MonthlyTarget:        req.MonthlyTarget,
		IsActive:             true,
	}

	if err := s.staffRepo.Create(ctx, staff); err != nil {
		return nil, err
	}

	return staff, nil
}

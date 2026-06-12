// Copyright 2026. Kimjibeom. All rights reserved.
package service

import (
	"context"
	"errors"
	"time"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/repository"
)

// MembershipService handles business logic for memberships.
type MembershipService struct {
	membershipRepo *repository.MembershipRepository
}

func NewMembershipService(membershipRepo *repository.MembershipRepository) *MembershipService {
	return &MembershipService{membershipRepo: membershipRepo}
}

func (s *MembershipService) Create(ctx context.Context, req *model.MembershipCreateRequest) (*model.Membership, error) {
	m := &model.Membership{
		CustomerID:      req.CustomerID,
		Type:            req.Type,
		Name:            req.Name,
		InitialBalance:  req.InitialBalance,
		Balance:         req.InitialBalance,
		InitialCount:    req.InitialCount,
		RemainingCount:  req.InitialCount,
		TargetTreatment: req.TargetTreatment,
		IsActive:        true,
	}

	if req.ExpiredAt != "" {
		t, err := time.Parse(time.RFC3339, req.ExpiredAt)
		if err == nil {
			m.ExpiredAt = &t
		}
	}

	if err := s.membershipRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MembershipService) GetByID(ctx context.Context, id string) (*model.Membership, error) {
	return s.membershipRepo.GetByID(ctx, id)
}

func (s *MembershipService) ListByCustomer(ctx context.Context, customerID string) ([]model.Membership, error) {
	return s.membershipRepo.ListByCustomer(ctx, customerID)
}

func (s *MembershipService) DeductMoney(ctx context.Context, membershipID string, amount float64, saleID string) error {
	m, err := s.membershipRepo.GetByID(ctx, membershipID)
	if err != nil {
		return err
	}
	if !m.IsActive {
		return errors.New("membership is not active")
	}
	if m.Type != "money" {
		return errors.New("membership is not money type")
	}
	if m.Balance < amount {
		return errors.New("insufficient balance")
	}
	return s.membershipRepo.DeductMoney(ctx, membershipID, amount, saleID, "시술 결제 차감")
}

func (s *MembershipService) DeductCount(ctx context.Context, membershipID string, saleID string) error {
	m, err := s.membershipRepo.GetByID(ctx, membershipID)
	if err != nil {
		return err
	}
	if !m.IsActive {
		return errors.New("membership is not active")
	}
	if m.Type != "count" {
		return errors.New("membership is not count type")
	}
	if m.RemainingCount <= 0 {
		return errors.New("no remaining count")
	}
	return s.membershipRepo.DeductCount(ctx, membershipID, saleID, "회원권 횟수 차감")
}

func (s *MembershipService) GetTransactions(ctx context.Context, membershipID string) ([]model.MembershipTransaction, error) {
	return s.membershipRepo.GetTransactions(ctx, membershipID)
}

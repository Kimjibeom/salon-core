// Copyright 2026. Kimjibeom. All rights reserved.
package service

import (
	"context"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/repository"
)

// SaleService handles business logic for sales/POS.
type SaleService struct {
	saleRepo        *repository.SaleRepository
	membershipRepo  *repository.MembershipRepository
	customerRepo    *repository.CustomerRepository
	reservationRepo *repository.ReservationRepository
}

func NewSaleService(saleRepo *repository.SaleRepository, membershipRepo *repository.MembershipRepository, customerRepo *repository.CustomerRepository, reservationRepo *repository.ReservationRepository) *SaleService {
	return &SaleService{saleRepo: saleRepo, membershipRepo: membershipRepo, customerRepo: customerRepo, reservationRepo: reservationRepo}
}

func (s *SaleService) Create(ctx context.Context, req *model.SaleCreateRequest) (*model.Sale, error) {
	// For reservation-linked sales, resolve the customer from the reservation's
	// phone number: match an existing customer or auto-register a new one.
	if req.ReservationID != "" && req.CustomerID == "" {
		if res, err := s.reservationRepo.GetByID(ctx, req.ReservationID); err == nil {
			if res.CustomerID != nil {
				req.CustomerID = *res.CustomerID
			} else if res.CustomerPhone != "" {
				if cust, cerr := s.customerRepo.FindOrCreateByPhone(ctx, res.CustomerName, res.CustomerPhone); cerr == nil {
					req.CustomerID = cust.ID
					_ = s.reservationRepo.SetCustomer(ctx, req.ReservationID, cust.ID)
				}
			}
		}
	}

	sale := &model.Sale{
		StaffID:            req.StaffID,
		ItemName:           req.ItemName,
		TotalAmount:        req.TotalAmount,
		Category:           req.Category,
		PaymentMethod:      req.PaymentMethod,
		CardAmount:         req.CardAmount,
		CashAmount:         req.CashAmount,
		MembershipAmount:   req.MembershipAmount,
		CardCommissionRate: req.CardCommissionRate,
		Memo:               req.Memo,
	}

	if req.ReservationID != "" {
		sale.ReservationID = &req.ReservationID
	}
	if req.CustomerID != "" {
		sale.CustomerID = &req.CustomerID
	}
	if req.MembershipID != "" {
		sale.MembershipID = &req.MembershipID
	}

	if err := s.saleRepo.Create(ctx, sale); err != nil {
		return nil, err
	}

	// Update customer visit history for direct POS sales.
	// (Reservation-linked sales already update it when the reservation is completed.)
	if req.CustomerID != "" && req.ReservationID == "" {
		_ = s.customerRepo.UpdateLastVisit(ctx, req.CustomerID)
	}

	// Auto-deduct membership if used
	if req.MembershipID != "" && req.MembershipAmount > 0 {
		m, err := s.membershipRepo.GetByID(ctx, req.MembershipID)
		if err == nil {
			switch m.Type {
			case "money":
				_ = s.membershipRepo.DeductMoney(ctx, req.MembershipID, req.MembershipAmount, sale.ID, "영업 결제 차감")
			case "count":
				_ = s.membershipRepo.DeductCount(ctx, req.MembershipID, sale.ID, "회원권 횟수 차감")
			}
		}
	}

	return sale, nil
}

func (s *SaleService) GetByID(ctx context.Context, id string) (*model.Sale, error) {
	return s.saleRepo.GetByID(ctx, id)
}

func (s *SaleService) ListByDateRange(ctx context.Context, startDate, endDate string) ([]model.Sale, error) {
	return s.saleRepo.ListByDateRange(ctx, startDate, endDate)
}

func (s *SaleService) ListByStaff(ctx context.Context, staffID, startDate, endDate string) ([]model.Sale, error) {
	return s.saleRepo.ListByStaff(ctx, staffID, startDate, endDate)
}

func (s *SaleService) Delete(ctx context.Context, id string) error {
	return s.saleRepo.Delete(ctx, id)
}

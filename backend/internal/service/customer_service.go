// Copyright 2026. Kimjibeom. All rights reserved.
package service

import (
	"context"
	"time"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/repository"
)

// CustomerService handles business logic for customer management.
type CustomerService struct {
	customerRepo *repository.CustomerRepository
	chartRepo    *repository.ChartRepository
	memberRepo   *repository.MembershipRepository
}

func NewCustomerService(customerRepo *repository.CustomerRepository, chartRepo *repository.ChartRepository, memberRepo *repository.MembershipRepository) *CustomerService {
	return &CustomerService{customerRepo: customerRepo, chartRepo: chartRepo, memberRepo: memberRepo}
}

func (s *CustomerService) Create(ctx context.Context, req *model.CustomerCreateRequest) (*model.Customer, error) {
	c := &model.Customer{
		Name:  req.Name,
		Phone: req.Phone,
		Email: req.Email,
		Memo:  req.Memo,
		Tags:  req.Tags,
	}
	if req.BirthDate != "" {
		t, err := time.Parse("2006-01-02", req.BirthDate)
		if err == nil {
			c.BirthDate = &t
		}
	}
	if c.Tags == nil {
		c.Tags = []string{}
	}
	if err := s.customerRepo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CustomerService) GetByID(ctx context.Context, id string) (*model.Customer, error) {
	return s.customerRepo.GetByID(ctx, id)
}

func (s *CustomerService) GetWithHistory(ctx context.Context, id string) (*model.CustomerWithHistory, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	charts, _ := s.chartRepo.ListByCustomer(ctx, id)
	memberships, _ := s.memberRepo.ListByCustomer(ctx, id)

	return &model.CustomerWithHistory{
		Customer:    *customer,
		Charts:      charts,
		Memberships: memberships,
	}, nil
}

func (s *CustomerService) Search(ctx context.Context, query string, limit, offset int) ([]model.Customer, error) {
	return s.customerRepo.Search(ctx, query, limit, offset)
}

func (s *CustomerService) GetByPhone(ctx context.Context, phone string) (*model.Customer, error) {
	return s.customerRepo.GetByPhone(ctx, phone)
}

func (s *CustomerService) List(ctx context.Context, limit, offset int) ([]model.Customer, int, error) {
	return s.customerRepo.List(ctx, limit, offset)
}

func (s *CustomerService) Update(ctx context.Context, id string, req *model.CustomerUpdateRequest) error {
	return s.customerRepo.Update(ctx, id, req)
}

func (s *CustomerService) Delete(ctx context.Context, id string) error {
	return s.customerRepo.Delete(ctx, id)
}

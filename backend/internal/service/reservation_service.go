// Copyright 2026. Kimjibeom. All rights reserved.
package service

import (
	"context"
	"errors"
	"time"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/repository"
	"github.com/Kimjibeom/salon-core/backend/internal/websocket"
)

// ReservationService handles business logic for reservations.
type ReservationService struct {
	reservationRepo *repository.ReservationRepository
	customerRepo    *repository.CustomerRepository
	hub             *websocket.Hub
}

func NewReservationService(reservationRepo *repository.ReservationRepository, customerRepo *repository.CustomerRepository, hub *websocket.Hub) *ReservationService {
	return &ReservationService{reservationRepo: reservationRepo, customerRepo: customerRepo, hub: hub}
}

func (s *ReservationService) Create(ctx context.Context, req *model.ReservationCreateRequest) (*model.Reservation, error) {
	// Check for conflicts if staff is assigned
	if req.StaffID != "" {
		conflict, err := s.reservationRepo.CheckConflict(ctx, req.StaffID, req.Date, req.StartTime, req.EndTime, "")
		if err != nil {
			return nil, err
		}
		if conflict {
			return nil, errors.New("reservation conflicts with existing booking")
		}
	}

	status := req.Status
	if status == "" {
		status = "reserved"
	}

	res := &model.Reservation{
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		TreatmentName: req.TreatmentName,
		Date:          req.Date,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Status:        status,
		Source:        req.Source,
		Memo:          req.Memo,
	}

	if req.CustomerID != "" {
		res.CustomerID = &req.CustomerID
	}
	if req.StaffID != "" {
		res.StaffID = &req.StaffID
	}
	if req.ServiceID != "" {
		res.ServiceID = &req.ServiceID
	}

	if err := s.reservationRepo.Create(ctx, res); err != nil {
		return nil, err
	}

	// Broadcast new reservation event
	eventType := websocket.EventNewReservation
	if req.Source == "online" || req.Source == "naver" {
		eventType = websocket.EventOnlineBooking
	}
	s.hub.BroadcastEvent(eventType, res)

	return res, nil
}

func (s *ReservationService) GetByID(ctx context.Context, id string) (*model.Reservation, error) {
	return s.reservationRepo.GetByID(ctx, id)
}

func (s *ReservationService) ListByDate(ctx context.Context, date string) ([]model.Reservation, error) {
	return s.reservationRepo.ListByDate(ctx, date)
}

func (s *ReservationService) ListByDateRange(ctx context.Context, startDate, endDate string) ([]model.Reservation, error) {
	return s.reservationRepo.ListByDateRange(ctx, startDate, endDate)
}

func (s *ReservationService) GetWaitingQueue(ctx context.Context) ([]model.WaitingQueueEntry, error) {
	return s.reservationRepo.GetWaitingQueue(ctx)
}

func (s *ReservationService) AddToWaitingQueue(ctx context.Context, req *model.WaitingQueueCreateRequest) (*model.Reservation, error) {
	customerName := req.CustomerName
	if customerName == "" {
		customerName = "워크인 고객"
	}
	// Walk-ins have no scheduled slot; record the arrival time (start_time/end_time are NOT NULL).
	now := time.Now().Format("15:04")

	res := &model.Reservation{
		CustomerName:  customerName,
		CustomerPhone: req.CustomerPhone,
		TreatmentName: req.TreatmentName,
		StartTime:     now,
		EndTime:       now,
		Source:        "offline",
		Memo:          req.Memo,
	}

	if req.CustomerID != "" {
		res.CustomerID = &req.CustomerID
	}
	if req.StaffID != "" {
		res.StaffID = &req.StaffID
	}
	if req.ServiceID != "" {
		res.ServiceID = &req.ServiceID
	}

	if err := s.reservationRepo.AddToWaitingQueue(ctx, res); err != nil {
		return nil, err
	}

	s.hub.BroadcastEvent(websocket.EventWaitingQueueUpdate, res)
	return res, nil
}

func (s *ReservationService) UpdateStatus(ctx context.Context, id string, status string) error {
	if err := s.reservationRepo.UpdateStatus(ctx, id, status); err != nil {
		return err
	}

	// If completed, update customer's last visit
	if status == "completed" {
		res, err := s.reservationRepo.GetByID(ctx, id)
		if err == nil && res.CustomerID != nil {
			_ = s.customerRepo.UpdateLastVisit(ctx, *res.CustomerID)
		}
	}

	res, _ := s.reservationRepo.GetByID(ctx, id)
	s.hub.BroadcastEvent(websocket.EventReservationUpdate, res)
	return nil
}

func (s *ReservationService) Update(ctx context.Context, id string, req *model.ReservationUpdateRequest) error {
	if err := s.reservationRepo.Update(ctx, id, req); err != nil {
		return err
	}
	res, _ := s.reservationRepo.GetByID(ctx, id)
	s.hub.BroadcastEvent(websocket.EventReservationUpdate, res)
	return nil
}

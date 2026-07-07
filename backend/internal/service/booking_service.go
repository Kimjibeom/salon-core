// Copyright 2026. Kimjibeom. All rights reserved.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/Kimjibeom/salon-core/backend/internal/repository"
	"github.com/Kimjibeom/salon-core/backend/internal/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BookingService handles public booking logic and availability calculation.
type BookingService struct {
	pool            *pgxpool.Pool
	reservationRepo *repository.ReservationRepository
	staffRepo       *repository.StaffRepository
	serviceRepo     *repository.ServiceMenuRepository
	settingRepo     *repository.SettingRepository
	hub             *websocket.Hub
}

func NewBookingService(
	pool *pgxpool.Pool,
	reservationRepo *repository.ReservationRepository,
	staffRepo *repository.StaffRepository,
	serviceRepo *repository.ServiceMenuRepository,
	settingRepo *repository.SettingRepository,
	hub *websocket.Hub,
) *BookingService {
	return &BookingService{
		pool:            pool,
		reservationRepo: reservationRepo,
		staffRepo:       staffRepo,
		serviceRepo:     serviceRepo,
		settingRepo:     settingRepo,
		hub:             hub,
	}
}

// GetAvailability calculates available time slots for a given date, staff, and service.
func (s *BookingService) GetAvailability(ctx context.Context, date, staffID, serviceID string) (*model.AvailabilityResponse, error) {
	// 1. Parse the requested date
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	// 2. Get staff details and check day off
	staff, err := s.staffRepo.GetByID(ctx, staffID)
	if err != nil {
		return nil, errors.New("designer not found")
	}

	dayOfWeek := int(parsedDate.Weekday()) // 0=Sunday, 1=Monday, ...
	for _, dOff := range staff.DayOff {
		if dOff == dayOfWeek {
			return &model.AvailabilityResponse{
				Date:      date,
				StaffID:   staffID,
				StaffName: staff.Name,
				Slots:     []model.TimeSlot{},
			}, nil
		}
	}

	// 3. Get service duration
	svc, err := s.serviceRepo.GetByID(ctx, serviceID)
	if err != nil {
		return nil, errors.New("service not found")
	}
	durationMinutes := svc.Duration

	// 4. Get shop open/close times from settings
	openTime := "10:00"
	closeTime := "20:00"
	slotInterval := 30 // default 30 minutes

	if setting, err := s.settingRepo.GetByKey(ctx, "shop_open_time"); err == nil {
		openTime = setting.Value
	}
	if setting, err := s.settingRepo.GetByKey(ctx, "shop_close_time"); err == nil {
		closeTime = setting.Value
	}
	if setting, err := s.settingRepo.GetByKey(ctx, "booking_slot_interval"); err == nil {
		fmt.Sscanf(setting.Value, "%d", &slotInterval)
	}
	if slotInterval <= 0 {
		slotInterval = 30
	}

	// 5. Get existing reservations for this staff on this date
	existingReservations, err := s.getStaffReservationsForDate(ctx, staffID, date)
	if err != nil {
		return nil, err
	}

	// 6. Build time slots
	slots, err := s.buildTimeSlots(openTime, closeTime, slotInterval, durationMinutes, existingReservations, parsedDate)
	if err != nil {
		return nil, err
	}

	return &model.AvailabilityResponse{
		Date:      date,
		StaffID:   staffID,
		StaffName: staff.Name,
		Slots:     slots,
	}, nil
}

// getStaffReservationsForDate gets all active reservations for a specific staff on a specific date.
func (s *BookingService) getStaffReservationsForDate(ctx context.Context, staffID, date string) ([]model.Reservation, error) {
	allReservations, err := s.reservationRepo.ListByDate(ctx, date)
	if err != nil {
		return nil, err
	}

	var staffReservations []model.Reservation
	for _, r := range allReservations {
		if r.StaffID != nil && *r.StaffID == staffID &&
			r.Status != "canceled" && r.Status != "no_show" {
			staffReservations = append(staffReservations, r)
		}
	}
	return staffReservations, nil
}

// buildTimeSlots generates time slot windows from open to close, checking for overlaps.
func (s *BookingService) buildTimeSlots(openTime, closeTime string, slotInterval, duration int, reservations []model.Reservation, date time.Time) ([]model.TimeSlot, error) {
	open, err := parseHHMM(openTime)
	if err != nil {
		return nil, fmt.Errorf("invalid open time: %w", err)
	}
	close, err := parseHHMM(closeTime)
	if err != nil {
		return nil, fmt.Errorf("invalid close time: %w", err)
	}

	now := time.Now()
	isToday := date.Year() == now.Year() && date.Month() == now.Month() && date.Day() == now.Day()

	var slots []model.TimeSlot

	for slotStart := open; slotStart+duration <= close; slotStart += slotInterval {
		slotEnd := slotStart + duration
		if slotEnd > close {
			break
		}

		startStr := minutesToHHMM(slotStart)
		endStr := minutesToHHMM(slotEnd)

		available := true

		// If today, don't allow past slots
		if isToday {
			currentMinutes := now.Hour()*60 + now.Minute()
			if slotStart <= currentMinutes {
				available = false
			}
		}

		// Check for overlap with existing reservations
		if available {
			for _, r := range reservations {
				rStart, err1 := parseHHMM(r.StartTime)
				rEnd, err2 := parseHHMM(r.EndTime)
				if err1 != nil || err2 != nil {
					continue
				}
				// Overlap check: two intervals [A, B) and [C, D) overlap if A < D && C < B
				if slotStart < rEnd && rStart < slotEnd {
					available = false
					break
				}
			}
		}

		slots = append(slots, model.TimeSlot{
			StartTime: startStr,
			EndTime:   endStr,
			Available: available,
		})
	}

	return slots, nil
}

// CreatePublicBooking creates a reservation from the public booking site with concurrency control.
func (s *BookingService) CreatePublicBooking(ctx context.Context, req *model.PublicBookingRequest) (*model.Reservation, error) {
	// 1. Get service to calculate end time
	svc, err := s.serviceRepo.GetByID(ctx, req.ServiceID)
	if err != nil {
		return nil, errors.New("service not found")
	}

	// Calculate end time from start_time + duration
	startMinutes, err := parseHHMM(req.StartTime)
	if err != nil {
		return nil, errors.New("invalid start time format")
	}
	endMinutes := startMinutes + svc.Duration
	endTime := minutesToHHMM(endMinutes)

	// 2. Use a database transaction for concurrency control (prevent overbooking)
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check for conflicts within the transaction (using advisory lock on the staff+date combo)
	var exists bool
	conflictQuery := `SELECT EXISTS(
		SELECT 1 FROM reservations
		WHERE staff_id = $1 AND date = $2
		AND status NOT IN ('canceled', 'no_show')
		AND start_time < $4::time AND end_time > $3::time
	)`
	err = tx.QueryRow(ctx, conflictQuery, req.StaffID, req.Date, req.StartTime, endTime).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("conflict check failed: %w", err)
	}
	if exists {
		return nil, errors.New("해당 시간대에 이미 예약이 있습니다. 다른 시간을 선택해주세요.")
	}

	// 3. Insert the reservation
	res := &model.Reservation{
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		TreatmentName: svc.Name,
		Date:          req.Date,
		StartTime:     req.StartTime,
		EndTime:       endTime,
		Status:        "reserved",
		Source:        "online",
		Memo:          req.Memo,
	}
	staffID := req.StaffID
	res.StaffID = &staffID
	serviceID := req.ServiceID
	res.ServiceID = &serviceID

	insertQuery := `INSERT INTO reservations (customer_id, staff_id, service_id, customer_name, customer_phone, treatment_name, date, start_time, end_time, status, source, memo)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, created_at, updated_at`
	err = tx.QueryRow(ctx, insertQuery,
		nil, res.StaffID, res.ServiceID, res.CustomerName, res.CustomerPhone,
		res.TreatmentName, res.Date, res.StartTime, res.EndTime,
		res.Status, res.Source, res.Memo,
	).Scan(&res.ID, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	// 4. Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 5. Broadcast event
	s.hub.BroadcastEvent(websocket.EventOnlineBooking, res)

	return res, nil
}

// GetPublicStaffList returns a simplified staff list for the public booking page.
func (s *BookingService) GetPublicStaffList(ctx context.Context) ([]model.PublicStaffResponse, error) {
	staffs, err := s.staffRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	var result []model.PublicStaffResponse
	for _, st := range staffs {
		if st.IsActive && (st.Role == "designer" || st.Role == "admin") {
			result = append(result, model.PublicStaffResponse{
				ID:     st.ID,
				Name:   st.Name,
				DayOff: st.DayOff,
			})
		}
	}
	return result, nil
}

// GetPublicServiceList returns a simplified active service list for the public booking page.
func (s *BookingService) GetPublicServiceList(ctx context.Context) ([]model.PublicServiceResponse, error) {
	services, err := s.serviceRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []model.PublicServiceResponse
	for _, svc := range services {
		if svc.IsActive {
			result = append(result, model.PublicServiceResponse{
				ID:       svc.ID,
				Category: svc.Category,
				Name:     svc.Name,
				Price:    svc.Price,
				Duration: svc.Duration,
			})
		}
	}
	return result, nil
}

// GetShopInfo returns basic shop information for the booking page.
func (s *BookingService) GetShopInfo(ctx context.Context) (map[string]string, error) {
	info := map[string]string{
		"shop_name":       "",
		"shop_phone":      "",
		"shop_address":    "",
		"shop_open_time":  "10:00",
		"shop_close_time": "20:00",
	}
	keys := []string{"shop_name", "shop_phone", "shop_address", "shop_open_time", "shop_close_time"}
	for _, k := range keys {
		if setting, err := s.settingRepo.GetByKey(ctx, k); err == nil {
			info[k] = setting.Value
		}
	}
	return info, nil
}

// parseHHMM parses "HH:MM" into total minutes from midnight.
func parseHHMM(s string) (int, error) {
	var h, m int
	_, err := fmt.Sscanf(s, "%d:%d", &h, &m)
	if err != nil {
		return 0, err
	}
	return h*60 + m, nil
}

// minutesToHHMM converts total minutes from midnight to "HH:MM".
func minutesToHHMM(minutes int) string {
	return fmt.Sprintf("%02d:%02d", minutes/60, minutes%60)
}

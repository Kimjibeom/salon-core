// Copyright 2026. Kimjibeom. All rights reserved.
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ReservationRepository handles database operations for reservations.
type ReservationRepository struct {
	pool *pgxpool.Pool
}

// NewReservationRepository creates a new ReservationRepository.
func NewReservationRepository(pool *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{pool: pool}
}

// Create inserts a new reservation record.
func (r *ReservationRepository) Create(ctx context.Context, res *model.Reservation) error {
	query := `INSERT INTO reservations (customer_id, staff_id, customer_name, customer_phone, treatment_name, date, start_time, end_time, status, source, memo)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		res.CustomerID, res.StaffID, res.CustomerName, res.CustomerPhone,
		res.TreatmentName, res.Date, res.StartTime, res.EndTime,
		res.Status, res.Source, res.Memo,
	).Scan(&res.ID, &res.CreatedAt, &res.UpdatedAt)
}

// GetByID fetches a reservation by ID.
func (r *ReservationRepository) GetByID(ctx context.Context, id string) (*model.Reservation, error) {
	res := &model.Reservation{}
	query := `SELECT r.id, r.customer_id, r.staff_id, r.customer_name, r.customer_phone,
		COALESCE(s.name, ''), r.treatment_name, r.date, r.start_time, r.end_time,
		r.status, r.source, r.waiting_number, r.waiting_started_at, r.memo, r.created_at, r.updated_at
		FROM reservations r LEFT JOIN staffs s ON r.staff_id = s.id WHERE r.id = $1`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&res.ID, &res.CustomerID, &res.StaffID, &res.CustomerName, &res.CustomerPhone,
		&res.StaffName, &res.TreatmentName, &res.Date, &res.StartTime, &res.EndTime,
		&res.Status, &res.Source, &res.WaitingNumber, &res.WaitingStartedAt, &res.Memo,
		&res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ListByDate fetches all reservations for a given date.
func (r *ReservationRepository) ListByDate(ctx context.Context, date string) ([]model.Reservation, error) {
	query := `SELECT r.id, r.customer_id, r.staff_id, r.customer_name, r.customer_phone,
		COALESCE(s.name, ''), r.treatment_name, r.date, r.start_time, r.end_time,
		r.status, r.source, r.waiting_number, r.waiting_started_at, r.memo, r.created_at, r.updated_at
		FROM reservations r LEFT JOIN staffs s ON r.staff_id = s.id
		WHERE r.date = $1
		ORDER BY r.start_time`
	rows, err := r.pool.Query(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanReservations(rows)
}

// ListByDateRange fetches reservations within a date range.
func (r *ReservationRepository) ListByDateRange(ctx context.Context, startDate, endDate string) ([]model.Reservation, error) {
	query := `SELECT r.id, r.customer_id, r.staff_id, r.customer_name, r.customer_phone,
		COALESCE(s.name, ''), r.treatment_name, r.date, r.start_time, r.end_time,
		r.status, r.source, r.waiting_number, r.waiting_started_at, r.memo, r.created_at, r.updated_at
		FROM reservations r LEFT JOIN staffs s ON r.staff_id = s.id
		WHERE r.date BETWEEN $1 AND $2
		ORDER BY r.date, r.start_time`
	rows, err := r.pool.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanReservations(rows)
}

// GetWaitingQueue returns the current walk-in waiting queue for today.
func (r *ReservationRepository) GetWaitingQueue(ctx context.Context) ([]model.WaitingQueueEntry, error) {
	today := time.Now().Format("2006-01-02")
	query := `SELECT r.id, r.customer_id, r.staff_id, r.customer_name, r.customer_phone,
		COALESCE(s.name, ''), r.treatment_name, r.date, r.start_time, r.end_time,
		r.status, r.source, r.waiting_number, r.waiting_started_at, r.memo, r.created_at, r.updated_at
		FROM reservations r LEFT JOIN staffs s ON r.staff_id = s.id
		WHERE r.date = $1 AND r.status = 'waiting'
		ORDER BY r.waiting_number NULLS LAST, r.created_at`
	rows, err := r.pool.Query(ctx, query, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var queue []model.WaitingQueueEntry
	position := 1
	for rows.Next() {
		var res model.Reservation
		if err := rows.Scan(
			&res.ID, &res.CustomerID, &res.StaffID, &res.CustomerName, &res.CustomerPhone,
			&res.StaffName, &res.TreatmentName, &res.Date, &res.StartTime, &res.EndTime,
			&res.Status, &res.Source, &res.WaitingNumber, &res.WaitingStartedAt, &res.Memo,
			&res.CreatedAt, &res.UpdatedAt,
		); err != nil {
			return nil, err
		}

		waitTime := 0
		if res.WaitingStartedAt != nil {
			waitTime = int(time.Since(*res.WaitingStartedAt).Minutes())
		}

		queue = append(queue, model.WaitingQueueEntry{
			Reservation:     res,
			WaitTimeMinutes: waitTime,
			Position:        position,
		})
		position++
	}
	return queue, rows.Err()
}

// AddToWaitingQueue adds a walk-in customer to the waiting queue.
func (r *ReservationRepository) AddToWaitingQueue(ctx context.Context, res *model.Reservation) error {
	today := time.Now().Format("2006-01-02")
	// Get next waiting number
	var maxNum *int
	err := r.pool.QueryRow(ctx,
		`SELECT MAX(waiting_number) FROM reservations WHERE date = $1 AND status = 'waiting'`, today,
	).Scan(&maxNum)
	if err != nil {
		return err
	}

	nextNum := 1
	if maxNum != nil {
		nextNum = *maxNum + 1
	}

	res.Date = today
	res.Status = "waiting"
	now := time.Now()
	res.WaitingStartedAt = &now
	res.WaitingNumber = &nextNum

	query := `INSERT INTO reservations (customer_id, staff_id, customer_name, customer_phone, treatment_name, date, start_time, end_time, status, source, waiting_number, waiting_started_at, memo)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		res.CustomerID, res.StaffID, res.CustomerName, res.CustomerPhone,
		res.TreatmentName, res.Date, res.StartTime, res.EndTime,
		res.Status, res.Source, res.WaitingNumber, res.WaitingStartedAt, res.Memo,
	).Scan(&res.ID, &res.CreatedAt, &res.UpdatedAt)
}

// UpdateStatus updates the status of a reservation.
func (r *ReservationRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE reservations SET status = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, status, id)
	return err
}

// Update modifies an existing reservation.
func (r *ReservationRepository) Update(ctx context.Context, id string, req *model.ReservationUpdateRequest) error {
	query := `UPDATE reservations SET `
	args := []interface{}{}
	argIdx := 1
	setClauses := []string{}

	if req.StaffID != nil {
		setClauses = append(setClauses, fmt.Sprintf("staff_id = $%d", argIdx))
		args = append(args, *req.StaffID)
		argIdx++
	}
	if req.TreatmentName != nil {
		setClauses = append(setClauses, fmt.Sprintf("treatment_name = $%d", argIdx))
		args = append(args, *req.TreatmentName)
		argIdx++
	}
	if req.Date != nil {
		setClauses = append(setClauses, fmt.Sprintf("date = $%d", argIdx))
		args = append(args, *req.Date)
		argIdx++
	}
	if req.StartTime != nil {
		setClauses = append(setClauses, fmt.Sprintf("start_time = $%d", argIdx))
		args = append(args, *req.StartTime)
		argIdx++
	}
	if req.EndTime != nil {
		setClauses = append(setClauses, fmt.Sprintf("end_time = $%d", argIdx))
		args = append(args, *req.EndTime)
		argIdx++
	}
	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, *req.Status)
		argIdx++
	}
	if req.Memo != nil {
		setClauses = append(setClauses, fmt.Sprintf("memo = $%d", argIdx))
		args = append(args, *req.Memo)
		argIdx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	for i, clause := range setClauses {
		if i > 0 {
			query += ", "
		}
		query += clause
	}
	query += fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, id)

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// CheckConflict checks if a reservation conflicts with existing ones.
func (r *ReservationRepository) CheckConflict(ctx context.Context, staffID, date, startTime, endTime string, excludeID string) (bool, error) {
	query := `SELECT EXISTS(
		SELECT 1 FROM reservations
		WHERE staff_id = $1 AND date = $2
		AND status NOT IN ('canceled', 'no_show')
		AND start_time < $4 AND end_time > $3
		AND ($5 = '' OR id != $5::uuid)
	)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, staffID, date, startTime, endTime, excludeID).Scan(&exists)
	return exists, err
}

type pgxRows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
	Close()
}

func scanReservations(rows pgxRows) ([]model.Reservation, error) {
	var reservations []model.Reservation
	for rows.Next() {
		var res model.Reservation
		if err := rows.Scan(
			&res.ID, &res.CustomerID, &res.StaffID, &res.CustomerName, &res.CustomerPhone,
			&res.StaffName, &res.TreatmentName, &res.Date, &res.StartTime, &res.EndTime,
			&res.Status, &res.Source, &res.WaitingNumber, &res.WaitingStartedAt, &res.Memo,
			&res.CreatedAt, &res.UpdatedAt,
		); err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}
	return reservations, rows.Err()
}

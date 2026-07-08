// Copyright 2026. Kimjibeom. All rights reserved.
package repository

import (
	"context"
	"fmt"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StaffRepository handles database operations for staffs.
type StaffRepository struct {
	pool *pgxpool.Pool
}

// NewStaffRepository creates a new StaffRepository.
func NewStaffRepository(pool *pgxpool.Pool) *StaffRepository {
	return &StaffRepository{pool: pool}
}

// Create inserts a new staff record.
func (r *StaffRepository) Create(ctx context.Context, s *model.Staff) error {
	query := `INSERT INTO staffs (name, email, password_hash, role, phone, service_incentive_rate, product_incentive_rate, monthly_target)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		s.Name, s.Email, s.PasswordHash, s.Role, s.Phone,
		s.ServiceIncentiveRate, s.ProductIncentiveRate, s.MonthlyTarget,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

// GetByID fetches a staff member by ID.
func (r *StaffRepository) GetByID(ctx context.Context, id string) (*model.Staff, error) {
	s := &model.Staff{}
	query := `SELECT id, name, COALESCE(email, ''), password_hash, role, COALESCE(phone, ''), service_incentive_rate, product_incentive_rate, monthly_target, COALESCE(day_off, ARRAY[]::integer[]), is_active, created_at, updated_at
		FROM staffs WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.Name, &s.Email, &s.PasswordHash, &s.Role, &s.Phone,
		&s.ServiceIncentiveRate, &s.ProductIncentiveRate, &s.MonthlyTarget, &s.DayOff,
		&s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetByEmail fetches a staff member by email for authentication.
func (r *StaffRepository) GetByEmail(ctx context.Context, email string) (*model.Staff, error) {
	s := &model.Staff{}
	query := `SELECT id, name, COALESCE(email, ''), password_hash, role, COALESCE(phone, ''), service_incentive_rate, product_incentive_rate, monthly_target, COALESCE(day_off, ARRAY[]::integer[]), is_active, created_at, updated_at
		FROM staffs WHERE email = $1 AND is_active = true`
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&s.ID, &s.Name, &s.Email, &s.PasswordHash, &s.Role, &s.Phone,
		&s.ServiceIncentiveRate, &s.ProductIncentiveRate, &s.MonthlyTarget, &s.DayOff,
		&s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// List retrieves all active staff members.
func (r *StaffRepository) List(ctx context.Context) ([]model.Staff, error) {
	query := `SELECT id, name, COALESCE(email, ''), role, COALESCE(phone, ''), service_incentive_rate, product_incentive_rate, monthly_target, COALESCE(day_off, ARRAY[]::integer[]), is_active, created_at, updated_at
		FROM staffs WHERE is_active = true ORDER BY name`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var staffs []model.Staff
	for rows.Next() {
		var s model.Staff
		if err := rows.Scan(
			&s.ID, &s.Name, &s.Email, &s.Role, &s.Phone,
			&s.ServiceIncentiveRate, &s.ProductIncentiveRate, &s.MonthlyTarget, &s.DayOff,
			&s.IsActive, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		staffs = append(staffs, s)
	}
	return staffs, rows.Err()
}

// Update modifies an existing staff record.
func (r *StaffRepository) Update(ctx context.Context, id string, req *model.StaffUpdateRequest) error {
	query := `UPDATE staffs SET `
	args := []interface{}{}
	argIdx := 1
	setClauses := []string{}

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *req.Name)
		argIdx++
	}
	if req.Role != nil {
		setClauses = append(setClauses, fmt.Sprintf("role = $%d", argIdx))
		args = append(args, *req.Role)
		argIdx++
	}
	if req.Phone != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argIdx))
		args = append(args, *req.Phone)
		argIdx++
	}
	if req.ServiceIncentiveRate != nil {
		setClauses = append(setClauses, fmt.Sprintf("service_incentive_rate = $%d", argIdx))
		args = append(args, *req.ServiceIncentiveRate)
		argIdx++
	}
	if req.ProductIncentiveRate != nil {
		setClauses = append(setClauses, fmt.Sprintf("product_incentive_rate = $%d", argIdx))
		args = append(args, *req.ProductIncentiveRate)
		argIdx++
	}
	if req.MonthlyTarget != nil {
		setClauses = append(setClauses, fmt.Sprintf("monthly_target = $%d", argIdx))
		args = append(args, *req.MonthlyTarget)
		argIdx++
	}
	if req.DayOff != nil {
		setClauses = append(setClauses, fmt.Sprintf("day_off = $%d::int4[]", argIdx))
		args = append(args, *req.DayOff)
		argIdx++
	}
	if req.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *req.IsActive)
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

// Delete soft-deactivates a staff member.
func (r *StaffRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE staffs SET is_active = false WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

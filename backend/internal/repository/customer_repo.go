// Copyright 2026. Kimjibeom. All rights reserved.
package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CustomerRepository handles database operations for customers.
type CustomerRepository struct {
	pool *pgxpool.Pool
}

// NewCustomerRepository creates a new CustomerRepository.
func NewCustomerRepository(pool *pgxpool.Pool) *CustomerRepository {
	return &CustomerRepository{pool: pool}
}

// Create inserts a new customer record.
func (r *CustomerRepository) Create(ctx context.Context, c *model.Customer) error {
	query := `INSERT INTO customers (name, phone, email, birth_date, memo, tags)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query,
		c.Name, c.Phone, c.Email, c.BirthDate, c.Memo, c.Tags,
	).Scan(&c.ID, &c.CreatedAt)
}

// GetByID fetches a customer by ID.
func (r *CustomerRepository) GetByID(ctx context.Context, id string) (*model.Customer, error) {
	c := &model.Customer{}
	query := `SELECT id, name, phone, email, birth_date, memo, tags, created_at, last_visited_at, visit_count
		FROM customers WHERE id = $1 AND is_deleted = false`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Phone, &c.Email, &c.BirthDate, &c.Memo, &c.Tags,
		&c.CreatedAt, &c.LastVisitedAt, &c.VisitCount,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Search finds customers by name or phone suffix with indexed search.
func (r *CustomerRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.Customer, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	searchQuery := `SELECT id, name, phone, email, birth_date, memo, tags, created_at, last_visited_at, visit_count
		FROM customers WHERE is_deleted = false AND (name ILIKE $1 OR phone LIKE $2)
		ORDER BY last_visited_at DESC NULLS LAST, name
		LIMIT $3 OFFSET $4`

	namePattern := "%" + sanitizeLikePattern(query) + "%"
	phonePattern := "%" + sanitizeLikePattern(query)

	rows, err := r.pool.Query(ctx, searchQuery, namePattern, phonePattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []model.Customer
	for rows.Next() {
		var c model.Customer
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Phone, &c.Email, &c.BirthDate, &c.Memo, &c.Tags,
			&c.CreatedAt, &c.LastVisitedAt, &c.VisitCount,
		); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, rows.Err()
}

// GetByPhone fetches a customer by phone number (for CID lookup).
func (r *CustomerRepository) GetByPhone(ctx context.Context, phone string) (*model.Customer, error) {
	c := &model.Customer{}
	query := `SELECT id, name, phone, email, birth_date, memo, tags, created_at, last_visited_at, visit_count
		FROM customers WHERE phone = $1 AND is_deleted = false`
	err := r.pool.QueryRow(ctx, query, phone).Scan(
		&c.ID, &c.Name, &c.Phone, &c.Email, &c.BirthDate, &c.Memo, &c.Tags,
		&c.CreatedAt, &c.LastVisitedAt, &c.VisitCount,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// List retrieves all customers with pagination.
func (r *CustomerRepository) List(ctx context.Context, limit, offset int) ([]model.Customer, int, error) {
	if limit <= 0 {
		limit = 20
	}

	countQuery := `SELECT COUNT(*) FROM customers WHERE is_deleted = false`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, name, phone, email, birth_date, memo, tags, created_at, last_visited_at, visit_count
		FROM customers WHERE is_deleted = false
		ORDER BY last_visited_at DESC NULLS LAST, name
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var customers []model.Customer
	for rows.Next() {
		var c model.Customer
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Phone, &c.Email, &c.BirthDate, &c.Memo, &c.Tags,
			&c.CreatedAt, &c.LastVisitedAt, &c.VisitCount,
		); err != nil {
			return nil, 0, err
		}
		customers = append(customers, c)
	}
	return customers, total, rows.Err()
}

// Update modifies an existing customer record.
func (r *CustomerRepository) Update(ctx context.Context, id string, req *model.CustomerUpdateRequest) error {
	query := `UPDATE customers SET `
	args := []interface{}{}
	argIdx := 1
	setClauses := []string{}

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *req.Name)
		argIdx++
	}
	if req.Phone != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argIdx))
		args = append(args, *req.Phone)
		argIdx++
	}
	if req.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argIdx))
		args = append(args, *req.Email)
		argIdx++
	}
	if req.BirthDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("birth_date = $%d", argIdx))
		args = append(args, *req.BirthDate)
		argIdx++
	}
	if req.Memo != nil {
		setClauses = append(setClauses, fmt.Sprintf("memo = $%d", argIdx))
		args = append(args, *req.Memo)
		argIdx++
	}
	if req.Tags != nil {
		setClauses = append(setClauses, fmt.Sprintf("tags = $%d", argIdx))
		args = append(args, req.Tags)
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
	query += fmt.Sprintf(" WHERE id = $%d AND is_deleted = false", argIdx)
	args = append(args, id)

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// UpdateLastVisit updates the last_visited_at and increments visit_count.
func (r *CustomerRepository) UpdateLastVisit(ctx context.Context, id string) error {
	query := `UPDATE customers SET last_visited_at = NOW(), visit_count = visit_count + 1 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// Delete soft-deletes a customer.
func (r *CustomerRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE customers SET is_deleted = true WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// GetBirthdayCustomers returns customers with birthdays today.
func (r *CustomerRepository) GetBirthdayCustomers(ctx context.Context) ([]model.Customer, error) {
	query := `SELECT id, name, phone, email, birth_date, memo, tags, created_at, last_visited_at, visit_count
		FROM customers
		WHERE is_deleted = false AND birth_date IS NOT NULL
		AND EXTRACT(MONTH FROM birth_date) = EXTRACT(MONTH FROM CURRENT_DATE)
		AND EXTRACT(DAY FROM birth_date) = EXTRACT(DAY FROM CURRENT_DATE)`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []model.Customer
	for rows.Next() {
		var c model.Customer
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Phone, &c.Email, &c.BirthDate, &c.Memo, &c.Tags,
			&c.CreatedAt, &c.LastVisitedAt, &c.VisitCount,
		); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, rows.Err()
}

// sanitizeLikePattern escapes special LIKE pattern characters.
func sanitizeLikePattern(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}

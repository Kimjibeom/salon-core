// Copyright 2026. Kimjibeom. All rights reserved.
package repository

import (
	"context"
	"errors"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceMenuRepository struct {
	db *pgxpool.Pool
}

func NewServiceMenuRepository(db *pgxpool.Pool) *ServiceMenuRepository {
	return &ServiceMenuRepository{db: db}
}

func (r *ServiceMenuRepository) Create(ctx context.Context, req *model.ServiceMenuCreateRequest) (*model.ServiceMenu, error) {
	query := `INSERT INTO services (category, name, price, duration, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, category, name, price, duration, is_active, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, req.Category, req.Name, req.Price, req.Duration, req.IsActive)
	return scanServiceMenu(row)
}

func (r *ServiceMenuRepository) FindAll(ctx context.Context) ([]model.ServiceMenu, error) {
	query := `SELECT id, category, name, price, duration, is_active, created_at, updated_at
		FROM services
		ORDER BY category ASC, name ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []model.ServiceMenu
	for rows.Next() {
		svc, err := scanServiceMenu(rows)
		if err != nil {
			return nil, err
		}
		services = append(services, *svc)
	}

	return services, nil
}

// GetAll is an alias for FindAll, used by the booking service.
func (r *ServiceMenuRepository) GetAll(ctx context.Context) ([]model.ServiceMenu, error) {
	return r.FindAll(ctx)
}

// GetByID fetches a single service menu by ID.
func (r *ServiceMenuRepository) GetByID(ctx context.Context, id string) (*model.ServiceMenu, error) {
	query := `SELECT id, category, name, price, duration, is_active, created_at, updated_at
		FROM services WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	return scanServiceMenu(row)
}

func (r *ServiceMenuRepository) Update(ctx context.Context, id string, req *model.ServiceMenuUpdateRequest) (*model.ServiceMenu, error) {
	query := `UPDATE services SET
		category = COALESCE($1, category),
		name = COALESCE($2, name),
		price = COALESCE($3, price),
		duration = COALESCE($4, duration),
		is_active = COALESCE($5, is_active)
		WHERE id = $6
		RETURNING id, category, name, price, duration, is_active, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, req.Category, req.Name, req.Price, req.Duration, req.IsActive, id)
	return scanServiceMenu(row)
}

func (r *ServiceMenuRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM services WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("service menu not found")
	}
	return nil
}

func scanServiceMenu(row pgx.Row) (*model.ServiceMenu, error) {
	var s model.ServiceMenu
	err := row.Scan(&s.ID, &s.Category, &s.Name, &s.Price, &s.Duration, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("service menu not found")
		}
		return nil, err
	}
	return &s, nil
}

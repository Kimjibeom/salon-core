// Copyright 2026. Kimjibeom. All rights reserved.
package repository

import (
	"context"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ChartRepository handles database operations for charts.
type ChartRepository struct {
	pool *pgxpool.Pool
}

// NewChartRepository creates a new ChartRepository.
func NewChartRepository(pool *pgxpool.Pool) *ChartRepository {
	return &ChartRepository{pool: pool}
}

// Create inserts a new chart record.
func (r *ChartRepository) Create(ctx context.Context, c *model.Chart) error {
	query := `INSERT INTO charts (customer_id, staff_id, recipe, treatment_name, treatment_duration, notes, before_img_url, after_img_url, consent_doc_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query,
		c.CustomerID, c.StaffID, c.Recipe, c.TreatmentName, c.TreatmentDuration,
		c.Notes, c.BeforeImgURL, c.AfterImgURL, c.ConsentDocURL,
	).Scan(&c.ID, &c.CreatedAt)
}

// GetByID fetches a chart by ID.
func (r *ChartRepository) GetByID(ctx context.Context, id string) (*model.Chart, error) {
	c := &model.Chart{}
	query := `SELECT c.id, c.customer_id, c.staff_id, s.name, c.recipe, c.treatment_name, c.treatment_duration, c.notes, c.before_img_url, c.after_img_url, c.consent_doc_url, c.created_at
		FROM charts c LEFT JOIN staffs s ON c.staff_id = s.id WHERE c.id = $1`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.CustomerID, &c.StaffID, &c.StaffName, &c.Recipe, &c.TreatmentName,
		&c.TreatmentDuration, &c.Notes, &c.BeforeImgURL, &c.AfterImgURL, &c.ConsentDocURL, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ListByCustomer fetches all charts for a given customer.
func (r *ChartRepository) ListByCustomer(ctx context.Context, customerID string) ([]model.Chart, error) {
	query := `SELECT c.id, c.customer_id, c.staff_id, COALESCE(s.name, ''), c.recipe, c.treatment_name, c.treatment_duration, c.notes, c.before_img_url, c.after_img_url, c.consent_doc_url, c.created_at
		FROM charts c LEFT JOIN staffs s ON c.staff_id = s.id
		WHERE c.customer_id = $1
		ORDER BY c.created_at DESC`
	rows, err := r.pool.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var charts []model.Chart
	for rows.Next() {
		var c model.Chart
		if err := rows.Scan(
			&c.ID, &c.CustomerID, &c.StaffID, &c.StaffName, &c.Recipe, &c.TreatmentName,
			&c.TreatmentDuration, &c.Notes, &c.BeforeImgURL, &c.AfterImgURL, &c.ConsentDocURL, &c.CreatedAt,
		); err != nil {
			return nil, err
		}
		charts = append(charts, c)
	}
	return charts, rows.Err()
}

// Delete removes a chart record.
func (r *ChartRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM charts WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// Copyright 2026. Kimjibeom. All rights reserved.
package repository

import (
	"context"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NaverMappingRepository handles database operations for naver_mappings.
type NaverMappingRepository struct {
	pool *pgxpool.Pool
}

func NewNaverMappingRepository(pool *pgxpool.Pool) *NaverMappingRepository {
	return &NaverMappingRepository{pool: pool}
}

// Create inserts a new naver mapping.
func (r *NaverMappingRepository) Create(ctx context.Context, m *model.NaverMapping) error {
	query := `INSERT INTO naver_mappings (internal_type, internal_id, naver_id)
		VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query,
		m.InternalType, m.InternalID, m.NaverID,
	).Scan(&m.ID, &m.CreatedAt)
}

// GetAll returns all naver mappings.
func (r *NaverMappingRepository) GetAll(ctx context.Context) ([]model.NaverMapping, error) {
	query := `SELECT id, internal_type, internal_id, naver_id, created_at FROM naver_mappings ORDER BY internal_type, created_at`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []model.NaverMapping
	for rows.Next() {
		var m model.NaverMapping
		if err := rows.Scan(&m.ID, &m.InternalType, &m.InternalID, &m.NaverID, &m.CreatedAt); err != nil {
			return nil, err
		}
		mappings = append(mappings, m)
	}
	return mappings, rows.Err()
}

// GetByNaverID finds a mapping by naver ID and type.
func (r *NaverMappingRepository) GetByNaverID(ctx context.Context, internalType, naverID string) (*model.NaverMapping, error) {
	query := `SELECT id, internal_type, internal_id, naver_id, created_at FROM naver_mappings WHERE internal_type = $1 AND naver_id = $2`
	var m model.NaverMapping
	err := r.pool.QueryRow(ctx, query, internalType, naverID).Scan(&m.ID, &m.InternalType, &m.InternalID, &m.NaverID, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// Delete removes a naver mapping by ID.
func (r *NaverMappingRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM naver_mappings WHERE id = $1`, id)
	return err
}

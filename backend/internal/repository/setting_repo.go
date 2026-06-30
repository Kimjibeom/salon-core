package repository

import (
	"context"
	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SettingRepository struct {
	db *pgxpool.Pool
}

func NewSettingRepository(db *pgxpool.Pool) *SettingRepository {
	return &SettingRepository{db: db}
}

func (r *SettingRepository) GetAll(ctx context.Context) ([]model.Setting, error) {
	query := `SELECT key, value, description, updated_at FROM settings ORDER BY key ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []model.Setting
	for rows.Next() {
		var s model.Setting
		if err := rows.Scan(&s.Key, &s.Value, &s.Description, &s.UpdatedAt); err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, nil
}

func (r *SettingRepository) GetByKey(ctx context.Context, key string) (*model.Setting, error) {
	query := `SELECT key, value, description, updated_at FROM settings WHERE key = $1`
	var s model.Setting
	err := r.db.QueryRow(ctx, query, key).Scan(&s.Key, &s.Value, &s.Description, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SettingRepository) Update(ctx context.Context, s *model.Setting) error {
	query := `
		INSERT INTO settings (key, value, description, updated_at) 
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (key) DO UPDATE SET 
			value = EXCLUDED.value,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(ctx, query, s.Key, s.Value, s.Description)
	return err
}

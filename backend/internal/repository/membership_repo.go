// Copyright 2026. Kimjibeom. All rights reserved.
package repository

import (
	"context"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MembershipRepository handles database operations for memberships.
type MembershipRepository struct {
	pool *pgxpool.Pool
}

// NewMembershipRepository creates a new MembershipRepository.
func NewMembershipRepository(pool *pgxpool.Pool) *MembershipRepository {
	return &MembershipRepository{pool: pool}
}

// Create inserts a new membership record.
func (r *MembershipRepository) Create(ctx context.Context, m *model.Membership) error {
	query := `INSERT INTO memberships (customer_id, type, name, initial_balance, balance, initial_count, remaining_count, target_treatment, expired_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		m.CustomerID, m.Type, m.Name, m.InitialBalance, m.Balance,
		m.InitialCount, m.RemainingCount, m.TargetTreatment, m.ExpiredAt,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

// GetByID fetches a membership by ID.
func (r *MembershipRepository) GetByID(ctx context.Context, id string) (*model.Membership, error) {
	m := &model.Membership{}
	query := `SELECT id, customer_id, type, name, initial_balance, balance, initial_count, remaining_count, target_treatment, expired_at, is_active, created_at, updated_at
		FROM memberships WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.CustomerID, &m.Type, &m.Name, &m.InitialBalance, &m.Balance,
		&m.InitialCount, &m.RemainingCount, &m.TargetTreatment, &m.ExpiredAt,
		&m.IsActive, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// ListByCustomer fetches all active memberships for a customer.
func (r *MembershipRepository) ListByCustomer(ctx context.Context, customerID string) ([]model.Membership, error) {
	query := `SELECT id, customer_id, type, name, initial_balance, balance, initial_count, remaining_count, target_treatment, expired_at, is_active, created_at, updated_at
		FROM memberships WHERE customer_id = $1 AND is_active = true
		ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []model.Membership
	for rows.Next() {
		var m model.Membership
		if err := rows.Scan(
			&m.ID, &m.CustomerID, &m.Type, &m.Name, &m.InitialBalance, &m.Balance,
			&m.InitialCount, &m.RemainingCount, &m.TargetTreatment, &m.ExpiredAt,
			&m.IsActive, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		memberships = append(memberships, m)
	}
	return memberships, rows.Err()
}

// DeductMoney deducts balance from a money-type membership and records the transaction.
func (r *MembershipRepository) DeductMoney(ctx context.Context, membershipID string, amount float64, saleID string, description string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Get current balance
	var balanceBefore float64
	err = tx.QueryRow(ctx,
		`SELECT balance FROM memberships WHERE id = $1 AND is_active = true FOR UPDATE`, membershipID,
	).Scan(&balanceBefore)
	if err != nil {
		return err
	}

	balanceAfter := balanceBefore - amount
	if balanceAfter < 0 {
		balanceAfter = 0
	}

	// Update balance
	_, err = tx.Exec(ctx,
		`UPDATE memberships SET balance = $1 WHERE id = $2`, balanceAfter, membershipID)
	if err != nil {
		return err
	}

	// Record transaction
	_, err = tx.Exec(ctx,
		`INSERT INTO membership_transactions (membership_id, sale_id, amount, balance_before, balance_after, description)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		membershipID, saleID, amount, balanceBefore, balanceAfter, description)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// DeductCount deducts one count from a count-type membership and records the transaction.
func (r *MembershipRepository) DeductCount(ctx context.Context, membershipID string, saleID string, description string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var countBefore int
	err = tx.QueryRow(ctx,
		`SELECT remaining_count FROM memberships WHERE id = $1 AND is_active = true FOR UPDATE`, membershipID,
	).Scan(&countBefore)
	if err != nil {
		return err
	}

	countAfter := countBefore - 1
	if countAfter < 0 {
		countAfter = 0
	}

	_, err = tx.Exec(ctx,
		`UPDATE memberships SET remaining_count = $1 WHERE id = $2`, countAfter, membershipID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO membership_transactions (membership_id, sale_id, count_change, count_before, count_after, description)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		membershipID, saleID, -1, countBefore, countAfter, description)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetExpiringMemberships returns memberships expiring within the given number of days.
func (r *MembershipRepository) GetExpiringMemberships(ctx context.Context, withinDays int) ([]model.Membership, error) {
	query := `SELECT m.id, m.customer_id, m.type, m.name, m.initial_balance, m.balance, m.initial_count, m.remaining_count, m.target_treatment, m.expired_at, m.is_active, m.created_at, m.updated_at
		FROM memberships m
		WHERE m.is_active = true AND m.expired_at IS NOT NULL
		AND m.expired_at BETWEEN NOW() AND NOW() + MAKE_INTERVAL(days => $1)
		ORDER BY m.expired_at`
	rows, err := r.pool.Query(ctx, query, withinDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []model.Membership
	for rows.Next() {
		var m model.Membership
		if err := rows.Scan(
			&m.ID, &m.CustomerID, &m.Type, &m.Name, &m.InitialBalance, &m.Balance,
			&m.InitialCount, &m.RemainingCount, &m.TargetTreatment, &m.ExpiredAt,
			&m.IsActive, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		memberships = append(memberships, m)
	}
	return memberships, rows.Err()
}

// GetTransactions returns membership transaction history.
func (r *MembershipRepository) GetTransactions(ctx context.Context, membershipID string) ([]model.MembershipTransaction, error) {
	query := `SELECT id, membership_id, sale_id, amount, count_change, balance_before, balance_after, count_before, count_after, description, created_at
		FROM membership_transactions WHERE membership_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, membershipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []model.MembershipTransaction
	for rows.Next() {
		var t model.MembershipTransaction
		var saleID *string
		if err := rows.Scan(
			&t.ID, &t.MembershipID, &saleID, &t.Amount, &t.CountChange,
			&t.BalanceBefore, &t.BalanceAfter, &t.CountBefore, &t.CountAfter,
			&t.Description, &t.CreatedAt,
		); err != nil {
			return nil, err
		}
		if saleID != nil {
			t.SaleID = *saleID
		}
		txns = append(txns, t)
	}
	return txns, rows.Err()
}

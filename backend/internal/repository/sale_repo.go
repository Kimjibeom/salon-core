// Copyright 2026. Kimjibeom. All rights reserved.
package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Kimjibeom/salon-core/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SaleRepository handles database operations for sales.
type SaleRepository struct {
	pool *pgxpool.Pool
}

// NewSaleRepository creates a new SaleRepository.
func NewSaleRepository(pool *pgxpool.Pool) *SaleRepository {
	return &SaleRepository{pool: pool}
}

// Create inserts a new sale record.
func (r *SaleRepository) Create(ctx context.Context, s *model.Sale) error {
	query := `INSERT INTO sales (reservation_id, customer_id, staff_id, service_id, item_name, total_amount, category, payment_method, card_amount, cash_amount, membership_amount, card_commission_rate, membership_id, memo)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query,
		s.ReservationID, s.CustomerID, s.StaffID, s.ServiceID, s.ItemName, s.TotalAmount,
		s.Category, s.PaymentMethod, s.CardAmount, s.CashAmount,
		s.MembershipAmount, s.CardCommissionRate, s.MembershipID, s.Memo,
	).Scan(&s.ID, &s.CreatedAt)
}

// GetByID fetches a sale by ID.
func (r *SaleRepository) GetByID(ctx context.Context, id string) (*model.Sale, error) {
	s := &model.Sale{}
	query := `SELECT s.id, s.reservation_id, s.customer_id, s.staff_id, COALESCE(st.name, ''), COALESCE(c.name, ''),
		s.service_id, s.item_name, s.total_amount, s.category, s.payment_method,
		s.card_amount, s.cash_amount, s.membership_amount, s.card_commission_rate, s.membership_id, s.memo, s.created_at
		FROM sales s
		LEFT JOIN staffs st ON s.staff_id = st.id
		LEFT JOIN customers c ON s.customer_id = c.id
		WHERE s.id = $1`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.ReservationID, &s.CustomerID, &s.StaffID, &s.StaffName, &s.CustomerName,
		&s.ServiceID, &s.ItemName, &s.TotalAmount, &s.Category, &s.PaymentMethod,
		&s.CardAmount, &s.CashAmount, &s.MembershipAmount, &s.CardCommissionRate,
		&s.MembershipID, &s.Memo, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// ListByDateRange retrieves sales within a date range.
func (r *SaleRepository) ListByDateRange(ctx context.Context, startDate, endDate string) ([]model.Sale, error) {
	query := `SELECT s.id, s.reservation_id, s.customer_id, s.staff_id, COALESCE(st.name, ''), COALESCE(c.name, ''),
		s.service_id, s.item_name, s.total_amount, s.category, s.payment_method,
		s.card_amount, s.cash_amount, s.membership_amount, s.card_commission_rate, s.membership_id, s.memo, s.created_at
		FROM sales s
		LEFT JOIN staffs st ON s.staff_id = st.id
		LEFT JOIN customers c ON s.customer_id = c.id
		WHERE s.created_at::date BETWEEN $1 AND $2
		ORDER BY s.created_at DESC`
	rows, err := r.pool.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSales(rows)
}

// ListByStaff retrieves sales for a specific staff within a date range.
func (r *SaleRepository) ListByStaff(ctx context.Context, staffID, startDate, endDate string) ([]model.Sale, error) {
	query := `SELECT s.id, s.reservation_id, s.customer_id, s.staff_id, COALESCE(st.name, ''), COALESCE(c.name, ''),
		s.service_id, s.item_name, s.total_amount, s.category, s.payment_method,
		s.card_amount, s.cash_amount, s.membership_amount, s.card_commission_rate, s.membership_id, s.memo, s.created_at
		FROM sales s
		LEFT JOIN staffs st ON s.staff_id = st.id
		LEFT JOIN customers c ON s.customer_id = c.id
		WHERE s.staff_id = $1 AND s.created_at::date BETWEEN $2 AND $3
		ORDER BY s.created_at DESC`
	rows, err := r.pool.Query(ctx, query, staffID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSales(rows)
}

// GetDailySummary returns aggregated sales data for a given date.
func (r *SaleRepository) GetDailySummary(ctx context.Context, date string) (*model.DailySummary, error) {
	ds := &model.DailySummary{Date: date}
	query := `SELECT
		COALESCE(SUM(total_amount), 0),
		COALESCE(SUM(CASE WHEN category = 'service' THEN total_amount ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN category = 'product' THEN total_amount ELSE 0 END), 0),
		COALESCE(SUM(card_amount), 0),
		COALESCE(SUM(cash_amount), 0),
		COALESCE(SUM(membership_amount), 0),
		COUNT(*),
		COUNT(DISTINCT customer_id)
		FROM sales WHERE created_at::date = $1`
	err := r.pool.QueryRow(ctx, query, date).Scan(
		&ds.TotalRevenue, &ds.ServiceRevenue, &ds.ProductRevenue,
		&ds.CardRevenue, &ds.CashRevenue, &ds.MembershipRevenue,
		&ds.TransactionCount, &ds.CustomerCount,
	)
	return ds, err
}

// GetMonthlyTrend returns daily aggregates for a given month.
func (r *SaleRepository) GetMonthlyTrend(ctx context.Context, year, month int) ([]model.DailySummary, error) {
	query := `SELECT
		created_at::date as sale_date,
		COALESCE(SUM(total_amount), 0),
		COALESCE(SUM(CASE WHEN category = 'service' THEN total_amount ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN category = 'product' THEN total_amount ELSE 0 END), 0),
		COALESCE(SUM(card_amount), 0),
		COALESCE(SUM(cash_amount), 0),
		COALESCE(SUM(membership_amount), 0),
		COUNT(*),
		COUNT(DISTINCT customer_id)
		FROM sales
		WHERE EXTRACT(YEAR FROM created_at) = $1 AND EXTRACT(MONTH FROM created_at) = $2
		GROUP BY sale_date ORDER BY sale_date`
	rows, err := r.pool.Query(ctx, query, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []model.DailySummary
	for rows.Next() {
		var ds model.DailySummary
		if err := rows.Scan(
			&ds.Date, &ds.TotalRevenue, &ds.ServiceRevenue, &ds.ProductRevenue,
			&ds.CardRevenue, &ds.CashRevenue, &ds.MembershipRevenue,
			&ds.TransactionCount, &ds.CustomerCount,
		); err != nil {
			return nil, err
		}
		summaries = append(summaries, ds)
	}
	return summaries, rows.Err()
}

// GetAnalyticsSummary returns monthly aggregated analytics for a given period.
func (r *SaleRepository) GetAnalyticsSummary(ctx context.Context, months int) ([]model.AnalyticsSummary, error) {
	query := `SELECT
		TO_CHAR(created_at, 'YYYY-MM') as period,
		COALESCE(SUM(total_amount), 0),
		COALESCE(SUM(CASE WHEN category = 'service' THEN total_amount ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN category = 'product' THEN total_amount ELSE 0 END), 0),
		COALESCE(AVG(CASE WHEN category = 'service' THEN total_amount END), 0),
		COALESCE(AVG(CASE WHEN category = 'product' THEN total_amount END), 0),
		COALESCE(SUM(card_amount), 0),
		COALESCE(SUM(cash_amount), 0),
		COALESCE(SUM(membership_amount), 0)
		FROM sales
		WHERE created_at >= NOW() - MAKE_INTERVAL(months => $1)
		GROUP BY period ORDER BY period`
	rows, err := r.pool.Query(ctx, query, months)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []model.AnalyticsSummary
	for rows.Next() {
		var as model.AnalyticsSummary
		if err := rows.Scan(
			&as.Period, &as.TotalRevenue, &as.ServiceRevenue, &as.ProductRevenue,
			&as.AvgServicePrice, &as.AvgProductPrice,
			&as.CardRevenue, &as.CashRevenue, &as.MembershipRevenue,
		); err != nil {
			return nil, err
		}
		summaries = append(summaries, as)
	}
	return summaries, rows.Err()
}

// GetPaymentBreakdown returns payment method distribution for a date range.
func (r *SaleRepository) GetPaymentBreakdown(ctx context.Context, startDate, endDate string) ([]model.PaymentBreakdown, error) {
	query := `WITH totals AS (
		SELECT
			COALESCE(SUM(card_amount), 0) as card_total,
			COALESCE(SUM(cash_amount), 0) as cash_total,
			COALESCE(SUM(membership_amount), 0) as membership_total,
			COALESCE(SUM(total_amount), 0) as grand_total
		FROM sales WHERE created_at::date BETWEEN $1 AND $2
	)
	SELECT * FROM (
		SELECT 'card' as method, card_total as amount,
			CASE WHEN grand_total > 0 THEN ROUND((card_total / grand_total * 100)::numeric, 1) ELSE 0 END as ratio FROM totals
		UNION ALL
		SELECT 'cash', cash_total,
			CASE WHEN grand_total > 0 THEN ROUND((cash_total / grand_total * 100)::numeric, 1) ELSE 0 END FROM totals
		UNION ALL
		SELECT 'membership', membership_total,
			CASE WHEN grand_total > 0 THEN ROUND((membership_total / grand_total * 100)::numeric, 1) ELSE 0 END FROM totals
	) sub ORDER BY amount DESC`
	rows, err := r.pool.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var breakdown []model.PaymentBreakdown
	for rows.Next() {
		var pb model.PaymentBreakdown
		if err := rows.Scan(&pb.Method, &pb.Amount, &pb.Ratio); err != nil {
			return nil, err
		}
		breakdown = append(breakdown, pb)
	}
	return breakdown, rows.Err()
}

// GetNewVsReturning returns new vs returning customer counts for a date range.
func (r *SaleRepository) GetNewVsReturning(ctx context.Context, startDate, endDate string) (int, int, error) {
	query := `SELECT
		COUNT(DISTINCT CASE WHEN c.visit_count <= 1 THEN c.id END) as new_customers,
		COUNT(DISTINCT CASE WHEN c.visit_count > 1 THEN c.id END) as returning_customers
		FROM sales s JOIN customers c ON s.customer_id = c.id
		WHERE s.created_at::date BETWEEN $1 AND $2`
	var newCount, returningCount int
	err := r.pool.QueryRow(ctx, query, startDate, endDate).Scan(&newCount, &returningCount)
	return newCount, returningCount, err
}

// GetChurnAnalysis returns churn statistics by visit number.
func (r *SaleRepository) GetChurnAnalysis(ctx context.Context, inactiveDays int) ([]model.ChurnAnalysis, error) {
	query := `WITH customer_visits AS (
		SELECT c.id, c.visit_count, c.last_visited_at,
			CASE WHEN c.last_visited_at < NOW() - MAKE_INTERVAL(days => $1) THEN true ELSE false END as is_churned
		FROM customers c WHERE c.is_deleted = false AND c.visit_count > 0
	)
	SELECT visit_count, COUNT(*) as total, SUM(CASE WHEN is_churned THEN 1 ELSE 0 END) as churned,
		ROUND(SUM(CASE WHEN is_churned THEN 1.0 ELSE 0 END) / NULLIF(COUNT(*), 0) * 100, 1) as churn_rate
	FROM customer_visits
	GROUP BY visit_count
	ORDER BY visit_count
	LIMIT 20`
	rows, err := r.pool.Query(ctx, query, inactiveDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []model.ChurnAnalysis
	for rows.Next() {
		var ca model.ChurnAnalysis
		if err := rows.Scan(&ca.VisitNumber, &ca.TotalCustomers, &ca.ChurnedCustomers, &ca.ChurnRate); err != nil {
			return nil, err
		}
		analyses = append(analyses, ca)
	}
	return analyses, rows.Err()
}

// GetStaffPerformance returns performance metrics for each staff member within a period.
func (r *SaleRepository) GetStaffPerformance(ctx context.Context, startDate, endDate string) ([]model.StaffPerformance, error) {
	query := `SELECT
		st.id, st.name, st.monthly_target,
		COALESCE(SUM(CASE WHEN s.category = 'service' THEN s.total_amount ELSE 0 END), 0) as service_rev,
		COALESCE(SUM(CASE WHEN s.category = 'product' THEN s.total_amount ELSE 0 END), 0) as product_rev,
		COALESCE(SUM(s.total_amount), 0) as total_rev,
		st.service_incentive_rate, st.product_incentive_rate,
		COALESCE(SUM(CASE WHEN s.category = 'service' THEN s.card_amount * s.card_commission_rate / 100 ELSE 0 END), 0) as service_commission,
		COALESCE(SUM(CASE WHEN s.category = 'product' THEN s.card_amount * s.card_commission_rate / 100 ELSE 0 END), 0) as product_commission,
		st.base_salary
		FROM staffs st
		LEFT JOIN sales s ON st.id = s.staff_id AND s.created_at::date BETWEEN $1 AND $2
		WHERE st.is_active = true
		GROUP BY st.id, st.name, st.monthly_target, st.service_incentive_rate, st.product_incentive_rate, st.base_salary
		ORDER BY total_rev DESC`
	rows, err := r.pool.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performances []model.StaffPerformance
	for rows.Next() {
		var sp model.StaffPerformance
		var serviceRate, productRate, serviceCommission, productCommission float64
		if err := rows.Scan(
			&sp.StaffID, &sp.StaffName, &sp.MonthlyTarget,
			&sp.ServiceRevenue, &sp.ProductRevenue, &sp.TotalRevenue,
			&serviceRate, &productRate, &serviceCommission, &productCommission,
			&sp.BaseSalary,
		); err != nil {
			return nil, err
		}

		// Calculate net revenue (after card commission deduction)
		sp.NetServiceRevenue = sp.ServiceRevenue - serviceCommission
		sp.NetProductRevenue = sp.ProductRevenue - productCommission

		// Calculate incentives based on net revenue
		sp.ServiceIncentive = sp.NetServiceRevenue * serviceRate / 100
		sp.ProductIncentive = sp.NetProductRevenue * productRate / 100
		sp.TotalIncentive = sp.ServiceIncentive + sp.ProductIncentive

		// Achievement rate
		if sp.MonthlyTarget > 0 {
			sp.AchievementRate = sp.TotalRevenue / sp.MonthlyTarget * 100
		}

		// Total payroll
		sp.TotalPayroll = sp.BaseSalary + sp.TotalIncentive

		performances = append(performances, sp)
	}
	return performances, rows.Err()
}

// TargetMarketingQuery builds a dynamic query for targeted customer extraction.
func (r *SaleRepository) TargetMarketingQuery(ctx context.Context, filters map[string]interface{}) ([]model.Customer, error) {
	baseQuery := `SELECT DISTINCT c.id, c.name, c.phone, c.email, c.birth_date, c.memo, c.tags, c.created_at, c.last_visited_at, c.visit_count
		FROM customers c
		LEFT JOIN sales s ON c.id = s.customer_id
		LEFT JOIN memberships m ON c.id = m.customer_id
		WHERE c.is_deleted = false`

	args := []interface{}{}
	argIdx := 1
	conditions := []string{}

	if treatment, ok := filters["treatment_name"].(string); ok && treatment != "" {
		conditions = append(conditions, fmt.Sprintf("s.item_name ILIKE $%d", argIdx))
		args = append(args, "%"+sanitizeTargetLike(treatment)+"%")
		argIdx++
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		conditions = append(conditions, fmt.Sprintf("s.category = $%d", argIdx))
		args = append(args, category)
		argIdx++
	}
	if hasMembership, ok := filters["has_membership"].(bool); ok && hasMembership {
		conditions = append(conditions, "m.is_active = true")
	}
	if membershipType, ok := filters["membership_type"].(string); ok && membershipType != "" {
		conditions = append(conditions, fmt.Sprintf("m.type = $%d AND m.is_active = true", argIdx))
		args = append(args, membershipType)
		argIdx++
	}
	if minVisits, ok := filters["min_visits"].(int); ok {
		conditions = append(conditions, fmt.Sprintf("c.visit_count >= $%d", argIdx))
		args = append(args, minVisits)
		argIdx++
	}
	if maxVisits, ok := filters["max_visits"].(int); ok {
		conditions = append(conditions, fmt.Sprintf("c.visit_count <= $%d", argIdx))
		args = append(args, maxVisits)
		argIdx++
	}
	if inactiveDays, ok := filters["inactive_days"].(int); ok {
		conditions = append(conditions, fmt.Sprintf("c.last_visited_at < NOW() - MAKE_INTERVAL(days => $%d)", argIdx))
		args = append(args, inactiveDays)
		argIdx++
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}
	baseQuery += " ORDER BY c.last_visited_at DESC NULLS LAST LIMIT 500"

	rows, err := r.pool.Query(ctx, baseQuery, args...)
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

// Delete removes a sale record. Only accessible by admin.
func (r *SaleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sales WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func sanitizeTargetLike(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}

type salesRows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
	Close()
}

func scanSales(rows salesRows) ([]model.Sale, error) {
	var sales []model.Sale
	for rows.Next() {
		var s model.Sale
		if err := rows.Scan(
			&s.ID, &s.ReservationID, &s.CustomerID, &s.StaffID, &s.StaffName, &s.CustomerName,
			&s.ServiceID, &s.ItemName, &s.TotalAmount, &s.Category, &s.PaymentMethod,
			&s.CardAmount, &s.CashAmount, &s.MembershipAmount, &s.CardCommissionRate,
			&s.MembershipID, &s.Memo, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		sales = append(sales, s)
	}
	return sales, rows.Err()
}

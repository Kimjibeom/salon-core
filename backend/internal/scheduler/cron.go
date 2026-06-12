// Copyright 2026. Kimjibeom. All rights reserved.
package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/Kimjibeom/salon-core/backend/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

// Scheduler handles recurring background tasks.
type Scheduler struct {
	cron           *cron.Cron
	pool           *pgxpool.Pool
	customerRepo   *repository.CustomerRepository
	membershipRepo *repository.MembershipRepository
}

// NewScheduler creates a new Scheduler.
func NewScheduler(pool *pgxpool.Pool, customerRepo *repository.CustomerRepository, membershipRepo *repository.MembershipRepository) *Scheduler {
	return &Scheduler{
		cron:           cron.New(cron.WithSeconds()),
		pool:           pool,
		customerRepo:   customerRepo,
		membershipRepo: membershipRepo,
	}
}

// Start begins the cron scheduler.
func (s *Scheduler) Start() {
	// Every day at midnight: process birthday customers, treatment cycle, membership expiry
	_, err := s.cron.AddFunc("0 0 0 * * *", s.dailyNotificationJob)
	if err != nil {
		log.Printf("Failed to add daily notification job: %v", err)
	}

	// Every hour: check for expiring memberships (more granular)
	_, err = s.cron.AddFunc("0 0 * * * *", s.membershipExpiryCheck)
	if err != nil {
		log.Printf("Failed to add membership expiry check: %v", err)
	}

	s.cron.Start()
	log.Println("Cron scheduler started")
}

// Stop gracefully stops the scheduler.
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("Cron scheduler stopped")
}

// dailyNotificationJob runs at midnight and queues notifications.
func (s *Scheduler) dailyNotificationJob() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	log.Println("Running daily notification job...")

	// 1. Birthday customers
	s.processBirthdayNotifications(ctx)

	// 2. Treatment cycle customers (visited 30+ days ago)
	s.processTreatmentCycleNotifications(ctx)

	// 3. Membership expiry (within 7 days)
	s.processMembershipExpiryNotifications(ctx)

	log.Println("Daily notification job completed")
}

func (s *Scheduler) processBirthdayNotifications(ctx context.Context) {
	customers, err := s.customerRepo.GetBirthdayCustomers(ctx)
	if err != nil {
		log.Printf("Failed to get birthday customers: %v", err)
		return
	}

	for _, c := range customers {
		_, err := s.pool.Exec(ctx,
			`INSERT INTO notification_queue (customer_id, type, title, message, channel, scheduled_at)
			VALUES ($1, 'birthday', $2, $3, 'email', NOW())`,
			c.ID,
			"🎂 생일 축하합니다!",
			c.Name+"님, 생일을 진심으로 축하드립니다! 특별한 할인 혜택을 준비했습니다.",
		)
		if err != nil {
			log.Printf("Failed to queue birthday notification for customer %s: %v", c.ID, err)
		}
	}
	log.Printf("Queued %d birthday notifications", len(customers))
}

func (s *Scheduler) processTreatmentCycleNotifications(ctx context.Context) {
	// Find customers who haven't visited in 30-90 days (treatment cycle reminder)
	query := `SELECT id, name, phone FROM customers
		WHERE is_deleted = false AND last_visited_at IS NOT NULL
		AND last_visited_at BETWEEN NOW() - INTERVAL '90 days' AND NOW() - INTERVAL '30 days'
		AND id NOT IN (
			SELECT customer_id FROM notification_queue
			WHERE type = 'treatment_cycle' AND created_at > NOW() - INTERVAL '7 days'
		)`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		log.Printf("Failed to get treatment cycle customers: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, name, phone string
		if err := rows.Scan(&id, &name, &phone); err != nil {
			continue
		}
		_, err := s.pool.Exec(ctx,
			`INSERT INTO notification_queue (customer_id, type, title, message, channel, scheduled_at)
			VALUES ($1, 'treatment_cycle', $2, $3, 'email', NOW())`,
			id,
			"✨ 시술 주기 알림",
			name+"님, 마지막 방문 이후 시간이 지났습니다. 새로운 스타일을 위해 방문해주세요!",
		)
		if err != nil {
			log.Printf("Failed to queue treatment cycle notification: %v", err)
		}
		count++
	}
	log.Printf("Queued %d treatment cycle notifications", count)
}

func (s *Scheduler) processMembershipExpiryNotifications(ctx context.Context) {
	memberships, err := s.membershipRepo.GetExpiringMemberships(ctx, 7)
	if err != nil {
		log.Printf("Failed to get expiring memberships: %v", err)
		return
	}

	for _, m := range memberships {
		// Check if we already sent a notification recently
		var exists bool
		err := s.pool.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM notification_queue WHERE customer_id = $1 AND type = 'membership_expiry' AND created_at > NOW() - INTERVAL '3 days')`,
			m.CustomerID,
		).Scan(&exists)
		if err != nil || exists {
			continue
		}

		expiryDate := ""
		if m.ExpiredAt != nil {
			expiryDate = m.ExpiredAt.Format("2006-01-02")
		}

		_, err = s.pool.Exec(ctx,
			`INSERT INTO notification_queue (customer_id, type, title, message, channel, scheduled_at)
			VALUES ($1, 'membership_expiry', $2, $3, 'email', NOW())`,
			m.CustomerID,
			"⚠️ 멤버십 만료 임박",
			m.Name+" 멤버십이 "+expiryDate+"에 만료됩니다. 갱신을 위해 방문해주세요.",
		)
		if err != nil {
			log.Printf("Failed to queue membership expiry notification: %v", err)
		}
	}
	log.Printf("Queued membership expiry notifications for %d memberships", len(memberships))
}

func (s *Scheduler) membershipExpiryCheck() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Deactivate expired memberships
	_, err := s.pool.Exec(ctx,
		`UPDATE memberships SET is_active = false WHERE expired_at < NOW() AND is_active = true`)
	if err != nil {
		log.Printf("Failed to deactivate expired memberships: %v", err)
	}
}

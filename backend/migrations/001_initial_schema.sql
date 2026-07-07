-- ============================================
-- Salon Core CRM - Initial Database Schema
-- Copyright 2026. Kimjibeom. All rights reserved.
-- ============================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- 1. staffs (직원)
-- ============================================
CREATE TABLE IF NOT EXISTS staffs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'staff' CHECK (role IN ('admin', 'designer', 'staff')),
    phone VARCHAR(20),
    service_incentive_rate NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    product_incentive_rate NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    base_salary NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    monthly_target NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_staffs_role ON staffs(role);
CREATE INDEX IF NOT EXISTS idx_staffs_is_active ON staffs(is_active);

-- ============================================
-- 2. customers (고객)
-- ============================================
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255),
    birth_date DATE,
    memo TEXT DEFAULT '',
    tags TEXT[] DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_visited_at TIMESTAMPTZ,
    visit_count INTEGER NOT NULL DEFAULT 0,
    is_deleted BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX IF NOT EXISTS idx_customers_phone ON customers(phone);
CREATE INDEX IF NOT EXISTS idx_customers_name ON customers(name);
CREATE INDEX IF NOT EXISTS idx_customers_birth_date ON customers(birth_date);
CREATE INDEX IF NOT EXISTS idx_customers_last_visited ON customers(last_visited_at);
CREATE INDEX IF NOT EXISTS idx_customers_tags ON customers USING GIN(tags);

-- ============================================
-- 3. charts (시술 차트 / 히스토리)
-- ============================================
CREATE TABLE IF NOT EXISTS charts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    staff_id UUID NOT NULL REFERENCES staffs(id) ON DELETE SET NULL,
    recipe TEXT NOT NULL DEFAULT '',
    treatment_name VARCHAR(200),
    treatment_duration INTEGER, -- minutes
    notes TEXT DEFAULT '',
    before_img_url TEXT,
    after_img_url TEXT,
    consent_doc_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_charts_customer ON charts(customer_id);
CREATE INDEX IF NOT EXISTS idx_charts_staff ON charts(staff_id);
CREATE INDEX IF NOT EXISTS idx_charts_created ON charts(created_at);

-- ============================================
-- 4. memberships (정액권 / 회원권)
-- ============================================
CREATE TABLE IF NOT EXISTS memberships (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    type VARCHAR(10) NOT NULL CHECK (type IN ('money', 'count')),
    name VARCHAR(200) NOT NULL DEFAULT '',
    initial_balance NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    balance NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    initial_count INTEGER NOT NULL DEFAULT 0,
    remaining_count INTEGER NOT NULL DEFAULT 0,
    target_treatment VARCHAR(200),
    expired_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_memberships_customer ON memberships(customer_id);
CREATE INDEX IF NOT EXISTS idx_memberships_type ON memberships(type);
CREATE INDEX IF NOT EXISTS idx_memberships_expired ON memberships(expired_at);
CREATE INDEX IF NOT EXISTS idx_memberships_active ON memberships(is_active);

-- ============================================
-- 5. reservations (예약)
-- ============================================
CREATE TABLE IF NOT EXISTS reservations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
    staff_id UUID REFERENCES staffs(id) ON DELETE SET NULL,
    customer_name VARCHAR(100),
    customer_phone VARCHAR(20),
    treatment_name VARCHAR(200),
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'reserved' CHECK (status IN ('reserved', 'waiting', 'in_progress', 'completed', 'canceled', 'no_show')),
    source VARCHAR(20) NOT NULL DEFAULT 'offline' CHECK (source IN ('online', 'offline', 'naver')),
    waiting_number INTEGER,
    waiting_started_at TIMESTAMPTZ,
    memo TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reservations_date ON reservations(date);
CREATE INDEX IF NOT EXISTS idx_reservations_customer ON reservations(customer_id);
CREATE INDEX IF NOT EXISTS idx_reservations_staff ON reservations(staff_id);
CREATE INDEX IF NOT EXISTS idx_reservations_status ON reservations(status);
CREATE INDEX IF NOT EXISTS idx_reservations_source ON reservations(source);
CREATE INDEX IF NOT EXISTS idx_reservations_date_staff ON reservations(date, staff_id);

-- ============================================
-- 6. sales (매출)
-- ============================================
CREATE TABLE IF NOT EXISTS sales (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reservation_id UUID REFERENCES reservations(id) ON DELETE SET NULL,
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
    staff_id UUID NOT NULL REFERENCES staffs(id) ON DELETE SET NULL,
    item_name VARCHAR(200) NOT NULL,
    total_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    category VARCHAR(20) NOT NULL CHECK (category IN ('service', 'product')),
    payment_method VARCHAR(20) NOT NULL CHECK (payment_method IN ('card', 'cash', 'membership', 'mixed')),
    card_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    cash_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    membership_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    card_commission_rate NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    membership_id UUID REFERENCES memberships(id) ON DELETE SET NULL,
    memo TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sales_staff ON sales(staff_id);
CREATE INDEX IF NOT EXISTS idx_sales_customer ON sales(customer_id);
CREATE INDEX IF NOT EXISTS idx_sales_reservation ON sales(reservation_id);
CREATE INDEX IF NOT EXISTS idx_sales_category ON sales(category);
CREATE INDEX IF NOT EXISTS idx_sales_created ON sales(created_at);
CREATE INDEX IF NOT EXISTS idx_sales_payment ON sales(payment_method);

-- ============================================
-- 7. notification_queue (알림 큐)
-- ============================================
CREATE TABLE IF NOT EXISTS notification_queue (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('birthday', 'treatment_cycle', 'membership_expiry', 'visit_thanks', 'membership_balance', 'marketing')),
    title VARCHAR(200) NOT NULL,
    message TEXT NOT NULL,
    channel VARCHAR(20) NOT NULL DEFAULT 'email' CHECK (channel IN ('email', 'sms', 'push')),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed', 'canceled')),
    scheduled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_status ON notification_queue(status);
CREATE INDEX IF NOT EXISTS idx_notifications_scheduled ON notification_queue(scheduled_at);
CREATE INDEX IF NOT EXISTS idx_notifications_customer ON notification_queue(customer_id);

-- ============================================
-- 8. membership_transactions (멤버십 차감 이력)
-- ============================================
CREATE TABLE IF NOT EXISTS membership_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    membership_id UUID NOT NULL REFERENCES memberships(id) ON DELETE CASCADE,
    sale_id UUID REFERENCES sales(id) ON DELETE SET NULL,
    amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    count_change INTEGER NOT NULL DEFAULT 0,
    balance_before NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    balance_after NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    count_before INTEGER NOT NULL DEFAULT 0,
    count_after INTEGER NOT NULL DEFAULT 0,
    description TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_membership_tx_membership ON membership_transactions(membership_id);
CREATE INDEX IF NOT EXISTS idx_membership_tx_sale ON membership_transactions(sale_id);

-- ============================================
-- Updated_at trigger function
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_staffs_updated_at') THEN
        CREATE TRIGGER update_staffs_updated_at BEFORE UPDATE ON staffs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_memberships_updated_at') THEN
        CREATE TRIGGER update_memberships_updated_at BEFORE UPDATE ON memberships FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_reservations_updated_at') THEN
        CREATE TRIGGER update_reservations_updated_at BEFORE UPDATE ON reservations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

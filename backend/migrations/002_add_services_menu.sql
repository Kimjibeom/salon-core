-- ============================================
-- Salon Core CRM - 002 Add Services Menu
-- ============================================

-- 1. Create services table
CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    price NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    duration INTEGER NOT NULL DEFAULT 30, -- minutes
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_services_category ON services(category);
CREATE INDEX idx_services_active ON services(is_active);

CREATE TRIGGER update_services_updated_at BEFORE UPDATE ON services FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 2. Add service_id to existing tables
ALTER TABLE reservations ADD COLUMN service_id UUID REFERENCES services(id) ON DELETE SET NULL;
ALTER TABLE sales ADD COLUMN service_id UUID REFERENCES services(id) ON DELETE SET NULL;
ALTER TABLE charts ADD COLUMN service_id UUID REFERENCES services(id) ON DELETE SET NULL;

CREATE INDEX idx_reservations_service ON reservations(service_id);
CREATE INDEX idx_sales_service ON sales(service_id);
CREATE INDEX idx_charts_service ON charts(service_id);

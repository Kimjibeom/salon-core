-- ============================================
-- 005: Booking Enhancements
-- Add naver_mappings table, staff day_off, shop settings
-- ============================================

-- 1. Staff day-off support (array of integers, 0=Sunday .. 6=Saturday)
ALTER TABLE staffs ADD COLUMN IF NOT EXISTS day_off INTEGER[] DEFAULT '{}';

-- 2. Naver Mappings: link internal IDs to Naver-side IDs
CREATE TABLE IF NOT EXISTS naver_mappings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    internal_type VARCHAR(50) NOT NULL CHECK (internal_type IN ('STAFF', 'SERVICE')),
    internal_id UUID NOT NULL,
    naver_id VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(internal_type, internal_id),
    UNIQUE(internal_type, naver_id)
);

CREATE INDEX IF NOT EXISTS idx_naver_mappings_type ON naver_mappings(internal_type);
CREATE INDEX IF NOT EXISTS idx_naver_mappings_naver_id ON naver_mappings(naver_id);

-- 3. Insert default shop settings (open/close times) if not exists
INSERT INTO settings (key, value, description) VALUES
    ('shop_name', '우리 미용실', '매장명'),
    ('shop_open_time', '10:00', '영업 시작 시간'),
    ('shop_close_time', '20:00', '영업 종료 시간'),
    ('shop_phone', '', '매장 전화번호'),
    ('shop_address', '', '매장 주소'),
    ('booking_slot_interval', '30', '예약 슬롯 간격 (분)')
ON CONFLICT (key) DO NOTHING;

-- 4. Add naver_external_id column to reservations (to track naver booking id)
ALTER TABLE reservations ADD COLUMN IF NOT EXISTS naver_external_id VARCHAR(100);
CREATE INDEX IF NOT EXISTS idx_reservations_naver_ext ON reservations(naver_external_id) WHERE naver_external_id IS NOT NULL;

-- 5. Add 'service_id' column to reservations if missing (already exists from 001 but ensure FK)
-- service_id column already exists from initial schema, just need to add the service table FK if missing
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'reservations_service_id_fkey'
        AND table_name = 'reservations'
    ) THEN
        BEGIN
            ALTER TABLE reservations ADD CONSTRAINT reservations_service_id_fkey
                FOREIGN KEY (service_id) REFERENCES service_menu(id) ON DELETE SET NULL;
        EXCEPTION WHEN undefined_table THEN
            -- service_menu table might not exist yet, skip
            NULL;
        END;
    END IF;
END $$;

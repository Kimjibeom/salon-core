-- Copyright 2026. Kimjibeom. All rights reserved.
-- Enforce phone-number-based customer identity.
-- 1. Deduplicate existing active customers sharing the same phone:
--    keep the record with the most visits (oldest as tiebreaker), soft-delete the rest.
UPDATE customers SET is_deleted = true
WHERE is_deleted = false
AND id NOT IN (
    SELECT DISTINCT ON (phone) id FROM customers
    WHERE is_deleted = false
    ORDER BY phone, visit_count DESC, created_at ASC
);

-- 2. One active customer per phone number.
CREATE UNIQUE INDEX IF NOT EXISTS idx_customers_phone_active_unique
    ON customers(phone) WHERE is_deleted = false;

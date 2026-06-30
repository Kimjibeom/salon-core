-- Create settings table for generic configurations
CREATE TABLE IF NOT EXISTS settings (
    key VARCHAR(255) PRIMARY KEY,
    value TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert initial values for naver sync
INSERT INTO settings (key, value, description) VALUES ('naver_place_id', '', '네이버 스마트플레이스 ID') ON CONFLICT DO NOTHING;
INSERT INTO settings (key, value, description) VALUES ('naver_webhook_secret', '', '네이버 예약 웹훅 시크릿 키') ON CONFLICT DO NOTHING;

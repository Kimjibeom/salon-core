// Copyright 2026. Kimjibeom. All rights reserved.
// Salon Core CRM - Configuration

package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Config holds all configuration values for the application.
type Config struct {
	// Database
	DatabaseURL string

	// Supabase
	SupabaseURL            string
	SupabaseAnonKey        string
	SupabaseServiceRoleKey string
	SupabaseStorageBucket  string

	// JWT
	JWTSecretKey string

	// Server
	BackendPort string
	FrontendURL string
	BackendURL  string

	// External integrations
	NaverBookingAPIKey    string
	NaverBookingAPISecret string
	SMSAPIKey             string
	SMSSenderNumber       string

	// SMTP
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string

	// Rate Limiting
	RateLimitRPS   int
	RateLimitBurst int
}

// Load reads configuration from environment variables.
// It validates required fields and applies secure defaults.
func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:            os.Getenv("DATABASE_URL"),
		SupabaseURL:            os.Getenv("SUPABASE_URL"),
		SupabaseAnonKey:        os.Getenv("SUPABASE_ANON_KEY"),
		SupabaseServiceRoleKey: os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
		SupabaseStorageBucket:  getEnvOrDefault("SUPABASE_STORAGE_BUCKET", "salon-media"),
		JWTSecretKey:           "",
		BackendPort:            getEnvOrDefault("BACKEND_PORT", "8080"),
		FrontendURL:            getEnvOrDefault("FRONTEND_URL", "http://localhost:3000"),
		BackendURL:             getEnvOrDefault("BACKEND_URL", "http://localhost:8080"),
		NaverBookingAPIKey:     os.Getenv("NAVER_BOOKING_API_KEY"),
		NaverBookingAPISecret:  os.Getenv("NAVER_BOOKING_API_SECRET"),
		SMSAPIKey:              os.Getenv("SMS_API_KEY"),
		SMSSenderNumber:        os.Getenv("SMS_SENDER_NUMBER"),
		SMTPHost:               os.Getenv("SMTP_HOST"),
		SMTPUser:               os.Getenv("SMTP_USER"),
		SMTPPassword:           os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:               getEnvOrDefault("SMTP_FROM", "noreply@salon-core.app"),
	}

	// SMTP port
	smtpPort, _ := strconv.Atoi(getEnvOrDefault("SMTP_PORT", "587"))
	cfg.SMTPPort = smtpPort

	// Rate limiting
	cfg.RateLimitRPS, _ = strconv.Atoi(getEnvOrDefault("RATE_LIMIT_RPS", "10"))
	cfg.RateLimitBurst, _ = strconv.Atoi(getEnvOrDefault("RATE_LIMIT_BURST", "20"))

	// JWT Secret: multi-tiered resolution (env -> file -> ephemeral random + warning)
	jwtSecret, err := resolveJWTSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve JWT secret: %w", err)
	}
	cfg.JWTSecretKey = jwtSecret

	return cfg, nil
}

// resolveJWTSecret resolves the JWT secret using a secure multi-tiered approach.
// Priority: Environment variable -> File -> Ephemeral random generation
func resolveJWTSecret() (string, error) {
	// Tier 1: Environment variable
	if secret := os.Getenv("JWT_SECRET_KEY"); secret != "" {
		return secret, nil
	}

	// Tier 2: Secret file
	if data, err := os.ReadFile("jwt_secret.txt"); err == nil {
		secret := string(data)
		if len(secret) > 0 {
			return secret, nil
		}
	}

	// Tier 3: Generate ephemeral random secret + severe warning
	log.Println("WARNING: Generating ephemeral JWT secret. This instance is isolated and sessions will not persist across restarts. Set JWT_SECRET_KEY environment variable for production use.")
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random secret: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

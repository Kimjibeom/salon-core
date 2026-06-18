// Copyright 2026. Kimjibeom. All rights reserved.
package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS configuration.
type CORSConfig struct {
	AllowedOrigins []string
}

// CORS returns a strict CORS middleware. Wildcard origins are NOT allowed.
func CORS(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Strict origin validation - no wildcard allowed
		allowed := false
		for _, allowedOrigin := range config.AllowedOrigins {
			if strings.EqualFold(origin, allowedOrigin) {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-CSRF-Token")
			c.Header("Access-Control-Max-Age", "86400")
			c.Header("Access-Control-Expose-Headers", "X-CSRF-Token")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimiter implements a simple token bucket rate limiter.
type RateLimiter struct {
	rps   int
	burst int
	// In production, use a distributed rate limiter (e.g., Redis-based)
	// TODO(security): Implement distributed rate limiting for horizontal scaling
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(rps, burst int) *RateLimiter {
	return &RateLimiter{rps: rps, burst: burst}
}

// Simple in-memory rate limiter using a fixed window counter
type rateLimitInfo struct {
	count int
	reset time.Time
}

var (
	ipLastRequest   = make(map[string]*rateLimitInfo)
	ipLastRequestMu sync.Mutex
)

// RateLimit returns a rate limiting middleware.
func RateLimit(rps int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		ipLastRequestMu.Lock()
		info, ok := ipLastRequest[ip]
		if !ok || now.After(info.reset) {
			ipLastRequest[ip] = &rateLimitInfo{
				count: 1,
				reset: now.Add(time.Second),
			}
		} else {
			if info.count >= rps {
				ipLastRequestMu.Unlock()
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
				return
			}
			info.count++
		}
		ipLastRequestMu.Unlock()
		c.Next()
	}
}

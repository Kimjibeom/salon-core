// Copyright 2026. Kimjibeom. All rights reserved.
package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security headers to all responses.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Strict CSP policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https://fonts.gstatic.com; connect-src 'self' wss: https:; frame-ancestors 'none'; base-uri 'self'; form-action 'self'")

		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// XSS protection (legacy browsers)
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Disable unnecessary browser features
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")

		// HTTPS enforcement
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Prevent caching of sensitive data
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
		c.Header("Pragma", "no-cache")

		c.Next()
	}
}

// Copyright 2026. Kimjibeom. All rights reserved.
package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims defines the JWT token claims.
type Claims struct {
	StaffID string `json:"staff_id"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

// InitJWT initializes the JWT secret key.
func InitJWT(secret string) {
	jwtSecret = []byte(secret)
}

// GenerateToken creates a new JWT token for a staff member.
func GenerateToken(staffID, email, role string) (string, error) {
	claims := Claims{
		StaffID: staffID,
		Email:   email,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "salon-core",
		},
	}

	// Hardcode HS256 algorithm - never derive from unverified token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// AuthMiddleware validates JWT tokens and extracts claims.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		tokenString := parts[1]

		// Parse with hardcoded HS256 algorithm - reject 'none' and other algorithms
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Reject 'none' algorithm and validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		}, jwt.WithValidMethods([]string{"HS256"}))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Validate exp claim
		if claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			return
		}

		// Set claims in context
		c.Set("staff_id", claims.StaffID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RequireRole returns a middleware that restricts access to specific roles.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		for _, allowed := range allowedRoles {
			if roleStr == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
	}
}

// GetJWTSecret resolves the JWT secret using a secure multi-tiered approach.
func GetJWTSecret() string {
	// Tier 1: Environment variable
	if secret := os.Getenv("JWT_SECRET_KEY"); secret != "" {
		return secret
	}

	// Tier 2: Secret file
	if data, err := os.ReadFile("jwt_secret.txt"); err == nil {
		secret := strings.TrimSpace(string(data))
		if len(secret) > 0 {
			return secret
		}
	}

	// Tier 3: Ephemeral random + warning
	log.Println("WARNING: Generating ephemeral JWT secret. Instance-isolated!")
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate random JWT secret: %v", err)
	}
	return hex.EncodeToString(bytes)
}

// Copyright 2026. Kimjibeom. All rights reserved.
// Salon Core CRM - Backend Entry Point

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kimjibeom/salon-core/backend/internal/config"
	"github.com/Kimjibeom/salon-core/backend/internal/database"
	"github.com/Kimjibeom/salon-core/backend/internal/handler"
	"github.com/Kimjibeom/salon-core/backend/internal/middleware"
	"github.com/Kimjibeom/salon-core/backend/internal/repository"
	"github.com/Kimjibeom/salon-core/backend/internal/scheduler"
	"github.com/Kimjibeom/salon-core/backend/internal/service"
	"github.com/Kimjibeom/salon-core/backend/internal/websocket"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize JWT
	middleware.InitJWT(cfg.JWTSecretKey)

	// Initialize database pool
	ctx := context.Background()
	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Initialize repositories
	staffRepo := repository.NewStaffRepository(pool)
	customerRepo := repository.NewCustomerRepository(pool)
	chartRepo := repository.NewChartRepository(pool)
	membershipRepo := repository.NewMembershipRepository(pool)
	reservationRepo := repository.NewReservationRepository(pool)
	saleRepo := repository.NewSaleRepository(pool)
	serviceRepo := repository.NewServiceMenuRepository(pool)
	settingRepo := repository.NewSettingRepository(pool)

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize services
	authService := service.NewAuthService(staffRepo)
	staffService := service.NewStaffService(staffRepo)
	customerService := service.NewCustomerService(customerRepo, chartRepo, membershipRepo)
	chartService := service.NewChartService(chartRepo)
	membershipService := service.NewMembershipService(membershipRepo)
	reservationService := service.NewReservationService(reservationRepo, customerRepo, hub)
	saleService := service.NewSaleService(saleRepo, membershipRepo, customerRepo)
	serviceService := service.NewServiceMenuService(serviceRepo)
	analyticsService := service.NewAnalyticsService(saleRepo)
	marketingService := service.NewMarketingService(saleRepo)
	settingService := service.NewSettingService(settingRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	staffHandler := handler.NewStaffHandler(staffService)
	customerHandler := handler.NewCustomerHandler(customerService)
	chartHandler := handler.NewChartHandler(chartService)
	membershipHandler := handler.NewMembershipHandler(membershipService)
	reservationHandler := handler.NewReservationHandler(reservationService)
	saleHandler := handler.NewSaleHandler(saleService)
	serviceHandler := handler.NewServiceMenuHandler(serviceService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
	marketingHandler := handler.NewMarketingHandler(marketingService)
	naverSyncHandler := handler.NewNaverSyncHandler()
	settingHandler := handler.NewSettingHandler(settingService)

	// Initialize cron scheduler
	cronScheduler := scheduler.NewScheduler(pool, customerRepo, membershipRepo)
	cronScheduler.Start()
	defer cronScheduler.Stop()

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Apply middleware
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS(middleware.CORSConfig{
		AllowedOrigins: []string{cfg.FrontendURL},
	}))
	r.Use(middleware.RateLimit(cfg.RateLimitRPS))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Format(time.RFC3339)})
	})

	// WebSocket endpoint
	r.GET("/ws", hub.HandleWebSocket)

	// API routes
	api := r.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// Webhook routes (public/verified internally)
		syncPublic := api.Group("/sync")
		{
			syncPublic.POST("/naver/webhook", naverSyncHandler.HandleWebhook)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Auth
			protected.GET("/auth/me", authHandler.Me)

			// Staff management (admin only for create/delete)
			staff := protected.Group("/staffs")
			{
				staff.GET("", staffHandler.List)
				staff.GET("/:id", staffHandler.GetByID)
				staff.PUT("/:id", middleware.RequireRole("admin"), staffHandler.Update)
				staff.DELETE("/:id", middleware.RequireRole("admin"), staffHandler.Delete)
			}

			// Staff registration (admin only)
			protected.POST("/auth/register", middleware.RequireRole("admin"), authHandler.Register)

			// Service Menu management
			services := protected.Group("/services")
			{
				services.GET("", serviceHandler.GetAll)
				services.POST("", middleware.RequireRole("admin"), serviceHandler.Create)
				services.PUT("/:id", middleware.RequireRole("admin"), serviceHandler.Update)
				services.DELETE("/:id", middleware.RequireRole("admin"), serviceHandler.Delete)
			}

			// Customer management
			customers := protected.Group("/customers")
			{
				customers.POST("", customerHandler.Create)
				customers.GET("", customerHandler.List)
				customers.GET("/search", customerHandler.Search)
				customers.GET("/lookup", customerHandler.LookupByPhone)
				customers.GET("/:id", customerHandler.GetByID)
				customers.PUT("/:id", customerHandler.Update)
				customers.DELETE("/:id", customerHandler.Delete)
			}

			// Chart management
			charts := protected.Group("/charts")
			{
				charts.POST("", chartHandler.Create)
				charts.GET("/:id", chartHandler.GetByID)
				charts.GET("/customer/:customerId", chartHandler.ListByCustomer)
				charts.DELETE("/:id", chartHandler.Delete)
			}

			// Membership management
			memberships := protected.Group("/memberships")
			{
				memberships.POST("", membershipHandler.Create)
				memberships.GET("/:id", membershipHandler.GetByID)
				memberships.GET("/customer/:customerId", membershipHandler.ListByCustomer)
				memberships.GET("/:id/transactions", membershipHandler.GetTransactions)
			}

			// Reservation management
			reservations := protected.Group("/reservations")
			{
				reservations.POST("", reservationHandler.Create)
				reservations.GET("/:id", reservationHandler.GetByID)
				reservations.GET("", reservationHandler.ListByDate)
				reservations.GET("/range", reservationHandler.ListByDateRange)
				reservations.GET("/waiting", reservationHandler.GetWaitingQueue)
				reservations.POST("/waiting", reservationHandler.AddToWaitingQueue)
				reservations.PATCH("/:id/status", reservationHandler.UpdateStatus)
				reservations.PUT("/:id", reservationHandler.Update)
			}

			// Sales / POS
			sales := protected.Group("/sales")
			{
				sales.POST("", saleHandler.Create)
				sales.GET("/:id", saleHandler.GetByID)
				sales.GET("", saleHandler.ListByDateRange)
				sales.GET("/staff/:staffId", saleHandler.ListByStaff)
				sales.DELETE("/:id", middleware.RequireRole("admin"), saleHandler.Delete)
			}

			// Analytics (admin/designer)
			analytics := protected.Group("/analytics")
			analytics.Use(middleware.RequireRole("admin", "designer"))
			{
				analytics.GET("/daily", analyticsHandler.GetDailySummary)
				analytics.GET("/monthly", analyticsHandler.GetMonthlyTrend)
				analytics.GET("/summary", analyticsHandler.GetAnalyticsSummary)
				analytics.GET("/payments", analyticsHandler.GetPaymentBreakdown)
				analytics.GET("/customers", analyticsHandler.GetCustomerStats)
				analytics.GET("/churn", analyticsHandler.GetChurnAnalysis)
				analytics.GET("/staff-performance", analyticsHandler.GetStaffPerformance)
			}

			// Marketing (admin only)
			marketing := protected.Group("/marketing")
			marketing.Use(middleware.RequireRole("admin"))
			{
				marketing.POST("/targets", marketingHandler.ExtractTargets)
			}

			// External System Sync
			sync := protected.Group("/sync")
			{
				sync.POST("/naver/manual", middleware.RequireRole("admin"), naverSyncHandler.TriggerManualSync)
			}

			// Settings (admin only)
			settings := protected.Group("/settings")
			settings.Use(middleware.RequireRole("admin"))
			{
				settings.GET("", settingHandler.GetSettings)
				settings.PUT("", settingHandler.UpdateSettingsBatch)
			}
		}
	}

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.BackendPort,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Salon Core API server starting on port %s", cfg.BackendPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}

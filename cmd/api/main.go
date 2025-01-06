// @openapi 3.0.0
// @title           Restaurant Management API
// @version         1.0
// @description     A Restaurant Management System API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/KNLopez/restaurant-api/docs"
	"github.com/KNLopez/restaurant-api/internal/config"
	"github.com/KNLopez/restaurant-api/internal/handler"
	"github.com/KNLopez/restaurant-api/internal/repository/postgres"
	"github.com/KNLopez/restaurant-api/internal/router"
	"github.com/KNLopez/restaurant-api/internal/service"
	"github.com/KNLopez/restaurant-api/internal/utils"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	// Initialize database
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	restaurantRepo := postgres.NewRestaurantRepository(db)
	menuRepo := postgres.NewMenuRepository(db)
	orderRepo := postgres.NewOrderRepository(db)
	tableRepo := postgres.NewTableRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	restaurantService := service.NewRestaurantService(restaurantRepo)
	menuService := service.NewMenuService(menuRepo)
	orderService := service.NewOrderService(orderRepo)
	tableService := service.NewTableService(tableRepo)

	// Initialize Cloudinary
	cloudinary, err := utils.NewCloudinaryService(
		cfg.Cloudinary.CloudName,
		cfg.Cloudinary.APIKey,
		cfg.Cloudinary.APISecret,
	)
	if err != nil {
		log.Fatal("failed to initialize cloudinary:", err)
	}

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	restaurantHandler := handler.NewRestaurantHandler(restaurantService, cloudinary)
	menuHandler := handler.NewMenuHandler(menuService, cloudinary)
	orderHandler := handler.NewOrderHandler(orderService)
	tableHandler := handler.NewTableHandler(tableService, cfg)

	// Setup router
	router := router.NewRouter(userHandler, restaurantHandler, menuHandler, orderHandler, tableHandler)

	// Create server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("Server starting on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed to start:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Println("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}

	log.Println("Server exited properly")
}

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

	_ "github.com/lib/pq"
	"github.com/yourusername/restaurant-api/internal/config"
	"github.com/yourusername/restaurant-api/internal/handler"
	"github.com/yourusername/restaurant-api/internal/repository/postgres"
	"github.com/yourusername/restaurant-api/internal/router"
	"github.com/yourusername/restaurant-api/internal/service"
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

	// Initialize services
	userService := service.NewUserService(userRepo)
	restaurantService := service.NewRestaurantService(restaurantRepo)
	menuService := service.NewMenuService(menuRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	restaurantHandler := handler.NewRestaurantHandler(restaurantService)
	menuHandler := handler.NewMenuHandler(menuService)

	// Setup router
	router := router.NewRouter(userHandler, restaurantHandler, menuHandler)

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

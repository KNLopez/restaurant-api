package router

import (
	"net/http"

	"github.com/yourusername/restaurant-api/internal/handler"
	"github.com/yourusername/restaurant-api/internal/middleware"
)

func NewRouter(
	userHandler *handler.UserHandler,
	restaurantHandler *handler.RestaurantHandler,
	menuHandler *handler.MenuHandler,
) http.Handler {
	// Initialize router
	mux := http.NewServeMux()

	// Apply global middleware
	handler := middleware.Chain(
		middleware.Logger,
		middleware.Recoverer,
		middleware.CORS,
	)

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// User routes
	mux.HandleFunc("POST /api/v1/users", userHandler.Create)
	mux.HandleFunc("GET /api/v1/users/{id}", userHandler.Get)

	// Restaurant routes
	mux.HandleFunc("POST /api/v1/restaurants", restaurantHandler.Create)
	mux.HandleFunc("GET /api/v1/restaurants/{id}", restaurantHandler.Get)

	// Menu routes
	mux.HandleFunc("POST /api/v1/restaurants/{id}/menu-items", menuHandler.Create)
	mux.HandleFunc("GET /api/v1/restaurants/{id}/menu-items", menuHandler.List)

	return handler(mux)
}

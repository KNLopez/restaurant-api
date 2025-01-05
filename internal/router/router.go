package router

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
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

	// Swagger documentation
	mux.HandleFunc("GET /swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// API v1 routes
	mux.HandleFunc("POST /api/v1/users", userHandler.Create)
	mux.HandleFunc("GET /api/v1/users/{id}", userHandler.Get)

	mux.HandleFunc("POST /api/v1/restaurants", restaurantHandler.Create)
	mux.HandleFunc("GET /api/v1/restaurants/{id}", restaurantHandler.Get)

	mux.HandleFunc("POST /api/v1/restaurants/{id}/menu-items", menuHandler.Create)
	mux.HandleFunc("GET /api/v1/restaurants/{id}/menu-items", menuHandler.List)

	return handler(mux)
}

package router

import (
	"net/http"

	"github.com/KNLopez/restaurant-api/internal/handler"
	"github.com/KNLopez/restaurant-api/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(
	userHandler *handler.UserHandler,
	restaurantHandler *handler.RestaurantHandler,
	menuHandler *handler.MenuHandler,
	orderHandler *handler.OrderHandler,
	tableHandler *handler.TableHandler,
) http.Handler {
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

	// Register routes by resource
	registerUserRoutes(mux, userHandler)
	registerRestaurantRoutes(mux, restaurantHandler)
	registerMenuRoutes(mux, menuHandler)
	registerOrderRoutes(mux, orderHandler)
	registerTableRoutes(mux, tableHandler)

	return handler(mux)
}

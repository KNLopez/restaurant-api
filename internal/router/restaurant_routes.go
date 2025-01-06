package router

import (
	"net/http"

	"github.com/KNLopez/restaurant-api/internal/constants"
	"github.com/KNLopez/restaurant-api/internal/handler"
)

func registerRestaurantRoutes(mux *http.ServeMux, h *handler.RestaurantHandler) {
	mux.HandleFunc("POST "+constants.RestaurantsRoute, h.Create)
	mux.HandleFunc("GET "+constants.RestaurantsRoute+"/{id}", h.Get)
	mux.HandleFunc("PUT "+constants.RestaurantsRoute+"/{id}", h.Update)
	mux.HandleFunc("DELETE "+constants.RestaurantsRoute+"/{id}", h.Delete)
}

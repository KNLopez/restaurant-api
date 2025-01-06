package router

import (
	"net/http"

	"github.com/KNLopez/restaurant-api/internal/constants"
	"github.com/KNLopez/restaurant-api/internal/handler"
)

func registerOrderRoutes(mux *http.ServeMux, h *handler.OrderHandler) {
	mux.HandleFunc("POST "+constants.OrdersRoute, h.Create)
	mux.HandleFunc("GET "+constants.OrdersRoute+"/{id}", h.Get)
	mux.HandleFunc("PUT "+constants.OrdersRoute+"/{id}", h.Update)
	mux.HandleFunc("PUT "+constants.OrdersRoute+"/{id}/status", h.UpdateStatus)
	mux.HandleFunc("DELETE "+constants.OrdersRoute+"/{id}", h.Delete)
}

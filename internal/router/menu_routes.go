package router

import (
	"net/http"

	"github.com/KNLopez/restaurant-api/internal/constants"
	"github.com/KNLopez/restaurant-api/internal/handler"
)

func registerMenuRoutes(mux *http.ServeMux, h *handler.MenuHandler) {
	base := constants.RestaurantsRoute + "/{id}/menu-items"
	mux.HandleFunc("POST "+base, h.Create)
	mux.HandleFunc("GET "+base+"/{item_id}", h.Get)
	mux.HandleFunc("GET "+base, h.List)
	mux.HandleFunc("PUT "+base+"/{item_id}", h.Update)
	mux.HandleFunc("DELETE "+base+"/{item_id}", h.Delete)
}

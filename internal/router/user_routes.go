package router

import (
	"net/http"

	"github.com/KNLopez/restaurant-api/internal/constants"
	"github.com/KNLopez/restaurant-api/internal/handler"
)

func registerUserRoutes(mux *http.ServeMux, h *handler.UserHandler) {
	mux.HandleFunc("POST "+constants.UsersRoute, h.Create)
	mux.HandleFunc("GET "+constants.UsersRoute+"/{id}", h.Get)
}

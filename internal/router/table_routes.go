package router

import (
	"net/http"

	"github.com/KNLopez/restaurant-api/internal/constants"
	"github.com/KNLopez/restaurant-api/internal/handler"
)

func registerTableRoutes(mux *http.ServeMux, h *handler.TableHandler) {
	base := constants.RestaurantsRoute + "/{id}/tables"
	mux.HandleFunc("POST "+base, h.Create)
	mux.HandleFunc("GET "+base+"/{table_id}", h.Get)
	mux.HandleFunc("GET "+constants.TablesRoute+"/qr/{qr_code}", h.GetByQR)
	mux.HandleFunc("PUT "+base+"/{table_id}/status", h.UpdateStatus)
	mux.HandleFunc("DELETE "+base+"/{table_id}", h.Delete)
}

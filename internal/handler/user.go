package handler

import (
	"net/http"

	"github.com/yourusername/restaurant-api/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement user creation
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement user retrieval
	w.WriteHeader(http.StatusNotImplemented)
}

package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/KNLopez/restaurant-api/internal/models"
	"github.com/KNLopez/restaurant-api/internal/service"
	"github.com/KNLopez/restaurant-api/internal/utils"
	"github.com/google/uuid"
)

type RestaurantHandler struct {
	restaurantService *service.RestaurantService
	cloudinary        *utils.CloudinaryService
}

func NewRestaurantHandler(restaurantService *service.RestaurantService, cloudinary *utils.CloudinaryService) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantService: restaurantService,
		cloudinary:        cloudinary,
	}
}

// Create godoc
// @Summary Create restaurant
// @Description Create a new restaurant
// @Tags restaurants
// @Accept json
// @Produce json
// @Param restaurant body models.Restaurant true "Restaurant object"
// @Success 201 {object} models.Restaurant
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/restaurants [post]
func (h *RestaurantHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	var restaurant models.Restaurant
	if err := json.NewDecoder(strings.NewReader(r.FormValue("data"))).Decode(&restaurant); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Handle logo upload
	file, _, err := r.FormFile("logo")
	if err == nil {
		defer file.Close()
		logoURL, err := h.cloudinary.UploadImage(r.Context(), file, "restaurants/logos")
		if err != nil {
			http.Error(w, "Failed to upload logo", http.StatusInternalServerError)
			return
		}
		restaurant.LogoURL = logoURL
	}

	if err := h.restaurantService.Create(r.Context(), &restaurant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(restaurant)
}

// Get godoc
// @Summary Get restaurant by ID
// @Description Get restaurant details by ID
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path string true "Restaurant ID"
// @Success 200 {object} models.Restaurant
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/restaurants/{id} [get]
func (h *RestaurantHandler) Get(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := uuid.Parse(path[len(path)-1])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	restaurant, err := h.restaurantService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if restaurant == nil {
		http.Error(w, "Restaurant not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurant)
}

// Update godoc
// @Summary Update restaurant
// @Description Update restaurant details
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path string true "Restaurant ID"
// @Param restaurant body models.Restaurant true "Restaurant object"
// @Success 200 {object} models.Restaurant
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/restaurants/{id} [put]
func (h *RestaurantHandler) Update(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := uuid.Parse(path[len(path)-1])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	var restaurant models.Restaurant
	if err := json.NewDecoder(r.Body).Decode(&restaurant); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	restaurant.ID = id

	if err := h.restaurantService.Update(r.Context(), &restaurant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurant)
}

// Delete godoc
// @Summary Delete restaurant
// @Description Delete a restaurant
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path string true "Restaurant ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/restaurants/{id} [delete]
func (h *RestaurantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := uuid.Parse(path[len(path)-1])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	if err := h.restaurantService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

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

type MenuHandler struct {
	menuService *service.MenuService
	cloudinary  *utils.CloudinaryService
}

func NewMenuHandler(menuService *service.MenuService, cloudinary *utils.CloudinaryService) *MenuHandler {
	return &MenuHandler{
		menuService: menuService,
		cloudinary:  cloudinary,
	}
}

// Create godoc
// @Summary Create menu item
// @Description Create a new menu item
// @Tags menu
// @Accept json
// @Produce json
// @Param restaurant_id path string true "Restaurant ID"
// @Param item body models.MenuItem true "Menu item object"
// @Success 201 {object} models.MenuItem
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/restaurants/{restaurant_id}/menu-items [post]
func (h *MenuHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Extract restaurant ID from URL
	path := strings.Split(r.URL.Path, "/")
	restaurantID, err := uuid.Parse(path[len(path)-2])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	var item models.MenuItem
	if err := json.NewDecoder(strings.NewReader(r.FormValue("data"))).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item.RestaurantID = restaurantID

	// Handle multiple image uploads
	form := r.MultipartForm
	files := form.File["images"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to process image", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		imageURL, err := h.cloudinary.UploadImage(r.Context(), file, "menu-items")
		if err != nil {
			http.Error(w, "Failed to upload image", http.StatusInternalServerError)
			return
		}
		item.ImageURLs = append(item.ImageURLs, imageURL)
	}

	if err := h.menuService.Create(r.Context(), &item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// Get godoc
// @Summary Get menu item
// @Description Get menu item by ID
// @Tags menu
// @Accept json
// @Produce json
// @Param restaurant_id path string true "Restaurant ID"
// @Param id path string true "Menu Item ID"
// @Success 200 {object} models.MenuItem
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/restaurants/{restaurant_id}/menu-items/{id} [get]
func (h *MenuHandler) Get(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := uuid.Parse(path[len(path)-1])
	if err != nil {
		http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}

	item, err := h.menuService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if item == nil {
		http.Error(w, "Menu item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// List godoc
// @Summary List menu items
// @Description List all menu items for a restaurant
// @Tags menu
// @Accept json
// @Produce json
// @Param restaurant_id path string true "Restaurant ID"
// @Success 200 {array} models.MenuItem
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/restaurants/{restaurant_id}/menu-items [get]
func (h *MenuHandler) List(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	restaurantID, err := uuid.Parse(path[len(path)-2])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	items, err := h.menuService.List(r.Context(), restaurantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// Update godoc
// @Summary Update menu item
// @Description Update menu item details
// @Tags menu
// @Accept json
// @Produce json
// @Param restaurant_id path string true "Restaurant ID"
// @Param id path string true "Menu Item ID"
// @Param item body models.MenuItem true "Menu item object"
// @Success 200 {object} models.MenuItem
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/restaurants/{restaurant_id}/menu-items/{id} [put]
func (h *MenuHandler) Update(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := uuid.Parse(path[len(path)-1])
	if err != nil {
		http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}

	restaurantID, err := uuid.Parse(path[len(path)-3])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item.ID = id
	item.RestaurantID = restaurantID

	if err := h.menuService.Update(r.Context(), &item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Delete godoc
// @Summary Delete menu item
// @Description Delete a menu item
// @Tags menu
// @Accept json
// @Produce json
// @Param restaurant_id path string true "Restaurant ID"
// @Param id path string true "Menu Item ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/restaurants/{restaurant_id}/menu-items/{id} [delete]
func (h *MenuHandler) Delete(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := uuid.Parse(path[len(path)-1])
	if err != nil {
		http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}

	if err := h.menuService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"encoding/base64"

	"github.com/KNLopez/restaurant-api/internal/config"
	"github.com/KNLopez/restaurant-api/internal/models"
	"github.com/KNLopez/restaurant-api/internal/service"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

type TableHandler struct {
	tableService *service.TableService
	config       *config.Config
}

func NewTableHandler(tableService *service.TableService, config *config.Config) *TableHandler {
	return &TableHandler{
		tableService: tableService,
		config:       config,
	}
}

// Create godoc
// @Summary Create table
// @Description Create a new table with QR code
// @Tags tables
// @Accept json
// @Produce json
// @Param restaurant_id path string true "Restaurant ID"
// @Param table body models.Table true "Table object"
// @Success 201 {object} models.Table
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/restaurants/{restaurant_id}/tables [post]
func (h *TableHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Extract restaurant ID from URL
	path := strings.Split(r.URL.Path, "/")
	restaurantID, err := uuid.Parse(path[len(path)-2])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	var table models.Table
	if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	table.RestaurantID = restaurantID

	// Generate table URL
	baseURL := h.config.BaseURL // You'll need to inject this from config
	table.GenerateTableURL(baseURL)

	// Generate QR code
	qr, err := qrcode.New(table.TableURL, qrcode.Medium)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	// Convert QR code to PNG and then base64
	png, err := qr.PNG(256) // 256x256 pixels
	if err != nil {
		http.Error(w, "Failed to generate QR code image", http.StatusInternalServerError)
		return
	}
	table.QRCode = base64.StdEncoding.EncodeToString(png)

	if err := h.tableService.Create(r.Context(), &table); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(table)
}

// Get godoc
// @Summary Get table by ID
// @Description Get table details by ID
// @Tags tables
// @Accept json
// @Produce json
// @Param restaurant_id path string true "Restaurant ID"
// @Param id path string true "Table ID"
// @Success 200 {object} models.Table
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/restaurants/{restaurant_id}/tables/{id} [get]
func (h *TableHandler) Get(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := uuid.Parse(path[len(path)-1])
	if err != nil {
		http.Error(w, "Invalid table ID", http.StatusBadRequest)
		return
	}

	table, err := h.tableService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if table == nil {
		http.Error(w, "Table not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(table)
}

// GetByQR godoc
// @Summary Get table by QR code
// @Description Get table details by scanning QR code
// @Tags tables
// @Accept json
// @Produce json
// @Param qr_code path string true "QR Code"
// @Success 200 {object} models.Table
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/tables/qr/{qr_code} [get]
func (h *TableHandler) GetByQR(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	qrCode := path[len(path)-1]

	table, err := h.tableService.GetByQRCode(r.Context(), qrCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if table == nil {
		http.Error(w, "Table not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(table)
}

// UpdateStatus godoc
// @Summary Update table status
// @Description Update table status (available/occupied/reserved)
// @Tags tables
// @Accept json
// @Produce json
// @Param restaurant_id path string true "Restaurant ID"
// @Param id path string true "Table ID"
// @Param status body models.TableStatus true "New status"
// @Success 200 {object} models.Table
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/restaurants/{restaurant_id}/tables/{id}/status [put]
func (h *TableHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := uuid.Parse(path[len(path)-2])
	if err != nil {
		http.Error(w, "Invalid table ID", http.StatusBadRequest)
		return
	}

	var status struct {
		Status models.TableStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.tableService.UpdateStatus(r.Context(), id, status.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete godoc
// @Summary Delete table
// @Description Delete a table
// @Tags tables
// @Accept json
// @Produce json
// @Param restaurant_id path string true "Restaurant ID"
// @Param id path string true "Table ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/restaurants/{restaurant_id}/tables/{id} [delete]
func (h *TableHandler) Delete(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := uuid.Parse(path[len(path)-1])
	if err != nil {
		http.Error(w, "Invalid table ID", http.StatusBadRequest)
		return
	}

	if err := h.tableService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/JerkyTreats/template-goapi/internal/logging"
)

// HealthResponse represents the JSON response for health checks
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() (*HealthHandler, error) {
	return &HealthHandler{}, nil
}

// ServeHTTP handles health check requests and returns JSON status
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logging.Debug("Processing health check request")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status: "HEALTHY",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logging.Error("Failed to encode health response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logging.Debug("Health check completed successfully")
}
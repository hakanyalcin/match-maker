package api

import (
	"encoding/json"
	"log"
	"net/http"

	"matchmaking-httpapi/pkg/matchmaker"
	"matchmaking-httpapi/pkg/metrics"

	"github.com/gorilla/mux"
)

// JoinRequest represents the request body for the JOIN endpoint
type JoinRequest struct {
	ID string `json:"id"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// Handler handles HTTP requests for the matchmaking service
type Handler struct {
	matchmaker *matchmaker.Matchmaker
	metrics    *metrics.Metrics
}

// NewHandler creates a new handler with the given matchmaker
func NewHandler(m *matchmaker.Matchmaker, metrics *metrics.Metrics) *Handler {
	return &Handler{
		matchmaker: m,
		metrics:    metrics,
	}
}

// RegisterRoutes registers the routes for the handler
func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/join", h.JoinHandler).Methods("POST")
	r.HandleFunc("/status/{match_id}", h.StatusHandler).Methods("GET")
}

// JoinHandler handles the POST /join endpoint
func (h *Handler) JoinHandler(w http.ResponseWriter, r *http.Request) {
	h.metrics.IncrementRequestCount("join")
	start := h.metrics.StartTimer("join")
	defer h.metrics.StopTimer("join", start)

	// Parse request body
	var req JoinRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.ID == "" {
		respondWithError(w, http.StatusBadRequest, "Player ID is required")
		return
	}

	// Add player to matchmaking
	match, err := h.matchmaker.AddPlayer(req.ID)
	if err != nil {
		log.Printf("Error adding player to matchmaking: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Respond with match info
	respondWithJSON(w, http.StatusOK, match)
}

// StatusHandler handles the GET /status/{match_id} endpoint
func (h *Handler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	h.metrics.IncrementRequestCount("status")
	start := h.metrics.StartTimer("status")
	defer h.metrics.StopTimer("status", start)

	// Get match ID from path
	vars := mux.Vars(r)
	matchID := vars["match_id"]

	// Get match
	match, err := h.matchmaker.GetMatch(matchID)
	if err != nil {
		if err == matchmaker.ErrMatchNotFound {
			respondWithError(w, http.StatusNotFound, "match not found")
			return
		}
		log.Printf("Error getting match: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Respond with match info
	respondWithJSON(w, http.StatusOK, match)
}

// Helper function to respond with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Helper function to respond with an error
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: message})
} 
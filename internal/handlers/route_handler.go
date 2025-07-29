package handlers

import (
	"net/http"
	"strings"
)

// HealthHandler handles health check endpoint
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		SendMethodNotAllowedResponse(w)
		return
	}

	response := map[string]interface{}{
		"status":  "ok",
		"service": "nft-platform-sample",
		"version": "1.0.0",
	}

	SendJSONResponse(w, http.StatusOK, response)
}

// RouteHandler helps with routing logic
type RouteHandler struct {
	AuthHandler   *AuthHandler
	PaperHandler  *PaperHandler
	ReviewHandler *ReviewHandler
	HealthHandler *HealthHandler
}

func NewRouteHandler(
	authHandler *AuthHandler,
	paperHandler *PaperHandler,
	reviewHandler *ReviewHandler,
) *RouteHandler {
	return &RouteHandler{
		AuthHandler:   authHandler,
		PaperHandler:  paperHandler,
		ReviewHandler: reviewHandler,
		HealthHandler: NewHealthHandler(),
	}
}

// HandlePapers routes paper-related requests
func (h *RouteHandler) HandlePapers(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/papers")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.PaperHandler.ListPapers(w, r)
		case http.MethodPost:
			h.PaperHandler.CreatePaper(w, r)
		default:
			SendMethodNotAllowedResponse(w)
		}
		return
	}

	// Handle papers/{id} routes
	parts := strings.Split(path, "/")
	if len(parts) >= 1 {
		switch r.Method {
		case http.MethodGet:
			h.PaperHandler.GetPaper(w, r)
		case http.MethodPut:
			h.PaperHandler.UpdatePaper(w, r)
		case http.MethodDelete:
			h.PaperHandler.DeletePaper(w, r)
		default:
			SendMethodNotAllowedResponse(w)
		}
	}
}

// HandleReviews routes review-related requests
func (h *RouteHandler) HandleReviews(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/reviews")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			// Handle GetReviews here if needed
			SendMethodNotAllowedResponse(w)
		case http.MethodPost:
			h.ReviewHandler.CreateReview(w, r)
		default:
			SendMethodNotAllowedResponse(w)
		}
		return
	}

	// Handle reviews/{id} routes
	parts := strings.Split(path, "/")
	if len(parts) >= 1 {
		switch r.Method {
		case http.MethodGet:
			h.ReviewHandler.GetReview(w, r)
		case http.MethodPut:
			h.ReviewHandler.UpdateReview(w, r)
		case http.MethodDelete:
			h.ReviewHandler.DeleteReview(w, r)
		default:
			SendMethodNotAllowedResponse(w)
		}
	}
}

// HandlePaperReviews routes paper review-related requests
func (h *RouteHandler) HandlePaperReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.ReviewHandler.GetPaperReviews(w, r)
	} else {
		SendMethodNotAllowedResponse(w)
	}
}

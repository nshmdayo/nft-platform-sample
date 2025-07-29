package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nshmdayo/nft-platform-sample/internal/dto"
	"github.com/nshmdayo/nft-platform-sample/internal/errors"
	"github.com/nshmdayo/nft-platform-sample/pkg/logger"
)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct{}

// SendResponse sends a standardized JSON response
func (h *BaseHandler) SendResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	response := dto.APIResponse{
		Success: statusCode >= 200 && statusCode < 300,
		Data:    data,
		Meta: &dto.MetaInfo{
			Timestamp: time.Now(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode JSON response", "error", err)
	}
}

// SendError sends a standardized error response
func (h *BaseHandler) SendError(w http.ResponseWriter, err error) {
	var appErr *errors.AppError
	if errors.IsAppError(err) {
		appErr = err.(*errors.AppError)
	} else {
		appErr = errors.Internal("Internal server error", err)
	}

	response := dto.APIResponse{
		Success: false,
		Error: &dto.ErrorInfo{
			Code:    string(appErr.Code),
			Message: appErr.Message,
			Details: appErr.Details,
		},
		Meta: &dto.MetaInfo{
			Timestamp: time.Now(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.StatusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode error response", "error", err)
	}
}

// DecodeJSON decodes JSON request body into the provided struct
func (h *BaseHandler) DecodeJSON(r *http.Request, dest interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.BadRequest("Content-Type must be application/json")
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dest); err != nil {
		return errors.BadRequest("Invalid JSON format")
	}

	return nil
}

// ExtractIDFromPath extracts ID from URL path
func (h *BaseHandler) ExtractIDFromPath(path string, resourceName string) (uint, error) {
	// This is a simplified implementation
	// In a real application, you might want to use a proper router
	// that handles path parameters more elegantly
	return 0, errors.BadRequest("Path parameter extraction not implemented")
}

// ValidateMethod checks if the request method is allowed
func (h *BaseHandler) ValidateMethod(r *http.Request, allowedMethods ...string) error {
	for _, method := range allowedMethods {
		if r.Method == method {
			return nil
		}
	}
	return errors.New(errors.ErrBadRequest, "Method not allowed")
}

// ParseQueryParams parses common query parameters
func (h *BaseHandler) ParseQueryParams(r *http.Request) *dto.QueryParams {
	query := r.URL.Query()

	params := &dto.QueryParams{
		Search:   query.Get("search"),
		Category: query.Get("category"),
		Status:   query.Get("status"),
		SortBy:   query.Get("sort_by"),
		SortDir:  query.Get("sort_dir"),
	}

	// Parse page
	if pageStr := query.Get("page"); pageStr != "" {
		if page := parseInt(pageStr, 1); page > 0 {
			params.Page = page
		}
	}

	// Parse page size
	if pageSizeStr := query.Get("page_size"); pageSizeStr != "" {
		if pageSize := parseInt(pageSizeStr, 10); pageSize > 0 && pageSize <= 100 {
			params.PageSize = pageSize
		}
	}

	params.SetDefaults()
	return params
}

// Helper function to parse integer with default value
func parseInt(s string, defaultVal int) int {
	// Simplified implementation - in production use strconv.Atoi
	if s == "" {
		return defaultVal
	}
	// TODO: Implement proper integer parsing
	return defaultVal
}

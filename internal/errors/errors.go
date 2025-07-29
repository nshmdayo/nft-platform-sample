package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents application error codes
type ErrorCode string

const (
	// Generic errors
	ErrInternal   ErrorCode = "INTERNAL_ERROR"
	ErrBadRequest ErrorCode = "BAD_REQUEST"
	ErrNotFound   ErrorCode = "NOT_FOUND"
	ErrConflict   ErrorCode = "CONFLICT"

	// Authentication errors
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrForbidden    ErrorCode = "FORBIDDEN"
	ErrInvalidToken ErrorCode = "INVALID_TOKEN"
	ErrTokenExpired ErrorCode = "TOKEN_EXPIRED"

	// Validation errors
	ErrValidation    ErrorCode = "VALIDATION_ERROR"
	ErrMissingField  ErrorCode = "MISSING_FIELD"
	ErrInvalidFormat ErrorCode = "INVALID_FORMAT"

	// Business logic errors
	ErrUserExists      ErrorCode = "USER_EXISTS"
	ErrUserNotFound    ErrorCode = "USER_NOT_FOUND"
	ErrPaperNotFound   ErrorCode = "PAPER_NOT_FOUND"
	ErrReviewNotFound  ErrorCode = "REVIEW_NOT_FOUND"
	ErrInvalidOwner    ErrorCode = "INVALID_OWNER"
	ErrAlreadyReviewed ErrorCode = "ALREADY_REVIEWED"
)

// AppError represents an application error
type AppError struct {
	Code       ErrorCode   `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	Cause      error       `json:"-"`
	StatusCode int         `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// New creates a new application error
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: getStatusCode(code),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Cause:      err,
		StatusCode: getStatusCode(code),
	}
}

// WithDetails adds details to an error
func (e *AppError) WithDetails(details interface{}) *AppError {
	e.Details = details
	return e
}

// WithStatusCode sets a custom status code
func (e *AppError) WithStatusCode(code int) *AppError {
	e.StatusCode = code
	return e
}

// getStatusCode maps error codes to HTTP status codes
func getStatusCode(code ErrorCode) int {
	switch code {
	case ErrBadRequest, ErrValidation, ErrMissingField, ErrInvalidFormat:
		return http.StatusBadRequest
	case ErrUnauthorized, ErrInvalidToken, ErrTokenExpired:
		return http.StatusUnauthorized
	case ErrForbidden, ErrInvalidOwner:
		return http.StatusForbidden
	case ErrNotFound, ErrUserNotFound, ErrPaperNotFound, ErrReviewNotFound:
		return http.StatusNotFound
	case ErrConflict, ErrUserExists, ErrAlreadyReviewed:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// Common error constructors

// BadRequest creates a bad request error
func BadRequest(message string) *AppError {
	return New(ErrBadRequest, message)
}

// NotFound creates a not found error
func NotFound(resource string) *AppError {
	return New(ErrNotFound, fmt.Sprintf("%s not found", resource))
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string) *AppError {
	if message == "" {
		message = "Authentication required"
	}
	return New(ErrUnauthorized, message)
}

// Forbidden creates a forbidden error
func Forbidden(message string) *AppError {
	if message == "" {
		message = "Access denied"
	}
	return New(ErrForbidden, message)
}

// Internal creates an internal server error
func Internal(message string, cause error) *AppError {
	return Wrap(cause, ErrInternal, message)
}

// Validation creates a validation error
func Validation(message string, details interface{}) *AppError {
	return New(ErrValidation, message).WithDetails(details)
}

// Conflict creates a conflict error
func Conflict(message string) *AppError {
	return New(ErrConflict, message)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// AsAppError converts an error to AppError if possible
func AsAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return Internal("Internal server error", err)
}

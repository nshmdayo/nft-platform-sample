package dto

import "time"

// APIResponse represents a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

// ErrorInfo represents error information
type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// MetaInfo represents metadata for responses
type MetaInfo struct {
	Timestamp  time.Time   `json:"timestamp"`
	RequestID  string      `json:"request_id,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Auth DTOs
type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Institution string `json:"institution,omitempty"`
	Role        string `json:"role,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string    `json:"token"`
	User  *UserInfo `json:"user"`
}

type UserInfo struct {
	ID          uint      `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Role        string    `json:"role"`
	Institution string    `json:"institution"`
	CreatedAt   time.Time `json:"created_at"`
}

// Paper DTOs
type CreatePaperRequest struct {
	Title    string   `json:"title" validate:"required,min=5,max=200"`
	Abstract string   `json:"abstract" validate:"required,min=50,max=2000"`
	Authors  []string `json:"authors" validate:"required,min=1"`
	Keywords []string `json:"keywords,omitempty"`
	Category string   `json:"category" validate:"required"`
}

type UpdatePaperRequest struct {
	Title    string   `json:"title,omitempty"`
	Abstract string   `json:"abstract,omitempty"`
	Authors  []string `json:"authors,omitempty"`
	Keywords []string `json:"keywords,omitempty"`
	Category string   `json:"category,omitempty"`
}

type PaperInfo struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Abstract  string    `json:"abstract"`
	Authors   []string  `json:"authors"`
	Keywords  []string  `json:"keywords"`
	Category  string    `json:"category"`
	Status    string    `json:"status"`
	IPFSHash  string    `json:"ipfs_hash,omitempty"`
	OwnerID   uint      `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaperListResponse struct {
	Papers     []PaperInfo `json:"papers"`
	Pagination *Pagination `json:"pagination"`
}

// Review DTOs
type CreateReviewRequest struct {
	PaperID        uint   `json:"paper_id" validate:"required"`
	Comment        string `json:"comment" validate:"required,min=10,max=2000"`
	Score          int    `json:"score" validate:"required,min=1,max=10"`
	Recommendation string `json:"recommendation" validate:"required,oneof=accept reject revision"`
}

type UpdateReviewRequest struct {
	Comment        string `json:"comment,omitempty"`
	Score          int    `json:"score,omitempty"`
	Recommendation string `json:"recommendation,omitempty"`
}

type ReviewInfo struct {
	ID             uint      `json:"id"`
	PaperID        uint      `json:"paper_id"`
	ReviewerID     uint      `json:"reviewer_id"`
	Comment        string    `json:"comment"`
	Score          int       `json:"score"`
	Recommendation string    `json:"recommendation"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ReviewListResponse struct {
	Reviews    []ReviewInfo `json:"reviews"`
	Pagination *Pagination  `json:"pagination"`
}

type PaperReviewsResponse struct {
	Reviews     []ReviewInfo `json:"reviews"`
	PaperScore  float64      `json:"paper_score"`
	ReviewCount int          `json:"review_count"`
	AvgScore    float64      `json:"average_score"`
}

// Query DTOs
type QueryParams struct {
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
	Search   string `json:"search,omitempty"`
	Category string `json:"category,omitempty"`
	Status   string `json:"status,omitempty"`
	SortBy   string `json:"sort_by,omitempty"`
	SortDir  string `json:"sort_dir,omitempty" validate:"omitempty,oneof=asc desc"`
}

// SetDefaults sets default values for query parameters
func (q *QueryParams) SetDefaults() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 || q.PageSize > 100 {
		q.PageSize = 10
	}
	if q.SortDir == "" {
		q.SortDir = "desc"
	}
}

// Health DTOs
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

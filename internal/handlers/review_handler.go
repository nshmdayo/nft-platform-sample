package handlers

import (
	"net/http"
	"strconv"

	"github.com/nshmdayo/nft-platform-sample/internal/service"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendMethodNotAllowedResponse(w)
		return
	}

	userID, err := GetUserIDFromContext(r)
	if err != nil {
		SendUnauthorizedResponse(w)
		return
	}

	var req service.CreateReviewRequest
	if err := DecodeJSONRequest(r, &req); err != nil {
		SendValidationErrorResponse(w, err.Error())
		return
	}

	// Validate required fields
	if err := ValidateRequiredFields(map[string]interface{}{
		"comment":        req.Comment,
		"recommendation": req.Recommendation,
	}); err != nil {
		SendValidationErrorResponse(w, err.Error())
		return
	}

	// Validate score range
	if req.Score < 1 || req.Score > 10 {
		SendValidationErrorResponse(w, "Score must be between 1 and 10")
		return
	}

	// Check eligibility
	if err := h.reviewService.CheckReviewEligibility(req.PaperID, userID); err != nil {
		SendErrorResponse(w, http.StatusForbidden, err.Error())
		return
	}

	review, err := h.reviewService.CreateReview(&req, userID)
	if err != nil {
		SendInternalServerErrorResponse(w, err)
		return
	}

	SendJSONResponse(w, http.StatusCreated, review)
}

func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		SendMethodNotAllowedResponse(w)
		return
	}

	reviewID, err := ExtractIDFromPath(r.URL.Path, "reviews")
	if err != nil {
		SendValidationErrorResponse(w, "Invalid review ID")
		return
	}

	review, err := h.reviewService.GetReview(reviewID)
	if err != nil {
		SendInternalServerErrorResponse(w, err)
		return
	}

	SendJSONResponse(w, http.StatusOK, review)
}

func (h *ReviewHandler) GetPaperReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		SendMethodNotAllowedResponse(w)
		return
	}

	paperID, err := ExtractIDFromPath(r.URL.Path, "papers")
	if err != nil {
		SendValidationErrorResponse(w, "Invalid paper ID")
		return
	}

	reviews, err := h.reviewService.GetPaperReviews(paperID)
	if err != nil {
		SendInternalServerErrorResponse(w, err)
		return
	}

	score, err := h.reviewService.CalculatePaperScore(paperID)
	if err != nil {
		SendInternalServerErrorResponse(w, err)
		return
	}

	response := map[string]interface{}{
		"reviews":     reviews,
		"paper_score": score,
	}

	SendJSONResponse(w, http.StatusOK, response)
}

func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		SendMethodNotAllowedResponse(w)
		return
	}

	userID, err := GetUserIDFromContext(r)
	if err != nil {
		SendUnauthorizedResponse(w)
		return
	}

	reviewID, err := ExtractIDFromPath(r.URL.Path, "reviews")
	if err != nil {
		SendValidationErrorResponse(w, "Invalid review ID")
		return
	}

	var req service.CreateReviewRequest
	if err := DecodeJSONRequest(r, &req); err != nil {
		SendValidationErrorResponse(w, err.Error())
		return
	}

	// Validate score range if provided
	if req.Score < 1 || req.Score > 10 {
		SendValidationErrorResponse(w, "Score must be between 1 and 10")
		return
	}

	review, err := h.reviewService.UpdateReview(reviewID, &req, userID)
	if err != nil {
		SendInternalServerErrorResponse(w, err)
		return
	}

	SendJSONResponse(w, http.StatusOK, review)
}

func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		SendMethodNotAllowedResponse(w)
		return
	}

	userID, err := GetUserIDFromContext(r)
	if err != nil {
		SendUnauthorizedResponse(w)
		return
	}

	reviewID, err := ExtractIDFromPath(r.URL.Path, "reviews")
	if err != nil {
		SendValidationErrorResponse(w, "Invalid review ID")
		return
	}

	if err := h.reviewService.DeleteReview(reviewID, userID); err != nil {
		SendInternalServerErrorResponse(w, err)
		return
	}

	SendJSONResponse(w, http.StatusNoContent, nil)
}

func (h *ReviewHandler) GetMyReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		SendMethodNotAllowedResponse(w)
		return
	}

	userID, err := GetUserIDFromContext(r)
	if err != nil {
		SendUnauthorizedResponse(w)
		return
	}

	// Parse pagination parameters with defaults
	page := 1
	limit := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	reviews, err := h.reviewService.GetReviewerReviews(userID, page, limit)
	if err != nil {
		SendInternalServerErrorResponse(w, err)
		return
	}

	SendJSONResponse(w, http.StatusOK, reviews)
}

func (h *ReviewHandler) GetPendingReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		SendMethodNotAllowedResponse(w)
		return
	}

	userID, err := GetUserIDFromContext(r)
	if err != nil {
		SendUnauthorizedResponse(w)
		return
	}

	// Parse pagination parameters with defaults
	page := 1
	limit := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	papers, err := h.reviewService.GetPendingReviews(userID, page, limit)
	if err != nil {
		SendInternalServerErrorResponse(w, err)
		return
	}

	SendJSONResponse(w, http.StatusOK, papers)
}

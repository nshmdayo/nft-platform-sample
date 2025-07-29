package handlers

import (
	"net/http"
	"strconv"

	"github.com/nshmdayo/nft-platform-sample/internal/service"
)

type PaperHandler struct {
	paperService *service.PaperService
}

func NewPaperHandler(paperService *service.PaperService) *PaperHandler {
	return &PaperHandler{
		paperService: paperService,
	}
}

func (h *PaperHandler) CreatePaper(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendMethodNotAllowedResponse(w)
		return
	}

	userID, err := GetUserIDFromContext(r)
	if err != nil {
		SendUnauthorizedResponse(w)
		return
	}

	var req service.CreatePaperRequest
	if err := DecodeJSONRequest(r, &req); err != nil {
		SendValidationErrorResponse(w, err.Error())
		return
	}

	// Validate required fields
	if err := ValidateRequiredFields(map[string]interface{}{
		"title":    req.Title,
		"abstract": req.Abstract,
		"category": req.Category,
	}); err != nil {
		SendValidationErrorResponse(w, err.Error())
		return
	}

	paper, err := h.paperService.CreatePaper(&req, userID)
	if err != nil {
		SendInternalServerErrorResponse(w, err)
		return
	}

	SendJSONResponse(w, http.StatusCreated, paper)
}

func (h *PaperHandler) GetPaper(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		SendMethodNotAllowedResponse(w)
		return
	}

	paperID, err := ExtractIDFromPath(r.URL.Path, "papers")
	if err != nil {
		SendValidationErrorResponse(w, "Invalid paper ID")
		return
	}

	paper, err := h.paperService.GetPaper(paperID)
	if err != nil {
		SendInternalServerErrorResponse(w, err)
		return
	}

	SendJSONResponse(w, http.StatusOK, paper)
}

func (h *PaperHandler) ListPapers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		SendMethodNotAllowedResponse(w)
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

	papers, err := h.paperService.ListPapers(page, limit)
	if err != nil {
		SendInternalServerErrorResponse(w, err)
		return
	}

	SendJSONResponse(w, http.StatusOK, papers)
}

func (h *PaperHandler) UpdatePaper(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		SendMethodNotAllowedResponse(w)
		return
	}

	userID, err := GetUserIDFromContext(r)
	if err != nil {
		SendUnauthorizedResponse(w)
		return
	}

	paperID, err := ExtractIDFromPath(r.URL.Path, "papers")
	if err != nil {
		SendValidationErrorResponse(w, "Invalid paper ID")
		return
	}

	var req service.UpdatePaperRequest
	if err := DecodeJSONRequest(r, &req); err != nil {
		SendValidationErrorResponse(w, err.Error())
		return
	}

	paper, err := h.paperService.UpdatePaper(paperID, &req, userID)
	if err != nil {
		SendInternalServerErrorResponse(w, err)
		return
	}

	SendJSONResponse(w, http.StatusOK, paper)
}

func (h *PaperHandler) DeletePaper(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		SendMethodNotAllowedResponse(w)
		return
	}

	userID, err := GetUserIDFromContext(r)
	if err != nil {
		SendUnauthorizedResponse(w)
		return
	}

	paperID, err := ExtractIDFromPath(r.URL.Path, "papers")
	if err != nil {
		SendValidationErrorResponse(w, "Invalid paper ID")
		return
	}

	if err := h.paperService.DeletePaper(paperID, userID); err != nil {
		SendInternalServerErrorResponse(w, err)
		return
	}

	SendJSONResponse(w, http.StatusNoContent, nil)
}

func (h *PaperHandler) GetMyPapers(w http.ResponseWriter, r *http.Request) {
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

	papers, err := h.paperService.GetUserPapers(userID, page, limit)
	if err != nil {
		SendInternalServerErrorResponse(w, err)
		return
	}

	SendJSONResponse(w, http.StatusOK, papers)
}

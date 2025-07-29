package service

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/nshmdayo/nft-platform-sample/internal/models"
	"github.com/nshmdayo/nft-platform-sample/internal/repository"
)

type PaperService struct {
	paperRepo *repository.PaperRepository
}

type CreatePaperRequest struct {
	Title    string   `json:"title" binding:"required"`
	Abstract string   `json:"abstract"`
	Authors  []string `json:"authors" binding:"required"`
	Keywords []string `json:"keywords"`
	Category string   `json:"category"`
}

type UpdatePaperRequest struct {
	Title    string   `json:"title"`
	Abstract string   `json:"abstract"`
	Authors  []string `json:"authors"`
	Keywords []string `json:"keywords"`
	Category string   `json:"category"`
}

func NewPaperService(paperRepo *repository.PaperRepository) *PaperService {
	return &PaperService{
		paperRepo: paperRepo,
	}
}

func (s *PaperService) CreatePaper(req *CreatePaperRequest, ownerID uint) (*models.Paper, error) {
	// Convert slices to JSON
	authorsJSON, err := json.Marshal(req.Authors)
	if err != nil {
		return nil, err
	}

	keywordsJSON, err := json.Marshal(req.Keywords)
	if err != nil {
		return nil, err
	}

	paper := &models.Paper{
		Title:    req.Title,
		Abstract: req.Abstract,
		Authors:  authorsJSON,
		Keywords: keywordsJSON,
		Category: req.Category,
		OwnerID:  ownerID,
		Status:   "draft",
	}

	if err := s.paperRepo.Create(paper); err != nil {
		return nil, err
	}

	return paper, nil
}

func (s *PaperService) GetPaper(id uint) (*models.Paper, error) {
	return s.paperRepo.GetByID(id)
}

func (s *PaperService) UpdatePaper(id uint, req *UpdatePaperRequest, userID uint) (*models.Paper, error) {
	paper, err := s.paperRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user owns the paper
	if paper.OwnerID != userID {
		return nil, errors.New("unauthorized to update this paper")
	}

	// Update fields if provided
	if req.Title != "" {
		paper.Title = req.Title
	}
	if req.Abstract != "" {
		paper.Abstract = req.Abstract
	}
	if len(req.Authors) > 0 {
		authorsJSON, err := json.Marshal(req.Authors)
		if err != nil {
			return nil, err
		}
		paper.Authors = authorsJSON
	}
	if len(req.Keywords) > 0 {
		keywordsJSON, err := json.Marshal(req.Keywords)
		if err != nil {
			return nil, err
		}
		paper.Keywords = keywordsJSON
	}
	if req.Category != "" {
		paper.Category = req.Category
	}

	if err := s.paperRepo.Update(paper); err != nil {
		return nil, err
	}

	return paper, nil
}

func (s *PaperService) DeletePaper(id, userID uint) error {
	paper, err := s.paperRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if user owns the paper
	if paper.OwnerID != userID {
		return errors.New("unauthorized to delete this paper")
	}

	return s.paperRepo.Delete(id)
}

func (s *PaperService) ListPapers(page, limit int) ([]models.Paper, error) {
	offset := (page - 1) * limit
	return s.paperRepo.List(limit, offset)
}

func (s *PaperService) GetUserPapers(userID uint, page, limit int) ([]models.Paper, error) {
	offset := (page - 1) * limit
	return s.paperRepo.GetByOwnerID(userID, limit, offset)
}

func (s *PaperService) SearchPapers(query string, page, limit int) ([]models.Paper, error) {
	offset := (page - 1) * limit
	return s.paperRepo.Search(query, limit, offset)
}

func (s *PaperService) SubmitForReview(id, userID uint) (*models.Paper, error) {
	paper, err := s.paperRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user owns the paper
	if paper.OwnerID != userID {
		return nil, errors.New("unauthorized to submit this paper")
	}

	// Check if paper is in draft status
	if paper.Status != "draft" {
		return nil, errors.New("paper is not in draft status")
	}

	paper.Status = "submitted"
	if err := s.paperRepo.Update(paper); err != nil {
		return nil, err
	}

	return paper, nil
}

func parseInt(s string, defaultValue int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return defaultValue
}

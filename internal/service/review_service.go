package service

import (
	"encoding/json"
	"errors"

	"github.com/nshmdayo/nft-platform-sample/internal/models"
	"github.com/nshmdayo/nft-platform-sample/internal/repository"
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	paperRepo  *repository.PaperRepository
}

type CreateReviewRequest struct {
	PaperID        uint   `json:"paper_id" binding:"required"`
	Comment        string `json:"comment" binding:"required"`
	Score          int    `json:"score" binding:"required,min=1,max=10"`
	Recommendation string `json:"recommendation" binding:"required,oneof=accept reject revision"`
}

type ReviewResponse struct {
	ID             uint                   `json:"id"`
	PaperID        uint                   `json:"paper_id"`
	ReviewerID     uint                   `json:"reviewer_id"`
	Comment        string                 `json:"comment"`
	Score          int                    `json:"score"`
	Recommendation string                 `json:"recommendation"`
	Metadata       *models.ReviewMetadata `json:"metadata,omitempty"`
	CreatedAt      string                 `json:"created_at"`
	UpdatedAt      string                 `json:"updated_at"`
}

func NewReviewService(reviewRepo *repository.ReviewRepository, paperRepo *repository.PaperRepository) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		paperRepo:  paperRepo,
	}
}

func (s *ReviewService) CreateReview(req *CreateReviewRequest, reviewerID uint) (*models.Review, error) {
	// Check if paper exists and is in submitted status
	paper, err := s.paperRepo.GetByID(req.PaperID)
	if err != nil {
		return nil, err
	}

	if paper.Status != "submitted" && paper.Status != "under_review" {
		return nil, errors.New("paper is not available for review")
	}

	// Check if reviewer already reviewed this paper
	existingReview, _ := s.reviewRepo.GetByPaperAndReviewer(req.PaperID, reviewerID)
	if existingReview != nil {
		return nil, errors.New("you have already reviewed this paper")
	}

	// Create review metadata
	metadata := &models.ReviewMetadata{
		ReviewCriteria: map[string]interface{}{
			"originality":  0,
			"methodology":  0,
			"clarity":      0,
			"significance": 0,
		},
		ReviewType: "peer_review",
		Version:    1,
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	review := &models.Review{
		PaperID:        req.PaperID,
		ReviewerID:     reviewerID,
		Comment:        req.Comment,
		Score:          req.Score,
		Recommendation: req.Recommendation,
		Metadata:       metadataJSON,
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}

	// Update paper status to under_review if it was submitted
	if paper.Status == "submitted" {
		paper.Status = "under_review"
		s.paperRepo.Update(paper)
	}

	return review, nil
}

func (s *ReviewService) GetReview(id uint) (*models.Review, error) {
	return s.reviewRepo.GetByID(id)
}

func (s *ReviewService) GetPaperReviews(paperID uint) ([]models.Review, error) {
	return s.reviewRepo.GetByPaperID(paperID)
}

func (s *ReviewService) GetReviewerReviews(reviewerID uint, page, limit int) ([]models.Review, error) {
	offset := (page - 1) * limit
	return s.reviewRepo.GetByReviewerID(reviewerID, limit, offset)
}

func (s *ReviewService) UpdateReview(id uint, req *CreateReviewRequest, reviewerID uint) (*models.Review, error) {
	review, err := s.reviewRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user owns the review
	if review.ReviewerID != reviewerID {
		return nil, errors.New("unauthorized to update this review")
	}

	// Update fields
	review.Comment = req.Comment
	review.Score = req.Score
	review.Recommendation = req.Recommendation

	if err := s.reviewRepo.Update(review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *ReviewService) DeleteReview(id, reviewerID uint) error {
	review, err := s.reviewRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if user owns the review
	if review.ReviewerID != reviewerID {
		return errors.New("unauthorized to delete this review")
	}

	return s.reviewRepo.Delete(id)
}

func (s *ReviewService) CalculatePaperScore(paperID uint) (float64, error) {
	reviews, err := s.reviewRepo.GetByPaperID(paperID)
	if err != nil {
		return 0, err
	}

	if len(reviews) == 0 {
		return 0, nil
	}

	totalScore := 0
	for _, review := range reviews {
		totalScore += review.Score
	}

	return float64(totalScore) / float64(len(reviews)), nil
}

func (s *ReviewService) CheckReviewEligibility(paperID, userID uint) error {
	paper, err := s.paperRepo.GetByID(paperID)
	if err != nil {
		return err
	}

	// Authors cannot review their own papers
	if paper.OwnerID == userID {
		return errors.New("authors cannot review their own papers")
	}

	// Check if paper is in a reviewable status
	if paper.Status != "submitted" && paper.Status != "under_review" {
		return errors.New("paper is not available for review")
	}

	// Check if user already reviewed this paper
	existingReview, _ := s.reviewRepo.GetByPaperAndReviewer(paperID, userID)
	if existingReview != nil {
		return errors.New("you have already reviewed this paper")
	}

	return nil
}

func (s *ReviewService) GetPendingReviews(reviewerID uint, page, limit int) ([]models.Paper, error) {
	// This would typically involve a more complex query
	// For now, we'll get papers that are submitted or under_review
	// and exclude papers already reviewed by this reviewer
	offset := (page - 1) * limit
	return s.paperRepo.GetPendingReviews(reviewerID, limit, offset)
}

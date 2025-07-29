package repository

import (
	"github.com/nshmdayo/nft-platform-sample/internal/models"
	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewRepository) GetByID(id uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Preload("Paper").Preload("Reviewer").First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

func (r *ReviewRepository) Delete(id uint) error {
	return r.db.Delete(&models.Review{}, id).Error
}

func (r *ReviewRepository) GetByPaperID(paperID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("paper_id = ?", paperID).Preload("Reviewer").Find(&reviews).Error
	return reviews, err
}

func (r *ReviewRepository) GetByReviewerID(reviewerID uint, limit, offset int) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("reviewer_id = ?", reviewerID).
		Preload("Paper").Limit(limit).Offset(offset).Find(&reviews).Error
	return reviews, err
}

func (r *ReviewRepository) GetByStatus(status string, limit, offset int) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("status = ?", status).
		Preload("Paper").Preload("Reviewer").Limit(limit).Offset(offset).Find(&reviews).Error
	return reviews, err
}

func (r *ReviewRepository) GetByPaperAndReviewer(paperID, reviewerID uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("paper_id = ? AND reviewer_id = ?", paperID, reviewerID).First(&review).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}

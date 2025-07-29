package repository

import (
	"github.com/nshmdayo/nft-platform-sample/internal/models"
	"gorm.io/gorm"
)

type PaperRepository struct {
	db *gorm.DB
}

func NewPaperRepository(db *gorm.DB) *PaperRepository {
	return &PaperRepository{db: db}
}

func (r *PaperRepository) Create(paper *models.Paper) error {
	return r.db.Create(paper).Error
}

func (r *PaperRepository) GetByID(id uint) (*models.Paper, error) {
	var paper models.Paper
	err := r.db.Preload("Owner").Preload("Reviews").First(&paper, id).Error
	if err != nil {
		return nil, err
	}
	return &paper, nil
}

func (r *PaperRepository) Update(paper *models.Paper) error {
	return r.db.Save(paper).Error
}

func (r *PaperRepository) Delete(id uint) error {
	return r.db.Delete(&models.Paper{}, id).Error
}

func (r *PaperRepository) List(limit, offset int) ([]models.Paper, error) {
	var papers []models.Paper
	err := r.db.Preload("Owner").Limit(limit).Offset(offset).Find(&papers).Error
	return papers, err
}

func (r *PaperRepository) GetByOwnerID(ownerID uint, limit, offset int) ([]models.Paper, error) {
	var papers []models.Paper
	err := r.db.Where("owner_id = ?", ownerID).Limit(limit).Offset(offset).Find(&papers).Error
	return papers, err
}

func (r *PaperRepository) GetByStatus(status string, limit, offset int) ([]models.Paper, error) {
	var papers []models.Paper
	err := r.db.Where("status = ?", status).Preload("Owner").Limit(limit).Offset(offset).Find(&papers).Error
	return papers, err
}

func (r *PaperRepository) Search(query string, limit, offset int) ([]models.Paper, error) {
	var papers []models.Paper
	searchQuery := "%" + query + "%"
	err := r.db.Where("title ILIKE ? OR abstract ILIKE ?", searchQuery, searchQuery).
		Preload("Owner").Limit(limit).Offset(offset).Find(&papers).Error
	return papers, err
}

func (r *PaperRepository) GetPendingReviews(reviewerID uint, limit, offset int) ([]models.Paper, error) {
	var papers []models.Paper

	// Get papers that are submitted or under_review
	// and exclude papers already reviewed by this reviewer
	subQuery := r.db.Select("paper_id").Where("reviewer_id = ?", reviewerID).Table("reviews")

	err := r.db.Where("status IN (?, ?) AND id NOT IN (?)", "submitted", "under_review", subQuery).
		Preload("Owner").Limit(limit).Offset(offset).Find(&papers).Error

	return papers, err
}

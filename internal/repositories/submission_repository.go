package repositories

import (
	"context"

	"autograder/internal/models"

	"gorm.io/gorm"
)

type SubmissionRepository interface {
	Create(ctx context.Context, submission *models.Submission) error
	FindByID(ctx context.Context, id uint) (*models.Submission, error)
	// ListByUser mengembalikan maksimal 50 submission terbaru. quizID == 0 berarti semua quiz.
	ListByUser(ctx context.Context, userID, quizID uint) ([]models.Submission, error)
}

type submissionRepository struct{ db *gorm.DB }

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}

func (r *submissionRepository) Create(ctx context.Context, submission *models.Submission) error {
	return translate(r.db.WithContext(ctx).Create(submission).Error)
}

func (r *submissionRepository) FindByID(ctx context.Context, id uint) (*models.Submission, error) {
	var submission models.Submission
	if err := r.db.WithContext(ctx).First(&submission, id).Error; err != nil {
		return nil, translate(err)
	}
	return &submission, nil
}

func (r *submissionRepository) ListByUser(ctx context.Context, userID, quizID uint) ([]models.Submission, error) {
	var list []models.Submission
	q := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if quizID != 0 {
		q = q.Where("quiz_id = ?", quizID)
	}
	err := q.Order("id DESC").Limit(50).Find(&list).Error
	return list, translate(err)
}

package repositories

import (
	"context"
	"encoding/json"

	"autograder/internal/models"

	"gorm.io/gorm"
)

type SubmissionRepository interface {
	Create(ctx context.Context, submission *models.Submission) error
	FindByID(ctx context.Context, id uint) (*models.Submission, error)
	// ListByUser mengembalikan maksimal 50 submission terbaru. quizID == 0 berarti semua quiz.
	ListByUser(ctx context.Context, userID, quizID uint) ([]models.Submission, error)
	// ListByCourse mengembalikan semua submission milik user dalam suatu course.
	ListByCourse(ctx context.Context, userID, courseID uint) ([]models.Submission, error)
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
	if submission.ResultsJson != "" {
		_ = json.Unmarshal([]byte(submission.ResultsJson), &submission.Results)
	}
	return &submission, nil
}

func (r *submissionRepository) ListByCourse(ctx context.Context, userID, courseID uint) ([]models.Submission, error) {
	var list []models.Submission
	// Query: join submissions dengan quizzes, filter berdasarkan user_id dan quiz.course_id
	err := r.db.WithContext(ctx).
		Joins("JOIN quizzes ON quizzes.id = submissions.quiz_id").
		Where("submissions.user_id = ? AND quizzes.course_id = ?", userID, courseID).
		Find(&list).Error
	return list, translate(err)
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

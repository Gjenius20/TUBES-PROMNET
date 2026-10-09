package repositories

import (
	"context"

	"autograder/internal/models"

	"gorm.io/gorm"
)

type QuizRepository interface {
	FindByID(ctx context.Context, id uint) (*models.Quiz, error)
	FindByIDWithTestCases(ctx context.Context, id uint) (*models.Quiz, error)
	Create(ctx context.Context, quiz *models.Quiz) error
	Update(ctx context.Context, id uint, fields map[string]interface{}) error
	Delete(ctx context.Context, id uint) error
}

type quizRepository struct{ db *gorm.DB }

func NewQuizRepository(db *gorm.DB) QuizRepository { return &quizRepository{db: db} }

func (r *quizRepository) FindByID(ctx context.Context, id uint) (*models.Quiz, error) {
	var quiz models.Quiz
	if err := r.db.WithContext(ctx).First(&quiz, id).Error; err != nil {
		return nil, translate(err)
	}
	return &quiz, nil
}

func (r *quizRepository) FindByIDWithTestCases(ctx context.Context, id uint) (*models.Quiz, error) {
	var quiz models.Quiz
	err := r.db.WithContext(ctx).
		Preload("TestCases", func(db *gorm.DB) *gorm.DB { return db.Order("test_cases.id ASC") }).
		First(&quiz, id).Error
	if err != nil {
		return nil, translate(err)
	}
	return &quiz, nil
}

func (r *quizRepository) Create(ctx context.Context, quiz *models.Quiz) error {
	return translate(r.db.WithContext(ctx).Create(quiz).Error)
}

func (r *quizRepository) Update(ctx context.Context, id uint, fields map[string]interface{}) error {
	return translate(r.db.WithContext(ctx).Model(&models.Quiz{}).Where("id = ?", id).Updates(fields).Error)
}

// Delete menghapus quiz; TestCases & Submissions ikut terhapus lewat FK ON DELETE CASCADE.
func (r *quizRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&models.Quiz{}, id)
	if res.Error != nil {
		return translate(res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

package repositories

import (
	"context"

	"autograder/internal/models"

	"gorm.io/gorm"
)

type TestCaseRepository interface {
	ListByQuiz(ctx context.Context, quizID uint) ([]models.TestCase, error)
	Create(ctx context.Context, testCase *models.TestCase) error
	CreateMany(ctx context.Context, testCases []models.TestCase) error
	Delete(ctx context.Context, id uint) error
}

type testCaseRepository struct{ db *gorm.DB }

func NewTestCaseRepository(db *gorm.DB) TestCaseRepository { return &testCaseRepository{db: db} }

func (r *testCaseRepository) ListByQuiz(ctx context.Context, quizID uint) ([]models.TestCase, error) {
	var cases []models.TestCase
	err := r.db.WithContext(ctx).Where("quiz_id = ?", quizID).Order("id ASC").Find(&cases).Error
	return cases, translate(err)
}

func (r *testCaseRepository) Create(ctx context.Context, testCase *models.TestCase) error {
	return translate(r.db.WithContext(ctx).Create(testCase).Error)
}

// CreateMany menyimpan banyak testcase dalam satu transaksi (all-or-nothing).
func (r *testCaseRepository) CreateMany(ctx context.Context, testCases []models.TestCase) error {
	if len(testCases) == 0 {
		return nil
	}
	return translate(r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(&testCases).Error
	}))
}

func (r *testCaseRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&models.TestCase{}, id)
	if res.Error != nil {
		return translate(res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

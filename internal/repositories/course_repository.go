package repositories

import (
	"context"

	"autograder/internal/models"

	"gorm.io/gorm"
)

type CourseRepository interface {
	List(ctx context.Context) ([]models.Course, error)
	FindByID(ctx context.Context, id uint) (*models.Course, error)
	Create(ctx context.Context, course *models.Course) error
	Update(ctx context.Context, id uint, fields map[string]interface{}) error
	Delete(ctx context.Context, id uint) error
}

type courseRepository struct{ db *gorm.DB }

func NewCourseRepository(db *gorm.DB) CourseRepository { return &courseRepository{db: db} }

func orderQuizzes(db *gorm.DB) *gorm.DB { return db.Order("quizzes.id ASC") }

func (r *courseRepository) List(ctx context.Context) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.WithContext(ctx).Preload("Quizzes", orderQuizzes).Order("courses.id ASC").Find(&courses).Error
	return courses, translate(err)
}

func (r *courseRepository) FindByID(ctx context.Context, id uint) (*models.Course, error) {
	var course models.Course
	if err := r.db.WithContext(ctx).Preload("Quizzes", orderQuizzes).First(&course, id).Error; err != nil {
		return nil, translate(err)
	}
	return &course, nil
}

func (r *courseRepository) Create(ctx context.Context, course *models.Course) error {
	return translate(r.db.WithContext(ctx).Create(course).Error)
}

func (r *courseRepository) Update(ctx context.Context, id uint, fields map[string]interface{}) error {
	return translate(r.db.WithContext(ctx).Model(&models.Course{}).Where("id = ?", id).Updates(fields).Error)
}

func (r *courseRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&models.Course{}, id)
	if res.Error != nil {
		return translate(res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

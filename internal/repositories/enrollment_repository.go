package repositories

import (
	"context"

	"autograder/internal/models"

	"gorm.io/gorm"
)

type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *models.Enrollment) error
	FindByID(ctx context.Context, id uint) (*models.Enrollment, error)
	FindByUserAndCourse(ctx context.Context, userID, courseID uint) (*models.Enrollment, error)
	ListByUser(ctx context.Context, userID uint) ([]models.Enrollment, error)
	ListByCourse(ctx context.Context, courseID uint) ([]models.Enrollment, error)
	ListAll(ctx context.Context, status string) ([]models.Enrollment, error)
	Update(ctx context.Context, id uint, fields map[string]interface{}) error
}

type enrollmentRepository struct{ db *gorm.DB }

func NewEnrollmentRepository(db *gorm.DB) EnrollmentRepository {
	return &enrollmentRepository{db: db}
}

func (r *enrollmentRepository) Create(ctx context.Context, enrollment *models.Enrollment) error {
	return translate(r.db.WithContext(ctx).Create(enrollment).Error)
}

func (r *enrollmentRepository) FindByID(ctx context.Context, id uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.WithContext(ctx).Preload("User").Preload("Course").First(&enrollment, id).Error; err != nil {
		return nil, translate(err)
	}
	return &enrollment, nil
}

func (r *enrollmentRepository) FindByUserAndCourse(ctx context.Context, userID, courseID uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.WithContext(ctx).Where("user_id = ? AND course_id = ?", userID, courseID).First(&enrollment).Error; err != nil {
		return nil, translate(err)
	}
	return &enrollment, nil
}

func (r *enrollmentRepository) ListByUser(ctx context.Context, userID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Course").Order("created_at DESC").Find(&enrollments).Error
	return enrollments, translate(err)
}

func (r *enrollmentRepository) ListByCourse(ctx context.Context, courseID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	err := r.db.WithContext(ctx).Where("course_id = ?", courseID).Preload("User").Order("created_at DESC").Find(&enrollments).Error
	return enrollments, translate(err)
}

func (r *enrollmentRepository) ListAll(ctx context.Context, status string) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	q := r.db.WithContext(ctx).Preload("User").Preload("Course").Order("created_at DESC")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	err := q.Find(&enrollments).Error
	return enrollments, translate(err)
}

func (r *enrollmentRepository) Update(ctx context.Context, id uint, fields map[string]interface{}) error {
	return translate(r.db.WithContext(ctx).Model(&models.Enrollment{}).Where("id = ?", id).Updates(fields).Error)
}
package repositories

import (
	"context"

	"autograder/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *models.Enrollment) error
	FindByID(ctx context.Context, id uint) (*models.Enrollment, error)
	FindByUserAndCourse(ctx context.Context, userID, courseID uint) (*models.Enrollment, error)
	ListByUser(ctx context.Context, userID uint) ([]models.Enrollment, error)
	ListByCourse(ctx context.Context, courseID uint) ([]models.Enrollment, error)
	// Admin bulk operations
	BulkUpsertApproved(ctx context.Context, courseID uint, studentIDs []uint) error
	DeleteByUserAndCourse(ctx context.Context, userID, courseID uint) error
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

func (r *enrollmentRepository) BulkUpsertApproved(ctx context.Context, courseID uint, studentIDs []uint) error {
	if len(studentIDs) == 0 {
		return nil
	}
	enrollments := make([]models.Enrollment, len(studentIDs))
	for i, sid := range studentIDs {
		enrollments[i] = models.Enrollment{
			UserID:   sid,
			CourseID: courseID,
			Status:   models.EnrollmentApproved,
		}
	}
	// Use OnConflict to upsert (update status to APPROVED if exists)
	return translate(r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "course_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status"}),
	}).Create(&enrollments).Error)
}

func (r *enrollmentRepository) DeleteByUserAndCourse(ctx context.Context, userID, courseID uint) error {
	return translate(r.db.WithContext(ctx).Where("user_id = ? AND course_id = ?", userID, courseID).Delete(&models.Enrollment{}).Error)
}
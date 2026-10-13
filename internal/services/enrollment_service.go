package services

import (
	"context"
	"errors"

	"autograder/internal/models"
	"autograder/internal/repositories"
)

var (
	ErrEnrollmentNotFound = errors.New("enrollment not found")
	ErrAlreadyEnrolled    = errors.New("already enrolled in this course")
	ErrEnrollmentPending  = errors.New("enrollment request already pending")
	ErrNotApproved        = errors.New("enrollment not approved")
	ErrInvalidStatus      = errors.New("invalid enrollment status")
)

type EnrollmentService struct {
	enrollments repositories.EnrollmentRepository
	courses     repositories.CourseRepository
}

func NewEnrollmentService(enrollments repositories.EnrollmentRepository, courses repositories.CourseRepository) *EnrollmentService {
	return &EnrollmentService{enrollments: enrollments, courses: courses}
}

// EnrollRequest daftar kursus baru (status PENDING).
func (s *EnrollmentService) EnrollRequest(ctx context.Context, userID, courseID uint) (*models.Enrollment, error) {
	// Cek apakah course ada
	if _, err := s.courses.FindByID(ctx, courseID); err != nil {
		return nil, err
	}

	// Cek apakah sudah pernah daftar
	existing, err := s.enrollments.FindByUserAndCourse(ctx, userID, courseID)
	if err == nil {
		switch existing.Status {
		case models.EnrollmentApproved:
			return nil, ErrAlreadyEnrolled
		case models.EnrollmentPending:
			return nil, ErrEnrollmentPending
		case models.EnrollmentRejected:
			// Update status menjadi PENDING lagi
			existing.Status = models.EnrollmentPending
			if err := s.enrollments.Update(ctx, existing.ID, map[string]interface{}{"status": models.EnrollmentPending}); err != nil {
				return nil, err
			}
			return existing, nil
		}
	} else if !errors.Is(err, repositories.ErrNotFound) {
		return nil, err
	}

	// Buat enrollment baru
	enrollment := &models.Enrollment{
		UserID:   userID,
		CourseID: courseID,
		Status:   models.EnrollmentPending,
	}
	if err := s.enrollments.Create(ctx, enrollment); err != nil {
		return nil, err
	}
	return enrollment, nil
}

// GetMyEnrollments ambil daftar enrollment user.
func (s *EnrollmentService) GetMyEnrollments(ctx context.Context, userID uint) ([]models.Enrollment, error) {
	return s.enrollments.ListByUser(ctx, userID)
}

// ListEnrollments admin list semua enrollment (filter status).
func (s *EnrollmentService) ListEnrollments(ctx context.Context, status string) ([]models.Enrollment, error) {
	return s.enrollments.ListAll(ctx, status)
}

// UpdateStatus admin update status enrollment.
func (s *EnrollmentService) UpdateStatus(ctx context.Context, id uint, status models.EnrollmentStatus) (*models.Enrollment, error) {
	if status != models.EnrollmentApproved && status != models.EnrollmentRejected {
		return nil, ErrInvalidStatus
	}

	enrollment, err := s.enrollments.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	enrollment.Status = status
	if err := s.enrollments.Update(ctx, id, map[string]interface{}{"status": status}); err != nil {
		return nil, err
	}
	return enrollment, nil
}

// CheckAccess cek apakah user approved untuk course.
func (s *EnrollmentService) CheckAccess(ctx context.Context, userID, courseID uint) (bool, error) {
	enrollment, err := s.enrollments.FindByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return enrollment.Status == models.EnrollmentApproved, nil
}
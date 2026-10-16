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
)

type EnrollmentService struct {
	enrollments repositories.EnrollmentRepository
	courses     repositories.CourseRepository
	users       repositories.UserRepository
}

func NewEnrollmentService(enrollments repositories.EnrollmentRepository, courses repositories.CourseRepository, users repositories.UserRepository) *EnrollmentService {
	return &EnrollmentService{enrollments: enrollments, courses: courses, users: users}
}

// GetMyEnrollments ambil daftar enrollment user (hanya APPROVED untuk student view).
func (s *EnrollmentService) GetMyEnrollments(ctx context.Context, userID uint) ([]models.Enrollment, error) {
	return s.enrollments.ListByUser(ctx, userID)
}

// ListEnrollmentsByCourse admin list enrollment untuk satu course.
func (s *EnrollmentService) ListEnrollmentsByCourse(ctx context.Context, courseID uint) ([]models.Enrollment, error) {
	return s.enrollments.ListByCourse(ctx, courseID)
}

// BulkEnrollStudents admin bulk enroll siswa ke kursus (Direct Assignment).
func (s *EnrollmentService) BulkEnrollStudents(ctx context.Context, courseID uint, studentIDs []uint, enrolledBy uint) error {
	return s.enrollments.BulkUpsertApproved(ctx, courseID, studentIDs, enrolledBy)
}

// RemoveStudent admin hapus siswa dari kursus.
func (s *EnrollmentService) RemoveStudent(ctx context.Context, courseID, studentID uint) error {
	return s.enrollments.DeleteByUserAndCourse(ctx, studentID, courseID)
}

// ListStudentsPaginated admin list siswa dengan search & pagination.
func (s *EnrollmentService) ListStudentsPaginated(ctx context.Context, search string, page, limit int) ([]models.User, int64, error) {
	return s.users.ListStudentsPaginated(ctx, search, page, limit)
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


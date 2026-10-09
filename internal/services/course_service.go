package services

import (
	"context"
	"strings"

	"autograder/internal/models"
	"autograder/internal/repositories"
)

type CourseService struct {
	courses repositories.CourseRepository
}

func NewCourseService(courses repositories.CourseRepository) *CourseService {
	return &CourseService{courses: courses}
}

func (s *CourseService) List(ctx context.Context) ([]models.Course, error) {
	return s.courses.List(ctx)
}

func (s *CourseService) Get(ctx context.Context, id uint) (*models.Course, error) {
	return s.courses.FindByID(ctx, id)
}

func (s *CourseService) Create(ctx context.Context, title, description string) (*models.Course, error) {
	course := &models.Course{Title: strings.TrimSpace(title), Description: strings.TrimSpace(description)}
	if err := s.courses.Create(ctx, course); err != nil {
		return nil, err
	}
	return course, nil
}

func (s *CourseService) Update(ctx context.Context, id uint, title, description string) (*models.Course, error) {
	if _, err := s.courses.FindByID(ctx, id); err != nil {
		return nil, err
	}
	fields := map[string]interface{}{
		"title":       strings.TrimSpace(title),
		"description": strings.TrimSpace(description),
	}
	if err := s.courses.Update(ctx, id, fields); err != nil {
		return nil, err
	}
	return s.courses.FindByID(ctx, id)
}

func (s *CourseService) Delete(ctx context.Context, id uint) error {
	return s.courses.Delete(ctx, id)
}

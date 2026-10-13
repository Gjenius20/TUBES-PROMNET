package services

import (
	"context"
	"autograder/internal/models"
	"autograder/internal/repositories"
)

type UserService struct {
	users repositories.UserRepository
}

func NewUserService(users repositories.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) ListStudents(ctx context.Context, search string, page, limit int) ([]models.User, int64, error) {
	return s.users.ListStudentsPaginated(ctx, search, page, limit)
}

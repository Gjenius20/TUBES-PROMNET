package services

import (
	"context"
	"errors"
	"strings"

	"autograder/internal/models"
	"autograder/internal/repositories"
	"autograder/internal/utils"
)

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	users  repositories.UserRepository
	tokens *utils.TokenManager
}

func NewAuthService(users repositories.UserRepository, tokens *utils.TokenManager) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}

func normalizeEmail(email string) string { return strings.ToLower(strings.TrimSpace(email)) }

// Register selalu membuat akun STUDENT. Akun ADMIN dibuat lewat cmd/createadmin.
func (s *AuthService) Register(ctx context.Context, name, email, password string) (*models.User, *AuthTokens, error) {
	email = normalizeEmail(email)

	if _, err := s.users.FindByEmail(ctx, email); err == nil {
		return nil, nil, ErrEmailTaken
	} else if !errors.Is(err, repositories.ErrNotFound) {
		return nil, nil, err
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	user := &models.User{
		Name:         strings.TrimSpace(name),
		Email:        email,
		PasswordHash: hash,
		Role:         models.RoleStudent,
	}
	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, repositories.ErrDuplicate) {
			return nil, nil, ErrEmailTaken
		}
		return nil, nil, err
	}

	tokens, err := s.issueTokens(user)
	if err != nil {
		return nil, nil, err
	}
	return user, tokens, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, *AuthTokens, error) {
	user, err := s.users.FindByEmail(ctx, normalizeEmail(email))
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, err
	}
	if !utils.CheckPassword(user.PasswordHash, password) {
		return nil, nil, ErrInvalidCredentials
	}

	tokens, err := s.issueTokens(user)
	if err != nil {
		return nil, nil, err
	}
	return user, tokens, nil
}

// Refresh memvalidasi refresh token lalu menerbitkan pasangan token baru (rotasi).
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*models.User, *AuthTokens, error) {
	claims, err := s.tokens.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, nil, ErrInvalidToken
	}
	user, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, nil, ErrInvalidToken
		}
		return nil, nil, err
	}
	tokens, err := s.issueTokens(user)
	if err != nil {
		return nil, nil, err
	}
	return user, tokens, nil
}

func (s *AuthService) Me(ctx context.Context, userID uint) (*models.User, error) {
	return s.users.FindByID(ctx, userID)
}

func (s *AuthService) issueTokens(user *models.User) (*AuthTokens, error) {
	access, err := s.tokens.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}
	refresh, err := s.tokens.GenerateRefreshToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}
	return &AuthTokens{AccessToken: access, RefreshToken: refresh}, nil
}

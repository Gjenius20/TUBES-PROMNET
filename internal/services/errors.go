package services

import (
	"errors"

	"autograder/internal/repositories"
)

var (
	ErrNotFound           = repositories.ErrNotFound
	ErrEmailTaken         = errors.New("email already registered")
	ErrNIMTaken           = errors.New("NIM already registered")
	ErrInvalidCredentials = errors.New("invalid identifier or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrNoTestCases        = errors.New("quiz has no test cases yet")
)

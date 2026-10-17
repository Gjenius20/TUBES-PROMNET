package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"autograder/internal/models"
	"autograder/internal/repositories"
	"autograder/internal/services/autograder"
	"autograder/internal/services/docker"
)

// Mock repositories and runner for testing
type mockQuizRepository struct {
	mock.Mock
}
func (m *mockQuizRepository) FindByID(ctx context.Context, id uint) (*models.Quiz, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Quiz), args.Error(1)
}
func (m *mockQuizRepository) Create(ctx context.Context, quiz *models.Quiz) (*models.Quiz, error) {
	panic("not implemented")
}
func (m *mockQuizRepository) Update(ctx context.Context, quiz *models.Quiz) (*models.Quiz, error) {
	panic("not implemented")
}
func (m *mockQuizRepository) Delete(ctx context.Context, id uint) error {
	panic("not implemented")
}
func (m *mockQuizRepository) ListByCourse(ctx context.Context, courseID uint) ([]models.Quiz, error) {
	panic("not implemented")
}
func (m *mockQuizRepository) List(ctx context.Context) ([]models.Quiz, error) {
	panic("not implemented")
}

type mockTestCaseRepository struct {
	mock.Mock
}
func (m *mockTestCaseRepository) ListByQuiz(ctx context.Context, quizID uint) ([]models.TestCase, error) {
	args := m.Called(ctx, quizID)
	return args.Get(0).([]models.TestCase), args.Error(1)
}
func (m *mockTestCaseRepository) Create(ctx context.Context, testCase *models.TestCase) (*models.TestCase, error) {
	panic("not implemented")
}
func (m *mockTestCaseRepository) Update(ctx context.Context, testCase *models.TestCase) (*models.TestCase, error) {
	panic("not implemented")
}
func (m *mockTestCaseRepository) Delete(ctx context.Context, id uint) error {
	panic("not implemented")
}

type mockEnrollmentRepository struct {
	mock.Mock
}
func (m *mockEnrollmentRepository) FindByUserAndCourse(ctx context.Context, userID, courseID uint) (*models.Enrollment, error) {
	args := m.Called(ctx, userID, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Enrollment), args.Error(1)
}
func (m *mockEnrollmentRepository) BulkUpsertApproved(ctx context.Context, courseID uint, userIDs []uint) ([]*models.Enrollment, error) {
	panic("not implemented")
}
func (m *mockEnrollmentRepository) Update(ctx context.Context, enrollment *models.Enrollment) (*models.Enrollment, error) {
	panic("not implemented")
}
func (m *mockEnrollmentRepository) Delete(ctx context.Context, id uint) error {
	panic("not implemented")
}
func (m *mockEnrollmentRepository) ListByCourse(ctx context.Context, courseID uint) ([]*models.Enrollment, error) {
	panic("not implemented")
}
func (m *mockEnrollmentRepository) ListByUser(ctx context.Context, userID uint) ([]*models.Enrollment, error) {
	panic("not implemented")
}

type mockSubmissionRepository struct {
	mock.Mock
}
func (m *mockSubmissionRepository) Save(ctx context.Context, sub *models.Submission) (*models.Submission, error) {
	args := m.Called(ctx, sub)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Submission), args.Error(1)
}
func (m *mockSubmissionRepository) Create(ctx context.Context, sub *models.Submission) (*models.Submission, error) {
	panic("not implemented")
}
func (m *mockSubmissionRepository) Update(ctx context.Context, sub *models.Submission) (*models.Submission, error) {
	panic("not implemented")
}
func (m *mockSubmissionRepository) Delete(ctx context.Context, id uint) error {
	panic("not implemented")
}
func (m *mockSubmissionRepository) ListMine(ctx context.Context, userID uint, quizID uint) ([]models.Submission, error) {
	panic("not implemented")
}
func (m *mockSubmissionRepository) Get(ctx context.Context, id uint, userID uint, role models.Role) (*models.Submission, error) {
	panic("not implemented")
}

type mockRunner struct {
	mock.Mock
}
func (m *mockRunner) Run(ctx context.Context, sourceCode, stdin string) (*docker.Result, error) {
	args := m.Called(ctx, sourceCode, stdin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*docker.Result), args.Error(1)
}

func TestSubmissionService_Run(t *testing.T) {
	// Setup
	quizRepo := new(mockQuizRepository)
	testCaseRepo := new(mockTestCaseRepository)
	enrollmentRepo := new(mockEnrollmentRepository)
	submissionRepo := new(mockSubmissionRepository)
	mockRunner := new(mockRunner)

	service := NewSubmissionService(submissionRepo, quizRepo, testCaseRepo, enrollmentRepo, mockRunner)

	ctx := context.Background()
	userID := uint(1)
	quizID := uint(1)
	sourceCode := `#include <stdio.h>
int main(void) {
    printf("Hello, World!");
    return 0;
}`

	// Mock quiz
	quiz := &models.Quiz{ID: quizID, CourseID: 10}
	quizRepo.On("FindByID", ctx, quizID).Return(quiz, nil)

	// Mock enrollment
	enrollment := &models.Enrollment{UserID: userID, CourseID: 10, Status: models.EnrollmentApproved}
	enrollmentRepo.On("FindByUserAndCourse", ctx, userID, 10).Return(enrollment, nil)

	// Mock test cases
	testCases := []models.TestCase{
		{Stdin: "", ExpectedOutput: "Hello, World!", IsHidden: false},
	}
	testCaseRepo.On("ListByQuiz", ctx, quizID).Return(testCases, nil)

	// Mock runner result
	mockRunner.On("Run", ctx, sourceCode, "").
		Return(&docker.Result{Stdout: "Hello, World!", Stderr: "", StatusID: docker.StatusAccepted}, nil)

	// Mock submission repository save (not expected to be called in Run, but we set up expectation that it's not called)
	submissionRepo.On("Save", mock.AnythingOfType("*models.Submission")).Return(nil, nil).Times(0)

	// Execute
	results, err := service.Run(ctx, userID, quizID, sourceCode)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, uint(1), results[0].Index)
	assert.Equal(t, false, results[0].Hidden)
	assert.Equal(t, true, results[0].Passed)
	assert.Equal(t, "OK", results[0].Status)
	assert.Equal(t, "", results[0].Message)
	assert.Equal(t, "", results[0].Stdin)
	assert.Equal(t, "Hello, World!", results[0].ExpectedOutput)
	assert.Equal(t, "Hello, World!", results[0].ActualOutput)

	// Verify all expectations were met
	quizRepo.AssertExpectations(t)
	testCaseRepo.AssertExpectations(t)
	enrollmentRepo.AssertExpectations(t)
	submissionRepo.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

func TestSubmissionService_Run_WrongAnswer(t *testing.T) {
	// Setup
	quizRepo := new(mockQuizRepository)
	testCaseRepo := new(mockTestCaseRepository)
	enrollmentRepo := new(mockEnrollmentRepository)
	submissionRepo := new(mockSubmissionRepository)
	mockRunner := new(mockRunner)

	service := NewSubmissionService(submissionRepo, quizRepo, testCaseRepo, enrollmentRepo, mockRunner)

	ctx := context.Background()
	userID := uint(1)
	quizID := uint(1)
	sourceCode := `#include <stdio.h>
int main(void) {
    printf("Goodbye");
    return 0;
}`

	// Mock quiz
	quiz := &models.Quiz{ID: quizID, CourseID: 10}
	quizRepo.On("FindByID", ctx, quizID).Return(quiz, nil)

	// Mock enrollment
	enrollment := &models.Enrollment{UserID: userID, CourseID: 10, Status: models.EnrollmentApproved}
	enrollmentRepo.On("FindByUserAndCourse", ctx, userID, 10).Return(enrollment, nil)

	// Mock test cases
	testCases := []models.TestCase{
		{Stdin: "", ExpectedOutput: "Hello, World!", IsHidden: false},
	}
	testCaseRepo.On("ListByQuiz", ctx, quizID).Return(testCases, nil)

	// Mock runner result
	mockRunner.On("Run", ctx, sourceCode, "").
		Return(&docker.Result{Stdout: "Goodbye", Stderr: "", StatusID: docker.StatusAccepted}, nil)

	// Mock submission repository save (not expected to be called)
	submissionRepo.On("Save", mock.AnythingOfType("*models.Submission")).Return(nil, nil).Times(0)

	// Execute
	results, err := service.Run(ctx, userID, quizID, sourceCode)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, uint(1), results[0].Index)
	assert.Equal(t, false, results[0].Hidden)
	assert.Equal(t, false, results[0].Passed)
	assert.Equal(t, string(models.StatusWrongAnswer), results[0].Status)
	assert.Equal(t, "", results[0].Message)
	assert.Equal(t, "", results[0].Stdin)
	assert.Equal(t, "Hello, World!", results[0].ExpectedOutput)
	assert.Equal(t, "Goodbye", results[0].ActualOutput)

	// Verify all expectations were met
	quizRepo.AssertExpectations(t)
	testCaseRepo.AssertExpectations(t)
	enrollmentRepo.AssertExpectations(t)
	submissionRepo.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

func TestSubmissionService_Run_CompileError(t *testing.T) {
	// Setup
	quizRepo := new(mockQuizRepository)
	testCaseRepo := new(mockTestCaseRepository)
	enrollmentRepo := new(mockEnrollmentRepository)
	submissionRepo := new(mockSubmissionRepository)
	mockRunner := new(mockRunner)

	service := NewSubmissionService(submissionRepo, quizRepo, testCaseRepo, enrollmentRepo, mockRunner)

	ctx := context.Background()
	userID := uint(1)
	quizID := uint(1)
	sourceCode := `invalid code`

	// Mock quiz
	quiz := &models.Quiz{ID: quizID, CourseID: 10}
	quizRepo.On("FindByID", ctx, quizID).Return(quiz, nil)

	// Mock enrollment
	enrollment := &models.Enrollment{UserID: userID, CourseID: 10, Status: models.EnrollmentApproved}
	enrollmentRepo.On("FindByUserAndCourse", ctx, userID, 10).Return(enrollment, nil)

	// Mock test cases
	testCases := []models.TestCase{
		{Stdin: "", ExpectedOutput: "Hello, World!", IsHidden: false},
	}
	testCaseRepo.On("ListByQuiz", ctx, quizID).Return(testCases, nil)

	// Mock runner result
	mockRunner.On("Run", ctx, sourceCode, "").
		Return(&docker.Result{Stdout: "", Stderr: "", CompileOutput: "error: expected ';' before '}'", StatusID: docker.StatusCompilationError}, nil)

	// Mock submission repository save (not expected to be called)
	submissionRepo.On("Save", mock.AnythingOfType("*models.Submission")).Return(nil, nil).Times(0)

	// Execute
	results, err := service.Run(ctx, userID, quizID, sourceCode)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, uint(1), results[0].Index)
	assert.Equal(t, false, results[0].Hidden)
	assert.Equal(t, false, results[0].Passed)
	assert.Equal(t, string(models.StatusCompileError), results[0].Status)
	assert.Equal(t, "error: expected ';' before '}'", results[0].Message)
	assert.Empty(t, results[0].Stdin)
	assert.Empty(t, results[0].ExpectedOutput)
	assert.Empty(t, results[0].ActualOutput)

	// Verify all expectations were met
	quizRepo.AssertExpectations(t)
	testCaseRepo.AssertExpectations(t)
	enrollmentRepo.AssertExpectations(t)
	submissionRepo.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}
EOF
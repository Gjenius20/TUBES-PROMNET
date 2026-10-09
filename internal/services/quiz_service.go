package services

import (
	"context"
	"fmt"
	"strings"

	"autograder/internal/models"
	"autograder/internal/repositories"
	"autograder/internal/services/ai"
)

const (
	defaultGenerateCount = 5
	// Dua testcase pertama dari AI ditampilkan sebagai contoh; sisanya disembunyikan.
	visibleAISamples = 2
)

type QuizService struct {
	quizzes   repositories.QuizRepository
	testCases repositories.TestCaseRepository
	courses   repositories.CourseRepository
	generator ai.Generator
}

func NewQuizService(
	quizzes repositories.QuizRepository,
	testCases repositories.TestCaseRepository,
	courses repositories.CourseRepository,
	generator ai.Generator,
) *QuizService {
	return &QuizService{quizzes: quizzes, testCases: testCases, courses: courses, generator: generator}
}

// Get mengembalikan quiz beserta SEMUA testcase. Pemanggil bertanggung jawab
// menyaring testcase hidden sebelum dikirim ke siswa.
func (s *QuizService) Get(ctx context.Context, id uint) (*models.Quiz, error) {
	return s.quizzes.FindByIDWithTestCases(ctx, id)
}

func (s *QuizService) Create(ctx context.Context, courseID uint, title, description, constraints string) (*models.Quiz, error) {
	if _, err := s.courses.FindByID(ctx, courseID); err != nil {
		return nil, err
	}
	quiz := &models.Quiz{
		CourseID:    courseID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Constraints: strings.TrimSpace(constraints),
	}
	if err := s.quizzes.Create(ctx, quiz); err != nil {
		return nil, err
	}
	return quiz, nil
}

func (s *QuizService) Update(ctx context.Context, id uint, title, description, constraints string) (*models.Quiz, error) {
	if _, err := s.quizzes.FindByID(ctx, id); err != nil {
		return nil, err
	}
	fields := map[string]interface{}{
		"title":       strings.TrimSpace(title),
		"description": strings.TrimSpace(description),
		"constraints": strings.TrimSpace(constraints),
	}
	if err := s.quizzes.Update(ctx, id, fields); err != nil {
		return nil, err
	}
	return s.quizzes.FindByIDWithTestCases(ctx, id)
}

func (s *QuizService) Delete(ctx context.Context, id uint) error {
	return s.quizzes.Delete(ctx, id)
}

func (s *QuizService) AddTestCase(ctx context.Context, quizID uint, stdin, expectedOutput string, hidden bool) (*models.TestCase, error) {
	if _, err := s.quizzes.FindByID(ctx, quizID); err != nil {
		return nil, err
	}
	tc := &models.TestCase{
		QuizID:         quizID,
		Stdin:          stdin,
		ExpectedOutput: expectedOutput,
		IsHidden:       hidden,
		Source:         models.TestCaseSourceManual,
	}
	if err := s.testCases.Create(ctx, tc); err != nil {
		return nil, err
	}
	return tc, nil
}

func (s *QuizService) DeleteTestCase(ctx context.Context, id uint) error {
	return s.testCases.Delete(ctx, id)
}

// GenerateTestCases meminta AI membuat testcase lalu menyimpannya (setelah divalidasi
// oleh ai.ParseTestCases). problemOverride boleh kosong: bila kosong dipakai deskripsi quiz.
func (s *QuizService) GenerateTestCases(ctx context.Context, quizID uint, problemOverride string, count int) ([]models.TestCase, error) {
	quiz, err := s.quizzes.FindByID(ctx, quizID)
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		count = defaultGenerateCount
	}

	problem := strings.TrimSpace(problemOverride)
	if problem == "" {
		problem = fmt.Sprintf("%s\n\n%s", quiz.Title, quiz.Description)
		if quiz.Constraints != "" {
			problem += "\n\nConstraints:\n" + quiz.Constraints
		}
	}

	generated, err := s.generator.GenerateTestCases(ctx, problem, count)
	if err != nil {
		return nil, err
	}

	cases := make([]models.TestCase, 0, len(generated))
	for i, g := range generated {
		cases = append(cases, models.TestCase{
			QuizID:         quizID,
			Stdin:          g.Stdin,
			ExpectedOutput: g.ExpectedOutput,
			IsHidden:       i >= visibleAISamples,
			Source:         models.TestCaseSourceAI,
		})
	}
	if err := s.testCases.CreateMany(ctx, cases); err != nil {
		return nil, err
	}
	return cases, nil
}

package services

import (
	"context"
	"fmt"
	"strings"

	"autograder/internal/models"
	"autograder/internal/repositories"
	"autograder/internal/services/autograder"
	"autograder/internal/services/docker"
	"encoding/json"
)

const (
	maxOutputPreview = 2000
	maxErrorMessage  = 4000
)

type SubmissionService struct {
	submissions repositories.SubmissionRepository
	quizzes     repositories.QuizRepository
	testCases   repositories.TestCaseRepository
	enrollments repositories.EnrollmentRepository
	runner      docker.Runner
}

func NewSubmissionService(
	submissions repositories.SubmissionRepository,
	quizzes repositories.QuizRepository,
	testCases repositories.TestCaseRepository,
	enrollments repositories.EnrollmentRepository,
	runner docker.Runner,
) *SubmissionService {
	return &SubmissionService{
		submissions: submissions,
		quizzes:     quizzes,
		testCases:   testCases,
		enrollments: enrollments,
		runner:      runner,
	}
}

// Submit menjalankan kode siswa terhadap seluruh testcase lewat Docker,
// menilai hasilnya dengan tolerance matcher, lalu menyimpan submission.
func (s *SubmissionService) Submit(ctx context.Context, userID, quizID uint, sourceCode string) (*models.Submission, error) {
	quiz, err := s.quizzes.FindByID(ctx, quizID)
	if err != nil {
		return nil, err
	}

	// Validasi Enrollment: Pastikan siswa terdaftar di course terkait quiz
	enrolled, err := s.enrollments.FindByUserAndCourse(ctx, userID, quiz.CourseID)
	if err != nil || enrolled == nil || enrolled.Status != models.EnrollmentApproved {
		return nil, fmt.Errorf("akses ditolak: siswa belum terdaftar di kursus ini")
	}

	cases, err := s.testCases.ListByQuiz(ctx, quizID)
	if err != nil {
		return nil, err
	}
	if len(cases) == 0 {
		return nil, ErrNoTestCases
	}

	sub := &models.Submission{
		UserID:     userID,
		QuizID:     quizID,
		SourceCode: sourceCode,
		TotalCount: len(cases),
	}
	results := make([]models.TestResult, 0, len(cases))
	runtimeFailure := false

	for i, tc := range cases {
		res, err := s.runner.Run(ctx, sourceCode, tc.Stdin)
		if err != nil {
			return nil, err
		}

		if res.StatusID == docker.StatusCompilationError {
			sub.Status = models.StatusCompileError
			sub.ErrorMessage = truncate(firstNonEmpty(res.CompileOutput, res.Message, res.StatusDescription), maxErrorMessage)
			return s.save(ctx, sub, nil)
		}

		tr := models.TestResult{Index: i + 1, Hidden: tc.IsHidden}
		if !tc.IsHidden {
			tr.Stdin = tc.Stdin
			tr.ExpectedOutput = tc.ExpectedOutput
		}

		switch {
		case res.StatusID == docker.StatusAccepted:
			if autograder.Match(res.Stdout, tc.ExpectedOutput) {
				tr.Passed = true
				tr.Status = "OK"
				sub.PassedCount++
			} else {
				tr.Status = string(models.StatusWrongAnswer)
				if !tc.IsHidden {
					tr.ActualOutput = truncate(autograder.Normalize(res.Stdout), maxOutputPreview)
				}
			}
		case isRuntimeFailure(res.StatusID):
			runtimeFailure = true
			tr.Status = string(models.StatusRuntimeError)
			tr.Message = res.StatusDescription
			if !tc.IsHidden {
				tr.Message = truncate(firstNonEmpty(res.Stderr, res.Message, res.StatusDescription), maxOutputPreview)
			}
			if sub.ErrorMessage == "" {
				sub.ErrorMessage = truncate(firstNonEmpty(res.Stderr, res.Message, res.StatusDescription), maxErrorMessage)
			}
		default:
			// Status antrean/proses/internal error dari Docker: bukan kesalahan siswa.
			return nil, fmt.Errorf("%w: unexpected status %d (%s)", docker.ErrUnavailable, res.StatusID, res.StatusDescription)
		}
		results = append(results, tr)
	}

	sub.Score = sub.PassedCount * 100 / sub.TotalCount
	switch {
	case sub.PassedCount == sub.TotalCount:
		sub.Status = models.StatusPassed
	case runtimeFailure:
		sub.Status = models.StatusRuntimeError
	default:
		sub.Status = models.StatusWrongAnswer
	}
	return s.save(ctx, sub, results)
}

func (s *SubmissionService) Run(ctx context.Context, userID, quizID uint, sourceCode string) ([]models.TestResult, error) {
	quiz, err := s.quizzes.FindByID(ctx, quizID)
	if err != nil {
		return nil, err
	}

	enrolled, err := s.enrollments.FindByUserAndCourse(ctx, userID, quiz.CourseID)
	if err != nil || enrolled == nil || enrolled.Status != models.EnrollmentApproved {
		return nil, fmt.Errorf("akses ditolak: siswa belum terdaftar di kursus ini")
	}

	cases, err := s.testCases.ListByQuiz(ctx, quizID)
	if err != nil {
		return nil, err
	}
	if len(cases) == 0 {
		return nil, ErrNoTestCases
	}

	results := make([]models.TestResult, 0, len(cases))
	for i, tc := range cases {
		res, err := s.runner.Run(ctx, sourceCode, tc.Stdin)
		if err != nil {
			return nil, err
		}

		tr := models.TestResult{
			Index:  i + 1,
			Hidden: tc.IsHidden,
		}

		if res.StatusID == docker.StatusCompilationError {
			tr.Status = string(models.StatusCompileError)
			tr.Message = truncate(firstNonEmpty(res.CompileOutput, res.Message, res.StatusDescription), maxOutputPreview)
			results = append(results, tr)
			break
		}

		if !tc.IsHidden {
			tr.Stdin = tc.Stdin
			tr.ExpectedOutput = tc.ExpectedOutput
			tr.ActualOutput = truncate(autograder.Normalize(res.Stdout), maxOutputPreview)
		}

		switch {
		case res.StatusID == docker.StatusAccepted:
			if autograder.Match(res.Stdout, tc.ExpectedOutput) {
				tr.Passed = true
				tr.Status = "OK"
			} else {
				tr.Status = string(models.StatusWrongAnswer)
			}
		case isRuntimeFailure(res.StatusID):
			tr.Status = string(models.StatusRuntimeError)
			if !tc.IsHidden {
				tr.Message = truncate(firstNonEmpty(res.Stderr, res.Message, res.StatusDescription), maxOutputPreview)
			} else {
				tr.Message = res.StatusDescription
			}
		default:
			return nil, fmt.Errorf("%w: unexpected status %d (%s)", docker.ErrUnavailable, res.StatusID, res.StatusDescription)
		}
		results = append(results, tr)
	}

	return results, nil
}

func (s *SubmissionService) save(ctx context.Context, sub *models.Submission, results []models.TestResult) (*models.Submission, error) {
	if results != nil && len(results) > 0 {
		if b, err := json.Marshal(results); err == nil {
			sub.ResultsJson = string(b)
		}
	} else {
		// MySQL JSON column rejects empty string, store empty array instead
		sub.ResultsJson = "[]"
	}
	if err := s.submissions.Create(ctx, sub); err != nil {
		return nil, err
	}
	sub.Results = results
	return sub, nil
}

// Get mengembalikan submission milik requester. Admin boleh melihat semua.
// Submission milik orang lain dilaporkan sebagai ErrNotFound agar ID tidak bisa dienumerasi.
func (s *SubmissionService) Get(ctx context.Context, id, requesterID uint, role models.Role) (*models.Submission, error) {
	sub, err := s.submissions.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sub.UserID != requesterID && role != models.RoleAdmin {
		return nil, ErrNotFound
	}
	return sub, nil
}

func (s *SubmissionService) ListMine(ctx context.Context, userID, quizID uint) ([]models.Submission, error) {
	return s.submissions.ListByUser(ctx, userID, quizID)
}

// isRuntimeFailure: TLE (5), runtime error (7-12), exec format error (14).
func isRuntimeFailure(statusID int) bool {
	return statusID == docker.StatusTimeLimit || (statusID >= 7 && statusID <= 12) || statusID == 14
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return strings.ToValidUTF8(s[:max], "") + "..."
}

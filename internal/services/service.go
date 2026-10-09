package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"autograder/internal/config"
	"autograder/internal/models"
	"autograder/internal/repositories"
	"autograder/internal/services/ai"
	"autograder/internal/services/autograder"
	"autograder/internal/services/judge0"
	"autograder/internal/utils"

	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email is already registered")
	ErrNotFound           = errors.New("resource not found")
	ErrNoTestCases        = errors.New("this quiz has no test cases yet")
)

const (
	AccessTTL  = 15 * time.Minute
	RefreshTTL = 7 * 24 * time.Hour
)

type Service struct {
	repo  *repositories.Repo
	judge *judge0.Client
	ai    *ai.Client
	cfg   *config.Config
}

func New(r *repositories.Repo, j *judge0.Client, a *ai.Client, c *config.Config) *Service {
	return &Service{repo: r, judge: j, ai: a, cfg: c}
}

func notFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}

// ---- auth ----

func (s *Service) Register(name, email, password string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if _, err := s.repo.UserByEmail(email); err == nil {
		return nil, ErrEmailTaken
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	u := &models.User{Name: name, Email: email, Password: hash, Role: models.RoleStudent}
	return u, s.repo.CreateUser(u)
}

func (s *Service) issue(u *models.User) (access, refresh string, err error) {
	if access, err = utils.SignToken(u.ID, u.Role, s.cfg.JWTSecret, AccessTTL); err != nil {
		return
	}
	refresh, err = utils.SignToken(u.ID, u.Role, s.cfg.JWTRefreshSecret, RefreshTTL)
	return
}

func (s *Service) Login(email, password string) (*models.User, string, string, error) {
	u, err := s.repo.UserByEmail(strings.ToLower(strings.TrimSpace(email)))
	if err != nil || !utils.CheckPassword(u.Password, password) {
		return nil, "", "", ErrInvalidCredentials
	}
	a, r, err := s.issue(u)
	return u, a, r, err
}

func (s *Service) Refresh(refreshToken string) (*models.User, string, error) {
	claims, err := utils.ParseToken(refreshToken, s.cfg.JWTRefreshSecret)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}
	u, err := s.repo.UserByID(claims.UserID)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}
	a, err := utils.SignToken(u.ID, u.Role, s.cfg.JWTSecret, AccessTTL)
	return u, a, err
}

// ---- courses & quizzes ----

func (s *Service) Courses() ([]models.Course, error)       { return s.repo.ListCourses() }
func (s *Service) CreateCourse(c *models.Course) error      { return s.repo.CreateCourse(c) }
func (s *Service) Quizzes() ([]models.Quiz, error)          { return s.repo.ListQuizzes() }
func (s *Service) CreateQuiz(q *models.Quiz) error          { return s.repo.CreateQuiz(q) }
func (s *Service) DeleteQuiz(id uint) error                 { return notFound(s.repo.DeleteQuiz(id)) }
func (s *Service) Quiz(id uint) (*models.Quiz, error) {
	q, err := s.repo.QuizByID(id)
	return q, notFound(err)
}

func (s *Service) GenerateTestCases(ctx context.Context, quizID uint, n int) ([]models.TestCase, error) {
	q, err := s.repo.QuizByID(quizID)
	if err != nil {
		return nil, notFound(err)
	}
	generated, err := s.ai.GenerateTestCases(ctx, q.Statement, n)
	if err != nil {
		return nil, err
	}
	tcs := make([]models.TestCase, len(generated))
	for i, g := range generated {
		tcs[i] = models.TestCase{QuizID: quizID, Stdin: g.Stdin, ExpectedOutput: g.ExpectedOutput}
	}
	return tcs, s.repo.CreateTestCases(tcs)
}

// ---- grading ----

func (s *Service) Grade(ctx context.Context, userID, quizID uint, code string) (*models.Submission, error) {
	q, err := s.repo.QuizByID(quizID)
	if err != nil {
		return nil, notFound(err)
	}
	if len(q.TestCases) == 0 {
		return nil, ErrNoTestCases
	}
	sub := &models.Submission{UserID: userID, QuizID: quizID, SourceCode: code, Status: models.StatusPassed}
	passed := 0
	for _, tc := range q.TestCases {
		res, err := s.judge.Run(ctx, code, tc.Stdin)
		if err != nil {
			return nil, err
		}
		if res.StatusID == judge0.StatusCompileError {
			sub.Status, passed = models.StatusCompileError, 0
			break
		}
		if res.StatusID == 3 && autograder.Match(res.Stdout, tc.ExpectedOutput) {
			passed++
		} else {
			sub.Status = models.StatusWrongAnswer
		}
	}
	if sub.Status != models.StatusCompileError {
		sub.Score = float64(passed) / float64(len(q.TestCases)) * 100
	}
	return sub, s.repo.CreateSubmission(sub)
}

func (s *Service) MySubmissions(uid uint) ([]models.Submission, error) {
	return s.repo.SubmissionsByUser(uid)
}

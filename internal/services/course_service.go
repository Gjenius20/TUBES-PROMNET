package services

import (
	"context"
	"fmt"
	"io"

	"autograder/internal/models"
	"autograder/internal/repositories"

	"github.com/xuri/excelize/v2"
)

func (s *CourseService) List(ctx context.Context) ([]models.Course, error) {
	return s.courses.List(ctx)
}

// ListByStudentID mengembalikan kursus yang sudah di-enroll APPROVED oleh student.
func (s *CourseService) ListByStudentID(ctx context.Context, studentID uint) ([]models.Course, error) {
	return s.courses.FindByStudentID(ctx, studentID)
}

func (s *CourseService) Get(ctx context.Context, id uint) (*models.Course, error) {
	return s.courses.FindByID(ctx, id)
}

func (s *CourseService) Create(ctx context.Context, title, description string) (*models.Course, error) {
	course := &models.Course{Title: title, Description: description}
	if err := s.courses.Create(ctx, course); err != nil {
		return nil, err
	}
	return course, nil
}

func (s *CourseService) Update(ctx context.Context, id uint, title, description string) (*models.Course, error) {
	fields := map[string]interface{}{"title": title, "description": description}
	if err := s.courses.Update(ctx, id, fields); err != nil {
		return nil, err
	}
	return s.courses.FindByID(ctx, id)
}

func (s *CourseService) Delete(ctx context.Context, id uint) error {
	return s.courses.Delete(ctx, id)
}

func (s *CourseService) ExportCourseProgressToExcel(ctx context.Context, courseID uint, writer io.Writer) error {
	progresses, err := s.GetStudentProgress(ctx, courseID)
	if err != nil {
		return err
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Report"
	f.NewSheet(sheet)

	// Headers
	headers := []string{"NIM", "Nama", "Email", "Status", "Progress (%)", "Rata-rata Skor"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Data
	for i, p := range progresses {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), p.NIM)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), p.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), p.Email)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), p.Status)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), p.ProgressPct)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), p.AvgScore)
	}

	// Remove default sheet
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(0)

	return f.Write(writer)
}

type CourseService struct {
	courses     repositories.CourseRepository
	users       repositories.UserRepository
	submissions repositories.SubmissionRepository
	enrollments repositories.EnrollmentRepository
}

func NewCourseService(
	courses repositories.CourseRepository,
	users repositories.UserRepository,
	submissions repositories.SubmissionRepository,
	enrollments repositories.EnrollmentRepository,
) *CourseService {
	return &CourseService{courses: courses, users: users, submissions: submissions, enrollments: enrollments}
}

type StudentProgress struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	NIM           string  `json:"nim"`
	Email         string  `json:"email"`
	Status        string  `json:"status"`
	ProgressPct   float64 `json:"progress_pct"`
	AvgScore      float64 `json:"avg_score"`
	TotalQuizzes  int     `json:"total_quizzes"`
	SolvedQuizzes int     `json:"solved_quizzes"`
}

func (s *CourseService) GetStudentProgress(ctx context.Context, courseID uint) ([]StudentProgress, error) {
	course, err := s.courses.FindByID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	totalQuizzes := len(course.Quizzes)
	if totalQuizzes == 0 {
		return []StudentProgress{}, nil
	}

	// PERBAIKAN: Ambil enrollment untuk kursus ini saja, preload User
	enrollments, err := s.enrollments.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, err
	}

	var progresses []StudentProgress
	for _, enr := range enrollments {
		stu := enr.User
		submissions, err := s.submissions.ListByCourse(ctx, stu.ID, courseID)
		if err != nil {
			return nil, err
		}

		// Map quiz_id -> best status (PASSED if any submission passed)
		bestStatus := make(map[uint]models.SubmissionStatus)
		var totalScore float64
		var scoredCount int

		for _, sub := range submissions {
			if sub.Status == models.StatusPassed {
				bestStatus[sub.QuizID] = models.StatusPassed
			} else if _, exists := bestStatus[sub.QuizID]; !exists {
				bestStatus[sub.QuizID] = sub.Status
			}

			// Calculate average score from passed submissions
			if sub.Status == models.StatusPassed && sub.TotalCount > 0 {
				score := float64(sub.PassedCount) / float64(sub.TotalCount) * 100
				totalScore += score
				scoredCount++
			}
		}

		solvedCount := 0
		for _, status := range bestStatus {
			if status == models.StatusPassed {
				solvedCount++
			}
		}

		progressPct := 0.0
		if totalQuizzes > 0 {
			progressPct = float64(solvedCount) / float64(totalQuizzes) * 100
		}

		status := "BELUM_MULAI"
		if solvedCount > 0 {
			if progressPct >= 100 {
				status = "SELESAI"
			} else {
				status = "SEDANG_BERJALAN"
			}
		}

		avgScore := 0.0
		if scoredCount > 0 {
			avgScore = totalScore / float64(scoredCount)
		}

		progresses = append(progresses, StudentProgress{
			ID:            stu.ID,
			Name:          stu.Name,
			NIM:           stu.NIM,
			Email:         stu.Email,
			Status:        status,
			ProgressPct:   progressPct,
			AvgScore:      avgScore,
			TotalQuizzes:  totalQuizzes,
			SolvedQuizzes: solvedCount,
		})
	}

	return progresses, nil
}

func (s *CourseService) ExportCourseGradesToExcel(ctx context.Context, courseID uint, writer io.Writer) error {
	course, err := s.courses.FindByID(ctx, courseID)
	if err != nil {
		return err
	}

	// PERBAIKAN: Ambil enrollment untuk kursus ini saja, preload User
	enrollments, err := s.enrollments.ListByCourse(ctx, courseID)
	if err != nil {
		return err
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Grades"
	f.NewSheet(sheet)

	// Headers
	headers := []string{"NIM", "Nama", "Email"}
	for _, quiz := range course.Quizzes {
		headers = append(headers, quiz.Title)
	}
	headers = append(headers, "Total Skor", "Max Skor", "Persentase")

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Data
	for i, enr := range enrollments {
		stu := enr.User
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), stu.NIM)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), stu.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), stu.Email)

		totalScore := 0.0
		totalMaxScore := 0

		for j, quiz := range course.Quizzes {
			col := j + 4 // Start from column D (index 4)
			cell, _ := excelize.CoordinatesToCellName(col, row)

			submissions, err := s.submissions.ListByCourse(ctx, stu.ID, courseID)
			if err != nil {
				return err
			}

			// Find best score for this quiz
			var bestScore int
			for _, sub := range submissions {
				if sub.QuizID == quiz.ID && sub.Status == models.StatusPassed {
					if sub.Score > bestScore {
						bestScore = sub.Score
					}
				}
			}

			scorePct := 0.0
			if quiz.MaxScore > 0 {
				scorePct = float64(bestScore) / float64(quiz.MaxScore) * 100
			}

			f.SetCellValue(sheet, cell, scorePct)
			totalScore += float64(bestScore)
			totalMaxScore += quiz.MaxScore
		}

		// Total skor
		totalCell, _ := excelize.CoordinatesToCellName(len(course.Quizzes)+4, row)
		f.SetCellValue(sheet, totalCell, totalScore)

		// Max skor
		maxCell, _ := excelize.CoordinatesToCellName(len(course.Quizzes)+5, row)
		f.SetCellValue(sheet, maxCell, totalMaxScore)

		// Persentase
		pctCell, _ := excelize.CoordinatesToCellName(len(course.Quizzes)+6, row)
		pct := 0.0
		if totalMaxScore > 0 {
			pct = totalScore / float64(totalMaxScore) * 100
		}
		f.SetCellValue(sheet, pctCell, pct)
	}

	// Remove default sheet
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(0)

	return f.Write(writer)
}
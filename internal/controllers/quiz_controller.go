package controllers

import (
	"net/http"

	"autograder/internal/services"
	"autograder/internal/utils"

	"github.com/gin-gonic/gin"
)

type QuizController struct{ quizzes *services.QuizService }

func NewQuizController(quizzes *services.QuizService) *QuizController {
	return &QuizController{quizzes: quizzes}
}

// Get (siswa): hanya menampilkan testcase non-hidden sebagai contoh.
func (h *QuizController) Get(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	quiz, err := h.quizzes.Get(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}
	samples := make([]SampleTestCase, 0)
	for _, tc := range quiz.TestCases {
		if !tc.IsHidden {
			samples = append(samples, SampleTestCase{Stdin: tc.Stdin, ExpectedOutput: tc.ExpectedOutput})
		}
	}
	utils.Success(c, http.StatusOK, "quiz retrieved", gin.H{"quiz": StudentQuizResponse{
		ID:              quiz.ID,
		CourseID:        quiz.CourseID,
		Title:           quiz.Title,
		Description:     quiz.Description,
		Constraints:     quiz.Constraints,
		SampleTestCases: samples,
	}})
}

// AdminGet: quiz lengkap dengan semua testcase.
func (h *QuizController) AdminGet(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	quiz, err := h.quizzes.Get(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "quiz retrieved", gin.H{"quiz": quiz})
}

func (h *QuizController) Create(c *gin.Context) {
	var req QuizCreateRequest
	if !bindJSON(c, &req) {
		return
	}
	quiz, err := h.quizzes.Create(c.Request.Context(), req.CourseID, req.Title, req.Description, req.Constraints)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "quiz created", gin.H{"quiz": quiz})
}

func (h *QuizController) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req QuizUpdateRequest
	if !bindJSON(c, &req) {
		return
	}
	quiz, err := h.quizzes.Update(c.Request.Context(), id, req.Title, req.Description, req.Constraints)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "quiz updated", gin.H{"quiz": quiz})
}

func (h *QuizController) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.quizzes.Delete(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "quiz deleted", nil)
}

func (h *QuizController) AddTestCase(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req TestCaseRequest
	if !bindJSON(c, &req) {
		return
	}
	tc, err := h.quizzes.AddTestCase(c.Request.Context(), id, req.Stdin, req.ExpectedOutput, req.IsHidden)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "test case created", gin.H{"test_case": tc})
}

func (h *QuizController) DeleteTestCase(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.quizzes.DeleteTestCase(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "test case deleted", nil)
}

func (h *QuizController) GenerateTestCases(c *gin.Context) {
	var req GenerateTestCasesRequest
	if !bindJSON(c, &req) {
		return
	}
	cases, err := h.quizzes.GenerateTestCases(c.Request.Context(), req.QuizID, req.ProblemStatement, req.Count)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "test cases generated", gin.H{"test_cases": cases})
}

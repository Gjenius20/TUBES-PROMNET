package controllers

import (
	"net/http"
	"strconv"

	"autograder/internal/middleware"
	"autograder/internal/services"
	"autograder/internal/utils"

	"github.com/gin-gonic/gin"
)

type SubmissionController struct{ submissions *services.SubmissionService }

func NewSubmissionController(submissions *services.SubmissionService) *SubmissionController {
	return &SubmissionController{submissions: submissions}
}

func (h *SubmissionController) Create(c *gin.Context) {
	var req SubmissionRequest
	if !bindJSON(c, &req) {
		return
	}
	sub, err := h.submissions.Submit(c.Request.Context(), middleware.CurrentUserID(c), req.QuizID, req.SourceCode)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "submission evaluated", gin.H{"submission": sub})
}

func (h *SubmissionController) List(c *gin.Context) {
	var quizID uint
	if raw := c.Query("quiz_id"); raw != "" {
		parsed, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || parsed == 0 {
			utils.Error(c, http.StatusBadRequest, "invalid quiz_id")
			return
		}
		quizID = uint(parsed)
	}
	list, err := h.submissions.ListMine(c.Request.Context(), middleware.CurrentUserID(c), quizID)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "submissions retrieved", gin.H{"submissions": list})
}

func (h *SubmissionController) Get(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	sub, err := h.submissions.Get(c.Request.Context(), id, middleware.CurrentUserID(c), middleware.CurrentRole(c))
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "submission retrieved", gin.H{"submission": sub})
}

package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"autograder/internal/models"
	"autograder/internal/services"
	"autograder/internal/utils"

	"github.com/gin-gonic/gin"
)

// RequireEnrollment memverifikasi user sudah APPROVED di course dengan :id = course_id.
// Digunakan untuk route /courses/:id.
func RequireEnrollment(enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil || courseID == 0 {
			utils.Error(c, http.StatusBadRequest, "invalid course id")
			c.Abort()
			return
		}
		if !ensureEnrolled(c, enrollmentService, uint(courseID)) {
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireQuizEnrollment memverifikasi enrollment untuk route dengan :id = quiz_id.
// Course diambil dari quiz (bukan dari input client) agar tidak bisa dipalsukan.
func RequireQuizEnrollment(enrollmentService *services.EnrollmentService, quizService *services.QuizService) gin.HandlerFunc {
	return func(c *gin.Context) {
		quizID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil || quizID == 0 {
			utils.Error(c, http.StatusBadRequest, "invalid quiz id")
			c.Abort()
			return
		}
		quiz, err := quizService.Get(c.Request.Context(), uint(quizID))
		if err != nil {
			utils.Error(c, http.StatusNotFound, "quiz not found")
			c.Abort()
			return
		}
		if !ensureEnrolled(c, enrollmentService, quiz.CourseID) {
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireEnrollmentForSubmission memverifikasi enrollment untuk POST /submissions,
// yang hanya membawa quiz_id di body. Body dibaca lalu dipulihkan agar handler
// berikutnya tetap dapat membaca payload.
func RequireEnrollmentForSubmission(enrollmentService *services.EnrollmentService, quizService *services.QuizService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "failed to read request body")
			c.Abort()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var req struct {
			QuizID uint `json:"quiz_id"`
		}
		if err := json.Unmarshal(bodyBytes, &req); err != nil || req.QuizID == 0 {
			utils.Error(c, http.StatusBadRequest, "quiz_id is required")
			c.Abort()
			return
		}

		quiz, err := quizService.Get(c.Request.Context(), req.QuizID)
		if err != nil {
			utils.Error(c, http.StatusNotFound, "quiz not found")
			c.Abort()
			return
		}
		if !ensureEnrolled(c, enrollmentService, quiz.CourseID) {
			c.Abort()
			return
		}
		c.Next()
	}
}

// ensureEnrolled mengembalikan true bila requester adalah admin atau sudah APPROVED
// di course. Bila tidak, response error sudah ditulis ke context.
func ensureEnrolled(c *gin.Context, enrollmentService *services.EnrollmentService, courseID uint) bool {
	if CurrentRole(c) == models.RoleAdmin {
		return true
	}
	approved, err := enrollmentService.CheckAccess(c.Request.Context(), CurrentUserID(c), courseID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to check enrollment")
		return false
	}
	if !approved {
		utils.Error(c, http.StatusForbidden, "you are not enrolled or not approved in this course")
		return false
	}
	return true
}

package controllers

import (
	"fmt"
	"net/http"

	"autograder/internal/middleware"
	"autograder/internal/services"
	"autograder/internal/utils"

	"github.com/gin-gonic/gin"
)

type EnrollmentController struct {
	enrollments *services.EnrollmentService
}

func NewEnrollmentController(enrollments *services.EnrollmentService) *EnrollmentController {
	return &EnrollmentController{enrollments: enrollments}
}

type BulkEnrollRequest struct {
	StudentIDs []uint `json:"student_ids" binding:"required,min=1"`
}

// MyEnrollments mengembalikan enrollment milik user yang login (untuk student melihat kursus yang diikutinya)
func (h *EnrollmentController) MyEnrollments(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

	enrollments, err := h.enrollments.GetMyEnrollments(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "enrollments retrieved", gin.H{"enrollments": enrollments})
}

// ListByCourse admin: list enrollment untuk satu course (untuk panel progress)
func (h *EnrollmentController) ListByCourse(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	enrollments, err := h.enrollments.ListEnrollmentsByCourse(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "enrollments retrieved", gin.H{"enrollments": enrollments})
}

// ListStudents admin: list semua siswa (role STUDENT) dengan search & pagination
func (h *EnrollmentController) ListStudents(c *gin.Context) {
	search := c.Query("search")
	page := 1
	limit := 10
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	users, total, err := h.enrollments.ListStudentsPaginated(c.Request.Context(), search, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "students retrieved", gin.H{
		"users": users,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// BulkEnroll admin: enroll multiple students to a course (Direct Assignment)
func (h *EnrollmentController) BulkEnroll(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var req BulkEnrollRequest
	if !bindJSON(c, &req) {
		return
	}

	if err := h.enrollments.BulkEnrollStudents(c.Request.Context(), id, req.StudentIDs); err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "students enrolled successfully", nil)
}

// RemoveStudent admin: remove a student from a course
func (h *EnrollmentController) RemoveStudent(c *gin.Context) {
	courseID, ok := parseID(c, "id")
	if !ok {
		return
	}

	studentIDStr := c.Param("student_id")
	var studentID uint
	fmt.Sscanf(studentIDStr, "%d", &studentID)

	if err := h.enrollments.RemoveStudent(c.Request.Context(), courseID, studentID); err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "student removed from course", nil)
}
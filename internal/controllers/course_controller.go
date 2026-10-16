package controllers

import (
	"fmt"
	"net/http"

	"autograder/internal/middleware"
	"autograder/internal/models"
	"autograder/internal/services"
	"autograder/internal/utils"

	"github.com/gin-gonic/gin"
)

type CourseController struct {
	courses     *services.CourseService
	enrollments *services.EnrollmentService
}

func NewCourseController(courses *services.CourseService, enrollments *services.EnrollmentService) *CourseController {
	return &CourseController{courses: courses, enrollments: enrollments}
}

func (h *CourseController) List(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var (
		courses []models.Course
		err     error
	)
	if middleware.CurrentRole(c) == models.RoleAdmin {
		courses, err = h.courses.List(c.Request.Context())
	} else {
		// Siswa hanya boleh melihat kursus yang sudah di-enroll APPROVED.
		courses, err = h.courses.ListByStudentID(c.Request.Context(), userID)
	}
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "courses retrieved", gin.H{"courses": courses})
}

// ListMyCourses mengembalikan kursus yang sudah di-enroll APPROVED oleh student yang login.
func (h *CourseController) ListMyCourses(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	courses, err := h.courses.ListByStudentID(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "my courses retrieved", gin.H{"courses": courses})
}

func (h *CourseController) Get(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	// Jika role siswa, validasi enrollment
	if middleware.CurrentRole(c) == models.RoleStudent {
		enrolled, err := h.enrollments.CheckAccess(c.Request.Context(), middleware.CurrentUserID(c), id)
		if err != nil || !enrolled {
			utils.Error(c, http.StatusForbidden, "Anda tidak terdaftar di kursus ini")
			return
		}
	}

	course, err := h.courses.Get(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "course retrieved", gin.H{"course": course})
}

func (h *CourseController) Create(c *gin.Context) {
	var req CourseRequest
	if !bindJSON(c, &req) {
		return
	}
	course, err := h.courses.Create(c.Request.Context(), req.Title, req.Description)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "course created", gin.H{"course": course})
}

func (h *CourseController) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req CourseRequest
	if !bindJSON(c, &req) {
		return
	}
	course, err := h.courses.Update(c.Request.Context(), id, req.Title, req.Description)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "course updated", gin.H{"course": course})
}

func (h *CourseController) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.courses.Delete(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "course deleted", nil)
}

func (h *CourseController) GetProgress(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	progress, err := h.courses.GetStudentProgress(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "progress retrieved", gin.H{"progress": progress})
}

func (h *CourseController) ExportGradesExcel(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"course-grades-%d.xlsx\"", id))

	if err := h.courses.ExportCourseGradesToExcel(c.Request.Context(), id, c.Writer); err != nil {
		handleError(c, err)
	}
}

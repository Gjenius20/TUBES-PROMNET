package controllers

import (
	"fmt"
	"net/http"

	"autograder/internal/services"
	"autograder/internal/utils"

	"github.com/gin-gonic/gin"
)

type CourseController struct{ courses *services.CourseService }

func NewCourseController(courses *services.CourseService) *CourseController {
	return &CourseController{courses: courses}
}

func (h *CourseController) List(c *gin.Context) {
	courses, err := h.courses.List(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "courses retrieved", gin.H{"courses": courses})
}

func (h *CourseController) Get(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
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

func (h *CourseController) ExportExcel(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"course-report-%d.xlsx\"", id))

	if err := h.courses.ExportCourseProgressToExcel(c.Request.Context(), id, c.Writer); err != nil {
		handleError(c, err)
	}
}

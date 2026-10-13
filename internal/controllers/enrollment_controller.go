package controllers

import (
	"net/http"

	"autograder/internal/middleware"
	"autograder/internal/models"
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

type EnrollRequest struct {
	CourseID uint `uri:"id" binding:"required"`
}

type UpdateEnrollmentStatusRequest struct {
	ID     uint   `uri:"id" binding:"required"`
	Status string `json:"status" binding:"required,oneof=APPROVED REJECTED"`
}

func (h *EnrollmentController) Enroll(c *gin.Context) {
	var req EnrollRequest
	if !bindJSON(c, &req) {
		return
	}
	userID := middleware.CurrentUserID(c)

	enrollment, err := h.enrollments.EnrollRequest(c.Request.Context(), userID, req.CourseID)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "enrollment request submitted", gin.H{"enrollment": enrollment})
}

func (h *EnrollmentController) MyEnrollments(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

	enrollments, err := h.enrollments.GetMyEnrollments(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "enrollments retrieved", gin.H{"enrollments": enrollments})
}

func (h *EnrollmentController) ListAll(c *gin.Context) {
	status := c.Query("status")

	enrollments, err := h.enrollments.ListEnrollments(c.Request.Context(), status)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "enrollments retrieved", gin.H{"enrollments": enrollments})
}

func (h *EnrollmentController) UpdateStatus(c *gin.Context) {
	var req UpdateEnrollmentStatusRequest
	if !bindJSON(c, &req) {
		return
	}

	status := models.EnrollmentStatus(req.Status)
	enrollment, err := h.enrollments.UpdateStatus(c.Request.Context(), req.ID, status)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "enrollment status updated", gin.H{"enrollment": enrollment})
}
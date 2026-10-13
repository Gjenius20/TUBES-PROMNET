package middleware

import (
	"net/http"
	"strconv"

	"autograder/internal/models"
	"autograder/internal/services"
	"autograder/internal/utils"

	"github.com/gin-gonic/gin"
)

// RequireEnrollment memverifikasi user sudah approved di course tersebut.
func RequireEnrollment(enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := CurrentUserID(c)
		courseIDStr := c.Param("id")
		courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "invalid course id")
			c.Abort()
			return
		}

		// Admin bypass enrollment check
		if CurrentRole(c) == models.RoleAdmin {
			c.Next()
			return
		}

		approved, err := enrollmentService.CheckAccess(c.Request.Context(), userID, uint(courseID))
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to check enrollment")
			c.Abort()
			return
		}

		if !approved {
			utils.Error(c, http.StatusForbidden, "you are not enrolled or not approved in this course")
			c.Abort()
			return
		}

		c.Next()
	}
}
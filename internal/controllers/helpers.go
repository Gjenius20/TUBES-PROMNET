package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"autograder/internal/services"
	"autograder/internal/services/ai"
	"autograder/internal/services/docker"
	"autograder/internal/utils"

	"github.com/gin-gonic/gin"
)

// parseID membaca path param bertipe uint positif; bila salah, langsung membalas 400.
func parseID(c *gin.Context, name string) (uint, bool) {
	raw := c.Param(name)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		utils.Error(c, http.StatusBadRequest, "invalid "+name)
		return 0, false
	}
	return uint(id), true
}

// bindJSON memvalidasi body request; bila gagal, langsung membalas 400.
func bindJSON(c *gin.Context, dst interface{}) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		utils.Error(c, http.StatusBadRequest, "validation failed", utils.BindErrors(err)...)
		return false
	}
	return true
}

// handleError memetakan error service ke status HTTP. Error tak dikenal dicatat di log
// dan dibalas generik agar detail internal tidak bocor.
func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		utils.Error(c, http.StatusNotFound, "resource not found")
	case errors.Is(err, services.ErrEmailTaken):
		utils.Error(c, http.StatusConflict, "email already registered")
	case errors.Is(err, services.ErrInvalidCredentials):
		utils.Error(c, http.StatusUnauthorized, "invalid email or password")
	case errors.Is(err, services.ErrInvalidToken):
		utils.Error(c, http.StatusUnauthorized, "invalid or expired token")
	case errors.Is(err, services.ErrNoTestCases):
		utils.Error(c, http.StatusUnprocessableEntity, "this quiz has no test cases yet")
	case errors.Is(err, docker.ErrUnavailable):
		log.Printf("docker execution error: %v", err)
		utils.Error(c, http.StatusServiceUnavailable, "code execution service is unavailable, please try again later")
	case errors.Is(err, ai.ErrInvalidOutput):
		log.Printf("ai invalid output: %v", err)
		utils.Error(c, http.StatusBadGateway, "AI returned invalid test cases, please try again")
	case errors.Is(err, ai.ErrUnavailable):
		log.Printf("ai error: %v", err)
		utils.Error(c, http.StatusBadGateway, "AI service is unavailable, please try again later")
	default:
		log.Printf("unhandled error: %v", err)
		utils.Error(c, http.StatusInternalServerError, "internal server error")
	}
}

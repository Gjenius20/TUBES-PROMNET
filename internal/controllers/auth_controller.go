package controllers

import (
	"net/http"

	"autograder/internal/middleware"
	"autograder/internal/models"
	"autograder/internal/services"
	"autograder/internal/utils"

	"github.com/gin-gonic/gin"
)

const (
	refreshCookieName = "refresh_token"
	refreshCookiePath = "/api/v1/auth"
)

type AuthController struct {
	auth         *services.AuthService
	tokens       *utils.TokenManager
	cookieSecure bool
}

func NewAuthController(auth *services.AuthService, tokens *utils.TokenManager, cookieSecure bool) *AuthController {
	return &AuthController{auth: auth, tokens: tokens, cookieSecure: cookieSecure}
}

func (a *AuthController) setRefreshCookie(c *gin.Context, token string, maxAge int) {
	c.SetSameSite(http.SameSiteLaxMode)
	// HttpOnly=true: tidak bisa dibaca JavaScript. Path dibatasi ke /api/v1/auth.
	c.SetCookie(refreshCookieName, token, maxAge, refreshCookiePath, "", a.cookieSecure, true)
}

func (a *AuthController) respondWithSession(c *gin.Context, status int, message string, user *models.User, tokens *services.AuthTokens) {
	a.setRefreshCookie(c, tokens.RefreshToken, int(a.tokens.RefreshTTL().Seconds()))
	utils.Success(c, status, message, gin.H{
		"user":         user,
		"access_token": tokens.AccessToken,
	})
}

func (a *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if !bindJSON(c, &req) {
		return
	}
	user, tokens, err := a.auth.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}
	a.respondWithSession(c, http.StatusCreated, "registration successful", user, tokens)
}

func (a *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if !bindJSON(c, &req) {
		return
	}
	user, tokens, err := a.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}
	a.respondWithSession(c, http.StatusOK, "login successful", user, tokens)
}

func (a *AuthController) Refresh(c *gin.Context) {
	token, err := c.Cookie(refreshCookieName)
	if err != nil || token == "" {
		utils.Error(c, http.StatusUnauthorized, "missing refresh token")
		return
	}
	user, tokens, err := a.auth.Refresh(c.Request.Context(), token)
	if err != nil {
		a.setRefreshCookie(c, "", -1)
		handleError(c, err)
		return
	}
	a.respondWithSession(c, http.StatusOK, "token refreshed", user, tokens)
}

func (a *AuthController) Logout(c *gin.Context) {
	a.setRefreshCookie(c, "", -1)
	utils.Success(c, http.StatusOK, "logged out", nil)
}

func (a *AuthController) Me(c *gin.Context) {
	user, err := a.auth.Me(c.Request.Context(), middleware.CurrentUserID(c))
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "current user", gin.H{"user": user})
}

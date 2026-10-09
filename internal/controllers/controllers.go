package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"autograder/internal/models"
	"autograder/internal/services"
	"autograder/internal/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct{ svc *services.Service }

func New(s *services.Service) *Handler { return &Handler{svc: s} }

const refreshCookie = "refresh_token"

func bad(c *gin.Context, err error) {
	utils.Fail(c, http.StatusBadRequest, "validation failed", err.Error())
}

func fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		utils.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrInvalidCredentials):
		utils.Fail(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, services.ErrEmailTaken):
		utils.Fail(c, http.StatusConflict, err.Error())
	case errors.Is(err, services.ErrNoTestCases):
		utils.Fail(c, http.StatusUnprocessableEntity, err.Error())
	default:
		utils.Fail(c, http.StatusBadGateway, "request could not be completed", err.Error())
	}
}

func idParam(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return uint(id), true
}

func (h *Handler) setRefresh(c *gin.Context, tok string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshCookie, tok, int(services.RefreshTTL.Seconds()), "/api/v1/auth", "", gin.Mode() == gin.ReleaseMode, true)
}

// ---- auth ----

func (h *Handler) Register(c *gin.Context) {
	var in struct {
		Name     string `json:"name" binding:"required,min=2"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		bad(c, err)
		return
	}
	u, err := h.svc.Register(in.Name, in.Email, in.Password)
	if err != nil {
		fail(c, err)
		return
	}
	utils.OK(c, http.StatusCreated, "registered", gin.H{"user": u})
}

func (h *Handler) Login(c *gin.Context) {
	var in struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		bad(c, err)
		return
	}
	u, access, refresh, err := h.svc.Login(in.Email, in.Password)
	if err != nil {
		fail(c, err)
		return
	}
	h.setRefresh(c, refresh)
	utils.OK(c, http.StatusOK, "logged in", gin.H{"user": u, "access_token": access})
}

func (h *Handler) Refresh(c *gin.Context) {
	tok, err := c.Cookie(refreshCookie)
	if err != nil {
		utils.Fail(c, http.StatusUnauthorized, "no refresh token")
		return
	}
	u, access, err := h.svc.Refresh(tok)
	if err != nil {
		fail(c, err)
		return
	}
	utils.OK(c, http.StatusOK, "refreshed", gin.H{"user": u, "access_token": access})
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie(refreshCookie, "", -1, "/api/v1/auth", "", false, true)
	utils.OK(c, http.StatusOK, "logged out", nil)
}

// ---- courses & quizzes ----

func (h *Handler) Courses(c *gin.Context) {
	cs, err := h.svc.Courses()
	if err != nil {
		fail(c, err)
		return
	}
	utils.OK(c, http.StatusOK, "ok", cs)
}

func (h *Handler) CreateCourse(c *gin.Context) {
	var in struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		bad(c, err)
		return
	}
	course := &models.Course{Title: in.Title, Description: in.Description}
	if err := h.svc.CreateCourse(course); err != nil {
		fail(c, err)
		return
	}
	utils.OK(c, http.StatusCreated, "course created", course)
}

func (h *Handler) Quizzes(c *gin.Context) {
	qs, err := h.svc.Quizzes()
	if err != nil {
		fail(c, err)
		return
	}
	utils.OK(c, http.StatusOK, "ok", qs)
}

func (h *Handler) Quiz(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	q, err := h.svc.Quiz(id)
	if err != nil {
		fail(c, err)
		return
	}
	if role, _ := c.Get("role"); role != models.RoleAdmin {
		q.TestCases = nil // students never see expected outputs
	}
	utils.OK(c, http.StatusOK, "ok", q)
}

func (h *Handler) CreateQuiz(c *gin.Context) {
	var in struct {
		CourseID  *uint  `json:"course_id"`
		Title     string `json:"title" binding:"required"`
		Statement string `json:"statement" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		bad(c, err)
		return
	}
	q := &models.Quiz{CourseID: in.CourseID, Title: in.Title, Statement: in.Statement}
	if err := h.svc.CreateQuiz(q); err != nil {
		fail(c, err)
		return
	}
	utils.OK(c, http.StatusCreated, "quiz created", q)
}

func (h *Handler) DeleteQuiz(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteQuiz(id); err != nil {
		fail(c, err)
		return
	}
	utils.OK(c, http.StatusOK, "quiz deleted", nil)
}

func (h *Handler) GenerateTestCases(c *gin.Context) {
	var in struct {
		QuizID uint `json:"quiz_id" binding:"required"`
		Count  int  `json:"count" binding:"omitempty,min=1,max=10"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		bad(c, err)
		return
	}
	if in.Count == 0 {
		in.Count = 5
	}
	tcs, err := h.svc.GenerateTestCases(c.Request.Context(), in.QuizID, in.Count)
	if err != nil {
		fail(c, err)
		return
	}
	utils.OK(c, http.StatusCreated, "test cases generated", tcs)
}

// ---- submissions ----

func (h *Handler) Submit(c *gin.Context) {
	var in struct {
		QuizID     uint   `json:"quiz_id" binding:"required"`
		SourceCode string `json:"source_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		bad(c, err)
		return
	}
	sub, err := h.svc.Grade(c.Request.Context(), c.GetUint("uid"), in.QuizID, in.SourceCode)
	if err != nil {
		fail(c, err)
		return
	}
	utils.OK(c, http.StatusCreated, "graded", sub)
}

func (h *Handler) MySubmissions(c *gin.Context) {
	ss, err := h.svc.MySubmissions(c.GetUint("uid"))
	if err != nil {
		fail(c, err)
		return
	}
	utils.OK(c, http.StatusOK, "ok", ss)
}

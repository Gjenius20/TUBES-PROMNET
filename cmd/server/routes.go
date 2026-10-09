package main

import (
	"time"

	"autograder/internal/config"
	"autograder/internal/controllers"
	"autograder/internal/middleware"
	"autograder/internal/models"
	"autograder/internal/utils"

	"github.com/gin-gonic/gin"
)

type handlers struct {
	auth        *controllers.AuthController
	course      *controllers.CourseController
	quiz        *controllers.QuizController
	submission  *controllers.SubmissionController
	tokenManager *utils.TokenManager
}

func newRouter(cfg *config.Config, h handlers) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), middleware.Recovery(), middleware.CORS(cfg.CORSAllowedOrigins))

	api := r.Group("/api/v1")

	// Rute publik (satu-satunya yang tidak memerlukan JWT): login, register, refresh.
	public := api.Group("/auth")
	public.Use(middleware.RateLimit(20, time.Minute, middleware.ByIP))
	public.POST("/register", h.auth.Register)
	public.POST("/login", h.auth.Login)
	public.POST("/refresh", h.auth.Refresh)

	// Semua rute lain wajib access token yang valid.
	protected := api.Group("")
	protected.Use(middleware.Auth(h.tokenManager))
	{
		protected.POST("/auth/logout", h.auth.Logout)
		protected.GET("/auth/me", h.auth.Me)

		protected.GET("/courses", h.course.List)
		protected.GET("/courses/:id", h.course.Get)
		protected.GET("/quizzes/:id", h.quiz.Get)

		protected.POST("/submissions", middleware.RateLimit(10, time.Minute, middleware.ByUser), h.submission.Create)
		protected.GET("/submissions", h.submission.List)
		protected.GET("/submissions/:id", h.submission.Get)
	}

	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole(models.RoleAdmin))
	{
		admin.POST("/courses", h.course.Create)
		admin.PUT("/courses/:id", h.course.Update)
		admin.DELETE("/courses/:id", h.course.Delete)

		admin.POST("/quizzes", h.quiz.Create)
		admin.POST("/quizzes/generate-testcases", middleware.RateLimit(5, time.Minute, middleware.ByUser), h.quiz.GenerateTestCases)
		admin.GET("/quizzes/:id", h.quiz.AdminGet)
		admin.PUT("/quizzes/:id", h.quiz.Update)
		admin.DELETE("/quizzes/:id", h.quiz.Delete)
		admin.POST("/quizzes/:id/testcases", h.quiz.AddTestCase)

		admin.DELETE("/testcases/:id", h.quiz.DeleteTestCase)
	}

	return r
}

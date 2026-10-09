package main

import (
	"log"
	"time"

	"autograder/internal/config"
	"autograder/internal/controllers"
	"autograder/internal/middleware"
	"autograder/internal/models"
	"autograder/internal/repositories"
	"autograder/internal/services"
	"autograder/internal/services/ai"
	"autograder/internal/services/judge0"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Course{}, &models.Quiz{}, &models.TestCase{}, &models.Submission{}); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	svc := services.New(repositories.New(db), judge0.New(cfg.Judge0URL), ai.New(cfg.GeminiKey, cfg.GeminiModel), cfg)
	h := controllers.New(svc)

	r := gin.New()
	r.Use(gin.Logger(), middleware.Recovery(), cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := r.Group("/api/v1")
	v1.POST("/auth/register", h.Register)
	v1.POST("/auth/login", h.Login)
	v1.POST("/auth/refresh", h.Refresh)
	v1.POST("/auth/logout", h.Logout)

	p := v1.Group("", middleware.Auth(cfg.JWTSecret))
	p.GET("/courses", h.Courses)
	p.GET("/quizzes", h.Quizzes)
	p.GET("/quizzes/:id", h.Quiz)
	p.POST("/submissions", h.Submit)
	p.GET("/submissions/me", h.MySubmissions)

	admin := p.Group("/admin", middleware.RequireRole(models.RoleAdmin))
	admin.POST("/courses", h.CreateCourse)
	admin.POST("/quizzes", h.CreateQuiz)
	admin.DELETE("/quizzes/:id", h.DeleteQuiz)
	admin.POST("/quizzes/generate-testcases", h.GenerateTestCases)

	log.Fatal(r.Run(":" + cfg.Port))
}

package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"autograder/internal/config"
	"autograder/internal/controllers"
	"autograder/internal/repositories"
	"autograder/internal/services"
	"autograder/internal/services/ai"
	"autograder/internal/services/docker"
	"autograder/internal/utils"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	db, err := config.ConnectDatabase(cfg.DBDSN)
	if err != nil {
		log.Fatalf("%v", err)
	}

	tokenManager := utils.NewTokenManager(cfg.JWTSecret, cfg.JWTRefreshSecret)

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	quizRepo := repositories.NewQuizRepository(db)
	testCaseRepo := repositories.NewTestCaseRepository(db)
	submissionRepo := repositories.NewSubmissionRepository(db)

	// Services & integrasi eksternal
	generator := ai.NewGeminiGenerator(cfg.GeminiAPIKey, cfg.GeminiModel)
	runner := docker.NewClient("gcc:12")

	authService := services.NewAuthService(userRepo, tokenManager)
	courseService := services.NewCourseService(courseRepo)
	quizService := services.NewQuizService(quizRepo, testCaseRepo, courseRepo, generator)
	submissionService := services.NewSubmissionService(submissionRepo, quizRepo, testCaseRepo, runner)

	router := newRouter(cfg, handlers{
		auth:         controllers.NewAuthController(authService, tokenManager, cfg.CookieSecure),
		course:       controllers.NewCourseController(courseService),
		quiz:         controllers.NewQuizController(quizService),
		submission:   controllers.NewSubmissionController(submissionService),
		tokenManager: tokenManager,
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		// Submission bisa memakan waktu (banyak testcase x Judge0), beri ruang cukup.
		WriteTimeout: 180 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
}

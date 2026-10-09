package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port, DBDSN, JWTSecret, JWTRefreshSecret string
	GeminiKey, GeminiModel, Judge0URL        string
	FrontendOrigin                           string
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	c := &Config{
		Port:             get("PORT", "8080"),
		DBDSN:            os.Getenv("DB_DSN"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
		GeminiKey:        os.Getenv("GEMINI_API_KEY"),
		GeminiModel:      get("GEMINI_MODEL", "gemini-2.5-flash"),
		Judge0URL:        os.Getenv("JUDGE0_API_URL"),
		FrontendOrigin:   get("FRONTEND_ORIGIN", "http://localhost:5173"),
	}
	required := map[string]string{"DB_DSN": c.DBDSN, "JWT_SECRET": c.JWTSecret, "JWT_REFRESH_SECRET": c.JWTRefreshSecret, "GEMINI_API_KEY": c.GeminiKey, "JUDGE0_API_URL": c.Judge0URL}
	for k, v := range required {
		if v == "" {
			return nil, fmt.Errorf("missing required env var %s", k)
		}
	}
	return c, nil
}

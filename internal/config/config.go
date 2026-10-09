package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config menampung seluruh konfigurasi aplikasi yang dibaca dari environment.
type Config struct {
	Port               string
	DBDSN              string
	JWTSecret          string
	JWTRefreshSecret   string
	GeminiAPIKey       string
	GeminiModel        string
	Judge0APIURL       string
	CORSAllowedOrigins []string
	CookieSecure       bool
}

// Load membaca file .env (jika ada) lalu memvalidasi variabel yang wajib.
func Load() (*Config, error) {
	// .env bersifat opsional: di production variabel biasanya di-inject langsung.
	_ = godotenv.Load()

	cfg := &Config{
		Port:             getEnv("PORT", "8080"),
		DBDSN:            os.Getenv("DB_DSN"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
		GeminiAPIKey:     os.Getenv("GEMINI_API_KEY"),
		GeminiModel:      getEnv("GEMINI_MODEL", "gemini-2.5-flash"),
		Judge0APIURL:     os.Getenv("JUDGE0_API_URL"),
	}

	required := []struct{ key, value string }{
		{"DB_DSN", cfg.DBDSN},
		{"JWT_SECRET", cfg.JWTSecret},
		{"JWT_REFRESH_SECRET", cfg.JWTRefreshSecret},
		{"GEMINI_API_KEY", cfg.GeminiAPIKey},
		{"JUDGE0_API_URL", cfg.Judge0APIURL},
	}
	var missing []string
	for _, r := range required {
		if strings.TrimSpace(r.value) == "" {
			missing = append(missing, r.key)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	if cfg.JWTSecret == cfg.JWTRefreshSecret {
		return nil, errors.New("JWT_SECRET and JWT_REFRESH_SECRET must be different values")
	}

	origins := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")
	for _, o := range strings.Split(origins, ",") {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		if o == "*" {
			return nil, errors.New("CORS_ALLOWED_ORIGINS must not contain \"*\" because credentials (cookies) are enabled")
		}
		cfg.CORSAllowedOrigins = append(cfg.CORSAllowedOrigins, o)
	}
	if len(cfg.CORSAllowedOrigins) == 0 {
		return nil, errors.New("CORS_ALLOWED_ORIGINS must contain at least one origin")
	}

	if raw := os.Getenv("COOKIE_SECURE"); raw != "" {
		secure, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid COOKIE_SECURE value %q: %w", raw, err)
		}
		cfg.CookieSecure = secure
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

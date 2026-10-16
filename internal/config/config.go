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
	ExecDockerImage    string
	ExecTimeoutSeconds int
	ExecMemoryLimitMB  int
	ExecCPULimit       string
	CORSAllowedOrigins []string
	CookieSecure       bool
}

// Load membaca file .env (jika ada) lalu memvalidasi variabel yang wajib.
func Load() (*Config, error) {
	// .env bersifat opsional: di production variabel biasanya di-inject langsung.
	_ = godotenv.Load()

	cfg := &Config{
		Port:             getEnv("PORT", "1111"),
		DBDSN:            os.Getenv("DB_DSN"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
		GeminiAPIKey:     os.Getenv("GEMINI_API_KEY"),
		GeminiModel:      getEnv("GEMINI_MODEL", "gemini-2.5-flash"),
		ExecDockerImage:  getEnv("EXEC_DOCKER_IMAGE", "gcc:12"),
		ExecCPULimit:     getEnv("EXEC_CPU_LIMIT", "0.5"),
	}

	timeout, err := loadIntEnv("EXEC_TIMEOUT_SECONDS", 5)
	if err != nil {
		return nil, err
	}
	memory, err := loadIntEnv("EXEC_MEMORY_LIMIT_MB", 128)
	if err != nil {
		return nil, err
	}
	cfg.ExecTimeoutSeconds = timeout
	cfg.ExecMemoryLimitMB = memory

	if cfg.ExecTimeoutSeconds <= 0 {
		return nil, errors.New("EXEC_TIMEOUT_SECONDS must be a positive integer")
	}
	if cfg.ExecMemoryLimitMB <= 0 {
		return nil, errors.New("EXEC_MEMORY_LIMIT_MB must be a positive integer")
	}
	if _, err := strconv.ParseFloat(cfg.ExecCPULimit, 64); err != nil {
		return nil, fmt.Errorf("invalid EXEC_CPU_LIMIT value %q: must be a number", cfg.ExecCPULimit)
	}

	required := []struct{ key, value string }{
		{"DB_DSN", cfg.DBDSN},
		{"JWT_SECRET", cfg.JWTSecret},
		{"JWT_REFRESH_SECRET", cfg.JWTRefreshSecret},
		{"GEMINI_API_KEY", cfg.GeminiAPIKey},
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

// loadIntEnv membaca variabel integer opsional; bila kosong memakai fallback,
// bila tidak valid mengembalikan error.
func loadIntEnv(key string, fallback int) (int, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s value %q: must be an integer", key, raw)
	}
	return n, nil
}

package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// Config는 애플리케이션 전체 설정입니다. 모든 값은 환경변수(.env 포함)에서 읽습니다.
type Config struct {
	Port               string
	DB                 DBConfig
	JWTSecret          string
	JWTTTL             time.Duration
	CORSAllowedOrigins []string
	RateLimitRPS       float64
	RateLimitBurst     int
	QueryTimeout       time.Duration
	ShutdownTimeout    time.Duration
}

// Load는 .env 파일(있으면)과 환경변수에서 설정을 읽고 검증합니다.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("warning: .env 파일을 읽지 못했습니다: %v", err)
	}

	var errs []error
	cfg := &Config{
		Port: getEnv("PORT", "8080"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "gdgdb"),
		},
		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTTTL:             getDuration("JWT_TTL", 24*time.Hour, &errs),
		CORSAllowedOrigins: splitList(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:8080")),
		RateLimitRPS:       getFloat("RATE_LIMIT_RPS", 10, &errs),
		RateLimitBurst:     getInt("RATE_LIMIT_BURST", 20, &errs),
		QueryTimeout:       getDuration("DB_QUERY_TIMEOUT", 5*time.Second, &errs),
		ShutdownTimeout:    getDuration("SHUTDOWN_TIMEOUT", 10*time.Second, &errs),
	}

	if len(cfg.JWTSecret) < 16 {
		errs = append(errs, errors.New("JWT_SECRET 은 16자 이상이어야 합니다"))
	}
	return cfg, errors.Join(errs...)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getDuration(key string, fallback time.Duration, errs *[]error) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s: %w", key, err))
		return fallback
	}
	return d
}

func getInt(key string, fallback int, errs *[]error) int {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s: %w", key, err))
		return fallback
	}
	return n
}

func getFloat(key string, fallback float64, errs *[]error) float64 {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s: %w", key, err))
		return fallback
	}
	return f
}

func splitList(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	GitHub   GitHubConfig
	Weather  WeatherConfig
	FCM      FCMConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Env  string
	Port string
}

type DatabaseConfig struct {
	URL      string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type GitHubConfig struct {
	Token  string
	Owner  string
	Repo   string
	Branch string
}

type WeatherConfig struct {
	APIKey string
	APIURL string
}

type FCMConfig struct {
	ServerKey string
}

type CORSConfig struct {
	AllowedOrigins string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	accessTTL, err := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"))
	if err != nil {
		accessTTL = 15 * time.Minute
	}
	refreshTTL, err := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))
	if err != nil {
		refreshTTL = 168 * time.Hour
	}

	return &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("PORT", "5005"),
		},
		Database: DatabaseConfig{
			URL:      getEnv("DATABASE_URL", ""),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "sitelogix"),
			User:     getEnv("DB_USER", "sitelogix"),
			Password: getEnv("DB_PASSWORD", "secret"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "change-me"),
			AccessTokenTTL:  accessTTL,
			RefreshTokenTTL: refreshTTL,
		},
		GitHub: GitHubConfig{
			Token:  getEnv("GITHUB_TOKEN", ""),
			Owner:  getEnv("GITHUB_OWNER", ""),
			Repo:   getEnv("GITHUB_REPO", "sitelogix-media"),
			Branch: getEnv("GITHUB_BRANCH", "main"),
		},
		Weather: WeatherConfig{
			APIKey: getEnv("WEATHER_API_KEY", ""),
			APIURL: getEnv("WEATHER_API_URL", "https://api.openweathermap.org/data/2.5/weather"),
		},
		FCM: FCMConfig{
			ServerKey: getEnv("FCM_SERVER_KEY", ""),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

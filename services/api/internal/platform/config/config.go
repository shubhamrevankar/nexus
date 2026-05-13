package config

import "os"

type Config struct {
	Port            string
	LogLevel        string
	DatabaseURL     string
	MigrationsDir   string
	AllowedOrigins  string
	GitHubAPIURL    string
	SessionTTLHours int
}

func Load() Config {
	return Config{
		Port:            envOrDefault("API_PORT", "8080"),
		LogLevel:        envOrDefault("API_LOG_LEVEL", "info"),
		DatabaseURL:     envOrDefault("DATABASE_URL", "postgres://nexus:nexus_dev_password@localhost:5432/nexus?sslmode=disable"),
		MigrationsDir:   envOrDefault("API_MIGRATIONS_DIR", "migrations"),
		AllowedOrigins:  envOrDefault("API_ALLOWED_ORIGINS", "http://localhost:3000"),
		GitHubAPIURL:    envOrDefault("GITHUB_API_URL", "https://api.github.com"),
		SessionTTLHours: 24 * 7,
	}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

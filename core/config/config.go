package config

import "os"

type Config struct {
	Environment  string
	Host         string
	Port         string
	LogLevel     string
	DatabasePath string
}

func Load() Config {
	return Config{
		Environment: getEnv(
			"FJK_ENV",
			"development",
		),
		Host: getEnv(
			"FJK_CORE_HOST",
			"127.0.0.1",
		),
		Port: getEnv(
			"FJK_CORE_PORT",
			"17890",
		),
		LogLevel: getEnv(
			"FJK_LOG_LEVEL",
			"debug",
		),
		DatabasePath: getEnv(
			"FJK_DATABASE_PATH",
			"../data/fjk-agent.db",
		),
	}
}

func getEnv(
	key string,
	fallback string,
) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

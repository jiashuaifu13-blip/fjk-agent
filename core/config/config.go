package config

import "os"

type Config struct {
	Environment  string
	Host         string
	Port         string
	LogLevel     string
	DatabasePath string

	AIProvider string
	AIModel    string

	DeepSeekAPIKey  string
	DeepSeekBaseURL string
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

		AIProvider: getEnv(
			"FJK_AI_PROVIDER",
			"deepseek",
		),
		AIModel: getEnv(
			"FJK_AI_MODEL",
			"deepseek-chat",
		),

		DeepSeekAPIKey: getEnv(
			"DEEPSEEK_API_KEY",
			"",
		),
		DeepSeekBaseURL: getEnv(
			"DEEPSEEK_BASE_URL",
			"https://api.deepseek.com",
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

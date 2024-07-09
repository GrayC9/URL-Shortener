package config

import (
	"os"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
}

type ServerConfig struct {
	Address string
}

type DBConfig struct {
	DSN string
}

func LoadConfig() Config {
	return Config{
		Server: ServerConfig{
			Address: getEnv("SERVER_ADDRESS", ":8080"),
		},
		DB: DBConfig{
			DSN: getEnv("DB_DSN", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

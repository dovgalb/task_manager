package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type DatabaseConfig struct {
	DbName        string
	DbUser        string
	DbPassword    string
	DbHost        string
	DbPort        string
	DbMaxAttempts int
}

type Config struct {
	DatabaseConfig
}

func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No env file")
	}
	return &Config{
		DatabaseConfig{
			DbName:        getEnv("POSTGRES_DB", "postgres_db"),
			DbUser:        getEnv("POSTGRES_USER", "postgres"),
			DbPassword:    getEnv("POSTGRES_PASSWORD", "postgres"),
			DbHost:        getEnv("POSTGRES_HOST", ""),
			DbPort:        getEnv("POSTGRES_PORT", ""),
			DbMaxAttempts: getEnvInt("POSTGRES_MAX_ATTEMPTS", 5),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Ошибка конвертации %s: %v. Используется значение по умолчанию: %d", key, err, defaultValue)
			return defaultValue
		}
		return intValue
	}
	return defaultValue
}

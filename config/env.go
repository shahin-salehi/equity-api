package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ConnectionString      string
	LocalConnectionString string
}

func InitConfig() Config {
	godotenv.Load()
	return Config{
		ConnectionString:      getEnv("pg_connection", ""),
		LocalConnectionString: getEnv("pg_local", ""),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

package config

import "os"

type Config struct {
	ConnectionString string
}

func InitConfig() Config {
	return Config{
		ConnectionString: getEnv("pg_connection", ""),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

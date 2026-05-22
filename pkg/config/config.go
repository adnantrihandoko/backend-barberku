package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPass        string
	DBName        string
	ServerPort    string
	WsOrigin      string
	JWTSecret     string
	FCMCredPath   string
	FCMProjectID  string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "barber"),
		DBPass:       getEnv("DB_PASS", "secret"),
		DBName:       getEnv("DB_NAME", "barbershop"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		WsOrigin:     getEnv("WS_ORIGIN_ALLOWED", "*"),
		JWTSecret:    getEnv("JWT_SECRET", "barberku-secret-key-change-in-production"),
		FCMCredPath:  getEnv("FCM_CREDENTIAL_PATH", "firebase-credentials.json"),
		FCMProjectID: getEnv("FCM_PROJECT_ID", ""),
	}
}

func (c *Config) DBConnectionString() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPass +
		" dbname=" + c.DBName +
		" sslmode=disable"
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

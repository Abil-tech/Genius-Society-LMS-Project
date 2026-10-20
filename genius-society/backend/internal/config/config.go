package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv         string
	Port           string
	MongoURI       string
	MongoDB        string
	JWTSecret      string
	CookieDomain   string
	CookieSameSite string
	CORSOrigin     string
	// Backblaze B2
	B2AccountID        string
	B2ApplicationKeyID string
	B2ApplicationKey   string
	B2BucketName       string
	B2S3Endpoint       string
	B2PublicURL        string
}

func Load() (*Config, error) {
	// Coba load .env di working directory saat ini atau di parent directory (jika dijalankan dari cmd/)
	_ = godotenv.Load(".env", "../.env")

	cfg := &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		Port:               getEnv("PORT", "8080"),
		MongoURI:           mustEnv("MONGODB_URI"),
		MongoDB:            getEnv("MONGODB_DATABASE", "genius_society"),
		JWTSecret:          mustEnv("JWT_SECRET"),
		CookieDomain:       getEnv("COOKIE_DOMAIN", ""),
		CookieSameSite:     getEnv("COOKIE_SAMESITE", "Lax"),
		CORSOrigin:         getEnv("CORS_ORIGIN", "http://localhost:5173"),
		B2AccountID:        getEnv("B2_ACCOUNT_ID", ""),
		B2ApplicationKeyID: getEnv("B2_APPLICATION_KEY_ID", ""),
		B2ApplicationKey:   getEnv("B2_APPLICATION_KEY", ""),
		B2BucketName:       getEnv("B2_BUCKET_NAME", ""),
		B2S3Endpoint:       getEnv("B2_S3_ENDPOINT", ""),
		B2PublicURL:        getEnv("B2_PUBLIC_URL", ""),
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required env var %s is not set", key))
	}
	return v
}

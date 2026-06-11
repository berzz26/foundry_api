package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	// "strconv"
)

type Config struct {
	AppEnv string

	HTTPPort string

	DatabaseURL string

	RedisURL string

	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
	HostURL        string

	GoogleClientID string
	GoogleSecret   string

	GithubClientID string
	GithubSecret   string

	GeminiAPI string
}

func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(".env file not found")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	port := os.Getenv("HTTP_PORT")
	// minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	// minioKey := os.Getenv("MINIO_ACCESS_KEY")
	// minioSecret := os.Getenv("MINIO_SECRET_KEY")
	// minioBucket := os.Getenv("MINIO_BUCKET")
	// minioUseSSL := os.Getenv("MINIO_USE_SSL")
	// hostURL := os.Getenv("HOST_URL")
	// googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	// googleSecret := os.Getenv("GOOGLE_SECRET")
	// githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	// githubSecret := os.Getenv("GITHUB_SECRET")
	// geminiAPI := os.Getenv("GEMINI_API")
	// appEnv := os.Getenv("APP_ENV")
	// redisURL := os.Getenv("REDIS_URL")

	if dbUrl == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	if port == "" {
		log.Fatal("PORT is not set")
	}
	// if minioEndpoint == "" {
	// 	log.Fatal("MINIO_ENDPOINT is not set")
	// }
	// if minioKey == "" {
	// 	log.Fatal("MINIO_KEY is not set")
	// }
	// if minioSecret == "" {
	// 	log.Fatal("MINIO_SECRET is not set")
	// }
	// if minioBucket == "" {
	// 	log.Fatal("MINIO_BUCKET is not set")
	// }
	// if minioUseSSL == "" {
	// 	log.Fatal("MINIO_USE_SSL is not set")
	// }
	// if hostURL == "" {
	// 	log.Fatal("HOST_URL is not set")
	// }
	// if googleClientID == "" {
	// 	log.Fatal("GOOGLE_CLIENT_ID is not set")
	// }
	// if googleSecret == "" {
	// 	log.Fatal("GOOGLE_SECRET is not set")
	// }
	// if githubClientID == "" {
	// 	log.Fatal("GITHUB_CLIENT_ID is not set")
	// }
	// if githubSecret == "" {
	// 	log.Fatal("GITHUB_SECRET is not set")
	// }
	// if geminiAPI == "" {
	// 	log.Fatal("GEMINI_API is not set")
	// }
	// useSSL := false
	// if minioUseSSL != "" {
	// 	var parseErr error
	// 	useSSL, parseErr = strconv.ParseBool(minioUseSSL)
	// 	if parseErr != nil {
	// 		log.Fatalf("MINIO_USE_SSL must be a boolean value: %v", parseErr)
	// 	}
	// }
	return &Config{
		// AppEnv:         appEnv,
		HTTPPort:       port,
		DatabaseURL:    dbUrl,
		// RedisURL:       redisURL,
		// MinioEndpoint:  minioEndpoint,
		// MinioAccessKey: minioKey,
		// MinioSecretKey: minioSecret,
		// MinioBucket:    minioBucket,
		// MinioUseSSL:    useSSL,
		// HostURL:        hostURL,
		// GoogleClientID: googleClientID,
		// GoogleSecret:   googleSecret,
		// GithubClientID: githubClientID,
		// GithubSecret:   githubSecret,
		// GeminiAPI:      geminiAPI,
	}

}

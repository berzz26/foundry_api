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
	hostURL := os.Getenv("HOST_URL")
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleSecret := os.Getenv("GOOGLE_SECRET")
	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	githubSecret := os.Getenv("GITHUB_SECRET")
	geminiAPI := os.Getenv("GEMINI_API")
	appEnv := os.Getenv("APP_ENV")
	redisURL := os.Getenv("REDIS_URL")

	if dbUrl == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	if port == "" {
		log.Fatal("PORT is not set")
	}
	if hostURL == "" {
		hostURL = "http://localhost:3000"
	}

	return &Config{
		AppEnv:         appEnv,
		HTTPPort:       port,
		DatabaseURL:    dbUrl,
		RedisURL:       redisURL,
		HostURL:        hostURL,
		GoogleClientID: googleClientID,
		GoogleSecret:   googleSecret,
		GithubClientID: githubClientID,
		GithubSecret:   githubSecret,
		GeminiAPI:      geminiAPI,
	}

}

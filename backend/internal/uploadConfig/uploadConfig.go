package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort      string
	JWTSecret       string
	AccessTokenExp  int
	RefreshTokenExp int
	UploadPath      string
	MaxUploadSize   int64
	MaxUploadCount  int
	AllowedTypes    []string
	OpenTweet       bool
}

func Load() *Config {
	_ = godotenv.Load()

	maxUploadSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "5242880"), 10, 64)

	return &Config{
		ServerPort:      getEnv("SERVER_PORT", ":8080"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		AccessTokenExp:  getEnvInt("ACCESS_TOKEN_EXP", 15),
		RefreshTokenExp: getEnvInt("REFRESH_TOKEN_EXP", 7),
		UploadPath:      getEnv("UPLOAD_PATH", "./uploads/images"),
		MaxUploadSize:   maxUploadSize,
		MaxUploadCount:  getEnvInt("MAX_UPLOAD_COUNT", 4),
		AllowedTypes:    []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
		OpenTweet:       getEnv("OPEN_TWEET", "false") == "true",
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type RateLimiterConfig struct {
	IpMaxRequest    int64
	IpBlockDuration int64
}
type RedisConfig struct {
	Port     string
	Host     string
	Password string
	Prefix   string
	DB       int
}

type ServerConfig struct {
	Port string
}

type AuthConfig struct {
	JwtSecret string
}

type Config struct {
	Redis       RedisConfig
	Server      ServerConfig
	RateLimiter RateLimiterConfig
	Auth        AuthConfig
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		redisDB = 0
	}

	iPmaxRequest, err := strconv.ParseInt(getEnv("IP_MAX_REQUESTS_PER_SECOND", "10"), 10, 64)
	if err != nil {
		iPmaxRequest = 10
	}

	iPblockDuration, err := strconv.ParseInt(getEnv("IP_BLOCK_DURATION", "60"), 10, 64)
	if err != nil {
		iPblockDuration = 60
	}

	return &Config{
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			Prefix:   getEnv("REDIS_PREFIX", ""),
			DB:       redisDB,
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		RateLimiter: RateLimiterConfig{
			IpMaxRequest:    iPmaxRequest,
			IpBlockDuration: iPblockDuration,
		},
		Auth: AuthConfig{
			JwtSecret: getEnv("JWT_SECRET_KEY", "defaultsecret"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

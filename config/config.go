package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type RateLimiterConfig struct {
	IpMaxRequest       int
	IpBlockDuration    int
	TokenMaxRequest    int
	TokenBlockDuration int
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

type Config struct {
	Redis       RedisConfig
	Server      ServerConfig
	RateLimiter RateLimiterConfig
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

	iPmaxRequest, err := strconv.Atoi(getEnv("IP_MAX_REQUESTS_PER_SECOND", "10"))
	if err != nil {
		iPmaxRequest = 0
	}

	iPblockDuration, err := strconv.Atoi(getEnv("IP_BLOCK_DURATION", "60"))
	if err != nil {
		iPblockDuration = 0
	}

	tokenMaxRequest, err := strconv.Atoi(getEnv("TOKEN_MAX_REQUESTS_PER_SECOND", "10"))
	if err != nil {
		tokenMaxRequest = 0
	}

	tokenBlockDuration, err := strconv.Atoi(getEnv("TOKEN_BLOCK_DURATION", "3000"))
	if err != nil {
		tokenBlockDuration = 0
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
			IpMaxRequest:       iPmaxRequest,
			IpBlockDuration:    iPblockDuration,
			TokenMaxRequest:    tokenMaxRequest,
			TokenBlockDuration: tokenBlockDuration,
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

package blogConfig

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type redisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Timeout  time.Duration
}

func loadRedisConfig() *redisConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("REDIS: Error loading .env file")
	}
	timeoutStr := getEnv("REDIS_TIMEOUT", "5s")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Println("REDIS: Error parsing timeout duration")
		timeout = 5 * time.Second
	}
	dbStr := getEnv("REDIS_DB", "0")
	var db int
	_, err = fmt.Sscanf(dbStr, "%d", &db)
	if err != nil {
		return nil
	}
	return &redisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       db,
		Timeout:  timeout,
	}
}

func ConnectRedis() *redis.Client {
	config := loadRedisConfig()
	if config == nil {
		return nil
	}
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  config.Timeout,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		PoolSize:     10,
	})

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("REDIS: Failed to connect to Redis at %s: %v", addr, err))
	}

	log.Println("REDIS: Connected to Redis")
	return client
}

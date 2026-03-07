package blogConfig

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func loadDBConfig() map[string]string {
	// 可选：加载 .env 文件（如果存在）
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	return map[string]string{
		"host":     getEnv("DB_HOST", "localhost"),
		"user":     getEnv("DB_USER", "postgres"),
		"password": getEnv("DB_PASS", ""),
		"dbname":   getEnv("DB_NAME", "myapp"),
		"port":     getEnv("DB_PORT", "5432"),
		"sslmode":  getEnv("DB_SSLMODE", "disable"),
		"TimeZone": getEnv("DB_TIMEZONE", "UTC"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func buildDSN(config map[string]string) string {
	var dsn string
	for k, v := range config {
		dsn += fmt.Sprintf("%s=%s ", k, url.QueryEscape(v))
	}
	return dsn
}

func ConnectDB() *gorm.DB {
	config := loadDBConfig()
	dsn := buildDSN(config)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get sql.DB")
	}
	if err := sqlDB.Ping(); err != nil {
		panic("failed to ping database")
	}

	// SetMaxIdleConns 设置空闲连接池中最大连接数
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库的最大连接数
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("🐘 Connected to PostgreSQL")
	return db
}

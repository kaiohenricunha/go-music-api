package config

import (
	"fmt"
	"os"

	"github.com/kaiohenricunha/go-music-k8s/backend/db" // Adjust import path as necessary
	"gorm.io/gorm"
)

type Config struct {
	DbHost     string
	DbName     string
	DbPass     string
	DbUser     string
	DB         *gorm.DB
	ServerPort string
}

func NewConfig() (*Config, error) {
	cfg := &Config{
		DbHost:     getEnv("CONFIG_DBHOST", "localhost:3306"),
		DbName:     getEnv("CONFIG_DBNAME", "infnet_music_db"),
		DbPass:     getEnv("CONFIG_DBPASS", "secret"),
		DbUser:     getEnv("CONFIG_DBUSER", "root"),
		ServerPort: getEnv("CONFIG_SERVER_PORT", "8081"),
	}

	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbName)
	var err error
	cfg.DB, err = db.InitDB(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return cfg, nil
}

// getEnv retrieves environment variables or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

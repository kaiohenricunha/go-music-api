package db

import (
	"fmt"
	"log"
	"strings"
	// "net/http"
	// "io/ioutil"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB initializes the database and seeds it with initial data.
func InitDB(dsn string) (*gorm.DB, error) {
	db := connectDB(dsn)

	if err := migrateSchema(db); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	if err := seedData(db); err != nil {
		return nil, fmt.Errorf("failed to seed data: %w", err)
	}

	return db, nil
}

// connectDB handles the database connection and creation if it doesn't exist.
func connectDB(dsn string) *gorm.DB {
	// Split DSN to extract database name and base DSN for initial connection.
	baseDSN, dbName := splitDSN(dsn)

	// Create database if it doesn't exist.
	createDatabase(baseDSN, dbName)

	// Connect to the database using the full DSN.
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

// splitDSN splits the DSN into base DSN and database name.
func splitDSN(dsn string) (baseDSN string, dbName string) {
	idx := strings.LastIndex(dsn, "/")
	baseDSN, dbName = dsn[:idx], dsn[idx+1:]
	if qIdx := strings.Index(dbName, "?"); qIdx != -1 {
		dbName = dbName[:qIdx]
	}
	baseDSN += "/"
	return
}

// createDatabase connects without a specific database and attempts to create it if it doesn't exist.
func createDatabase(baseDSN string, dbName string) {
	db, err := gorm.Open(mysql.Open(baseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
	if err != nil {
		log.Fatalf("Failed to create database '%s': %v", dbName, err)
	}
}

// migrateSchema auto-migrates the database schema using GORM's AutoMigrate.
func migrateSchema(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{}, &model.Song{})
}

// seedData seeds the database with initial data if necessary.
func seedData(db *gorm.DB) error {
	if err := seedUsers(db); err != nil {
		return err
	}

	return nil
}

// seedUsers seeds the Users table with initial data.
func seedUsers(db *gorm.DB) error {
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count == 0 {
		users := []model.User{
			{Username: "admin", Password: "admin"},
			{Username: "user2", Password: "user2"},
		}
		for _, user := range users {
			if err := db.Create(&user).Error; err != nil {
				return err
			}
		}
		log.Println("Seeded Users successfully")
	}
	return nil
}

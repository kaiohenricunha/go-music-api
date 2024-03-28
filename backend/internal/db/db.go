package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB returns a new db connection, with automatic database creation if not exists
func InitDB(dsn string) *gorm.DB {
	// Correctly extract the base DSN (without the database name) for the initial connection.
	// The goal is to end up with something like username:password@protocol(address)/
	baseDSN := dsn
	dbName := ""
	if idx := strings.LastIndex(dsn, "/"); idx != -1 {
		baseDSN = dsn[:idx+1] // Include the slash
		dbName = dsn[idx+1:]  // Extract database name
	}
	if idx := strings.Index(dbName, "?"); idx != -1 {
		dbName = dbName[:idx] // Remove parameters from the database name if present
	}

	// Ensure the base DSN ends with a slash but does not include the database name for the initial connection.
	if !strings.HasSuffix(baseDSN, "/") {
		baseDSN += "/"
	}

	// Connect to MySQL without specifying a database, to check/create the database.
	sqlDB, err := sql.Open("mysql", baseDSN)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer sqlDB.Close()

	// Try to create the database if it doesn't exist.
	_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
	if err != nil {
		log.Fatalf("Failed to create database '%s': %v", dbName, err)
	}

	// Now, connect with GORM using the full DSN including the database name.
	fullDSN := dsn // This already includes the database name
	db, err := gorm.Open(mysql.Open(fullDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database with GORM: %v", err)
	}

	// Set connection pool parameters.
	gormDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object from GORM: %v", err)
	}
	gormDB.SetMaxIdleConns(10)
	gormDB.SetMaxOpenConns(100)

	return db
}

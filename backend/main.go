package main

import (
	"log"
	"net/http"

	"github.com/kaiohenricunha/go-music-k8s/backend/config"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/api/routes"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations
	if err := db.AutoMigrate(&model.User{}, &model.Song{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Setup DAOs
	userDAO := dao.NewGormUserDAO(db)
	songDAO := dao.NewGormSongDAO(db)

	// Setup Services
	userService := service.NewUserService(userDAO)
	songService := service.NewSongService(songDAO)

	// Setup API routes
	router := routes.SetupRoutes(userService, songService)

	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

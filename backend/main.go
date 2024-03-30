package main

import (
	"log"
	"net/http"

	"github.com/kaiohenricunha/go-music-k8s/backend/api/routes"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/config"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/service"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db := cfg.DB // Use the *gorm.DB instance from the configuration

	// Setup DAOs with the database connection
	userDAO := dao.NewGormDAO(db)
	songDAO := dao.NewGormDAO(db)
	playlistDAO := dao.NewGormDAO(db)

	// Setup Services with the DAOs
	userService := service.NewUserService(userDAO)
	songService := service.NewSongService(songDAO)
	playlistService := service.NewPlaylistService(playlistDAO)

	// Setup API routes with the services
	router := routes.SetupRoutes(userService, songService, playlistService)

	// Start the server
	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

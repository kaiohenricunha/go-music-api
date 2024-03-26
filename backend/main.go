package main

import (
	"log"

	"github.com/kaiohenricunha/go-music-k8s/backend/api/server"
	"github.com/kaiohenricunha/go-music-k8s/backend/config"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/musicapp"
)

var cfg *config.Config
var err error

func init() {
	log.Print("Welcome to music api...")

	// get a config
	cfg, err = config.NewConfig()
	if err != nil {
		log.Fatal("Config init failed", err)
	}

	// migrate db
	if err = musicapp.DbInit(cfg.DB); err != nil {
		log.Fatal("DB migration failed...")
	}
}

func main() {
	server.Start(cfg)
}

package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaiohenricunha/go-music-k8s/backend/config"
)

var gcfg *config.Config // global config for server

func Start(cfg *config.Config) {
	gcfg = cfg
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/music", GetSongHandler).Methods("GET")
	r.HandleFunc("/api/v1/music", PostSongHandler).Methods("POST")
	r.HandleFunc("/api/v1/music", UpdateSongHandler).Methods("PUT")
	r.HandleFunc("/api/v1/music/{id}", DeleteSongHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}

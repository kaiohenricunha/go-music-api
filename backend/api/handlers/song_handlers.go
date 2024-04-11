package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaiohenricunha/go-music-k8s/backend/api"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/service"
)

type SongHandlers struct {
	songService service.SongService
}

func NewSongHandlers(songService service.SongService) *SongHandlers {
	return &SongHandlers{
		songService: songService,
	}
}

func (h *SongHandlers) GetAllSongsHandler(w http.ResponseWriter, r *http.Request) {
	songs, err := h.songService.GetAllSongs()
	if err != nil {
		api.LogErrorWithDetails(w, "Failed to get songs", err, http.StatusInternalServerError)
		return
	}
	api.RespondWithJSON(w, http.StatusOK, songs)
}

func (h *SongHandlers) GetSongFromSpotifyByIDHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to get song from Spotify by ID")
	vars := mux.Vars(r)
	spotifyID, ok := vars["spotifyID"]
	if !ok || spotifyID == "" {
		api.LogErrorWithDetails(w, "Spotify ID is required", nil, http.StatusBadRequest)
		return
	}

	song, err := h.songService.GetSongFromSpotifyByID(spotifyID)
	if err != nil {
		api.LogErrorWithDetails(w, "Failed to get song from Spotify", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, song)
}

func (h *SongHandlers) SearchSongsFromSpotifyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to search songs...")
	songName, artistName := r.URL.Query().Get("songName"), r.URL.Query().Get("artistName")
	if songName == "" || artistName == "" {
		api.LogErrorWithDetails(w, "Song name and artist name are required", nil, http.StatusBadRequest)
		return
	}

	songs, err := h.songService.SearchSongsFromSpotify(songName, artistName)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrSongNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrFetchingFromSpotify) {
			status = http.StatusBadGateway
		}
		api.LogErrorWithDetails(w, "Failed to search for songs", err, status)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, songs)
}

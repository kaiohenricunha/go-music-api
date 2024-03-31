package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/kaiohenricunha/go-music-k8s/backend/api"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/service"
)

// SongHandlers encapsulates handlers for dealing with songs.
type SongHandlers struct {
	songService service.SongService
}

// NewSongHandlers creates an instance of SongHandlers.
func NewSongHandlers(songService service.SongService) *SongHandlers {
	return &SongHandlers{
		songService: songService,
	}
}

// GetAllSongs handles GET requests to retrieve all songs.
func (h *SongHandlers) GetAllSongsHandler(w http.ResponseWriter, r *http.Request) {
	songs, err := h.songService.GetAllSongs()
	if err != nil {
		api.LogErrorWithDetails(w, "Failed to get songs", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, songs)
}

// GetSongFromSpotifyByIdHandler handles GET requests to retrieve a song from Spotify by its ID.
func (h *SongHandlers) GetSongFromSpotifyByIDHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to get song from Spotify by ID")

	// Retrieve the Spotify ID from the path variables
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

// SearchSongsFromSpotifyHandler handles GET requests to search for songs on Spotify.
func (h *SongHandlers) SearchSongsFromSpotifyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to search for songs on Spotify")

	// log the whole request
	log.Printf("Request: %v", r.URL)

	// Retrieve the query parameters for track name and artist name
	songName := r.URL.Query().Get("songName")
	artistName := r.URL.Query().Get("artistName")

	if songName == "" || artistName == "" {
		api.LogErrorWithDetails(w, "Track name and artist name are required", nil, http.StatusBadRequest)
		return
	}

	songs, err := h.songService.SearchSongsFromSpotify(songName, artistName)
	if err != nil {
		api.LogErrorWithDetails(w, "Failed to search for songs on Spotify", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, songs)
}

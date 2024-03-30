package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kaiohenricunha/go-music-k8s/backend/api"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
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

// AddSong handles POST requests to add a new song.
func (h *SongHandlers) AddSong(w http.ResponseWriter, r *http.Request) {
	var song model.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		api.LogErrorWithDetails(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	if err := h.songService.AddSong(&song); err != nil {
		api.LogErrorWithDetails(w, "Failed to add song", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Song added successfully"})
}

// GetAllSongs handles GET requests to retrieve all songs.
func (h *SongHandlers) GetAllSongs(w http.ResponseWriter, r *http.Request) {
	songs, err := h.songService.GetAllSongs()
	if err != nil {
		api.LogErrorWithDetails(w, "Failed to get songs", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, songs)
}

// UpdateSong handles PUT requests to update a song.
func (h *SongHandlers) UpdateSong(w http.ResponseWriter, r *http.Request) {
	var song model.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		api.LogErrorWithDetails(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	if err := h.songService.UpdateSong(&song); err != nil {
		api.LogErrorWithDetails(w, "Failed to update song", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Song updated successfully"})
}

// DeleteSong handles DELETE requests to remove a song.
func (h *SongHandlers) DeleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		api.LogErrorWithDetails(w, "Invalid song ID", err, http.StatusBadRequest)
		return
	}

	if err := h.songService.DeleteSong(uint(id)); err != nil {
		api.LogErrorWithDetails(w, "Failed to delete song", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Song deleted successfully"})
}

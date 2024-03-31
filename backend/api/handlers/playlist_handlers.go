package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kaiohenricunha/go-music-k8s/backend/api"
	"github.com/kaiohenricunha/go-music-k8s/backend/api/middleware"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/service"

	"github.com/gorilla/mux"
)

// PlaylistHandlers encapsulates handlers for dealing with playlists.
type PlaylistHandlers struct {
	playlistService service.PlaylistService
}

// NewPlaylistHandlers creates an instance of PlaylistHandlers.
func NewPlaylistHandlers(playlistService service.PlaylistService) *PlaylistHandlers {
	return &PlaylistHandlers{
		playlistService: playlistService,
	}
}

// CreatePlaylist handles POST requests to create a new playlist.
func (h *PlaylistHandlers) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var playlist model.Playlist
	if err := json.NewDecoder(r.Body).Decode(&playlist); err != nil {
		api.LogErrorWithDetails(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	// Get authenticated user ID from context
	authUserID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		api.LogErrorWithDetails(w, "Failed to get authenticated user ID", nil, http.StatusInternalServerError)
		return
	}

	// Set the UserID for the playlist to the authenticated user's ID
	playlist.UserID = authUserID

	// Attempt to create the playlist
	if err := h.playlistService.CreatePlaylist(&playlist); err != nil {
		api.LogErrorWithDetails(w, "Failed to create playlist", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Playlist created successfully"})
}

func (h *PlaylistHandlers) AddSongToPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	// Get playlist ID from URL parameters
	vars := mux.Vars(r)
	playlistID, err := strconv.Atoi(vars["playlistID"])
	if err != nil {
		api.LogErrorWithDetails(w, "Invalid playlist ID", err, http.StatusBadRequest)
		return
	}

	// Get song ID from URL parameters
	songID, err := strconv.Atoi(vars["songID"])
	if err != nil {
		api.LogErrorWithDetails(w, "Invalid song ID", err, http.StatusBadRequest)
		return
	}

	// Attempt to add the song to the playlist
	if err := h.playlistService.AddSongToPlaylist(uint(playlistID), uint(songID)); err != nil {
		api.LogErrorWithDetails(w, "Failed to add song to playlist", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Song added to playlist successfully"})
}

// GetAllPlaylistsHandler handles GET requests to retrieve all playlists.
func (h *PlaylistHandlers) GetAllPlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	playlists, err := h.playlistService.GetAllPlaylists()
	if err != nil {
		// handle error, e.g., log it and return an appropriate HTTP error response
		api.LogErrorWithDetails(w, "Failed to fetch playlists", err, http.StatusInternalServerError)
		return
	}

	// Respond with the fetched playlists
	api.RespondWithJSON(w, http.StatusOK, playlists)
}

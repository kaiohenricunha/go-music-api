package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kaiohenricunha/go-music-k8s/backend/api"
	"github.com/kaiohenricunha/go-music-k8s/backend/api/middleware"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/service"
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

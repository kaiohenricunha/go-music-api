package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kaiohenricunha/go-music-k8s/backend/api"
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
	var playlist model.Playlist
	if err := json.NewDecoder(r.Body).Decode(&playlist); err != nil {
		api.LogErrorWithDetails(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	// get the authenticated user ID from the request content
	authUserID, ok := r.Context().Value("userID").(uint)
	if !ok {
		api.LogErrorWithDetails(w, "Failed to get authenticated user ID", nil, http.StatusInternalServerError)
		return
	}

	if err := h.playlistService.CreatePlaylist(authUserID, &playlist); err != nil {
		api.LogErrorWithDetails(w, "Failed to create playlist", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Playlist created successfully"})
}

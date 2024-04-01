package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/kaiohenricunha/go-music-k8s/backend/api"
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

// GetAllPlaylistsHandler handles GET requests to list all playlists.
func (h *PlaylistHandlers) GetAllPlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	playlists, err := h.playlistService.GetAllPlaylists()
	if err != nil {
		api.LogErrorWithDetails(w, "Failed to retrieve playlists", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, playlists)
}

// GetPlaylistByIDHandler handles GET requests to retrieve a playlist by ID.
func (h *PlaylistHandlers) GetPlaylistByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playlistID := vars["playlistID"]

	playlist, err := h.playlistService.GetPlaylistByID(playlistID)
	if err != nil {
		api.LogErrorWithDetails(w, "Failed to retrieve playlist", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, playlist)
}

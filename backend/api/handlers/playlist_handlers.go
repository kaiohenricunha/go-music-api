package handlers

import (
	"errors"
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

// AddSongToPlaylistHandler handles POST requests to add a song to a playlist.
func (h *PlaylistHandlers) AddSongToPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playlistID := vars["playlistID"]
	songID := vars["songID"] // Assuming the song ID is passed as a path parameter or you could choose to receive it in the request body.

	// Call the service method to add the song to the playlist
	err := h.playlistService.AddSongToPlaylist(playlistID, songID)
	if err != nil {
		if errors.Is(err, service.ErrPlaylistNotFound) || errors.Is(err, service.ErrSongNotFound) {
			api.LogErrorWithDetails(w, "Playlist or song not found", err, http.StatusNotFound)
			return
		}
		api.LogErrorWithDetails(w, "Failed to add song to playlist", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Song added to playlist successfully"})
}

// RemoveSongFromPlaylistHandler handles DELETE requests to remove a song from a playlist.
func (h *PlaylistHandlers) RemoveSongFromPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playlistID := vars["playlistID"]
	songID := vars["songID"] // Assuming the song ID is passed as a path parameter.

	// Call the service method to remove the song from the playlist
	err := h.playlistService.RemoveSongFromPlaylist(playlistID, songID)
	if err != nil {
		if errors.Is(err, service.ErrPlaylistNotFound) || errors.Is(err, service.ErrSongNotFound) {
			api.LogErrorWithDetails(w, "Playlist or song not found", err, http.StatusNotFound)
			return
		}
		api.LogErrorWithDetails(w, "Failed to remove song from playlist", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Song removed from playlist successfully"})
}

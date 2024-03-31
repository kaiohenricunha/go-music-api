package routes

import (
	"github.com/gorilla/mux"
	"github.com/kaiohenricunha/go-music-k8s/backend/api/handlers"
	"github.com/kaiohenricunha/go-music-k8s/backend/api/middleware"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/service"
)

// SetupRoutes configures and returns a new router with all routes defined.
func SetupRoutes(userService service.UserService, songService service.SongService, playlistService service.PlaylistService) *mux.Router {
	r := mux.NewRouter()

	// Middleware for basic auth, excluding specific routes
	authMiddleware := middleware.BasicAuthMiddleware(userService)

	// Apply global middleware directly
	r.Use(middleware.LoggingMiddleware)

	// Initialize handlers
	userHandlers := handlers.NewUserHandlers(userService)
	songHandlers := handlers.NewSongHandlers(songService)
	playlistHandlers := handlers.NewPlaylistHandlers(playlistService)

	// Public routes
	publicRouter := r.PathPrefix("/api/v1").Subrouter()
	publicRouter.HandleFunc("/register", userHandlers.RegisterUserHandler).Methods("POST")

	// User-specific routes
	authUserRouter := r.PathPrefix("/api/v1/users").Subrouter()
	authUserRouter.Use(authMiddleware)
	authUserRouter.HandleFunc("", userHandlers.ListUsersHandler).Methods("GET")
	authUserRouter.HandleFunc("/{id}", userHandlers.UpdateUserHandler).Methods("PUT")
	authUserRouter.HandleFunc("/{id}", userHandlers.DeleteUserHandler).Methods("DELETE")
	authUserRouter.HandleFunc("/{username}", userHandlers.FindUserByUsernameHandler).Methods("GET")

	// Song Routes
	authSongRouter := r.PathPrefix("/api/v1/songs").Subrouter()
	authSongRouter.Use(authMiddleware)
	// Place the more specific /search route before the more general /{spotifyID} route
	authSongRouter.HandleFunc("/search", songHandlers.SearchSongsFromSpotifyHandler).Methods("GET")
	authSongRouter.HandleFunc("", songHandlers.GetAllSongsHandler).Methods("GET")
	authSongRouter.HandleFunc("/{spotifyID}", songHandlers.GetSongFromSpotifyByIDHandler).Methods("GET")

	// Playlist Routes
	authPlaylistRouter := r.PathPrefix("/api/v1/playlists").Subrouter()
	authPlaylistRouter.Use(authMiddleware)
	authPlaylistRouter.HandleFunc("", playlistHandlers.CreatePlaylist).Methods("POST")
	authPlaylistRouter.HandleFunc("/{playlistID}/songs/{songID}", playlistHandlers.AddSongToPlaylistHandler).Methods("POST")
	authPlaylistRouter.HandleFunc("", playlistHandlers.GetAllPlaylistsHandler).Methods("GET")

	return r
}

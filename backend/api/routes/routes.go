package routes

import (
	"net/http"

	goHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kaiohenricunha/go-music-k8s/backend/api/handlers"
	"github.com/kaiohenricunha/go-music-k8s/backend/api/middleware"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/service"
)

func SetupRoutes(userService service.UserService, songService service.SongService, playlistService service.PlaylistService) http.Handler {
	r := mux.NewRouter()

	// Middleware for JWT Auth
	jwtMiddleware := middleware.JWTAuthMiddleware(userService)

	// Apply global middleware directly
	r.Use(middleware.LoggingMiddleware)

	// Initialize handlers
	userHandlers := handlers.NewUserHandlers(userService)
	songHandlers := handlers.NewSongHandlers(songService)
	playlistHandlers := handlers.NewPlaylistHandlers(playlistService)

	// Public routes (no auth needed)
	publicRouter := r.PathPrefix("/api/v1").Subrouter()
	publicRouter.HandleFunc("/register", userHandlers.RegisterUserHandler).Methods("POST")
	// The login route itself will handle basic authentication inside its handler
	publicRouter.HandleFunc("/login", userHandlers.UserLoginHandler).Methods("POST")

	// Protected routes (JWT Auth)
	protectedRouter := r.PathPrefix("/api/v1").Subrouter()
	protectedRouter.Use(jwtMiddleware) // Apply JWT middleware here

	// User-specific routes
	protectedRouter.HandleFunc("/users", userHandlers.ListUsersHandler).Methods("GET")
	protectedRouter.HandleFunc("/users/{username}", userHandlers.GetUserByUsername).Methods("GET")
	// TODO: implement a route to update and delete users

	// Song Routes
	protectedRouter.HandleFunc("/songs", songHandlers.GetAllSongsHandler).Methods("GET")
	protectedRouter.HandleFunc("/songs/search", songHandlers.SearchSongsFromSpotifyHandler).Methods("GET")
	protectedRouter.HandleFunc("/songs/{spotifyID}", songHandlers.GetSongFromSpotifyByIDHandler).Methods("GET")

	// Playlist Routes
	protectedRouter.HandleFunc("/playlists", playlistHandlers.GetAllPlaylistsHandler).Methods("GET")
	protectedRouter.HandleFunc("/playlists/{playlistID}", playlistHandlers.GetPlaylistByIDHandler).Methods("GET")
	protectedRouter.HandleFunc("/playlists/{playlistID}/songs/{songID}", playlistHandlers.AddSongToPlaylistHandler).Methods("POST")
	protectedRouter.HandleFunc("/playlists/{playlistID}/songs/{songID}", playlistHandlers.RemoveSongFromPlaylistHandler).Methods("DELETE")

	// Wrap the entire router with CORS middleware
	corsMiddleware := goHandlers.CORS(
		goHandlers.AllowedOrigins([]string{"http://localhost:3000"}),
		goHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		goHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		goHandlers.AllowCredentials(),
	)

	// Apply CORS middleware to the router and return
	return corsMiddleware(r)
}

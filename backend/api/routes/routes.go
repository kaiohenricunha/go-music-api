package routes

import (
	"github.com/gorilla/mux"
	"github.com/kaiohenricunha/go-music-k8s/backend/api/handlers"
	"github.com/kaiohenricunha/go-music-k8s/backend/api/middleware"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/service"
)

// SetupRoutes configures and returns a new router with all routes defined.
func SetupRoutes(userService service.UserService, songService service.SongService) *mux.Router {
	r := mux.NewRouter()

	// Apply global middleware directly
	r.Use(middleware.LoggingMiddleware)

	// Middleware for basic auth
	authMiddleware := middleware.BasicAuthMiddleware(userService)

	// User-specific routes with auth middleware
	userRouter := r.PathPrefix("/api/v1").Subrouter()
	userRouter.Use(authMiddleware)
	userHandlers := handlers.NewUserHandlers(userService)
	userRouter.HandleFunc("/register", userHandlers.RegisterUserHandler).Methods("POST")
	userRouter.HandleFunc("/users", userHandlers.ListUsersHandler).Methods("GET")
	userRouter.HandleFunc("/users/{id}", userHandlers.UpdateUserHandler).Methods("PUT")
	userRouter.HandleFunc("/users/{id}", userHandlers.DeleteUserHandler).Methods("DELETE")

	// Song Routes - assuming these routes also require authentication
	songRouter := r.PathPrefix("/api/v1/music").Subrouter()
	songRouter.Use(authMiddleware)
	songHandlers := handlers.NewSongHandlers(songService)
	songRouter.HandleFunc("", songHandlers.AddSong).Methods("POST")
	songRouter.HandleFunc("", songHandlers.GetAllSongs).Methods("GET")
	songRouter.HandleFunc("/{id}", songHandlers.UpdateSong).Methods("PUT")
	songRouter.HandleFunc("/{id}", songHandlers.DeleteSong).Methods("DELETE")

	return r
}

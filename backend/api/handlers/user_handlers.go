package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaiohenricunha/go-music-k8s/backend/api"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/service"
)

// UserHandlers encapsulates handlers related to user operations.
type UserHandlers struct {
	userService service.UserService
}

// NewUserHandlers creates a new instance of UserHandlers.
func NewUserHandlers(userService service.UserService) *UserHandlers {
	return &UserHandlers{
		userService: userService,
	}
}

// RegisterUserHandler handles the user registration requests.
func (h *UserHandlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		api.LogErrorWithDetails(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	err := h.userService.RegisterUser(&user)
	if err != nil {
		if err == service.ErrUsernameTaken {
			api.LogErrorWithDetails(w, "Username already taken", err, http.StatusConflict) // Use HTTP 409 Conflict for duplicate username
		} else {
			api.LogErrorWithDetails(w, "Failed to register user", err, http.StatusInternalServerError)
		}
		return
	}

	api.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

// ListUsersHandler handles requests to list all users.
func (h *UserHandlers) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		api.LogErrorWithDetails(w, "Failed to retrieve users", err, http.StatusInternalServerError)
		return
	}

	// Prepare the response by omitting the password from each user.
	var response []map[string]interface{}
	for _, user := range users {
		userMap := map[string]interface{}{
			"ID":        user.ID,
			"CreatedAt": user.CreatedAt,
			"UpdatedAt": user.UpdatedAt,
			"Username":  user.Username,
			"Playlists": user.Playlists,
			"Role":      user.Role,
			"FullName":  user.FullName,
			"Email":     user.Email,
		}
		response = append(response, userMap)
	}

	api.RespondWithJSON(w, http.StatusOK, response)
}

// GetUserByUsername handles requests to find a user by their username.
func (h *UserHandlers) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	user, err := h.userService.GetUserByUsername(username)
	if err != nil {
		if err == service.ErrUserNotFound {
			api.LogErrorAndRespond(w, "User not found", http.StatusNotFound)
		} else {
			api.LogErrorWithDetails(w, "Failed to retrieve user", err, http.StatusInternalServerError)
		}
		return
	}

	// Convert the user object to a map excluding the password
	userMap := map[string]interface{}{
		"ID":        user.ID,
		"CreatedAt": user.CreatedAt,
		"UpdatedAt": user.UpdatedAt,
		"Username":  user.Username,
		"Playlists": user.Playlists,
		"Role":      user.Role,
		"FullName":  user.FullName,
		"Email":     user.Email,
	}

	api.RespondWithJSON(w, http.StatusOK, userMap)
}

// UserLoginHandler handles the user login requests.
func (h *UserHandlers) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Missing credentials", http.StatusBadRequest)
		return
	}

	userID, valid := h.userService.ValidateUser(username, password)
	if !valid {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT for the user
	token, err := api.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

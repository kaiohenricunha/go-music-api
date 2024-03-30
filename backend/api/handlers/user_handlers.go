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
		}
		response = append(response, userMap)
	}

	api.RespondWithJSON(w, http.StatusOK, response)
}

// UpdateUserHandler handles requests to update a user.
func (h *UserHandlers) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestUserID, err := strconv.Atoi(vars["id"])
	if err != nil {
		api.LogErrorWithDetails(w, "Invalid user ID", err, http.StatusBadRequest)
		return
	}

	authUserID, ok := r.Context().Value("userID").(uint)
	if !ok || authUserID != uint(requestUserID) {
		api.LogErrorAndRespond(w, "Unauthorized to update this user", http.StatusUnauthorized)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		api.LogErrorWithDetails(w, "Error decoding user data", err, http.StatusBadRequest)
		return
	}

	user.ID = uint(requestUserID) // Ensuring the user ID is correctly set
	err = h.userService.UpdateUser(&user)
	if err != nil {
		api.LogErrorWithDetails(w, "Failed to update user", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

// DeleteUserHandler handles requests to delete a user.
func (h *UserHandlers) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		api.LogErrorWithDetails(w, "Invalid user ID", err, http.StatusBadRequest)
		return
	}

	err = h.userService.DeleteUser(uint(userID))
	if err != nil {
		api.LogErrorWithDetails(w, "Failed to delete user", err, http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// FindUserByUsernameHandler handles requests to find a user by their username.
func (h *UserHandlers) FindUserByUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	user, err := h.userService.FindUserByUsername(username)
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
	}

	api.RespondWithJSON(w, http.StatusOK, userMap)
}

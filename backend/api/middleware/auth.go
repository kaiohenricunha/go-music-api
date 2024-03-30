package middleware

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/service"
)

// Define a custom type for context keys to avoid collisions
type ContextKey string

// Export UserContextKey by capitalizing its first letter
const UserContextKey ContextKey = "userID"

func BasicAuthMiddleware(userService service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
				http.Error(w, "Authorization required", http.StatusUnauthorized)
				return
			}

			encodedPayload := strings.TrimPrefix(authHeader, "Basic ")
			payload, err := base64.StdEncoding.DecodeString(encodedPayload)
			if err != nil {
				http.Error(w, "Invalid authorization format", http.StatusBadRequest)
				return
			}

			pair := strings.SplitN(string(payload), ":", 2)
			if len(pair) != 2 {
				http.Error(w, "Invalid authorization format", http.StatusBadRequest)
				return
			}

			userID, valid := userService.ValidateUser(pair[0], pair[1])
			if !valid {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
				return
			}

			// Use the exported UserContextKey
			ctx := context.WithValue(r.Context(), UserContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

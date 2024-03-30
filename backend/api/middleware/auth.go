package middleware

import (
	"encoding/base64"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/service"
	"net/http"
	"strings"
)

func BasicAuthMiddleware(userService service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
				http.Error(w, "Authorization required", http.StatusUnauthorized)
				return
			}

			payload, _ := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
			pair := strings.SplitN(string(payload), ":", 2)

			if len(pair) != 2 || !userService.ValidateUser(pair[0], pair[1]) {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

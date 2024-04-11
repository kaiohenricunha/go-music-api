package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/service"
)

// contextKey used to avoid collisions
type contextKey string

const userContextKey contextKey = "userID"

// Replace the jwtKey initialization with this
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func JWTAuthMiddleware(userService service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				w.Header().Set("WWW-Authenticate", `Bearer realm="restricted"`)
				http.Error(w, "Authorization required", http.StatusUnauthorized)
				return
			}

			claims := &jwt.StandardClaims{}
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
			var err error
			tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Log the token's content if validation fails
				if err != nil {
					log.Printf("JWT Validation Error: %v. Token: %s", err, tokenString)
				}
				return jwtKey, nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					http.Error(w, "Invalid token signature", http.StatusUnauthorized)
					return
				}

				// Log any validation errors in detail
				log.Printf("JWT Validation Error: %v", err)
				http.Error(w, "Invalid token", http.StatusBadRequest)
				return
			}

			if !tkn.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// If the token was valid, set the user ID in the context
			ctx := context.WithValue(r.Context(), userContextKey, claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

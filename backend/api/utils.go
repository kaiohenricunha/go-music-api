package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// RespondWithJSON takes a payload, marshals it to JSON, and writes it to the response writer with the given status code.
func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		LogErrorWithDetails(w, "Error marshaling JSON", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(response); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

// LogErrorAndRespond logs the provided error message and sends an HTTP error response with the given status code.
func LogErrorAndRespond(w http.ResponseWriter, errMsg string, statusCode int) {
	log.Println(errMsg)
	http.Error(w, errMsg, statusCode)
}

// LogErrorWithDetails logs an error message along with additional details and sends an HTTP error response.
func LogErrorWithDetails(w http.ResponseWriter, errMsg string, err error, statusCode int) {
	log.Printf("%s: %v", errMsg, err)
	http.Error(w, errMsg, statusCode)
}

// GenerateJWT generates a new JWT token for a given user ID.
func GenerateJWT(userID uint) (string, error) {
	// Fetch the JWT secret key from the environment
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	// Define token claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Token expires in 1 day
		Issuer:    "go-music-k8s",
		Subject:   fmt.Sprintf("%d", userID),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

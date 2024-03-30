package api

import (
	"encoding/json"
	"log"
	"net/http"
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

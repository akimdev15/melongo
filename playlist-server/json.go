package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Martial the JSON data -> bytes
	data, err := json.Marshal(payload)

	if err != nil {
		slog.Error("Failed to marshal JSON response", "error", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		slog.Error("Failed to write to the response", "error", err)
		return
	}
}

// respondWithError constructs consistent error response
func respondWithError(w http.ResponseWriter, code int, msg string) {
	// error in 400's are all client errors
	// using the api in a weird way
	if code > 499 {
		slog.Error("Responding with 5XX error: ", "error", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

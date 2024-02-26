package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Martial the JSON data -> bytes
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to martial JSON resposne: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

// respondWithError constructs consistent error response
func respondWithError(w http.ResponseWriter, code int, msg string) {
	// error in 400's are all client errors
	// using the api in a weird way
	if code > 499 {
		log.Println("Responding with 5XX error: ", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

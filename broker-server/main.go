package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/akimdev15/melongo/broker/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const PORT = ":8080"

type UserCreatedPayload struct {
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	mux := http.NewServeMux()

	// define routes here
	mux.HandleFunc("GET /authorize", handleAuthorization)
	mux.HandleFunc("GET /callback", handleSpotifyCallback)
	mux.HandleFunc("GET /missedTracks", middlewareAuth(handleGetMissedTracks))
	mux.HandleFunc("POST /handle", middlewareAuth(handleSubmission))
	mux.HandleFunc("POST /createPlaylist", middlewareAuth(handleCreatePlaylist))
	mux.HandleFunc("POST /melonTop100/create", middlewareAuth(handleMelonTop100))
	mux.HandleFunc("POST /melonTop100/save", middlewareAuth(handleSaveMelonTop100DB))
	mux.HandleFunc("POST /resolveMissedTracks", middlewareAuth(handleResolveMissedTracks))

	corsHandler := corsMiddleware(mux)

	fmt.Printf("Listening on port%s", PORT)
	err = http.ListenAndServe(PORT, corsHandler)
	if err != nil {
		return
	}

}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Expose-Headers", "Link")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "300")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func handleSpotifyCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	// connect to server
	conn, err := grpc.Dial("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		fmt.Println("Error during gRPC dial")
		return
	}
	defer conn.Close()

	// create client
	client := proto.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// call AuthorizeUser method in the auth service
	response, err := client.AuthorizeUser(ctx, &proto.AuthCallbackRequest{
		Code: code,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Error in AuthorizeUser method")
		return
	}

	var payload UserCreatedPayload
	payload.Name = response.Name
	payload.AccessToken = response.AccessToken

	writeJSON(w, http.StatusAccepted, payload)
}

func handleSubmission(w http.ResponseWriter, r *http.Request, accessToken string, userID string) {
	fmt.Println("Hit Handle Submission")
	fmt.Printf("Access Token: %s\n", accessToken)
	fmt.Printf("UserID: %s\n", userID)
}

func handleAuthorization(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("ClientID")
	redirectURI := os.Getenv("RedirectURI")
	scopes := os.Getenv("Scopes")

	// Construct Authorization URL
	authURL := fmt.Sprintf(
		"https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s",
		clientID, redirectURI, scopes,
	)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// Write JSON
func writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

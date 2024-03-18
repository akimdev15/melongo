package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/akimdev15/melongo/auth-server/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

const (
	RedirectURI = "http://localhost:8080/callback"
	Scopes      = "user-read-email user-read-private playlist-modify-public playlist-modify-private playlist-read-collaborative playlist-read-private"
)

var (
	clientID     string
	clientSecret string
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type apiConfig struct {
	DB *database.Queries
}

const PORT = ":8081"

func main() {
	// Step 1: Get client info
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
	// Step 1.1: Setup Database
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the env file")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("dbUrl: %s\n", dbURL)
		log.Fatal("Can't connect to the database. err: ", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	clientID = os.Getenv("ClientID")
	clientSecret = os.Getenv("ClientSecret")
	if clientID == "" || clientSecret == "" {
		log.Fatal("")
	}

	// Redirect User to Spotify Authorization Page
	http.HandleFunc("GET /authorize", handleRedirect)

	// Handle Authorization Response
	http.HandleFunc("GET /callback", apiCfg.handleAuthorizationResponse)

	go apiCfg.grpcListen()

	// Start server
	fmt.Println("Server listening on port ", PORT)
	err = http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatalf("Error starting the server %v", err)
		return
	}
}

package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/akimdev15/melongo/auth-server/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var (
	clientID     string
	clientSecret string
	RedirectURI  string
)

type TokenResponse struct {
	AccessToken   string `json:"access_token"`
	Expires_In    int    `json:"expires_in"`
	Refresh_Token string `json:"refresh_token"`
}

type apiConfig struct {
	DB     *database.Queries
	DBConn *sql.DB
}

const PORT = ":8081"

func main() {
	// Step 1: Get client info
	err := godotenv.Load("config.env")
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}
	// Step 1.1: Setup Database
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		slog.Error("DB_URL is not found in the env file")
		os.Exit(1)
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		slog.Error("Can't connect to the database", "error", err)
		os.Exit(1)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB:     db,
		DBConn: conn,
	}

	clientID = os.Getenv("ClientID")
	clientSecret = os.Getenv("ClientSecret")
	RedirectURI = os.Getenv("RedirectURI")
	if clientID == "" || clientSecret == "" || RedirectURI == "" {
		slog.Error("ClientID, ClientSecret, RedirectURI are not found in the env file")
		os.Exit(1)
	}

	// Start gRPC server
	go apiCfg.grpcListen()

	// Start http server
	slog.Info("Listening on", "PORT", PORT)
	err = http.ListenAndServe(PORT, nil)
	if err != nil {
		slog.Error("Error starting the server", "error", err)
		os.Exit(1)
	}
}

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
		DB:     db,
		DBConn: conn,
	}

	clientID = os.Getenv("ClientID")
	clientSecret = os.Getenv("ClientSecret")
	RedirectURI = os.Getenv("RedirectURI")
	if clientID == "" || clientSecret == "" || RedirectURI == "" {
		log.Fatal("")
	}

	// Start gRPC server
	go apiCfg.grpcListen()

	// Start http server
	fmt.Println("Server listening on port ", PORT)
	err = http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatalf("Error starting the server %v", err)
		return
	}
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3/database"
)

type apiConfig struct {
}

func main() {

	godotenv.Load(".env")
	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT is not found in the ENV file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the env file")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to the database")
	}

	db := database.New(conn)
	// TODO - add the config

	// TODO - run the scraper here

	router := chi.NewRouter()
	// adding cors
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// TODO - create new router and add them here

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portStr,
	}
	fmt.Println("Server starting on the PORT: ", portStr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

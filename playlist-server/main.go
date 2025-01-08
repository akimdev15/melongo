package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/akimdev15/melongo/playlist-server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB     *database.Queries
	DBConn *sql.DB
}

const PORT = ":8082"

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
		return
	}

	// Step 1.1: Setup Database
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		slog.Error("DB_URL is not found in the env file")
		return
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		slog.Error("Can't connect to the database", "error", err)
		return
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB:     db,
		DBConn: conn,
	}

	// Start gRPC server
	go apiCfg.grpcListen()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /test", apiCfg.testHandler)
	mux.HandleFunc("GET /testAlbum", apiCfg.testNewAlbumsHandler)
	corsHandler := corsMiddleware(mux)

	slog.Info("Listening on", "PORT", PORT)

	err = http.ListenAndServe(PORT, corsHandler)
	if err != nil {
		slog.Error("Error starting the server", "error", err)
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

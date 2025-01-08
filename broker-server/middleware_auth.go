package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/akimdev15/melongo/broker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authHandler func(http.ResponseWriter, *http.Request, string, string)

func middlewareAuth(handler authHandler) http.HandlerFunc {
	// Creating a anonymous function
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("spotify_access_token")
		if err != nil || cookie == nil {
			slog.Error("No access token found in the cookie")
			respondWithError(w, 401, "Unauthorized. Please login again.")
			return
		}

		accessToken := cookie.Value

		conn, err := grpc.Dial("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			slog.Error("Error during gRPC dial", "error", err)
			respondWithError(w, 403, fmt.Sprintf("Error during gRPC dial: %v", err))
			return
		}
		defer conn.Close()

		// create client
		client := proto.NewAuthServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// call AuthorizeUser method in the auth service
		response, err := client.AuthenticateUser(ctx, &proto.AuthenticateRequest{
			AccessToken: accessToken,
		})

		if err != nil {
			slog.Error("Error in AuthenticateUser method", "error", err)
			respondWithError(w, 403, fmt.Sprintf("Error in AuthenticateUser method: %v", err))
			return
		}

		// If the access token exipred and was refreshed during AuthenticateUser, set the new access token in a cookie
		if response.IsRefreshed {
			accessToken = response.AccessToken

			http.SetCookie(w, &http.Cookie{
				Name:     "spotify_access_token",
				Value:    accessToken,
				HttpOnly: true, // Prevent JS access to the cookie
				Secure:   true, // Only send cookie over HTTPS (important for production)
				SameSite: http.SameSiteStrictMode,
				Path:     "/",
				MaxAge:   7200, // Expiry time (2 hour)
			})
		}

		handler(w, r, accessToken, response.UserID)
	}
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/akimdev15/melongo/broker/internal/auth"
	"github.com/akimdev15/melongo/broker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authHandler func(http.ResponseWriter, *http.Request, string, string)

func middlewareAuth(handler authHandler) http.HandlerFunc {
	// Creating a anonymous function
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := auth.GetAccessToken(r.Header)
		if err != nil {
			respondWithError(w, 403, "API key not found")
			return
		}

		conn, err := grpc.Dial("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Error during gRPC dial: %v", err))
			return
		}
		defer conn.Close()

		// create client
		client := proto.NewAuthServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// call AuthorizeUser method in the auth service
		token, err := client.AuthenticateUser(ctx, &proto.AuthenticateRequest{
			AccessToken: accessToken,
		})

		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Error in AuthenticateUser method: %v", err))
			return
		}

		handler(w, r, token.AccessToken, token.UserID)
	}
}

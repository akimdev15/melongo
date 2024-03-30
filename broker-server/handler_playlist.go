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

type CreatePlaylistPayload struct {
	SpotifyPlaylistID string `json:"spotifyPlaylistID"`
	ExternalUrl       string `json:"externalUrl"`
	Name              string `json:"name"`
}

func handleCreatePlaylist(w http.ResponseWriter, r *http.Request, accessToken string) {
	fmt.Println("HandleCreatePlaylist")

	apiKey, err := auth.GetAPIKey(r.Header)
	// connect to server
	conn, err := grpc.Dial("localhost:50002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		fmt.Println("Error during gRPC dial")
		return
	}
	defer conn.Close()

	// create client
	client := proto.NewPlaylistServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// call AuthorizeUser method in the auth service
	playlistResponse, err := client.CreatePlaylist(ctx, &proto.CreatePlaylistRequest{
		ApiKey:      apiKey,
		AccessToken: accessToken,
	})
	if err != nil {
		// TODO - need to return errorJSON
		fmt.Printf("Error in AuthorizeUser method: %v\n", err)
		return
	}
	fmt.Println("Playlist Response: ", playlistResponse)

	var responsePayload CreatePlaylistPayload
	responsePayload.SpotifyPlaylistID = playlistResponse.SpotifyPlaylistID
	responsePayload.ExternalUrl = playlistResponse.ExternalURL
	responsePayload.Name = playlistResponse.Name

	err = writeJSON(w, http.StatusOK, responsePayload)
	if err != nil {
		fmt.Println("Error writing JSON")
		return
	}
}

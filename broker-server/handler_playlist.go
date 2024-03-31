package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/akimdev15/melongo/broker/internal/auth"
	"github.com/akimdev15/melongo/broker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CreatePlaylistResponse struct {
	SpotifyPlaylistID string `json:"spotifyPlaylistID"`
	ExternalUrl       string `json:"externalUrl"`
	Name              string `json:"name"`
}

type CreatePlaylistPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

func handleCreatePlaylist(w http.ResponseWriter, r *http.Request, accessToken string, userID string) {
	fmt.Println("HandleCreatePlaylist")

	apiKey, err := auth.GetAPIKey(r.Header)
	// connect to server
	conn, err := grpc.Dial("localhost:50002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		fmt.Println("Error during gRPC dial")
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	// create client
	client := proto.NewPlaylistServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var payload CreatePlaylistPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Println("Error decoding payload")
		return
	}

	playlistResponse, err := client.CreatePlaylist(ctx, &proto.CreatePlaylistRequest{
		ApiKey:       apiKey,
		AccessToken:  accessToken,
		UserID:       userID,
		PlaylistName: payload.Name,
		Description:  payload.Description,
		IsPublic:     payload.Public,
	})

	// TODO -> ERROR HERE
	if err != nil {
		// TODO - need to return errorJSON
		fmt.Printf("Error in AuthorizeUser method: %v\n", err)
		return
	}
	fmt.Println("Playlist Response: ", playlistResponse)

	var responsePayload CreatePlaylistResponse
	responsePayload.SpotifyPlaylistID = playlistResponse.SpotifyPlaylistID
	responsePayload.ExternalUrl = playlistResponse.ExternalURL
	responsePayload.Name = playlistResponse.Name

	err = writeJSON(w, http.StatusOK, responsePayload)
	if err != nil {
		fmt.Println("Error writing JSON")
		return
	}
}

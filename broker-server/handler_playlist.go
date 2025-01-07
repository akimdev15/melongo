package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

type MelonTop100Request struct {
	PlaylistID string `json:"playlistID"`
	Date       string `json:"date"`
}

type MelonTop100Response struct {
	Status string `json:"status"`
}

type ResolveMissedTracksRequest struct {
	ResolvedTracks []*proto.ResolvedTrack `json:"resolvedTracks"`
}

func handleCreatePlaylist(w http.ResponseWriter, r *http.Request, accessToken string, userID string) {
	fmt.Println("HandleCreatePlaylist")
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
		AccessToken:  accessToken,
		UserID:       userID,
		PlaylistName: payload.Name,
		Description:  payload.Description,
		IsPublic:     payload.Public,
	})

	// TODO -> ERROR HERE
	if err != nil {
		// TODO - need to return errorJSON
		fmt.Printf("Error creating playlist: %v\n", err)
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

func handleMelonTop100(w http.ResponseWriter, r *http.Request, accessToken string, userID string) {
	fmt.Println("Hit HandleMelonTop100")

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

	var payload MelonTop100Request
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Println("Error decoding payload")
		return
	}

	response, err := client.CreateMelonTop100(ctx, &proto.CreateMelonTop100Request{
		AccessToken: accessToken,
		UserID:      userID,
		PlaylistID:  payload.PlaylistID,
		Date:        payload.Date,
	})

	if err != nil {
		// TODO - need to return errorJSON
		fmt.Printf("Error creating playlist: %v\n", err)
		return
	}
	fmt.Println("Playlist Response: ", response)

	var responsePayload MelonTop100Response
	responsePayload.Status = response.Status

	err = writeJSON(w, http.StatusOK, responsePayload)
	if err != nil {
		fmt.Println("Error writing JSON")
		return
	}
}

func handleSaveMelonTop100DB(w http.ResponseWriter, r *http.Request, accessToken string, userID string) {
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

	response, err := client.SaveMelonTop100DB(ctx, &proto.SaveMelonTop100DBRequest{
		AccessToken: accessToken,
	})

	if err != nil {
		fmt.Printf("Error in handleSaveMelonTop100DB: %v\n", err)
		return
	}

	var responsePayload MelonTop100Response
	responsePayload.Status = response.Status

	err = writeJSON(w, http.StatusOK, responsePayload)
	if err != nil {
		fmt.Println("Error writing JSON")
		return
	}
}

func handleGetMissedTracks(w http.ResponseWriter, r *http.Request, accessToken string, userID string) {
	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, "Missing date parameter", http.StatusBadRequest)
		return
	}

	conn, client, ctx, cancel, err := connectToGRPCServer("localhost:50002")
	if err != nil {
		fmt.Println("Error during gRPC connection setup:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("[handleGetMissedTracks] - Error closing connection")
		}
	}(conn)

	defer cancel()

	response, err := client.GetMissedTracks(ctx, &proto.GetMissedTracksRequest{
		AccessToken: accessToken,
		Date:        date,
	})

	if err != nil {
		fmt.Printf("Error in handleGetMissedTracks: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = writeJSON(w, http.StatusOK, response)
	if err != nil {
		fmt.Println("Error writing JSON")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func handleResolveMissedTracks(w http.ResponseWriter, r *http.Request, accessToken string, userID string) {

	conn, client, ctx, cancel, err := connectToGRPCServer("localhost:50002")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("[handleResolveMissedTracks] - Error closing connection")
		}
	}(conn)

	defer cancel()

	var requestPayload ResolveMissedTracksRequest
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response, err := client.ResolveMissedTracks(ctx, &proto.ResolveMissedTracksRequest{
		AccessToken:    accessToken,
		ResolvedTracks: requestPayload.ResolvedTracks,
	})

	if err != nil {
		fmt.Printf("Error in handleResolveMissedTracks: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = writeJSON(w, http.StatusOK, response)
	if err != nil {
		fmt.Println("Error writing JSON")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// ---------- HELPER FUNCTIONS ----------
func connectToGRPCServer(address string) (*grpc.ClientConn, proto.PlaylistServiceClient, context.Context, context.CancelFunc, error) {
	// Connect to the gRPC server
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	// Create the gRPC client
	client := proto.NewPlaylistServiceClient(conn)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return conn, client, ctx, cancel, nil
}

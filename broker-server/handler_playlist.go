package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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
		slog.Error("Error during gRPC dial", "error", err)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			slog.Error("Error closing connection", "error", err)
		}
	}(conn)

	// create client
	client := proto.NewPlaylistServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var payload CreatePlaylistPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		slog.Error("Error decoding payload", "error", err)
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
		slog.Error("Error creating playlist", "error", err)
		return
	}

	slog.Info("Playlist Response: ", "response", playlistResponse)

	var responsePayload CreatePlaylistResponse
	responsePayload.SpotifyPlaylistID = playlistResponse.SpotifyPlaylistID
	responsePayload.ExternalUrl = playlistResponse.ExternalURL
	responsePayload.Name = playlistResponse.Name

	err = writeJSON(w, http.StatusOK, responsePayload)
	if err != nil {
		slog.Error("Error writing JSON", "error", err)
		return
	}
}

func handleMelonTop100(w http.ResponseWriter, r *http.Request, accessToken string, userID string) {
	// connect to server
	conn, err := grpc.Dial("localhost:50002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		slog.Error("Error during gRPC dial", "error", err)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			slog.Error("Error closing connection", "error", err)
		}
	}(conn)

	// create client
	client := proto.NewPlaylistServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var payload MelonTop100Request
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		slog.Error("Error decoding payload", "error", err)
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
		slog.Error("Error creating playlist", "error", err)
		return
	}
	slog.Info("Melon Top 100 Response: ", "response", response)

	var responsePayload MelonTop100Response
	responsePayload.Status = response.Status

	err = writeJSON(w, http.StatusOK, responsePayload)
	if err != nil {
		slog.Error("Error writing JSON", "error", err)
		return
	}
}

func handleSaveMelonTop100DB(w http.ResponseWriter, r *http.Request, accessToken string, userID string) {
	// connect to server
	conn, err := grpc.Dial("localhost:50002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		slog.Error("Error during gRPC dial", "error", err)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			slog.Error("Error closing connection", "error", err)
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
		slog.Error("Error in handleSaveMelonTop100DB", "error", err)
		return
	}

	var responsePayload MelonTop100Response
	responsePayload.Status = response.Status

	err = writeJSON(w, http.StatusOK, responsePayload)
	if err != nil {
		slog.Error("Error writing JSON", "error", err)
		return
	}
}

func handleGetMissedTracks(w http.ResponseWriter, r *http.Request, accessToken string, userID string) {
	date := r.URL.Query().Get("date")
	if date == "" {
		slog.Error("Missing date parameter")
		http.Error(w, "Missing date parameter", http.StatusBadRequest)
		return
	}

	conn, client, ctx, cancel, err := connectToGRPCServer("localhost:50002")
	if err != nil {
		slog.Error("Error during gRPC connection setup", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			slog.Error("Error closing connection", "error", err)
		}
	}(conn)

	defer cancel()

	response, err := client.GetMissedTracks(ctx, &proto.GetMissedTracksRequest{
		AccessToken: accessToken,
		Date:        date,
	})

	if err != nil {
		slog.Error("Error in handleGetMissedTracks", "error", err)
		fmt.Printf("Error in handleGetMissedTracks: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = writeJSON(w, http.StatusOK, response)
	if err != nil {
		slog.Error("Error writing JSON", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func handleResolveMissedTracks(w http.ResponseWriter, r *http.Request, accessToken string, userID string) {

	conn, client, ctx, cancel, err := connectToGRPCServer("localhost:50002")
	if err != nil {
		slog.Error("Error during gRPC connection setup", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			slog.Error("Error closing connection", "error", err)
			fmt.Println("[handleResolveMissedTracks] - Error closing connection")
		}
	}(conn)

	defer cancel()

	var requestPayload ResolveMissedTracksRequest
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		slog.Error("Error decoding payload", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response, err := client.ResolveMissedTracks(ctx, &proto.ResolveMissedTracksRequest{
		AccessToken:    accessToken,
		ResolvedTracks: requestPayload.ResolvedTracks,
	})

	if err != nil {
		slog.Error("Error in handleResolveMissedTracks", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = writeJSON(w, http.StatusOK, response)
	if err != nil {
		slog.Error("Error writing JSON", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func handleGetPlaylists(w http.ResponseWriter, r *http.Request, accessToken string, userID string) {
	conn, client, ctx, cancel, err := connectToGRPCServer("localhost:50002")
	if err != nil {
		slog.Error("Error during gRPC connection setup", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			slog.Error("Error closing connection", "error", err)
		}
	}(conn)

	defer cancel()

	response, err := client.GetUserPlaylists(ctx, &proto.GetUserPlaylistsRequest{
		AccessToken: accessToken,
	})

	if err != nil {
		slog.Error("Error in handleGetPlaylists", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = writeJSON(w, http.StatusOK, response)
	if err != nil {
		slog.Error("Error writing JSON", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func handleGetPlaylistTracks(w http.ResponseWriter, r *http.Request, accessToken string, userID string) {
	tracksEnpoint := r.URL.Query().Get("endpoint")
	if tracksEnpoint == "" {
		slog.Error("Missing tracksEnpoint parameter")
		http.Error(w, "Missing tracksEnpoint parameter", http.StatusBadRequest)
		return
	}

	conn, client, ctx, cancel, err := connectToGRPCServer("localhost:50002")
	if err != nil {
		slog.Error("Error during gRPC connection setup", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			slog.Error("Error closing connection", "error", err)
		}
	}(conn)

	defer cancel()

	response, err := client.GetUserPlaylistTracks(ctx, &proto.GetUserPlaylistTracksRequest{
		AccessToken:    accessToken,
		TracksEndpoint: tracksEnpoint,
	})

	if err != nil {
		slog.Error("Error in handleGetPlaylistTracks", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = writeJSON(w, http.StatusOK, response)
	if err != nil {
		slog.Error("Error writing JSON", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// ---------- HELPER FUNCTIONS ----------
func connectToGRPCServer(address string) (*grpc.ClientConn, proto.PlaylistServiceClient, context.Context, context.CancelFunc, error) {
	// Connect to the gRPC server
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		slog.Error("Error during gRPC dial", "error", err)
		return nil, nil, nil, nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	// Create the gRPC client
	client := proto.NewPlaylistServiceClient(conn)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return conn, client, ctx, cancel, nil
}

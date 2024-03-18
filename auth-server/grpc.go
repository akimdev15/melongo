package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/akimdev15/melongo/auth-server/auth"
	"github.com/akimdev15/melongo/auth-server/internal/database"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	DB *database.Queries
}

const gRPCPORT = "50001"

func (authServer *AuthServer) AuthorizeUser(ctx context.Context, req *auth.AuthCallbackRequest) (*auth.AuthCallbackResponse, error) {
	code := req.GetCode()

	token, err := exchangeCodeForToken(ctx, code)
	if err != nil {
		fmt.Println("Error exchanging code for token:", err)
		return nil, err
	}

	// TODO Step 2: Make request to get user's info and create user
	userParams, err := createUserParams(token.AccessToken)
	if err != nil {
		// TODO - need to handle error
		fmt.Printf("Failed to create user parameter %v\n", err)
	}

	// Save to the database
	dbUser, err := authServer.DB.CreateUser(ctx, *userParams)
	if err != nil {
		// TODO - need to handle error here
		fmt.Printf("Error creating the user %v.\n", dbUser)
	}

	user := databaseUserToUser(dbUser)

	fmt.Printf("Successfully saved user to the database: %v\n", user)

	res := &auth.AuthCallbackResponse{
		ApiKey: user.APIKey,
		Name:   user.Name,
	}

	return res, nil
}

func (app *apiConfig) grpcListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPORT))
	if err != nil {
		log.Fatal("Failed to listen for grpc %v", err)
	}

	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, &AuthServer{DB: app.DB})
	log.Printf("gRPC Server started on port %s\n", gRPCPORT)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for grpc %v", err)
	}
}

func exchangeCodeForToken(ctx context.Context, code string) (*TokenResponse, error) {
	// Prepare request body
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", RedirectURI)

	// Prepare request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON response
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// createUserParams constructs parameter to save to Users database
// populates fileds by calling the spotify api
func createUserParams(accessToken string) (*database.CreateUserParams, error) {
	// Step 1: Get user's info from the Spotify API
	// Prepare request
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON response received from spotify
	var userInfoResponse struct {
		Email       string `json:"email"`
		UserID      string `json:"id"`
		DisplayName string `json:"display_name"`
	}

	err = json.Unmarshal(body, &userInfoResponse)

	if err != nil {
		fmt.Println("Error parsing JSON %s", err)
		return nil, err
	}

	fmt.Printf("User received; userID: %s, and user email: %s, displayName: %s\n",
		userInfoResponse.UserID,
		userInfoResponse.Email,
		userInfoResponse.DisplayName)

	createUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      userInfoResponse.DisplayName,
		UserID:    userInfoResponse.UserID,
		Email:     userInfoResponse.Email,
	}
	return &createUserParams, nil

}

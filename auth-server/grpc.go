package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/akimdev15/melongo/auth-server/internal/database"
	"github.com/akimdev15/melongo/auth-server/proto"
	"github.com/google/uuid"
)

type AuthServer struct {
	proto.UnimplementedAuthServiceServer
	DB     *database.Queries
	DBConn *sql.DB
}

// Parse JSON response received from spotify
type UserInfoResponse struct {
	UserID      string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

const gRPCPORT = "50001"

func (app *apiConfig) grpcListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPORT))
	if err != nil {
		log.Fatalf("Failed to listen for grpc %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterAuthServiceServer(grpcServer, &AuthServer{DB: app.DB, DBConn: app.DBConn})
	log.Printf("gRPC Server started on port %s\n", gRPCPORT)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for grpc %v", err)
	}
}

func (authServer *AuthServer) AuthenticateUser(ctx context.Context, req *proto.AuthenticateRequest) (*proto.AuthenticateResponse, error) {
	apiKey := req.GetApiKey()
	if apiKey == "" {
		return nil, errors.New("API_KEY is empty")
	}

	userToken, err := authServer.DB.GetUserTokenByAPIKey(ctx, apiKey)
	if err != nil {
		fmt.Printf("Error getting user token from the DB. error: %v\n", err)
		return nil, err
	}

	// TODO - put in a logic for refresh token

	res := &proto.AuthenticateResponse{
		AccessToken: userToken.AccessToken,
	}

	return res, nil
}

// AuthorizeUser handles authorization callback from spotify.
// Handles token generation and save/get user
// For initial login to the website
func (authServer *AuthServer) AuthorizeUser(ctx context.Context, req *proto.AuthCallbackRequest) (*proto.AuthCallbackResponse, error) {
	code := req.GetCode()
	fmt.Println("Received gRPC call from broker-service")

	token, err := exchangeCodeForToken(ctx, code)
	if err != nil {
		fmt.Println("Error exchanging code for token:", err)
		return nil, err
	}

	// TODO Step 2: Make request to get user's info and create user
	userInfoResponse, err := getUserFromSpotify(token.AccessToken)
	if err != nil {
		// TODO - need to handle error
		fmt.Println("Failed to get userResponse.")
	}
	fmt.Printf("UserInfoResponse: %v\n", userInfoResponse)

	// Check if user already exists
	savedUser, err := authServer.DB.GetUserById(ctx, userInfoResponse.UserID)
	fmt.Printf("Existing user: %v", savedUser)
	if err == sql.ErrNoRows {
		// Use transaction
		tx, err := authServer.DBConn.Begin()
		if err != nil {
			fmt.Println("Failed creating DB transaction")
			return nil, err
		}
		defer tx.Rollback()
		qtx := authServer.DB.WithTx(tx)

		savedUser, err = authServer.createAndSaveUser(ctx, userInfoResponse, qtx)
		if err != nil {
			// TODO - need to handle error
			fmt.Printf("Failed to create user. err: %v\n", err)
			return nil, err
		}

		// save the token in the user_tokens db
		err = authServer.createUserToken(ctx, token, savedUser.ApiKey, qtx)
		if err != nil {
			fmt.Printf("Failed to save token to db for the user: %s\n. err: %v\n", savedUser.ID, err)
			return nil, err
		}

		// transaction succeeded
		tx.Commit()
	} else if err != nil {
		fmt.Printf("Error checking for existing user. err: %v\n", err)
		return nil, err
	}

	res := &proto.AuthCallbackResponse{
		ApiKey: savedUser.ApiKey,
		Name:   savedUser.Name,
	}

	return res, nil
}

// createUserToken saves the token information in the DB
func (authServer *AuthServer) createUserToken(ctx context.Context, token *TokenResponse, apiKey string, qtx *database.Queries) error {
	expireSecond := int64(token.Expires_In)
	expirationTime := time.Now().Add(time.Second * time.Duration(expireSecond))
	_, err := qtx.CreateUserToken(ctx, database.CreateUserTokenParams{

		ID:           uuid.New(),
		ApiKey:       apiKey,
		AccessToken:  token.AccessToken,
		RefreshToken: token.Refresh_Token,
		ExpireTime:   expirationTime,
	})

	return err
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

func getUserFromSpotify(accessToken string) (UserInfoResponse, error) {
	// Step 1: Get user's info from the Spotify API
	// Prepare request
	userInfoResponse := UserInfoResponse{}
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return userInfoResponse, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return userInfoResponse, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return userInfoResponse, err
	}

	err = json.Unmarshal(body, &userInfoResponse)

	if err != nil {
		fmt.Printf("Error parsing JSON %s\n", err)
		return userInfoResponse, err
	}

	return userInfoResponse, err
}

func (authServer *AuthServer) createAndSaveUser(ctx context.Context, userInfoResponse UserInfoResponse, qtx *database.Queries) (database.User, error) {

	createUserParams := database.CreateUserParams{
		ID:        userInfoResponse.UserID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      userInfoResponse.DisplayName,
		Email:     userInfoResponse.Email,
	}

	// Save to the database
	dbUser, err := qtx.CreateUser(ctx, createUserParams)
	if err != nil {
		// TODO - need to handle error here
		fmt.Printf("Error creating the user %v.\n", dbUser)
		return database.User{}, err
	}

	fmt.Printf("Successfully saved user to the database: %v\n", dbUser)
	return dbUser, nil

}

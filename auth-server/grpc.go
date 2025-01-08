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
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/akimdev15/melongo/auth-server/internal/database"
	"github.com/akimdev15/melongo/auth-server/proto"
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

// Define a struct to represent the response from the Spotify token endpoint
type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

const gRPCPORT = "50001"

func (app *apiConfig) grpcListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPORT))
	if err != nil {
		slog.Error("Failed to listen for grpc", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterAuthServiceServer(grpcServer, &AuthServer{DB: app.DB, DBConn: app.DBConn})
	log.Printf("gRPC Server started on port %s\n", gRPCPORT)
	if err = grpcServer.Serve(lis); err != nil {
		slog.Error("Failed to listen for grpc", "error", err)
		os.Exit(1)
	}
}

func (authServer *AuthServer) AuthenticateUser(ctx context.Context, req *proto.AuthenticateRequest) (*proto.AuthenticateResponse, error) {
	accessToken := req.GetAccessToken()
	if accessToken == "" {
		return nil, errors.New("AccessToken is empty")
	}

	userToken, err := authServer.DB.GetUserTokenByAccessToken(ctx, accessToken)
	if err != nil {
		slog.Error("Error getting user token from the DB", "error", err)
		return nil, err
	}

	tokenRefreshed := false
	if userToken.ExpireTime.Before(time.Now().UTC()) {
		slog.Info("Access-Token expired. Getting a new token...")
		refreshToken, err := RefreshToken(userToken.RefreshToken, ctx)
		if err != nil {
			slog.Error("Error refreshing token", "error", err)
			return nil, err
		}

		// Update user token with the refresh token
		newExpireTime := time.Now().UTC().Add(time.Duration(refreshToken.ExpiresIn) * time.Second)
		userToken.AccessToken = refreshToken.AccessToken

		// Refresh token returned from the response might be empty. In this case use the existing refresh token (According to the documentation)
		if refreshToken.RefreshToken == "" {
			refreshToken.RefreshToken = userToken.RefreshToken
		}

		// Asynchronously update to the DB
		errCh := make(chan error)
		go func() {
			err = authServer.DB.UpdateToken(ctx, database.UpdateTokenParams{
				AccessToken:  refreshToken.AccessToken,
				RefreshToken: refreshToken.RefreshToken,
				ExpireTime:   newExpireTime,
				UpdatedAt:    time.Now(),
				ID:           userToken.ID,
			})
			// send err or nil if successful to the channel
			errCh <- err
		}()

		err = <-errCh

		if err != nil {
			slog.Error("Error updating the user token with the refresh token", "error", err)
			return nil, err
		}

		slog.Info("Successfully updated the refresh token")
		tokenRefreshed = true
	}

	res := &proto.AuthenticateResponse{
		AccessToken: userToken.AccessToken,
		UserID:      userToken.ID,
		IsRefreshed: tokenRefreshed,
	}

	return res, nil
}

// AuthorizeUser handles authorization callback from spotify.
// Handles token generation and save/get user
// For initial login to the website
func (authServer *AuthServer) AuthorizeUser(ctx context.Context, req *proto.AuthCallbackRequest) (*proto.AuthCallbackResponse, error) {
	code := req.GetCode()

	token, err := exchangeCodeForToken(ctx, code)
	if err != nil {
		slog.Error("Failed to exchange code for token", "error", err)
		return nil, err
	}

	// TODO Step 2: Make request to get user's info and create user
	userInfoResponse, err := getUserFromSpotify(token.AccessToken)
	if err != nil {
		// TODO - need to handle error
		slog.Error("Failed to get userResponse.", "error", err)
		return nil, err
	}

	// Check if user already exists
	savedUser, err := authServer.DB.GetUserById(ctx, userInfoResponse.UserID)
	slog.Info("[AuthorizeUser] - Saved user", "user", savedUser)
	if err == sql.ErrNoRows {
		// Use transaction
		tx, err := authServer.DBConn.Begin()
		if err != nil {
			slog.Error("Failed creating DB transaction", "error", err)
			return nil, err
		}
		defer tx.Rollback()
		qtx := authServer.DB.WithTx(tx)

		savedUser, err = authServer.createAndSaveUser(ctx, userInfoResponse, qtx)
		if err != nil {
			// TODO - need to handle error
			slog.Error("Failed to save user to the database", "error", err)
			return nil, err
		}

		// save the token in the user_tokens db
		err = authServer.createUserToken(ctx, token, savedUser.ApiKey, savedUser.ID, qtx)
		if err != nil {
			slog.Error("Failed to save token to the database", "error", err)
			return nil, err
		}

		// transaction succeeded
		tx.Commit()
	} else if err != nil {
		slog.Error("Error checking for existing user", "error", err)
		return nil, err
	} else {
		// Should update the token info in the DB for the existing user with new token
		err = authServer.DB.UpdateToken(ctx, database.UpdateTokenParams{
			AccessToken:  token.AccessToken,
			RefreshToken: token.Refresh_Token,
			ExpireTime:   time.Now().UTC().Add(time.Duration(token.Expires_In) * time.Second),
			UpdatedAt:    time.Now(),
			ID:           savedUser.ID,
		})
		if err != nil {
			slog.Error("Failed to update token in the database", "error", err)
			return nil, err
		}
	}

	res := &proto.AuthCallbackResponse{
		AccessToken: token.AccessToken,
		Name:        savedUser.Name,
	}

	return res, nil
}

// createUserToken saves the token information in the DB
func (authServer *AuthServer) createUserToken(ctx context.Context, token *TokenResponse, apiKey string, userID string, qtx *database.Queries) error {
	expireSecond := int64(token.Expires_In)
	expirationTime := time.Now().UTC().Add(time.Second * time.Duration(expireSecond))
	_, err := qtx.CreateUserToken(ctx, database.CreateUserTokenParams{
		ID:           userID,
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
		slog.Error("[exchangeCodeForToken] - Failed to create request", "error", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("[exchangeCodeForToken] - Failed to send request", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		slog.Error("[exchangeCodeForToken] - Unexpected status code", "status", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("[exchangeCodeForToken] - Failed to read response body", "error", err)
		return nil, err
	}

	// Parse JSON response
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		slog.Error("[exchangeCodeForToken] - Failed to parse JSON response", "error", err)
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
		slog.Error("[getUserFromSpotify] - Failed to create request", "error", err)
		return userInfoResponse, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("[getUserFromSpotify] - Failed to send request", "error", err)
		return userInfoResponse, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("[getUserFromSpotify] - Failed to read response body", "error", err)
		return userInfoResponse, err
	}

	err = json.Unmarshal(body, &userInfoResponse)

	if err != nil {
		slog.Error("[getUserFromSpotify] - Failed to parse JSON response", "error", err)
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
		slog.Error("Failed to save user to the database", "error", err)
		return database.User{}, err
	}

	slog.Info("Successfully saved user to the database", "user", dbUser)
	return dbUser, nil

}

func RefreshToken(refreshToken string, ctx context.Context) (SpotifyTokenResponse, error) {
	// Construct the request body
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	// Prepare request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return SpotifyTokenResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))

	// Send request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to send request", "error", err)
		return SpotifyTokenResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			err := resp.Body.Close()
			if err != nil {
				slog.Error("Failed to close response body", "error", err)
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Unexpected status code", "status", resp.StatusCode)
		return SpotifyTokenResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	//// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body", "error", err)
		return SpotifyTokenResponse{}, err
	}

	var tokenResponse SpotifyTokenResponse
	err = json.Unmarshal(body, &tokenResponse)

	if err != nil {
		slog.Error("Error parsing JSON", "error", err)
		return SpotifyTokenResponse{}, err
	}

	return tokenResponse, nil
}

package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/akimdev15/melongo/auth-server/internal/database"
	"github.com/google/uuid"
)

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	// Construct Authorization URL
	authURL := fmt.Sprintf(
		"https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s",
		clientID, RedirectURI, Scopes,
	)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (apiCfg *apiConfig) handleAuthorizationResponse(w http.ResponseWriter, r *http.Request) {
	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Extract authorization code from query parameters
	code := r.URL.Query().Get("code")

	token, err := exchangeCodeForToken(ctx, code)
	if err != nil {
		fmt.Println("Error exchanging code for token:", err)
		return
	}

	// TODO Step 2: Make request to get user's info and create user
	userParams, err := createUserParams(token.AccessToken)
	if err != nil {
		// TODO - need to handle error
		fmt.Printf("Failed to create user parameter %v\n", err)
	}

	// Save to the database
	dbUser, err := apiCfg.DB.CreateUser(r.Context(), *userParams)
	if err != nil {
		// TODO - need to handle error here
		fmt.Printf("Error creating the user %v.\n", dbUser)
	}

	user := databaseUserToUser(dbUser)

	fmt.Printf("Successfully saved user to the database: %v\n", user)
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

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
)

const (
	ClientID     = "60e1ce49166b49ebb3c2999beabe8ac5"
	RedirectURI  = "http://localhost:8080/callback"
	Scopes       = "user-read-email user-read-private playlist-modify-public playlist-modify-private playlist-read-collaborative playlist-read-private"
	ClientSecret = "9938440d200544f7beb21db96d29161e"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func main() {
	// Step 1: Construct Authorization URL
	authURL := fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s", ClientID, RedirectURI, Scopes)

	// Step 2: Redirect User to Spotify Authorization Page
	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	})

	// Step 3: Handle Authorization Response
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Create a context with a timeout of 10 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// Extract authorization code from query parameters
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		fmt.Println("State: ", state)
		fmt.Println("Authorization Code:", code)

		token, err := exchangeCodeForToken(ctx, code)
		if err != nil {
			fmt.Println("Error exchanging code for token:", err)
			return
		}

		// Step 2: Make a request to get user's playlists
		playlists, err := getUserPlaylists(token.AccessToken)
		if err != nil {
			fmt.Println("Error getting user playlists:", err)
			return
		}

		fmt.Println("Playlist OBJ: ", playlists)

		// Print playlists
		fmt.Println("User Playlists:")
		for _, playlist := range playlists {
			fmt.Println(playlist.Name)
		}
	})

	// Start server
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
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
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(ClientID+":"+ClientSecret)))
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

func getUserPlaylists(accessToken string) ([]Playlist, error) {
	// Prepare request
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/playlists", nil)
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

	// Parse JSON response
	var playlistsResponse struct {
		Items []Playlist `json:"items"`
	}
	err = json.Unmarshal(body, &playlistsResponse)
	if err != nil {
		return nil, err
	}

	return playlistsResponse.Items, nil
}

type Playlist struct {
	Name string `json:"name"`
	// Add other playlist fields you want to extract
}

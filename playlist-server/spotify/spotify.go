package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Playlist struct {
	Name string `json:"name"`
	// Add other playlist fields you want to extract
}

type ExternalURLs struct {
	Spotify string `json:"spotify"`
}

type Artist struct {
	ExternalURLs ExternalURLs `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

type ArtistItem struct {
	ID string `json:"id"`
	// Add other artist fields you want to extract
}

type ArtistsResponse struct {
	Artists struct {
		Items []ArtistItem `json:"items"`
	} `json:"artists"`
}

type Album struct {
	Artists []Artist `json:"artists"`
}

type Track struct {
	Album Album  `json:"album"`
	URI   string `json:"uri"`
	Name  string `json:"name"`
}

type TracksResponse struct {
	Tracks []Track `json:"tracks"`
}

// GetUserPlaylists gets all the user's playlist names
func GetUserPlaylists(accessToken string) ([]Playlist, error) {
	// Prepare request
	address := "https://api.spotify.com/v1/me/playlists"
	body, err := makeSpotifyRequest(address, accessToken)
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

func SearchArtistID(artistName string, accessToken string) (string, error) {
	encodedArtistName := url.QueryEscape(artistName)
	// Construct the search query for the artist
	address := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=artist&locale=ko_KR", encodedArtistName)

	body, err := makeSpotifyRequest(address, accessToken)
	if err != nil {
		fmt.Println("Error making the request to the spotify with the url: ", address)
		return "", err
	}

	// Unmarshal JSON data into TracksResponse struct
	var response ArtistsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	// Extract the Spotify ID of the first artist
	if len(response.Artists.Items) <= 0 {
		fmt.Println("Couldn't search for an Artist ID. response.Artist.Items is empty")
		return "", nil
	}

	artistID := response.Artists.Items[0].ID
	return artistID, nil
}

func SearchTracksByArtist(artistID string, accessToken string) ([]Track, error) {
	address := fmt.Sprintf("https://api.spotify.com/v1/artists/%s/top-tracks?country=KR", artistID)
	body, err := makeSpotifyRequest(address, accessToken)
	if err != nil {
		fmt.Println("Error making the request to the spotify with the url: ", address)
		return nil, err
	}
	var tracksResponse TracksResponse
	if err := json.Unmarshal(body, &tracksResponse); err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return tracksResponse.Tracks, nil
}

// searchTrack looks up a music by title and artist
func SearchTrack(title, artist, accessToken string) (TracksResponse, error) {
	encodedTitle := url.QueryEscape(title)
	encodedArtist := url.QueryEscape(artist)
	address := fmt.Sprintf("https://api.spotify.com/v1/search?q=track:%s+artist:%s&type=track", encodedTitle, encodedArtist)

	req, err := http.NewRequest("GET", address, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	if err != nil {
		fmt.Printf("Error creating the request. Error: %v\n", err)
		return TracksResponse{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error response. err: ", err)
		return TracksResponse{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error: ", err)
		}
	}(resp.Body)

	var bodyReader = resp.Body

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code:", resp.Status)
		return TracksResponse{}, err
	}

	// Read the response body
	body, err := io.ReadAll(bodyReader)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return TracksResponse{}, err
	}

	// Unmarshal JSON data into TracksResponse struct
	var response TracksResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error:", err)
		return TracksResponse{}, err
	}

	fmt.Println("Successfully searched the track")
	return response, nil
}

// MakeSpotifyRequest - With the given address, make a GET request to spotify
// returns response
func makeSpotifyRequest(address string, accessToken string) ([]byte, error) {
	// Create request
	req, err := http.NewRequest("GET", address, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	if err != nil {
		fmt.Printf("Error creating the request. Error: %v\n", err)
		return nil, err
	}

	// make a GET request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error response. err: ", err)
		return nil, err
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code:", resp.Status)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error: ", err)
		}
	}(resp.Body)

	var bodyReader = resp.Body

	// Read the response body
	body, err := io.ReadAll(bodyReader)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return nil, err
	}

	return body, nil
}

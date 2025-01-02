package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Playlist struct {
	Name string `json:"name"`
	ID   string `json:"id"`
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

type Album struct {
	Artists []Artist `json:"artists"`
}

type Track struct {
	Artist string
	URI    string `json:"uri"`
	Name   string `json:"name"`
}

// Contains artist names in an array (in case there are more than one)
type AlbumTrack struct {
	Artist []struct {
		Name string `json:"name"`
	} `json:"artists"`
	URI  string `json:"uri"`
	Name string `json:"name"`
}

// SearchResponse - result of query search by title and artist
type SearchResponse struct {
	Tracks struct {
		Items []struct {
			Name    string `json:"name"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			URI string `json:"uri"`
		} `json:"items"`
	} `json:"tracks"`
}

type SearchResponseAlbum struct {
	Albums struct {
		Items []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			URI  string `json:"uri"`
		} `json:"items"`
	} `json:"albums"`
}

// ArtistsResponse - result of artist search
type ArtistsResponse struct {
	Artists struct {
		Items []ArtistItem `json:"items"`
	} `json:"artists"`
}

type ArtistItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Add other artist fields you want to extract
}

// END Artist Response

// ArtistTracksResponse - result of artist's top tracks
type ArtistTracksResponse struct {
	Tracks []Track `json:"tracks"`
}

// TrackResponse - result of query search (URI of a track)
type TrackResponse struct {
	Items struct {
		Tracks []Track `json:"items"`
	} `json:"tracks"`
}

// TracksResponse - result of album search by id and returns all tracks
type AlbumTracksResponse struct {
	Items []AlbumTrack `json:"items"`
}

type NewPlayListRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

// NewPlaylistResponse - result of creating new empty playlist
type NewPlaylistResponse struct {
	Description       string `json:"description"`
	SpotifyPlaylistID string `json:"id"`
	Name              string `json:"name"`
	URI               string `json:"uri"`
}

// AddTrackRequest - request to add new track(s) to the playlist
type AddTrackRequest struct {
	URIs     []string `json:"uris"`
	Position int      `json:"position"`
}

// AddTrackResponse - returns new id of the playlist
type AddTrackResponse struct {
	SnapshotID string `json:"snapshot_id"`
}

// GetUserPlaylists gets all the user's playlist names
func GetUserPlaylists(accessToken string) ([]Playlist, error) {
	// Prepare request
	address := "https://api.spotify.com/v1/me/playlists"
	body, err := makeSpotifyGetRequest(address, accessToken)
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

func SearchArtistID(artistName string, accessToken string) (ArtistItem, error) {
	encodedArtistName := url.QueryEscape(artistName)
	// Construct the search query for the artist
	address := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=artist&locale=ko_KR", encodedArtistName)

	body, err := makeSpotifyGetRequest(address, accessToken)
	if err != nil {
		fmt.Println("Error making the request to the spotify with the url: ", address)
		return ArtistItem{}, err
	}

	// Unmarshal JSON data into TracksResponse struct
	var response ArtistsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error:", err)
		return ArtistItem{}, err
	}

	// Extract the Spotify ID of the first artist
	if len(response.Artists.Items) <= 0 {
		fmt.Println("Couldn't search for an Artist ID. response.Artist.Items is empty")
		return ArtistItem{}, nil
	}

	fmt.Printf("*** Artist Items: %v\n", response.Artists.Items[0])

	artist := response.Artists.Items[0]
	return artist, nil
}

// SearchTrack looks up a music by title and artist
// returns the Track which contains the URI of the track
// which can be used to add to the playlist
func SearchTrack(title, artist, accessToken string) (*Track, error) {
	if title == "" || artist == "" || accessToken == "" {
		return nil, fmt.Errorf("title, artist, or access token is empty")
	}

	// remove any brackets in the song title
	formattedTitle := formatTitle(title)

	const spotifyAPIURL = "https://api.spotify.com/v1/search"

	// Formulate the search query
	query := fmt.Sprintf("track:%s artist:%s", formattedTitle, artist)

	// URL-encode the query string
	encodedQuery := url.QueryEscape(query)

	// Construct the search URL
	searchURL := fmt.Sprintf("%s?q=%s&type=track", spotifyAPIURL, encodedQuery)

	body, err := makeSpotifyGetRequest(searchURL, accessToken)

	if err != nil {
		fmt.Println("Error making the request to the spotify")
		return nil, err
	}

	// Parse the JSON response into the struct
	var searchResp SearchResponse
	err = json.Unmarshal(body, &searchResp)
	if err != nil {
		fmt.Printf("Error parsing the JSON response for title: %s, %v\n", title, err)
		return nil, err
	}

	// Check if any tracks were found
	if len(searchResp.Tracks.Items) == 0 {
		return nil, fmt.Errorf("no tracks found for title: %s, artist: %s", title, artist)
	}

	// Return the URI of the first matching track
	trackURI := searchResp.Tracks.Items[0].URI

	return &Track{
		Artist: artist,
		Name:   formattedTitle,
		URI:    trackURI,
	}, nil
}

func SearchTracksFromAlbum(albumName, artistName, accessToken string) ([]AlbumTrack, error) {
	if albumName == "" || artistName == "" || accessToken == "" {
		return nil, fmt.Errorf("album name, artist name, or access token is empty")
	}

	// Format the album name to remove brackets (if any)
	formattedAlbumName := formatTitle(albumName)

	const spotifyAPIURL = "https://api.spotify.com/v1/search"

	// Formulate the search query
	query := fmt.Sprintf("album:%s artist:%s", formattedAlbumName, artistName)

	// URL-encode the query string
	encodedQuery := url.QueryEscape(query)

	// Construct the search URL
	searchURL := fmt.Sprintf("%s?q=%s&type=album", spotifyAPIURL, encodedQuery)

	// Make the GET request to Spotify to search for albums
	body, err := makeSpotifyGetRequest(searchURL, accessToken)
	if err != nil {
		return nil, fmt.Errorf("error making the request to Spotify: %v", err)
	}

	// Parse the search response
	var searchResp SearchResponseAlbum
	err = json.Unmarshal(body, &searchResp)
	if err != nil {
		return nil, fmt.Errorf("error parsing the JSON response for album: %v", err)
	}

	// Check if any albums were found
	if len(searchResp.Albums.Items) == 0 {
		// TODO - For sigle album, might want to search by track since it returns the result for track but not album for some
		return nil, fmt.Errorf("no albums found for album: %s, artist: %s", albumName, artistName)
	}

	// Get the album ID from the first result
	albumID := searchResp.Albums.Items[0].ID
	fmt.Println("Album ID: ", albumID)

	url := fmt.Sprintf("https://api.spotify.com/v1/albums/%s/tracks", albumID)
	body, err = makeSpotifyGetRequest(url, accessToken)
	if err != nil {
		return nil, fmt.Errorf("error fetching album tracks: %v", err)
	}

	// Unmarshal the JSON response into a TracksResponse struct
	var albumTracksResp AlbumTracksResponse
	err = json.Unmarshal(body, &albumTracksResp)
	if err != nil {
		fmt.Println("Error parsing the JSON response for album tracks: ", err)
		return nil, err
	}

	// Check if any tracks were found
	if len(albumTracksResp.Items) == 0 {
		return nil, fmt.Errorf("no tracks found for album: %s, artist: %s", albumName, artistName)
	}

	// TODO - Maybe convert AlbumTrack to Track by concatenating the artist names with a comma separated string

	return albumTracksResp.Items, nil

}

// CreateNewPlaylist - creates a empty new playlist for the user
func CreateNewPlaylist(name string, description string, isPublic bool, userId string, accessToken string) (NewPlaylistResponse, error) {
	address := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", userId)

	playlistRequest := NewPlayListRequest{
		Name:        name,
		Description: description,
		Public:      isPublic,
	}

	body, err := json.Marshal(playlistRequest)
	if err != nil {
		fmt.Println("Error during json.Marshal of playlistRequest. Err: ", err)
		return NewPlaylistResponse{}, err
	}

	body, err = makeSpotifyPostRequest(address, body, accessToken)

	if err != nil {
		fmt.Println("Error making the request to the spotify")
		return NewPlaylistResponse{}, err
	}

	// Unmarshal JSON data into TracksResponse struct
	var response NewPlaylistResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error:", err)
		return NewPlaylistResponse{}, err
	}

	fmt.Println("Successfully searched the track")
	return response, nil
}

// AddTrackToPlaylist - adds trackURI which is comma separated track uris to the playlist
// TODO - Try to prevent duplicate tracks by checking the existing playlist or maybe store something in the database
func AddTrackToPlaylist(playlistID string, trackURI []string, accessToken string) (AddTrackResponse, error) {
	address := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)

	addTrackRequest := AddTrackRequest{
		URIs: trackURI,
	}

	body, err := json.Marshal(addTrackRequest)
	if err != nil {
		fmt.Println("Error during json.Marshal")
		return AddTrackResponse{}, err
	}

	body, err = makeSpotifyPostRequest(address, body, accessToken)
	if err != nil {
		fmt.Println("Error making the request to the spotify")
		return AddTrackResponse{}, err
	}

	var response AddTrackResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error during AddTrackResponse. Err: ", err)
		return AddTrackResponse{}, err
	}

	return response, nil
}

// makeSpotifyGetRequest - With the given address, make a GET request to spotify
// returns response
func makeSpotifyGetRequest(address string, accessToken string) ([]byte, error) {
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

// makeSpotifyPostRequest - make a post request where data is the body
func makeSpotifyPostRequest(address string, data []byte, accessToken string) ([]byte, error) {
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error creating the request. Error: ", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error reading...")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	//// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// formatTitle - remove everything inside the brackets
// ex) "Cry Me A River - Justin Timberlake (Official Music Video)" -> "Cry Me A River - Justin Timberlake"
func formatTitle(title string) string {
	if title == "" {
		return ""
	}

	// Find the first index of the bracket
	idx := bytes.IndexByte([]byte(title), '(')
	if idx == -1 {
		return title
	}

	// Remove the bracket and the content inside
	return title[:idx]
}

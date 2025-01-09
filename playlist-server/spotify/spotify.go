package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"unicode"
)

type Image struct {
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

// Playlist - for getting detailed information of a single playlist
type Playlist struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ID          string  `json:"id"`
	Images      []Image `json:"images"`
	Tracks      struct {
		Total int `json:"total"`
		Item  []struct {
			Track struct {
				Artists    []Artist `json:"artists"`
				Name       string   `json:"name"`
				Popularity int      `json:"popularity"`
				URI        string   `json:"uri"`
			} `json:"track"`
		} `json:"items"`
	} `json:"tracks"`
}

// SimplifiedPlaylist - Used when fetching all the playlists of the user
type SimplifiedPlaylist struct {
	Next  string     `json:"next"`  // URL to the next page of items
	Total int        `json:"total"` // Total number of items available
	Items []struct { // Items is now a slice of playlists
		ExternalURLs struct {
			SpotifyURL string `json:"spotify"`
		} `json:"external_urls"`
		PlaylistEndpoint string  `json:"href"` // URL endpoint to get playlist details
		Name             string  `json:"name"`
		Description      string  `json:"description"`
		ID               string  `json:"id"`
		Images           []Image `json:"images"`
		Tracks           struct {
			TracksEndpoint string `json:"href"`  // URL endpoint to get full tracks of the playlist
			Total          int    `json:"total"` // Total number of tracks in the playlist
		} `json:"tracks"`
	} `json:"items"`
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
	Artist     string
	URI        string `json:"uri"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
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
			URI        string `json:"uri"`
			Popularity int    `json:"popularity"`
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

// TODO - need  to check
func GetPlaylistDetails(accessToken string, userId string) (*Playlist, error) {
	// Prepare request
	address := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", userId)
	body, err := makeSpotifyGetRequest(address, accessToken)
	if err != nil {
		return nil, err
	}

	var playlistResponse *Playlist
	err = json.Unmarshal(body, &playlistResponse)
	if err != nil {
		return nil, err
	}

	return playlistResponse, nil
}

// GetUserPlaylists gets all the current user's playlist
func GetUserPlaylists(accessToken string) (*SimplifiedPlaylist, error) {
	// Prepare request
	address := "https://api.spotify.com/v1/me/playlists?limit=50"
	body, err := makeSpotifyGetRequest(address, accessToken)
	if err != nil {
		return nil, err
	}

	var playlistResponse *SimplifiedPlaylist
	err = json.Unmarshal(body, &playlistResponse)
	if err != nil {
		return nil, err
	}

	return playlistResponse, nil
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
		slog.Error("Couldn't search for an Artist ID. response.Artist.Items is empty")
		return ArtistItem{}, nil
	}

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

	// get english artist name
	artist = formatArtistName(artist)

	const spotifyAPIURL = "https://api.spotify.com/v1/search"

	// Formulate the search query
	query := fmt.Sprintf("track:%s artist:%s", formattedTitle, artist)

	// URL-encode the query string
	encodedQuery := url.QueryEscape(query)

	// Construct the search URL
	searchURL := fmt.Sprintf("%s?q=%s&type=track", spotifyAPIURL, encodedQuery)

	body, err := makeSpotifyGetRequest(searchURL, accessToken)

	if err != nil {
		slog.Error("Error making the request to the spotify")
		return nil, err
	}

	// Parse the JSON response into the struct
	var searchResp SearchResponse
	err = json.Unmarshal(body, &searchResp)
	if err != nil {
		slog.Error("Error parsing the JSON response for", "title: ", title, "Error: ", err)
		return nil, err
	}

	// Check if any tracks were found
	if len(searchResp.Tracks.Items) == 0 {
		slog.Info("No tracks found for", "title", title, " artist: ", artist)
		return nil, fmt.Errorf("no tracks found for title: %s, artist: %s", title, artist)
	}

	// Return the URI of the first matching track
	trackURI := searchResp.Tracks.Items[0].URI
	popularity := searchResp.Tracks.Items[0].Popularity

	return &Track{
		Artist:     artist,
		Name:       formattedTitle,
		URI:        trackURI,
		Popularity: popularity,
	}, nil
}

func SearchTracksFromAlbum(albumName, artistName, accessToken string) ([]AlbumTrack, error) {
	if albumName == "" || artistName == "" || accessToken == "" {
		slog.Error("album name, artist name, or access token is empty")
		return nil, fmt.Errorf("album name, artist name, or access token is empty")
	}

	// Format the album name to remove brackets (if any)
	formattedAlbumName := formatTitle(albumName)

	// TODO - test if this still works (Otherwise remove it)
	// get english artist name
	artistName = formatArtistName(artistName)

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
		slog.Error("Error making the request to Spotify: ", err)
		return nil, fmt.Errorf("error making the request to Spotify: %v", err)
	}

	// Parse the search response
	var searchResp SearchResponseAlbum
	err = json.Unmarshal(body, &searchResp)
	if err != nil {
		slog.Error("Error parsing the JSON response for album: ", err)
		return nil, fmt.Errorf("error parsing the JSON response for album: %v", err)
	}

	// Check if any albums were found
	if len(searchResp.Albums.Items) == 0 {
		// TODO - For sigle album, might want to search by track since it returns the result for track but not album for some
		slog.Info("No albums found for", "album", albumName, " artist: ", artistName)
		return nil, fmt.Errorf("no albums found for album: %s, artist: %s", albumName, artistName)
	}

	// Get the album ID from the first result
	albumID := searchResp.Albums.Items[0].ID

	url := fmt.Sprintf("https://api.spotify.com/v1/albums/%s/tracks", albumID)
	body, err = makeSpotifyGetRequest(url, accessToken)
	if err != nil {
		return nil, fmt.Errorf("error fetching album tracks: %v", err)
	}

	// Unmarshal the JSON response into a TracksResponse struct
	var albumTracksResp AlbumTracksResponse
	err = json.Unmarshal(body, &albumTracksResp)
	if err != nil {
		slog.Error("Error parsing the JSON response for album tracks: ", err)
		return nil, err
	}

	// Check if any tracks were found
	if len(albumTracksResp.Items) == 0 {
		slog.Info("No tracks found for", "album", albumName, " artist: ", artistName)
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
		slog.Error("Error during json.Marshal of playlistRequest. Err: ", err)
		return NewPlaylistResponse{}, err
	}

	body, err = makeSpotifyPostRequest(address, body, accessToken)

	if err != nil {
		slog.Error("Error making the request to the spotify")
		fmt.Println("Error making the request to the spotify")
		return NewPlaylistResponse{}, err
	}

	// Unmarshal JSON data into TracksResponse struct
	var response NewPlaylistResponse
	if err := json.Unmarshal(body, &response); err != nil {
		slog.Error("Error:", err)
		return NewPlaylistResponse{}, err
	}

	slog.Info("Successfully created the new playlist")
	return response, nil
}

// AddTrackToPlaylist - adds trackURI which is comma separated track uris to the playlist
// TODO - Try to prevent duplicate tracks by checking the existing playlist or maybe store something in the database
func AddTrackToPlaylist(playlistID string, trackURI []string, accessToken string) (AddTrackResponse, error) {
	address := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)

	addTrackRequest := AddTrackRequest{
		URIs: trackURI,
	}

	slog.Info("TrackURI: ", len(trackURI))

	body, err := json.Marshal(addTrackRequest)
	if err != nil {
		slog.Error("Error during json.Marshal")
		return AddTrackResponse{}, err
	}

	body, err = makeSpotifyPostRequest(address, body, accessToken)
	if err != nil {
		slog.Error("Error making the request to the spotify")
		return AddTrackResponse{}, err
	}

	var response AddTrackResponse
	if err := json.Unmarshal(body, &response); err != nil {
		slog.Error("Error during AddTrackResponse. Err: ", err)
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
		slog.Error("Error creating the request. Error: ", err)
		return nil, err
	}

	// make a GET request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error response. err: ", err)
		return nil, err
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		slog.Error("Error: Unexpected status code:", resp.Status)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Error reading...")
		}
	}(resp.Body)

	var bodyReader = resp.Body

	// Read the response body
	body, err := io.ReadAll(bodyReader)
	if err != nil {
		slog.Error("Error reading the response body:", err)
		return nil, err
	}

	return body, nil
}

// makeSpotifyPostRequest - make a post request where data is the body
func makeSpotifyPostRequest(address string, data []byte, accessToken string) ([]byte, error) {
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(data))
	if err != nil {
		slog.Error("Error creating the request. Error: ", err)
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
			slog.Error("Error reading...")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		slog.Error("Error: Unexpected", "status code", resp.Status)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	//// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading the response body:", err)
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

// formatArtistName extracts the English name from the artist string
func formatArtistName(artist string) string {
	if artist == "" {
		return ""
	}

	// 1. Check if the artist name contains brackets (in case (여자)아이들)
	idx := bytes.IndexByte([]byte(artist), '(')
	if idx < 1 {
		return artist
	}

	// 2. Split the artist name to two parts by the first bracket
	parts := strings.SplitN(artist, "(", 2)

	if len(parts) < 2 {
		slog.Error("Error splitting the artist name", "Parts: ", parts)
		return artist
	}

	// 3. Remove the space at the end of the first part
	firstPart := strings.TrimSpace(parts[0])

	// 4. Remove the last brakcet at the end of the second part
	secondPart := strings.TrimSuffix(parts[1], ")")

	if containsEnglish(secondPart) {
		return secondPart
	}
	return firstPart
}

// containsEnglish checks if a string contains any English characters
func containsEnglish(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Latin, r) {
			return true
		}
	}
	return false
}

// containsKorean checks if a string contains any Korean characters
func containsKorean(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Hangul, r) {
			return true
		}
	}
	return false
}

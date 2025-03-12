package main

import (
	"log/slog"
	"net/http"

	"github.com/akimdev15/melongo/playlist-server/spotify"
	"github.com/akimdev15/mscraper"
)

// AccessToken TODO - only for testing purpose. Should be REMOVED!!
const AccessToken = ""

func getSongs() []mscraper.Song {
	songs := mscraper.GetNewestSongsMelon("0300")
	return songs
}

func getMelonToop100() []mscraper.Song {
	songs := mscraper.GetMelonTop100Songs()
	return songs
}

func (apiCfg *apiConfig) testHandler(w http.ResponseWriter, r *http.Request) {
	songs := getMelonToop100()
	uris := []string{}
	for _, song := range songs {
		track, err := spotify.SearchTrack(song.Title, song.Artist, AccessToken)
		if err != nil {
			slog.Error("Error searching track", "error", err)
		}
		if track != nil && track.URI != "" {
			uris = append(uris, track.URI)
		}
	}

	// spotify.AddTrackToPlaylist("2XCwgZm2ornbTdEvaDT1h9", uris, AccessToken)

	respondWithJSON(w, 200, uris)
}

func (apiCfg *apiConfig) testNewAlbumsHandler(w http.ResponseWriter, r *http.Request) {
	// albums := mscraper.GetNewestAlbumFromMelon()
	// tracks, err := spotify.SearchTracksFromAlbum(albums[0].Name, albums[0].Artist, AccessToken)

	// This one works
	tracks, err := spotify.SearchTracksFromAlbum("아무렇지 않게", "DK", AccessToken)

	// This works but artis name is empty
	// tracks, err := spotify.SearchTracksFromAlbum("#2024: 가장자리", "Minit", AccessToken)

	if err != nil {
		slog.Error("Error searching tracks", "error", err)
		return
	}

	uris := []string{}
	for _, track := range tracks {
		uris = append(uris, track.URI)
	}

	// spotify.AddTrackToPlaylist("0msfdSZz5ZKXibCW6uZlvU", uris, AccessToken)

	respondWithJSON(w, 200, tracks)
}

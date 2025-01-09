import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './Home.css';

// Type definitions based on your new structure
interface Track {
  title: string;
  artist: string;
  popularity: number;
  uri: string;
}

interface Playlist {
  next: string;
  total: number;
  playlistPageURL: string;
  detailedPlaylistEndpoint: string;
  name: string;
  description: string;
  spotifyPlaylistID: string;
  imageUrl: string;
  totalTracks: number;
  tracksEndpoint: string;
}

const Home: React.FC = () => {
  const [playlists, setPlaylists] = useState<Playlist[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');
  const [expandedPlaylist, setExpandedPlaylist] = useState<string | null>(null); // To handle expanding playlist cards

  // Fetch user's playlists
  useEffect(() => {
    const fetchPlaylists = async () => {
      setLoading(true);
      setError('');
      try {
        const response = await axios.get('http://localhost:8080/playlists', {
          withCredentials: true,
        });
        console.log(response);
        setPlaylists(response.data.playlists); // Assuming response is in the structure `GetUserPlaylistsResponse`
      } catch (err) {
        setError('Failed to fetch playlists');
      } finally {
        setLoading(false);
      }
    };

    fetchPlaylists();
  }, []);

  // Fetch detailed tracks for a playlist when button is clicked
  const fetchDetailedTracks = async (detailedPlaylistEndpoint: string) => {
    setLoading(true);
    setError('');
    try {
      const response = await axios.get(detailedPlaylistEndpoint, {
        withCredentials: true,
      });
      console.log('Detailed Tracks:', response);
      // Here you would handle the detailed tracks response
      // Maybe you want to add them to the playlist's tracks or display them elsewhere
    } catch (err) {
      setError('Failed to fetch detailed tracks');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="home-container">
      <h1 className="page-title">Your Playlists</h1>

      {loading ? (
        <p className="loading-message">Loading playlists...</p>
      ) : error ? (
        <p className="error-message">{error}</p>
      ) : (
        <div className="playlist-grid">
          {playlists.map((playlist, index) => (
            <div className="playlist-card" key={index}>
              <img
                className="playlist-image"
                src={playlist.imageUrl || 'https://via.placeholder.com/150'}
                alt={playlist.name}
              />
              <div className="playlist-info">
                <h3 className="playlist-name">{playlist.name}</h3>
                <p className="playlist-description">{playlist.description}</p>

                {/* Show more options when playlist is expanded */}
                {expandedPlaylist === playlist.spotifyPlaylistID ? (
                  <div className="expanded-details">
                    <p className="playlist-total-tracks">Total Tracks: {playlist.totalTracks}</p>
                    <p className="playlist-next">Next: {playlist.next}</p>
                    <a
                      href={playlist.playlistPageURL}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="playlist-url"
                    >
                      Open Playlist on Spotify
                    </a>
                    <button
                      onClick={() => fetchDetailedTracks(playlist.detailedPlaylistEndpoint)}
                      className="fetch-detailed-tracks-button"
                    >
                      Fetch Detailed Tracks
                    </button>
                  </div>
                ) : null}

                <button
                  className="expand-button"
                  onClick={() =>
                    setExpandedPlaylist(expandedPlaylist === playlist.spotifyPlaylistID ? null : playlist.spotifyPlaylistID)
                  }
                >
                  {expandedPlaylist === playlist.spotifyPlaylistID ? 'Hide Details' : 'Show Details'}
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Home;

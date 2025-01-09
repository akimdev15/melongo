import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './Home.css';

// Type definitions
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
  const [expandedPlaylist, setExpandedPlaylist] = useState<string | null>(null);
  const [nextPageUrl, setNextPageUrl] = useState<string | null>(null);
  const [detailedTracks, setDetailedTracks] = useState<Track[]>([]); // New state for tracks
  const [showModal, setShowModal] = useState<boolean>(false); // State to toggle modal visibility

  // Fetch user's playlists
  useEffect(() => {
    const fetchPlaylists = async () => {
      setLoading(true);
      setError('');
      try {
        const response = await axios.get('http://localhost:8080/playlists', {
          withCredentials: true,
        });
        setPlaylists(response.data.playlists);
        setNextPageUrl(response.data.nextPageUrl);
      } catch (err) {
        setError('Failed to fetch playlists');
      } finally {
        setLoading(false);
      }
    };

    fetchPlaylists();
  }, []);

  // Fetch detailed tracks for a playlist when button is clicked
  const fetchDetailedTracks = async (event: React.MouseEvent<HTMLButtonElement>, tracksEndpoint: string) => {
    event.preventDefault(); // Prevent page reload on button click
    setError('');
    try {
      const response = await axios.get(`http://localhost:8080/playlist/tracks?endpoint=${tracksEndpoint}`, {
        withCredentials: true,
      });
      setDetailedTracks(response.data.playlistTracks); // Store tracks in state
      setShowModal(true); // Show modal with tracks
    } catch (err) {
      setError('Failed to fetch detailed tracks');
    }
  };

  // Close the modal
  const closeModal = () => {
    setShowModal(false);
    setDetailedTracks([]); // Clear tracks when modal is closed
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
                    <a
                      href={playlist.playlistPageURL}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="playlist-url"
                    >
                      Open Playlist on Spotify
                    </a>
                    <button
                      onClick={(event) => fetchDetailedTracks(event, playlist.tracksEndpoint)}
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

      {/* Modal for detailed tracks */}
      {showModal && (
        <div className="modal-overlay" onClick={closeModal}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h2>Track Details</h2>
            <ul className="track-list">
              {detailedTracks.map((track, index) => (
                <li key={index} className="track-item">
                  <p><strong>{track.title}</strong> by {track.artist}</p>
                  <p>Popularity: {track.popularity}</p>
                </li>
              ))}
            </ul>
            <button className="close-modal-button" onClick={closeModal}>
              Close
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default Home;

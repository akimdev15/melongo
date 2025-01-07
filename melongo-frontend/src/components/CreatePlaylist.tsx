import React, { useState } from 'react';
import axios from 'axios';
import './CreatePlaylist.css';

// Type definition for the playlist
interface Playlist {
    name: string;
    description: string;
    public: boolean;  // Indicates if the playlist is public or private
}

interface CreatePlaylistResponse {
    spotifyPlaylistID: string;
    externalUrl: string;
    name: string;
}

const CreatePlaylist: React.FC = () => {
    // State for Create Playlist form
    const [playlist, setPlaylist] = useState<Playlist>({
        name: '',
        description: '',
        public: true, // Default to public playlist
    });
    const [loading, setLoading] = useState<boolean>(false);
    const [status, setStatus] = useState<string>('');
    const [playlistDetails, setPlaylistDetails] = useState<CreatePlaylistResponse | null>(null);

    // State for Add Melon Songs to Playlist form
    const [melonPlaylistId, setMelonPlaylistId] = useState<string>('');
    const [melonDate, setMelonDate] = useState<string>(''); // New state for melon date
    const [melonStatus, setMelonStatus] = useState<string>('');

    // Handle input changes for playlist name, description, and public toggle
    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>, field: 'name' | 'description') => {
        setPlaylist({
            ...playlist,
            [field]: e.target.value,
        });
    };

    // Handle toggle for public/private playlist
    const handlePublicChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setPlaylist({
            ...playlist,
            public: e.target.checked,  // Toggle between true and false
        });
    };

    // Handle submitting the new playlist
    const handleSubmitPlaylist = async () => {
        if (!playlist.name || !playlist.description) {
            setStatus('Please provide a name and description for the playlist.');
            return;
        }

        setLoading(true);
        setStatus('');
        try {
            const response = await axios.post<CreatePlaylistResponse>(
                'http://localhost:8080/createPlaylist', // Replace with your backend endpoint
                playlist,
                {
                    withCredentials: true,
                }
            );

            // Update the state with response data
            setPlaylistDetails(response.data);
            setStatus(`Playlist "${response.data.name}" created successfully!`);
        } catch (err) {
            setStatus('Failed to create playlist.');
        } finally {
            setLoading(false);
        }
    };

    // Handle adding Melon Top 100 songs to the playlist
    const handleAddMelonSongsToPlaylist = async () => {
        if (!melonPlaylistId) {
            setMelonStatus('Please enter a playlist ID.');
            return;
        }

        if (!melonDate) {
            setMelonStatus('Please select a date for Melon Top 100.');
            return;
        }

        setLoading(true);
        setMelonStatus('');
        try {
            // Send a request to add Melon Top 100 songs to the playlist
            const response = await axios.post(
                'http://localhost:8080/melonTop100/create', // Replace with your backend endpoint
                {
                    playlistID: melonPlaylistId, // The playlist ID entered by the user
                    date: melonDate, // Include the date in the request body
                },
                {
                    withCredentials: true,
                }
            );

            setMelonStatus('Melon Top 100 songs added to playlist successfully!');
        } catch (err) {
            setMelonStatus('Failed to add Melon Top 100 songs to playlist.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="create-playlist-container">
            <h2>Create New Playlist</h2>

            {/* Playlist Creation Form */}
            <div className="playlist-form">
                <h3>Playlist Information</h3>

                <div className="input-group">
                    <label>Playlist Name:</label>
                    <input
                        type="text"
                        value={playlist.name}
                        onChange={(e) => handleInputChange(e, 'name')}
                        placeholder="Enter playlist name"
                        disabled={loading}
                    />
                </div>

                <div className="input-group">
                    <label>Description:</label>
                    <input
                        type="text"
                        value={playlist.description}
                        onChange={(e) => handleInputChange(e, 'description')}
                        placeholder="Enter playlist description"
                        disabled={loading}
                    />
                </div>

                {/* Public / Private Toggle */}
                <div className="input-group checkbox-container">
                    <label className="checkbox-label">
                        <input
                            type="checkbox"
                            checked={playlist.public}
                            onChange={handlePublicChange}
                            disabled={loading}
                        />
                        Public Playlist
                    </label>
                </div>

                {/* Submit Playlist Button */}
                <button onClick={handleSubmitPlaylist} disabled={loading}>
                    {loading ? 'Creating Playlist...' : 'Create Playlist'}
                </button>

                {/* Status Message for Playlist Creation */}
                {status && <p className="status">{status}</p>}

                {/* Display Created Playlist Details */}
                {playlistDetails && (
                    <div className="playlist-details">
                        <h3>Playlist Created:</h3>
                        <p><strong>Playlist Name:</strong> {playlistDetails.name}</p>
                        <p><strong>Spotify Playlist ID:</strong> {playlistDetails.spotifyPlaylistID}</p>
                        <p><strong>External URL:</strong> <a href={playlistDetails.externalUrl} target="_blank" rel="noopener noreferrer">Open Playlist</a></p>
                    </div>
                )}
            </div>

            <hr />

            {/* Add Melon Songs to Playlist Form */}
            <div className="melon-form">
                <h3>Add Melon Top 100 Songs to Playlist</h3>

                <div className="input-group">
                    <label>Spotify Playlist ID:</label>
                    <input
                        type="text"
                        value={melonPlaylistId}
                        onChange={(e) => setMelonPlaylistId(e.target.value)}
                        placeholder="Enter playlist ID"
                        disabled={loading}
                    />
                </div>

                {/* Date Selector */}
                <div className="input-group">
                    <label>Select Date for Melon Top 100:</label>
                    <input
                        type="date"
                        value={melonDate}
                        onChange={(e) => setMelonDate(e.target.value)}
                        disabled={loading}
                    />
                </div>

                <button onClick={handleAddMelonSongsToPlaylist} disabled={loading}>
                    {loading ? 'Adding Songs...' : 'Add Melon Top 100 Songs'}
                </button>

                {/* Status Message for Melon Songs Addition */}
                {melonStatus && <p className="status">{melonStatus}</p>}
            </div>
        </div>
    );
};

export default CreatePlaylist;

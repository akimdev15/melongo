import React, { useState } from 'react';
import axios from 'axios';

// Type definition for missed track, now with both original and edited title and artist
interface MissedTrack {
    rank: number;
    title: string;  // The original title (from the backend)
    editedTitle: string; // The edited title (updated by the user)
    artist: string; // The original artist name (from the backend)
    editedArtist: string; // The edited artist name (updated by the user)
    date: string;
}

const getTodayDate = (): string => {
    const today = new Date();
    const year = today.getFullYear();
    const month = String(today.getMonth() + 1).padStart(2, '0'); // Months are zero-based
    const day = String(today.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
};

const MissedTracks: React.FC = () => {
    const [missedTracks, setMissedTracks] = useState<MissedTrack[]>([]);
    const [date, setDate] = useState<string>(getTodayDate()); // Default date
    const [loading, setLoading] = useState<boolean>(false);
    const [status, setStatus] = useState<string>('');
    const [apiKey, setApiKey] = useState<string>('0ecb6fc2fc51f6bfa015969446a0557bf242556ba9fcec9c4476ead4fdc74c37');

    // Handle fetching missed tracks
    const fetchMissedTracks = async () => {
        setLoading(true);
        setStatus('');
        try {
            const response = await axios.get(`http://localhost:8080/missedTracks`, {
                params: { date },
                headers: {
                    Authorization: `ApiKey ${apiKey}`,
                },
            });

            // Ensure the response contains the 'missedTracks' array
            if (response.data && Array.isArray(response.data.missedTracks)) {
                // Store original title, artist, and edited versions of each
                const tracksWithEditedFields = response.data.missedTracks.map((track: any) => ({
                    ...track,
                    editedTitle: track.title, // Initialize editedTitle with the original title
                    editedArtist: track.artist, // Initialize editedArtist with the original artist
                }));
                setMissedTracks(tracksWithEditedFields); // Set missed tracks with original title and artist
            } else if (response.data && Object.keys(response.data).length === 0) {
                setStatus(`No missed tracks found for the date: ${date}`);
                setMissedTracks([]); // Clear missed tracks
            } else {
                setStatus('Invalid response format');
                setMissedTracks([]); // Clear missed tracks
            }
        } catch (err) {
            setStatus('Failed to fetch missed tracks');
        } finally {
            setLoading(false);
        }
    };

    // Handle input changes for the edited title or artist
    const handleInputChange = (
        e: React.ChangeEvent<HTMLInputElement>,
        index: number,
        field: 'editedTitle' | 'editedArtist' // Either editedTitle or editedArtist
    ) => {
        const updatedTracks = [...missedTracks];
        updatedTracks[index] = {
            ...updatedTracks[index],
            [field]: e.target.value, // Dynamically update either editedTitle or editedArtist
        };
        setMissedTracks(updatedTracks);
    };

    // Handle submit of resolved tracks
    const handleSubmitResolvedTracks = async () => {
        setLoading(true);
        try {
            const resolvedTracksData = missedTracks.map((track) => ({
                rank: track.rank,
                missed_title: track.title, // Original title from backend (missed_title)
                missed_artist: track.artist, // Original artist from backend (missed_artist)
                title: track.editedTitle, // The edited title (title)
                artist: track.editedArtist, // The edited artist (artist)
                date,
            }));

            const response = await axios.post(
                'http://localhost:8080/resolveMissedTracks',
                { resolvedTracks: resolvedTracksData },
                {
                    headers: {
                        Authorization: `ApiKey ${apiKey}`,
                    },
                }
            );

            setStatus('Resolving missed tracks... Please check again');
            setMissedTracks([]); // Clear missed tracks
        } catch (err) {
            setStatus('Failed to resolve missed tracks');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="missed-tracks-container">
            <h2>Missed Tracks</h2>

            {/* Date input */}
            <div>
                <label>Date: </label>
                <input
                    type="date"
                    value={date}
                    onChange={(e) => setDate(e.target.value)}
                />
            </div>

            {/* Button to fetch missed tracks */}
            <button onClick={fetchMissedTracks} disabled={loading}>
                {loading ? 'Loading...' : 'Fetch Missed Tracks'}
            </button>

            {status && <p className="status">{status}</p>}

            {/* Display missed tracks */}
            <ul>
                {missedTracks.map((track, index) => (
                    <li key={index}>
                        <div>
                            <strong>Rank: {track.rank}</strong>
                        </div>
                        <div>
                            <label>Original Title:</label>
                            <input
                                type="text"
                                value={track.title}
                                disabled
                            />
                        </div>
                        <div>
                            <label>Edited Title:</label>
                            <input
                                type="text"
                                value={track.editedTitle} // Bind to the edited title
                                onChange={(e) => handleInputChange(e, index, 'editedTitle')} // Update the editedTitle field
                            />
                        </div>
                        <div>
                            <label>Original Artist:</label>
                            <input
                                type="text"
                                value={track.artist}
                                disabled
                            />
                        </div>
                        <div>
                            <label>Edited Artist:</label>
                            <input
                                type="text"
                                value={track.editedArtist} // Bind to the edited artist
                                onChange={(e) => handleInputChange(e, index, 'editedArtist')} // Update the editedArtist field
                            />
                        </div>
                    </li>
                ))}
            </ul>

            {/* Button to submit resolved tracks */}
            <button onClick={handleSubmitResolvedTracks} disabled={loading}>
                Resolve Tracks
            </button>
        </div>
    );
};

export default MissedTracks;

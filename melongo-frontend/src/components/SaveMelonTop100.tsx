import React, { useState } from 'react';
import axios from 'axios';
import './SaveMelonTop100.css';

const SaveMelonTop100: React.FC = () => {
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string>('');
    const [success, setSuccess] = useState<string>('');

    // Function to save the top 100 melon tracks
    const saveMelonTop100 = async () => {
        setLoading(true);
        setError('');
        setSuccess('');
        try {
            const response = await axios.post('http://localhost:8080/melonTop100/save', {}, {
                withCredentials: true,
            });
            setSuccess('Successfully saved the top 100 melon tracks.');
        } catch (err) {
            setError('Failed to save the top 100 melon tracks.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="save-melon-top100-container">
            <h1 className="page-title">Save Melon Top 100 Tracks</h1>

            {loading ? (
                <p className="loading-message">Saving tracks...</p>
            ) : error ? (
                <p className="error-message">{error}</p>
            ) : success ? (
                <p className="success-message">{success}</p>
            ) : (
                <button onClick={saveMelonTop100} className="save-button">
                    Save Melon Top 100 Tracks
                </button>
            )}
        </div>
    );
};

export default SaveMelonTop100;
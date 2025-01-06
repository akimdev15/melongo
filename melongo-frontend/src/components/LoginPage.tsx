import React, { useState } from 'react';
import './LoginPage.css';  // Importing the CSS file

const LoginPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Spotify login URL (backend will handle the redirect to Spotify)
    const handleLogin = async () => {
        setLoading(true);
        setError(null);

        try {
            const authURL = 'https://accounts.spotify.com/authorize?client_id=60e1ce49166b49ebb3c2999beabe8ac5&response_type=code&redirect_uri=http://localhost:8080/callback&scope=user-read-email%20user-read-private%20playlist-modify-public%20playlist-modify-private%20playlist-read-collaborative%20playlist-read-private';
            window.location.href = authURL;
        } catch (err) {
            setError('Error initiating login with Spotify.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="login-container">
            <h2 className="login-title">Login with Spotify</h2>
            <button className="login-button" onClick={handleLogin} disabled={loading}>
                {loading ? 'Loading...' : 'Login with Spotify'}
            </button>
            {error && <p className="error-message">{error}</p>}
        </div>
    );
};

export default LoginPage;

// src/main.tsx (or src/index.tsx)
import React from 'react';
import ReactDOM from 'react-dom/client';  // Updated import for React 18
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import App from './App';
import './styles/theme.css';

const Main: React.FC = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<App />} />
        {/* Future routes */}
        <Route path="/createPlaylist" element={<div>Create Playlist</div>} />
        <Route path="/resolveMissedTracks" element={<div>Resolve Missed Tracks</div>} />
      </Routes>
    </Router>
  );
};

// Create the root and render the app
const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);  // Type assertion for `root`
root.render(<Main />);

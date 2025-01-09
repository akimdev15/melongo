import React, { useEffect } from "react";
import ReactDOM from 'react-dom/client';  // Updated import for React 18
import { BrowserRouter as Router, Route, Routes, useNavigate } from 'react-router-dom';
import './styles/theme.css';
import Dashboard from "./components/Dashbooard";
import LoginPage from "./components/LoginPage";
import CreatePlaylist from "./components/CreatePlaylist";
import Home from "./pages/Home";

const Main: React.FC = () => {
  const navigate = useNavigate(); // useNavigate hook allows us to navigate programmatically

  // useEffect(() => {
  //   if (token) {
  //     // If the token exists, navigate to /dashboard
  //     navigate('/');
  //   } else {
  //     // If no token exists, navigate to /
  //     navigate('/login');
  //   }
  // }, [token, navigate]); // This effect will run when 'token' or 'navigate' changes

  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/createPlaylist" element={<CreatePlaylist />} />
      <Route path="/dashboard" element={<Dashboard />} />
      <Route path="/createPlaylist" element={<div>Create Playlist</div>} />
      <Route path="/resolveMissedTracks" element={<div>Resolve Missed Tracks</div>} />
    </Routes>
  );
};

// Create the root and render the app
const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);  // Type assertion for `root`

// Wrap Main component in Router here
root.render(
  <Router>
    <Main />
  </Router>
);

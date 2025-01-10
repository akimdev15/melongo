import React from "react";
import ReactDOM from "react-dom/client"; // Updated import for React 18
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import "./styles/theme.css"; // Import the CSS file
import LoginPage from "./components/LoginPage";
import CreatePlaylist from "./components/CreatePlaylist";
import Home from "./pages/Home";
import Navbar from "./components/NavBar";
import Dashboard from "./components/Dashbooard";
import './styles/main.css'

const Main: React.FC = () => {
  return (
    <div>
      <Navbar /> {/* Add the Navbar here */}
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/createPlaylist" element={<CreatePlaylist />} />
        <Route path="/dashboard" element={<Dashboard />} />
      </Routes>
    </div>
  );
};

// Create the root and render the app
const root = ReactDOM.createRoot(document.getElementById("root") as HTMLElement);

root.render(
  <Router>
    <Main />
  </Router>
);
``
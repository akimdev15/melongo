import React, { useState } from "react";
import { NavLink } from "react-router-dom"; // Import NavLink
import "./Navbar.css"; // Import the CSS file

const Navbar: React.FC = () => {
    const [isMenuOpen, setMenuOpen] = useState(false);

    const toggleMenu = () => {
        setMenuOpen(!isMenuOpen);
    };

    return (
        <nav className={`navbar ${isMenuOpen ? 'active' : ''}`}>
            <div className="navbar-toggle" onClick={toggleMenu}>
                <span className={isMenuOpen ? 'open' : ''}>☰</span>
            </div>
            <ul className="nav-list">
                <li className="nav-item">
                    <NavLink
                        to="/"
                        className={({ isActive }) => isActive ? "nav-link active-link" : "nav-link"}
                    >
                        Home
                    </NavLink>
                </li>
                <li className="nav-item">
                    <NavLink
                        to="/login"
                        className={({ isActive }) => isActive ? "nav-link active-link" : "nav-link"}
                    >
                        Login
                    </NavLink>
                </li>
                <li className="nav-item">
                    <NavLink
                        to="/createPlaylist"
                        className={({ isActive }) => isActive ? "nav-link active-link" : "nav-link"}
                    >
                        Create Playlist
                    </NavLink>
                </li>
                <li className="nav-item">
                    <NavLink
                        to="/dashboard"
                        className={({ isActive }) => isActive ? "nav-link active-link" : "nav-link"}
                    >
                        Dashboard
                    </NavLink>
                </li>
            </ul>
        </nav>
    );
};

export default Navbar;

// src/App.tsx
import React from 'react';
import './styles/theme.css';
import MissedTracks from './components/MissedTracks';

const App: React.FC = () => {
    return (
        <div className="App">
            <MissedTracks />
        </div>
    );
}

export default App;

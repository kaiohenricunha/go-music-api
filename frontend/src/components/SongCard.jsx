import React from 'react';

function SongCard({ name, artist }) {
  return (
    <div className="song-card" style={{ margin: '1rem', padding: '1rem', borderRadius: '8px', background: 'var(--color-darkblue-alpha)', color: 'var(--color-white)' }}>
      <h3>{name}</h3>
      <p>{artist}</p>
      {/* Placeholder for future edit and delete buttons */}
    </div>
  );
}

export default SongCard;

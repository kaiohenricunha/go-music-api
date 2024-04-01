package service

import (
	"fmt"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
)

type PlaylistService interface {
	GetAllPlaylists() ([]model.Playlist, error)
	GetPlaylistByID(playlistID string) (*model.Playlist, error)
	AddSongToPlaylist(playlistID, songID string) error
	RemoveSongFromPlaylist(playlistID, songID string) error
}

type playlistService struct {
	musicDAO dao.MusicDAO
}

func NewPlaylistService(musicDAO dao.MusicDAO) PlaylistService {
	return &playlistService{musicDAO: musicDAO}
}

var (
	ErrPlaylistNotFound = fmt.Errorf("playlist not found")
)

func (s *playlistService) GetAllPlaylists() ([]model.Playlist, error) {
	return s.musicDAO.GetAllPlaylists()
}

func (s *playlistService) GetPlaylistByID(playlistID string) (*model.Playlist, error) {
	playlist, err := s.musicDAO.GetPlaylistByID(playlistID)
	if err != nil {
		return nil, err
	}

	if playlist == nil {
		return nil, fmt.Errorf("playlist with ID %s not found", playlistID)
	}

	return playlist, nil
}

// AddSongToPlaylist adds a song to a playlist.
func (s *playlistService) AddSongToPlaylist(playlistID, songID string) error {
	return s.musicDAO.AddSongToPlaylist(playlistID, songID)
}

func (s *playlistService) RemoveSongFromPlaylist(playlistID, songID string) error {
	return s.musicDAO.RemoveSongFromPlaylist(playlistID, songID)
}

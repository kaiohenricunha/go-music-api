package service

import (
	"errors"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
)

type PlaylistService interface {
	CreatePlaylist(authenticatedUserID uint, playlist *model.Playlist) error
	AddSongToPlaylist(playlistID, songID uint) error
	RemoveSongFromPlaylist(playlistID, songID uint) error
	GetPlaylistsByUserID(userID uint) ([]model.Playlist, error)
	DeletePlaylist(playlistID uint) error
}

type playlistService struct {
	musicDAO dao.MusicDAO
}

func NewPlaylistService(musicDAO dao.MusicDAO) PlaylistService {
	return &playlistService{musicDAO: musicDAO}
}

func (s *playlistService) CreatePlaylist(authenticatedUserID uint, playlist *model.Playlist) error {
	// Check if the playlist name is empty
	if playlist.Name == "" {
		return errors.New("playlist name is required")
	}

	// Check if the playlist already exists
	existingPlaylist, err := s.musicDAO.GetPlaylistByName(authenticatedUserID, playlist.Name)
	if err != nil {
		return err
	}
	if existingPlaylist != nil {
		return errors.New("playlist already exists")
	}

	// Create the playlist
	return s.musicDAO.CreatePlaylist(playlist)
}

func (s *playlistService) AddSongToPlaylist(playlistID, songID uint) error {
	return s.musicDAO.AddSongToPlaylist(playlistID, songID)
}

func (s *playlistService) RemoveSongFromPlaylist(playlistID, songID uint) error {
	// Similarly, check if the playlist and song exist and if the song is actually in the playlist
	return s.musicDAO.RemoveSongFromPlaylist(playlistID, songID)
}

func (s *playlistService) GetPlaylistsByUserID(userID uint) ([]model.Playlist, error) {
	return s.musicDAO.GetPlaylistsByUserID(userID)
}

func (s *playlistService) DeletePlaylist(playlistID uint) error {
	return s.musicDAO.DeletePlaylist(playlistID)
}

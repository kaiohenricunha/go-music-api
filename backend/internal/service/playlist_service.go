package service

import (
	"errors"
	"fmt"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
)

type PlaylistService interface {
	CreatePlaylist(playlist *model.Playlist) error
	AddSongToPlaylist(playlistID, songID uint) error
	GetAllPlaylists() ([]model.Playlist, error)
}

type playlistService struct {
	musicDAO dao.MusicDAO
}

func NewPlaylistService(musicDAO dao.MusicDAO) PlaylistService {
	return &playlistService{musicDAO: musicDAO}
}

var (
	ErrPlaylistAlreadyExists = errors.New("playlist already exists for this user")
	ErrUserDoesNotExist      = errors.New("user does not exist")
)

func (s *playlistService) CreatePlaylist(playlist *model.Playlist) error {
	if playlist.Name == "" {
		return errors.New("playlist name cannot be empty")
	}

	existingPlaylist, err := s.musicDAO.GetPlaylistByNameAndUserID(playlist.Name, playlist.UserID)
	if err != nil {
		return fmt.Errorf("error checking for existing playlist: %w", err)
	}
	if existingPlaylist != nil {
		return errors.New("playlist already exists for this user")
	}

	// Proceed with creation if existingPlaylist is nil
	return s.musicDAO.CreatePlaylist(playlist)
}

func (s *playlistService) AddSongToPlaylist(playlistID, songID uint) error {
	return s.musicDAO.AddSongToPlaylist(playlistID, songID)
}

func (s *playlistService) GetAllPlaylists() ([]model.Playlist, error) {
	return s.musicDAO.GetAllPlaylists()
}

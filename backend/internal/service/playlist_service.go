package service

import (
	"fmt"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
)

type PlaylistService interface {
	GetAllPlaylists() ([]model.Playlist, error)
	GetPlaylistByID(playlistID string) (*model.Playlist, error)
}

type playlistService struct {
	musicDAO dao.MusicDAO
}

func NewPlaylistService(musicDAO dao.MusicDAO) PlaylistService {
	return &playlistService{musicDAO: musicDAO}
}

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

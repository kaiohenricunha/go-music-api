package service

import (
	"errors"
	"fmt"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
)

var (
	ErrSongNameRequired  = errors.New("song name and artist are required")
	ErrSongAlreadyExists = errors.New("a song with the same name by the same artist already exists")
	ErrSongNotFound      = errors.New("song not found")
)

type SongService interface {
	AddSong(song *model.Song) error
	GetAllSongs() ([]model.Song, error)
	UpdateSong(song *model.Song) error
	DeleteSong(songID uint) error
}

type songService struct {
	songDAO dao.MusicDAO
}

func NewSongService(songDAO dao.MusicDAO) SongService {
	return &songService{songDAO: songDAO}
}

func (s *songService) AddSong(song *model.Song) error {
	if song.Name == "" || song.Artist == "" {
		return ErrSongNameRequired
	}

	existingSong, err := s.songDAO.FindSongByName(song.Name)
	if err != nil {
		return fmt.Errorf("failed to check for existing song: %w", err)
	}
	if existingSong != nil {
		return ErrSongAlreadyExists
	}

	return s.songDAO.CreateSong(song)
}

func (s *songService) GetAllSongs() ([]model.Song, error) {
	return s.songDAO.GetAllSongs()
}

func (s *songService) UpdateSong(song *model.Song) error {
	if song.Name == "" {
		return ErrSongNameRequired
	}

	existingSong, err := s.songDAO.FindSongByName(song.Name)
	if err != nil || existingSong == nil {
		return ErrSongNotFound
	}

	return s.songDAO.UpdateSong(song)
}

func (s *songService) DeleteSong(songID uint) error {
	return s.songDAO.DeleteSong(songID)
}

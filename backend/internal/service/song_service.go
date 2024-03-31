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
	GetAllSongs() ([]model.Song, error)
}

type songService struct {
	songDAO dao.MusicDAO
}

func NewSongService(songDAO dao.MusicDAO) SongService {
	return &songService{songDAO: songDAO}
}

func (s *songService) GetAllSongs() ([]model.Song, error) {
	return s.songDAO.GetAllSongs()
}

package dao

import (
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
)

type MusicDAO interface {
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(userID uint) error
	GetAllUsers() ([]model.User, error)
	FindByUsername(username string) (*model.User, error)

	CreateSong(song *model.Song) error
	GetAllSongs() ([]model.Song, error)
	UpdateSong(song *model.Song) error
	DeleteSong(songID uint) error
	FindSongByName(songName string) (*model.Song, error)
}

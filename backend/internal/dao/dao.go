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
	FindUserByID(userID uint) (*model.User, error)

	CreateSong(song *model.Song) error
	GetAllSongs() ([]model.Song, error)
	UpdateSong(song *model.Song) error
	DeleteSong(songID uint) error
	FindSongByName(songName string) (*model.Song, error)
	FindSongByID(songID uint) (*model.Song, error)

	CreatePlaylist(playlist *model.Playlist) error
	AddSongToPlaylist(playlistID, songID uint) error
	GetPlaylistByNameAndUserID(playlistName string, userID uint) (*model.Playlist, error)
	GetAllPlaylists() ([]model.Playlist, error)
}

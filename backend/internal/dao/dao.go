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

	GetAllSongs() ([]model.Song, error)
	GetSongFromSpotifyByID(spotifyID string) (*model.Song, error)
	SearchSongsFromSpotify(trackName, artistName string) ([]model.Song, error)
	CreatePlaylist(playlist *model.Playlist) error
	AddSongToPlaylist(playlistID, songID uint) error
	GetPlaylistByNameAndUserID(playlistName string, userID uint) (*model.Playlist, error)
	GetAllPlaylists() ([]model.Playlist, error)
}

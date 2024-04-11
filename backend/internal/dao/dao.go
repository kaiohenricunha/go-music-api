package dao

import (
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
)

type MusicDAO interface {
	CreateUser(user *model.User) error
	GetAllUsers() ([]model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(userID uint) (*model.User, error)

	CreateSong(song *model.Song) error
	GetAllSongs() ([]model.Song, error)
	GetSongByID(songID string) (*model.Song, error)
	GetSongBySpotifyID(spotifyID string) (*model.Song, error)
	GetSongByNameAndArtist(songName, artistName string) (*model.Song, error)
	GetSongFromSpotifyByID(spotifyID string) (*model.Song, error)
	SearchSongsFromSpotify(trackName, artistName string) ([]model.Song, error)

	GetAllPlaylists() ([]model.Playlist, error)
	GetPlaylistByID(playlistID string) (*model.Playlist, error)
	AddSongToPlaylist(playlistID, songID string) error
	RemoveSongFromPlaylist(playlistID, songID string) error
}

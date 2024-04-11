package mocks

import (
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// MusicDAO is a mock type for the MusicDAO type
type MusicDAO struct {
	mock.Mock
}

////////////////////////////////
// USER METHODS //
////////////////////////////////

// CreateUser mocks the CreateUser method
func (_m *MusicDAO) CreateUser(user *model.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserByID mocks the GetUserByID method
func (_m *MusicDAO) GetUserByID(userID uint) (*model.User, error) {
	ret := _m.Called(userID)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(uint) *model.User); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllUsers mocks the GetAllUsers method
func (_m *MusicDAO) GetAllUsers() ([]model.User, error) {
	ret := _m.Called()

	var r0 []model.User
	if rf, ok := ret.Get(0).(func() []model.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByUsername mocks the GetUserByUsername method
func (_m *MusicDAO) GetUserByUsername(username string) (*model.User, error) {
	ret := _m.Called(username)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(string) *model.User); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

////////////////////////////////
// SONG METHODS //
////////////////////////////////

// CreateSong mocks the CreateSong method
func (_m *MusicDAO) CreateSong(song *model.Song) error {
	ret := _m.Called(song)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Song) error); ok {
		r0 = rf(song)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllSongs mocks the GetAllSongs method
func (_m *MusicDAO) GetAllSongs() ([]model.Song, error) {
	ret := _m.Called()

	var r0 []model.Song
	if rf, ok := ret.Get(0).(func() []model.Song); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Song)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSongByID mocks the GetSongByID method
func (_m *MusicDAO) GetSongByID(songID string) (*model.Song, error) {
	ret := _m.Called(songID)

	var r0 *model.Song
	if rf, ok := ret.Get(0).(func(string) *model.Song); ok {
		r0 = rf(songID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Song)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(songID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSongByNameAndArtist mocks the GetSongByNameAndArtist method
func (_m *MusicDAO) GetSongByNameAndArtist(songName, artistName string) (*model.Song, error) {
	ret := _m.Called(songName, artistName)

	var r0 *model.Song
	if rf, ok := ret.Get(0).(func(string, string) *model.Song); ok {
		r0 = rf(songName, artistName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Song)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(songName, artistName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSongFromSpotifyByID mocks the GetSongFromSpotifyByID method
func (_m *MusicDAO) GetSongFromSpotifyByID(spotifyID string) (*model.Song, error) {
	ret := _m.Called(spotifyID)

	var r0 *model.Song
	if rf, ok := ret.Get(0).(func(string) *model.Song); ok {
		r0 = rf(spotifyID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Song)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(spotifyID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchSongsFromSpotify mocks the SearchSongsFromSpotify method
func (_m *MusicDAO) SearchSongsFromSpotify(trackName, artistName string) ([]model.Song, error) {
	ret := _m.Called(trackName, artistName)

	var r0 []model.Song
	if rf, ok := ret.Get(0).(func(string, string) []model.Song); ok {
		r0 = rf(trackName, artistName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Song)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(trackName, artistName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

////////////////////////////////
// PLAYLIST METHODS //
////////////////////////////////

// GetPlaylistByID mocks the GetPlaylistByID method
func (_m *MusicDAO) GetPlaylistByID(playlistID string) (*model.Playlist, error) {
	ret := _m.Called(playlistID)

	var r0 *model.Playlist
	if rf, ok := ret.Get(0).(func(string) *model.Playlist); ok {
		r0 = rf(playlistID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Playlist)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(playlistID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllPlaylists mocks the GetAllPlaylists method
func (_m *MusicDAO) GetAllPlaylists() ([]model.Playlist, error) {
	ret := _m.Called()

	var r0 []model.Playlist
	if rf, ok := ret.Get(0).(func() []model.Playlist); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Playlist)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddSongToPlaylist mocks the AddSongToPlaylist method
func (_m *MusicDAO) AddSongToPlaylist(playlistID, songID string) error {
	ret := _m.Called(playlistID, songID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(playlistID, songID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveSongFromPlaylist mocks the RemoveSongFromPlaylist method
func (_m *MusicDAO) RemoveSongFromPlaylist(playlistID, songID string) error {
	ret := _m.Called(playlistID, songID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(playlistID, songID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetSongBySpotifyID mocks the GetSongBySpotifyID method
func (_m *MusicDAO) GetSongBySpotifyID(spotifyID string) (*model.Song, error) {
	ret := _m.Called(spotifyID)

	var r0 *model.Song
	if rf, ok := ret.Get(0).(func(string) *model.Song); ok {
		r0 = rf(spotifyID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Song)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(spotifyID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

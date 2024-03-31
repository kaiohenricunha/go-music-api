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

// UpdateUser mocks the UpdateUser method
func (_m *MusicDAO) UpdateUser(user *model.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteUser mocks the DeleteUser method
func (_m *MusicDAO) DeleteUser(userID uint) error {
	ret := _m.Called(userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindUserByID mocks the FindUserByID method
func (_m *MusicDAO) FindUserByID(userID uint) (*model.User, error) {
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

// FindByUsername mocks the FindByUsername method
func (_m *MusicDAO) FindByUsername(username string) (*model.User, error) {
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

// CreatePlaylist mocks the CreatePlaylist method
func (_m *MusicDAO) CreatePlaylist(playlist *model.Playlist) error {
	ret := _m.Called(playlist)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Playlist) error); ok {
		r0 = rf(playlist)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddSongToPlaylist mocks the AddSongToPlaylist method
func (_m *MusicDAO) AddSongToPlaylist(playlistID, songID uint) error {
	ret := _m.Called(playlistID, songID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint, uint) error); ok {
		r0 = rf(playlistID, songID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveSongFromPlaylist mocks the RemoveSongFromPlaylist method
func (_m *MusicDAO) RemoveSongFromPlaylist(playlistID, songID uint) error {
	ret := _m.Called(playlistID, songID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint, uint) error); ok {
		r0 = rf(playlistID, songID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetPlaylistByNameAndUserID mocks the GetPlaylistByNameAndUserID method
func (_m *MusicDAO) GetPlaylistByNameAndUserID(playlistName string, userID uint) (*model.Playlist, error) {
	ret := _m.Called(playlistName, userID)

	var r0 *model.Playlist
	if rf, ok := ret.Get(0).(func(string, uint) *model.Playlist); ok {
		r0 = rf(playlistName, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Playlist)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint) error); ok {
		r1 = rf(playlistName, userID)
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

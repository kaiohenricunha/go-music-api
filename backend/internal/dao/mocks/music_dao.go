package mocks

import (
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// MusicDAO is a mock type for the MusicDAO type
type MusicDAO struct {
	mock.Mock
}

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

// UpdateSong mocks the UpdateSong method
func (_m *MusicDAO) UpdateSong(song *model.Song) error {
	ret := _m.Called(song)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Song) error); ok {
		r0 = rf(song)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSong mocks the DeleteSong method
func (_m *MusicDAO) DeleteSong(songID uint) error {
	ret := _m.Called(songID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(songID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindSongByID mocks the FindSongByID method
func (_m *MusicDAO) FindSongByID(songID uint) (*model.Song, error) {
	ret := _m.Called(songID)

	var r0 *model.Song
	if rf, ok := ret.Get(0).(func(uint) *model.Song); ok {
		r0 = rf(songID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Song)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(songID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindSongByName mocks the FindSongByName method
func (_m *MusicDAO) FindSongByName(songName string) (*model.Song, error) {
	ret := _m.Called(songName)

	var r0 *model.Song
	if rf, ok := ret.Get(0).(func(string) *model.Song); ok {
		r0 = rf(songName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Song)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(songName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

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

// UpdatePlaylist mocks the UpdatePlaylist method
func (_m *MusicDAO) UpdatePlaylist(playlist *model.Playlist) error {
	ret := _m.Called(playlist)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Playlist) error); ok {
		r0 = rf(playlist)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePlaylist mocks the DeletePlaylist method
func (_m *MusicDAO) DeletePlaylist(playlistID uint) error {
	ret := _m.Called(playlistID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(playlistID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
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

// GetPlaylistByID mocks the GetPlaylistByID method
func (_m *MusicDAO) GetPlaylistByID(id uint) (*model.Playlist, error) {
	ret := _m.Called(id)

	var r0 *model.Playlist
	if rf, ok := ret.Get(0).(func(uint) *model.Playlist); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Playlist)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPlaylistsByUserID mocks the GetPlaylistsByUserID method
func (_m *MusicDAO) GetPlaylistsByUserID(userID uint) ([]model.Playlist, error) {
	ret := _m.Called(userID)

	var r0 []model.Playlist
	if rf, ok := ret.Get(0).(func(uint) []model.Playlist); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Playlist)
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

// GetPlaylistByName mocks the GetPlaylistByName method
func (_m *MusicDAO) GetPlaylistByName(playlistName string) (*model.Playlist, error) {
	ret := _m.Called(playlistName)

	var r0 *model.Playlist
	if rf, ok := ret.Get(0).(func(string) *model.Playlist); ok {
		r0 = rf(playlistName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Playlist)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(playlistName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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

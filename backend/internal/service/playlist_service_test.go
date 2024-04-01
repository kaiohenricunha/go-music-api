package service

import (
	"testing"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao/mocks"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGetAllPlaylists(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	ps := NewPlaylistService(mockDAO)
	mockPlaylists := []model.Playlist{{Name: "Chill Vibes"}, {Name: "Workout"}}

	mockDAO.On("GetAllPlaylists").Return(mockPlaylists, nil)

	playlists, err := ps.GetAllPlaylists()
	assert.NoError(t, err)
	assert.Equal(t, mockPlaylists, playlists)
}

func TestGetPlaylistByID(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	ps := NewPlaylistService(mockDAO)
	mockPlaylist := &model.Playlist{Name: "Chill Vibes"}

	mockDAO.On("GetPlaylistByID", "1").Return(mockPlaylist, nil)
	mockDAO.On("GetPlaylistByID", "2").Return(nil, ErrPlaylistNotFound)

	playlist, err := ps.GetPlaylistByID("1")
	assert.NoError(t, err)
	assert.Equal(t, mockPlaylist, playlist)

	_, err = ps.GetPlaylistByID("2")
	assert.Equal(t, ErrPlaylistNotFound, err)
}

func TestAddSongToPlaylist(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	ps := NewPlaylistService(mockDAO)

	mockDAO.On("AddSongToPlaylist", "1", "1").Return(nil)
	mockDAO.On("AddSongToPlaylist", "1", "2").Return(ErrPlaylistNotFound)

	err := ps.AddSongToPlaylist("1", "1")
	assert.NoError(t, err)

	err = ps.AddSongToPlaylist("1", "2")
	assert.Equal(t, ErrPlaylistNotFound, err)
}

func TestRemoveSongFromPlaylist(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	ps := NewPlaylistService(mockDAO)

	mockDAO.On("RemoveSongFromPlaylist", "1", "1").Return(nil)
	mockDAO.On("RemoveSongFromPlaylist", "1", "2").Return(ErrPlaylistNotFound)

	err := ps.RemoveSongFromPlaylist("1", "1")
	assert.NoError(t, err)

	err = ps.RemoveSongFromPlaylist("1", "2")
	assert.Equal(t, ErrPlaylistNotFound, err)
}

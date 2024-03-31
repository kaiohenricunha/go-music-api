package service

import (
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao/mocks"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePlaylist_AlreadyExists(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	playlistService := NewPlaylistService(mockDAO)

	playlist := &model.Playlist{
		Name:   "Chill Hits",
		UserID: 1,
	}

	mockDAO.On("GetPlaylistByNameAndUserID", playlist.Name, playlist.UserID).Return(playlist, nil)

	err := playlistService.CreatePlaylist(playlist)

	assert.EqualError(t, err, ErrPlaylistAlreadyExists.Error())
	mockDAO.AssertExpectations(t)
}

func TestCreatePlaylist_Success(t *testing.T) {
	mockDAO := &mocks.MusicDAO{}
	service := NewPlaylistService(mockDAO)

	playlist := &model.Playlist{
		Name:   "Workout Jams",
		UserID: 1,
	}

	mockDAO.On("GetPlaylistByNameAndUserID", playlist.Name, playlist.UserID).Return(nil, nil) // Simulate no existing playlist found
	mockDAO.On("CreatePlaylist", playlist).Return(nil)

	err := service.CreatePlaylist(playlist)

	assert.NoError(t, err)
	mockDAO.AssertExpectations(t)
}

func TestAddSongToPlaylist_Success(t *testing.T) {
	mockDAO := &mocks.MusicDAO{}
	service := NewPlaylistService(mockDAO)

	playlistID := uint(1)
	songID := uint(1)

	mockDAO.On("AddSongToPlaylist", playlistID, songID).Return(nil)

	err := service.AddSongToPlaylist(playlistID, songID)

	assert.NoError(t, err)
	mockDAO.AssertExpectations(t)
}

func TestAddSongToPlaylist_Error(t *testing.T) {
	mockDAO := &mocks.MusicDAO{}
	service := NewPlaylistService(mockDAO)

	playlistID := uint(1)
	songID := uint(1)

	mockDAO.On("AddSongToPlaylist", playlistID, songID).Return(ErrPlaylistAlreadyExists)

	err := service.AddSongToPlaylist(playlistID, songID)

	assert.EqualError(t, err, ErrPlaylistAlreadyExists.Error())
	mockDAO.AssertExpectations(t)
}

func TestGetAllPlaylists_Success(t *testing.T) {
	mockDAO := &mocks.MusicDAO{}
	service := NewPlaylistService(mockDAO)

	mockPlaylists := []model.Playlist{
		{Name: "Chill Hits", UserID: 1},
		{Name: "Workout Jams", UserID: 1},
	}

	mockDAO.On("GetAllPlaylists").Return(mockPlaylists, nil)

	playlists, err := service.GetAllPlaylists()

	assert.NoError(t, err)
	assert.Equal(t, mockPlaylists, playlists)
	mockDAO.AssertExpectations(t)
}

func TestGetAllPlaylists_AuthError(t *testing.T) {
	mockDAO := &mocks.MusicDAO{}
	service := NewPlaylistService(mockDAO)

	mockDAO.On("GetAllPlaylists").Return(nil, ErrUserDoesNotExist)

	playlists, err := service.GetAllPlaylists()

	assert.EqualError(t, err, ErrUserDoesNotExist.Error())
	assert.Nil(t, playlists)
	mockDAO.AssertExpectations(t)
}

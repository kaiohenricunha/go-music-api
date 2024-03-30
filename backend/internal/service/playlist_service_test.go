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

func TestCreatePlaylist_EmptyName(t *testing.T) {
	mockDAO := &mocks.MusicDAO{}
	service := NewPlaylistService(mockDAO)

	playlist := &model.Playlist{
		UserID: 1,
	}

	// No need to setup mock expectations here since the method should return early

	err := service.CreatePlaylist(playlist)

	assert.EqualError(t, err, "playlist name cannot be empty")
}

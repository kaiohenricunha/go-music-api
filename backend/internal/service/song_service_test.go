package service

import (
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao/mocks"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllSongs(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	songService := NewSongService(mockDAO)

	mockSongs := []model.Song{
		{Name: "Test Song 1", Artist: "Test Artist 1"},
		{Name: "Test Song 2", Artist: "Test Artist 2"},
	}
	mockDAO.On("GetAllSongs").Return(mockSongs, nil)

	songs, err := songService.GetAllSongs()
	assert.Nil(t, err)
	assert.Equal(t, mockSongs, songs)
	mockDAO.AssertExpectations(t)
}

func TestUpdateSong(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	songService := NewSongService(mockDAO)

	mockSong := &model.Song{Name: "Test Song", Artist: "Test Artist"}
	mockSong.ID = 1
	mockDAO.On("UpdateSong", mockSong).Return(nil)

	err := songService.UpdateSong(mockSong)
	assert.Nil(t, err)
	mockDAO.AssertExpectations(t)
}

func TestDeleteSong(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	songService := NewSongService(mockDAO)

	mockSongID := uint(1)
	mockDAO.On("DeleteSong", mockSongID).Return(nil)

	err := songService.DeleteSong(mockSongID)
	assert.Nil(t, err)
	mockDAO.AssertExpectations(t)
}

// TestAddSongError tests the AddSong function when an error occurs.
func TestAddSongError(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	songService := NewSongService(mockDAO)

	mockSong := &model.Song{Name: "Test Song", Artist: "Test Artist"}

	// Simulate finding an existing song to trigger an error
	mockDAO.On("FindSongByName", "Test Song").Return(mockSong, nil)

	err := songService.AddSong(mockSong)
	assert.Equal(t, ErrSongAlreadyExists, err)
	mockDAO.AssertExpectations(t)
}

// TestGetAllSongsError tests the GetAllSongs function when an error occurs.
func TestGetAllSongsError(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	songService := NewSongService(mockDAO)

	mockDAO.On("GetAllSongs").Return(nil, assert.AnError)

	songs, err := songService.GetAllSongs()
	assert.NotNil(t, err)
	assert.Nil(t, songs)
	mockDAO.AssertExpectations(t)
}

// TestUpdateSongError tests the UpdateSong function when an error occurs.
func TestUpdateSongError(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	songService := NewSongService(mockDAO)

	mockSong := &model.Song{Name: "Test Song", Artist: "Test Artist"}
	mockSong.ID = 1
	mockDAO.On("UpdateSong", mockSong).Return(assert.AnError)

	err := songService.UpdateSong(mockSong)
	assert.NotNil(t, err)
	mockDAO.AssertExpectations(t)
}

// TestDeleteSongError tests the DeleteSong function when an error occurs.
func TestDeleteSongError(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	songService := NewSongService(mockDAO)

	mockSongID := uint(1)
	mockDAO.On("DeleteSong", mockSongID).Return(assert.AnError)

	err := songService.DeleteSong(mockSongID)
	assert.NotNil(t, err)
	mockDAO.AssertExpectations(t)
}

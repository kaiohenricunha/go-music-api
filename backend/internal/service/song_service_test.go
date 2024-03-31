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

package service

import (
	"testing"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao/mocks"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSong(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	songService := NewSongService(mockDAO)
	testSong := &model.Song{Name: "Test Song", Artist: "Test Artist"}

	t.Run("success", func(t *testing.T) {
		mockDAO.On("GetSongByNameAndArtist", testSong.Name, testSong.Artist).Return(nil, nil).Once()
		mockDAO.On("CreateSong", mock.AnythingOfType("*model.Song")).Return(nil).Once()

		err := songService.CreateSong(testSong)
		assert.NoError(t, err)

		mockDAO.AssertExpectations(t)
	})

	t.Run("song exists", func(t *testing.T) {
		mockDAO.On("GetSongByNameAndArtist", testSong.Name, testSong.Artist).Return(testSong, nil).Once()

		err := songService.CreateSong(testSong)
		assert.Equal(t, ErrSongAlreadyExists, err)

		mockDAO.AssertExpectations(t)
	})

	t.Run("missing required fields", func(t *testing.T) {
		err := songService.CreateSong(&model.Song{Name: "", Artist: ""})
		assert.Equal(t, ErrSongNameRequired, err)
	})
}

func TestGetSongByID(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	songService := NewSongService(mockDAO)
	testID := "1"
	testSong := &model.Song{Name: "Test Song", Artist: "Test Artist"}

	mockDAO.On("GetSongByID", testID).Return(testSong, nil) // Song found
	resultSong, err := songService.GetSongByID(testID)
	assert.NoError(t, err)
	assert.Equal(t, testSong, resultSong)

	mockDAO.On("GetSongByID", "non-existent-id").Return(nil, ErrSongNotFound) // Song not found
	_, err = songService.GetSongByID("non-existent-id")
	assert.Equal(t, ErrSongNotFound, err)
}

func TestGetSongByNameAndArtist(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	songService := NewSongService(mockDAO)
	testSong := &model.Song{Name: "Test Song", Artist: "Test Artist"}

	mockDAO.On("GetSongByNameAndArtist", testSong.Name, testSong.Artist).Return(testSong, nil) // Song exists
	foundSong, err := songService.GetSongByNameAndArtist(testSong.Name, testSong.Artist)
	assert.NoError(t, err)
	assert.Equal(t, testSong, foundSong)

	mockDAO.On("GetSongByNameAndArtist", "unknown", "unknown").Return(nil, ErrSongNotFound) // Song does not exist
	_, err = songService.GetSongByNameAndArtist("unknown", "unknown")
	assert.Equal(t, ErrSongNotFound, err)
}

// // TestGetSongFromSpotifyByID tests the GetSongFromSpotifyByID method
// func TestGetSongFromSpotifyByID(t *testing.T) {
// 	mockDAO := new(mocks.MusicDAO)
// 	songService := NewSongService(mockDAO)
// 	testSong := &model.Song{Name: "Test Song", Artist: "Test Artist", SpotifyID: "test-id"}

// 	mockDAO.On("GetSongFromSpotifyByID", testSong.SpotifyID).Return(testSong, nil) // Song found
// 	foundSong, err := songService.GetSongFromSpotifyByID(testSong.SpotifyID)
// 	assert.NoError(t, err)
// 	assert.Equal(t, testSong, foundSong)

// 	mockDAO.On("GetSongFromSpotifyByID", "unknown-id").Return(nil, ErrSongNotFound) // Song not found
// 	_, err = songService.GetSongFromSpotifyByID("unknown-id")
// 	assert.Equal(t, ErrSongNotFound, err)
// }

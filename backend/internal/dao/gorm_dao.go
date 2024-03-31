package dao

import (
	"errors"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"gorm.io/gorm"
)

type GormDAO struct {
	DB *gorm.DB
}

// NewGormDAO creates a new instance of GormDAO.
func NewGormDAO(db *gorm.DB) *GormDAO {
	return &GormDAO{DB: db}
}

//////////////////////
// USER METHODS //
//////////////////////

var (
	ErrUserNotFound = errors.New("user not found")
	ErrSongNotFound = errors.New("song not found")
)

// CreateUser checks for a soft-deleted user with the same username and permanently deletes it before creating a new one.
func (g *GormDAO) CreateUser(user *model.User) error {
	var existingUser model.User
	// Check for a soft-deleted user with the same username.
	result := g.DB.Unscoped().Where("username = ?", user.Username).First(&existingUser)
	if result.Error == nil && existingUser.DeletedAt.Valid {
		// If a soft-deleted record exists, permanently delete it to free up the username.
		if err := g.DB.Unscoped().Delete(&existingUser).Error; err != nil {
			return err
		}
	} else if result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	return g.DB.Create(user).Error
}

// UpdateUser updates an existing user's information in the database.
func (g *GormDAO) UpdateUser(user *model.User) error {
	var existingUser model.User
	result := g.DB.Where("id = ?", user.ID).First(&existingUser)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrUserNotFound
	}

	// Update the user with new values where applicable
	if user.Username != "" {
		existingUser.Username = user.Username
	}

	// The password field should already be hashed if it was updated,
	// as handled in the UserService.UpdateUser method.
	if user.Password != "" {
		existingUser.Password = user.Password
	}

	// Save updates
	return g.DB.Save(&existingUser).Error
}

// FindUserByID retrieves a single user by ID.
func (g *GormDAO) FindUserByID(userID uint) (*model.User, error) {
	var user model.User
	err := g.DB.Where("id = ?", userID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

// DeleteUser leverages GORM's soft delete functionality, which is automatically applied if the model includes a `DeletedAt` field.
func (g *GormDAO) DeleteUser(userID uint) error {
	return g.DB.Delete(&model.User{}, userID).Error
}

// GetAllUsers retrieves all users with their playlists.
func (g *GormDAO) GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := g.DB.Preload("Playlists").Find(&users).Error
	return users, err
}

// FindByUsername finds a single user by username with playlists.
func (g *GormDAO) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := g.DB.Preload("Playlists").Where("username = ?", username).First(&user).Error
	return &user, err
}

//////////////////////
// SONG METHODS //
//////////////////////

// GetAllSongs retrieves all songs from the database.
func (g *GormDAO) GetAllSongs() ([]model.Song, error) {
	var songs []model.Song
	err := g.DB.Find(&songs).Error
	return songs, err
}

func (g *GormDAO) GetSongFromSpotifyByID(spotifyID string) (*model.Song, error) {
	var song model.Song
	err := g.DB.Where("spotify_id = ?", spotifyID).First(&song).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSongNotFound
	}
	return &song, err
}

func (g *GormDAO) SearchSongsFromSpotify(trackName, artistName string) ([]model.Song, error) {
	var songs []model.Song
	err := g.DB.Where("track_name = ? AND artist_name = ?", trackName, artistName).Find(&songs).Error
	return songs, err
}

//////////////////////
// PLAYLIST METHODS //
//////////////////////

// CreatePlaylist inserts a new playlist into the database.
func (g *GormDAO) CreatePlaylist(playlist *model.Playlist) error {
	return g.DB.Create(playlist).Error
}

// AddSongToPlaylist inserts a new record into the playlist_songs table to associate a song with a playlist.
func (g *GormDAO) AddSongToPlaylist(playlistID, songID uint) error {
	return g.DB.Exec("INSERT INTO playlist_songs (playlist_id, song_id) VALUES (?, ?)", playlistID, songID).Error
}

// GetPlaylistByNameAndUserID retrieves a single playlist by name and user ID.
func (g *GormDAO) GetPlaylistByNameAndUserID(playlistName string, userID uint) (*model.Playlist, error) {
	var playlist model.Playlist
	err := g.DB.Where("name = ? AND user_id = ?", playlistName, userID).First(&playlist).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &playlist, err
}

// GetAllPlaylists retrieves all playlists from the database.
func (g *GormDAO) GetAllPlaylists() ([]model.Playlist, error) {
	var playlists []model.Playlist
	err := g.DB.Find(&playlists).Error
	return playlists, err
}

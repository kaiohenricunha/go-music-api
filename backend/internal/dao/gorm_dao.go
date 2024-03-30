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

// GetAllUsers retrieves all users from the database.
func (g *GormDAO) GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := g.DB.Find(&users).Error
	return users, err
}

// FindByUsername finds a single user by username.
func (g *GormDAO) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := g.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// CreateSong inserts a new song into the database.
func (g *GormDAO) CreateSong(song *model.Song) error {
	return g.DB.Create(song).Error
}

// UpdateSong updates an existing song's information in the database.
func (g *GormDAO) UpdateSong(song *model.Song) error {
	return g.DB.Model(&model.Song{}).Where("id = ?", song.ID).Updates(song).Error
}

// DeleteSong removes a song from the database by ID.
func (g *GormDAO) DeleteSong(songID uint) error {
	return g.DB.Delete(&model.Song{}, songID).Error
}

// GetAllSongs retrieves all songs from the database.
func (g *GormDAO) GetAllSongs() ([]model.Song, error) {
	var songs []model.Song
	err := g.DB.Find(&songs).Error
	return songs, err
}

func (g *GormDAO) FindSongByName(songName string) (*model.Song, error) {
	var song model.Song
	err := g.DB.Where("name = ?", songName).First(&song).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &song, err
}

func (g *GormDAO) FindSongByID(songID uint) (*model.Song, error) {
	var song model.Song
	err := g.DB.Where("id = ?", songID).First(&song).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &song, err
}

func (g *GormDAO) CreatePlaylist(playlist *model.Playlist) error {
	return g.DB.Create(playlist).Error
}

func (g *GormDAO) UpdatePlaylist(playlist *model.Playlist) error {
	return g.DB.Model(&model.Playlist{}).Where("id = ?", playlist.ID).Updates(playlist).Error
}

func (g *GormDAO) DeletePlaylist(playlistID uint) error {
	return g.DB.Delete(&model.Playlist{}, playlistID).Error
}

func (g *GormDAO) GetAllPlaylists() ([]model.Playlist, error) {
	var playlists []model.Playlist
	err := g.DB.Find(&playlists).Error
	return playlists, err
}

// GetPlaylistByID retrieves a playlist by its ID
func (g *GormDAO) GetPlaylistByID(id uint) (*model.Playlist, error) {
	var playlist model.Playlist
	result := g.DB.Preload("Songs").First(&playlist, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &playlist, nil
}

func (g *GormDAO) GetPlaylistsByUserID(userID uint) ([]model.Playlist, error) {
	var playlists []model.Playlist
	err := g.DB.Where("user_id = ?", userID).Find(&playlists).Error
	return playlists, err
}

func (g *GormDAO) GetPlaylistByName(playlistName string) (*model.Playlist, error) {
	var playlist model.Playlist
	err := g.DB.Where("name = ?", playlistName).First(&playlist).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &playlist, err
}

func (dao *GormDAO) GetPlaylistByNameAndUserID(name string, userID uint) (*model.Playlist, error) {
	var playlist model.Playlist
	if err := dao.DB.Where("name = ? AND user_id = ?", name, userID).First(&playlist).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &playlist, nil
}

func (g *GormDAO) AddSongToPlaylist(playlistID, songID uint) error {
	return g.DB.Exec("INSERT INTO playlist_songs (playlist_id, song_id) VALUES (?, ?)", playlistID, songID).Error
}

func (g *GormDAO) RemoveSongFromPlaylist(playlistID, songID uint) error {
	return g.DB.Exec("DELETE FROM playlist_songs WHERE playlist_id = ? AND song_id = ?", playlistID, songID).Error
}

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
	ErrUserNotFound     = errors.New("user not found")
	ErrSongNotFound     = errors.New("song not found")
	ErrPlaylistNotFound = errors.New("playlist not found")
	ErrRecordNotFound   = gorm.ErrRecordNotFound
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

// GetUserByID retrieves a single user by ID.
func (g *GormDAO) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	err := g.DB.Where("id = ?", userID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

// GetAllUsers retrieves all users with their playlists.
func (g *GormDAO) GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := g.DB.Preload("Playlists").Find(&users).Error
	return users, err
}

// GetUserByUsername retrieves a single user by username.
func (g *GormDAO) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := g.DB.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

//////////////////////
// SONG METHODS //
//////////////////////

// CreateSong inserts a new song into the database.
func (g *GormDAO) CreateSong(song *model.Song) error {
	return g.DB.Create(song).Error
}

// GetAllSongs retrieves all songs from the database.
func (g *GormDAO) GetAllSongs() ([]model.Song, error) {
	var songs []model.Song
	err := g.DB.Find(&songs).Error
	return songs, err
}

// GetSongByID retrieves a single song by ID.
func (g *GormDAO) GetSongByID(songID string) (*model.Song, error) {
	var song model.Song
	err := g.DB.Where("id = ?", songID).First(&song).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSongNotFound
	}
	return &song, err
}

// GetSongByNameAndArtist retrieves a single song by name and artist.
func (g *GormDAO) GetSongByNameAndArtist(songName, artistName string) (*model.Song, error) {
	var song model.Song
	err := g.DB.Where("song_name = ? AND artist_name = ?", songName, artistName).First(&song).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSongNotFound
	}
	return &song, err
}

// GetSongFromSpotifyByID retrieves a single song by Spotify ID.
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
	err := g.DB.Where("song_name = ? AND artist_name = ?", trackName, artistName).Find(&songs).Error
	return songs, err
}

//////////////////////
// PLAYLIST METHODS //
//////////////////////

// GetPlaylistByID retrieves a single playlist by ID.
func (g *GormDAO) GetPlaylistByID(playlistID string) (*model.Playlist, error) {
	var playlist model.Playlist
	err := g.DB.Preload("Songs").Preload("Ratings").Where("id = ?", playlistID).First(&playlist).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPlaylistNotFound
	}
	return &playlist, err
}

// GetAllPlaylists retrieves all playlists from the database.
func (g *GormDAO) GetAllPlaylists() ([]model.Playlist, error) {
	var playlists []model.Playlist
	err := g.DB.Preload("Songs").Preload("Ratings").Find(&playlists).Error
	if err != nil {
		return nil, err
	}
	return playlists, nil
}

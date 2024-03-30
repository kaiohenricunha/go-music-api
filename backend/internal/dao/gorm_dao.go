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

// CreateUser inserts a new user into the database.
func (g *GormDAO) CreateUser(user *model.User) error {
	return g.DB.Create(user).Error
}

// UpdateUser updates an existing user's information in the database.
func (g *GormDAO) UpdateUser(user *model.User) error {
	return g.DB.Save(user).Error
}

// DeleteUser removes a user from the database by ID.
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
	return g.DB.Save(song).Error
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

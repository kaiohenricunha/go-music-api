package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string // Hashed password
}

type Song struct {
	gorm.Model
	Name   string `json:"name"`
	Artist string `json:"artist"`
}

type Playlist struct {
	gorm.Model
	Name    string
	UserID  uint
	Songs   []Song   `gorm:"many2many:playlist_songs;"`
	Ratings []Rating // Optional: direct access to ratings from a playlist
}

type Rating struct {
	gorm.Model
	PlaylistID uint
	UserID     uint
	Score      int // Numerical rating
}

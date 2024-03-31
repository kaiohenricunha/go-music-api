package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string     `gorm:"unique"`
	Password  string     // Hashed password
	Playlists []Playlist `gorm:"foreignKey:UserID"`
}

type Song struct {
	gorm.Model
	SpotifyID     string `json:"spotify_id"`
	Name          string `json:"name"`
	Artist        string `json:"artist"`
	AlbumName     string `json:"album_name"`
	AlbumImageURL string `json:"album_image_url"`
	PreviewURL    string `json:"preview_url"`
	ExternalURL   string `json:"external_url"` // URL to Spotify
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

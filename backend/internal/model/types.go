package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName  string     `json:"full_name"`
	Email     string     `gorm:"unique" json:"email"`
	Username  string     `gorm:"unique" json:"username"`
	Password  string     // Consider storing hashed passwords only
	Role      string     `json:"role"`
	Playlists []Playlist `gorm:"foreignKey:UserID" json:"playlists"`
}

type Song struct {
	gorm.Model
	SpotifyID     string `gorm:"column:spotify_id" json:"spotify_id"`
	Name          string `gorm:"column:song_name" json:"name"`
	Artist        string `gorm:"column:artist_name" json:"artist"`
	AlbumName     string `gorm:"column:album_name" json:"album_name"`
	AlbumImageURL string `gorm:"column:album_image_url" json:"album_image_url"`
	PreviewURL    string `gorm:"column:preview_url" json:"preview_url"`
	ExternalURL   string `gorm:"column:external_url" json:"external_url"`
}

type Playlist struct {
	gorm.Model
	Name             string   `gorm:"column:playlist_name" json:"playlist_name"`
	UserID           uint     `gorm:"column:user_id" json:"user_id"`
	PlaylistImageURL string   `gorm:"column:playlist_image_url" json:"playlist_image_url"`
	Songs            []Song   `gorm:"many2many:playlist_songs;"`
	Ratings          []Rating `gorm:"foreignKey:PlaylistID" json:"ratings"`
}

type Rating struct {
	gorm.Model
	PlaylistID string `json:"playlist_id"`
	UserID     string `json:"user_id"`
	Score      int    `json:"score" gorm:"column:score"`
}

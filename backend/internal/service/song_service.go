package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"gorm.io/gorm"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"golang.org/x/oauth2/clientcredentials"
)

// Error definitions
var (
	ErrSongNameRequired    = errors.New("song name and artist are required")
	ErrSongAlreadyExists   = errors.New("a song with the same name by the same artist already exists")
	ErrSongNotFound        = errors.New("song not found")
	ErrInvalidSongID       = errors.New("invalid song ID")
	ErrFetchingFromSpotify = errors.New("error fetching data from Spotify")
)

type SongService interface {
	CreateSong(song *model.Song) error
	GetAllSongs() ([]model.Song, error)
	GetSongByID(id string) (*model.Song, error)
	GetSongByNameAndArtist(name, artist string) (*model.Song, error)
	GetSongFromSpotifyByID(spotifyID string) (*model.Song, error)
	SearchSongsFromSpotify(trackName, artistName string) ([]*model.Song, error)
}

type songService struct {
	songDAO    dao.MusicDAO
	httpClient *http.Client // HTTP client for Spotify API requests
}

// SpotifyTrackResponse struct for unmarshalling Spotify API response.
type SpotifyTrackResponse struct {
	Tracks struct {
		Items []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Album struct {
				Name   string `json:"name"`
				Images []struct {
					URL string `json:"url"`
				} `json:"images"`
			} `json:"album"`
			PreviewURL   string `json:"preview_url"`
			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
		} `json:"items"`
	} `json:"tracks"`
}

func NewSongService(songDAO dao.MusicDAO) SongService {
	// Configure the client for Spotify API authentication
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"), // these are set in the environment or Kubernetes secrets
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		TokenURL:     "https://accounts.spotify.com/api/token",
	}
	httpClient := config.Client(context.Background())

	return &songService{songDAO: songDAO, httpClient: httpClient}
}

// CreateSong creates a new song in the database.
func (s *songService) CreateSong(song *model.Song) error {
	// Validate the song name and artist
	if song.Name == "" || song.Artist == "" {
		return ErrSongNameRequired
	}

	// Check if a song with the same name by the same artist already exists
	existingSong, err := s.songDAO.GetSongByNameAndArtist(song.Name, song.Artist)
	if err != nil {
		return err
	}
	if existingSong != nil {
		return ErrSongAlreadyExists
	}

	// Create the song in the database
	return s.songDAO.CreateSong(song)
}

// GetSongByID retrieves a song by its ID.
func (s *songService) GetSongByID(id string) (*model.Song, error) {
	if id == "" {
		return nil, ErrInvalidSongID
	}

	return s.songDAO.GetSongByID(id)
}

func (s *songService) GetAllSongs() ([]model.Song, error) {
	return s.songDAO.GetAllSongs()
}

func (s *songService) GetSongByNameAndArtist(name, artist string) (*model.Song, error) {
	return s.songDAO.GetSongByNameAndArtist(name, artist)
}

// New method to get a song from Spotify by its ID
func (s *songService) GetSongFromSpotifyByID(spotifyID string) (*model.Song, error) {
	// Construct the request to the Spotify Web API
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", spotifyID), nil)
	if err != nil {
		return nil, ErrFetchingFromSpotify
	}

	// Make the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, ErrFetchingFromSpotify
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrFetchingFromSpotify
	}

	// Assuming you have a function to parse the Spotify API response into your Song model
	var song model.Song
	if _, err := s.parseSpotifyResponse(resp); err != nil {
		return nil, err
	}

	return &song, nil
}

func (s *songService) SearchSongsFromSpotify(trackName, artistName string) ([]*model.Song, error) {
	songs, err := s.searchLocalDatabase(trackName, artistName)
	if err != nil {
		if err == ErrSongNotFound {
			return s.searchOnSpotify(trackName, artistName) // Proceed to search on Spotify
		}
		return nil, err // Return all other errors immediately
	}
	return songs, nil // Return songs found locally
}

func (s *songService) searchLocalDatabase(trackName, artistName string) ([]*model.Song, error) {
	log.Printf("Searching for song in local database: %s by %s", trackName, artistName)

	song, err := s.songDAO.GetSongByNameAndArtist(trackName, artistName)
	if err != nil {
		// Here, we need to specifically check if the error is due to the record not being found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No song found in local database for: %s by %s", trackName, artistName)
			return nil, ErrSongNotFound
		}
		log.Printf("Error searching song in local database: %v", err)
		return nil, err // Return other database-related errors
	}
	return []*model.Song{song}, nil
}

func (s *songService) searchOnSpotify(trackName, artistName string) ([]*model.Song, error) {
	log.Printf("Searching for song on Spotify: %s by %s", trackName, artistName)

	query := url.QueryEscape(fmt.Sprintf("track:%s artist:%s", trackName, artistName))
	requestURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", query)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Printf("Failed to create request for Spotify: %v", err)
		return nil, ErrFetchingFromSpotify
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := s.httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Failed to execute request on Spotify: %v, status code: %d", err, resp.StatusCode)
		return nil, ErrFetchingFromSpotify
	}
	defer resp.Body.Close()

	return s.parseSpotifyResponse(resp)
}

func (s *songService) parseSpotifyResponse(resp *http.Response) ([]*model.Song, error) {
	var spotifyResponse SpotifyTrackResponse
	if err := json.NewDecoder(resp.Body).Decode(&spotifyResponse); err != nil {
		log.Printf("Failed to decode Spotify response: %v", err)
		return nil, err
	}

	if len(spotifyResponse.Tracks.Items) == 0 {
		log.Printf("No tracks found on Spotify for the given query")
		return nil, ErrSongNotFound
	}

	var songs []*model.Song
	for _, item := range spotifyResponse.Tracks.Items {
		song := &model.Song{
			SpotifyID:     item.ID,
			Name:          item.Name,
			Artist:        item.Artists[0].Name,
			AlbumName:     item.Album.Name,
			AlbumImageURL: albumImageURLFrom([]struct{ URL string }(item.Album.Images)),
			PreviewURL:    item.PreviewURL,
			ExternalURL:   item.ExternalURLs.Spotify,
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func albumImageURLFrom(images []struct{ URL string }) string {
	if len(images) > 0 {
		return images[0].URL
	}
	return ""
}

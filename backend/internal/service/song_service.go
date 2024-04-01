package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"golang.org/x/oauth2/clientcredentials"
)

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
	SearchSongsFromSpotify(trackName, artistName string) ([]model.Song, error)
}

type songService struct {
	songDAO    dao.MusicDAO
	httpClient *http.Client // HTTP client for Spotify API requests
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
	err = parseSpotifyResponseToSongModel(resp.Body, &song)
	if err != nil {
		return nil, ErrFetchingFromSpotify
	}

	return &song, nil
}

func (s *songService) SearchSongsFromSpotify(trackName, artistName string) ([]model.Song, error) {
	// Attempt to find the song in the local database first.
	log.Printf("Searching for song: %s by artist: %s in the database", trackName, artistName)
	existingSong, err := s.songDAO.GetSongByNameAndArtist(trackName, artistName)
	if err == nil {
		// If the song is found, return it in a slice.
		log.Printf("Song found in the database: %s by %s", existingSong.Name, existingSong.Artist)
		return []model.Song{*existingSong}, nil
	} else if err != ErrSongNotFound {
		// If an unexpected error occurs, it means the song is not yet in the local database.
		// Log the error and proceed to search on Spotify.
		log.Printf("Error searching for song in the database: %v", err)
	}

	// Song not found locally, proceed to search on Spotify.
	log.Printf("Searching for song: %s by artist: %s on Spotify", trackName, artistName)
	query := url.QueryEscape(fmt.Sprintf("track:%s artist:%s", trackName, artistName))
	requestURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", query)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, ErrFetchingFromSpotify
	}

	resp, err := s.httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, ErrFetchingFromSpotify
	}
	defer resp.Body.Close()

	var spotifyResponse SpotifyTrackResponse
	if err := json.NewDecoder(resp.Body).Decode(&spotifyResponse); err != nil {
		return nil, err
	}

	var songs []model.Song
	for _, item := range spotifyResponse.Tracks.Items {
		// Assuming only the first artist is relevant.
		firstArtistName := ""
		if len(item.Artists) > 0 {
			firstArtistName = item.Artists[0].Name
		}

		song := model.Song{
			SpotifyID:     item.ID,
			Name:          item.Name,
			Artist:        firstArtistName,
			AlbumName:     item.Album.Name,
			AlbumImageURL: albumImageURLFrom([]struct{ URL string }(item.Album.Images)),
			PreviewURL:    item.PreviewURL,
			ExternalURL:   item.ExternalURLs.Spotify,
		}
		songs = append(songs, song)

		// Add song to the database, assuming it wasn't found earlier.
		log.Printf("Caching song '%s' by '%s' to the database", song.Name, song.Artist)
		if dbErr := s.songDAO.CreateSong(&song); dbErr != nil {
			log.Printf("Error adding song '%s' by '%s' to the database: %v", song.Name, song.Artist, dbErr)
			// Opt to continue processing other songs despite the error.
		}
	}

	return songs, nil
}

// Helper function to extract the album image URL from the Spotify response.
func albumImageURLFrom(images []struct{ URL string }) string {
	if len(images) > 0 {
		return images[0].URL
	}
	return ""
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

// GetSongByID retrieves a song from DB by its ID
func (s *songService) GetSongByID(id string) (*model.Song, error) {
	// Validate the song ID
	if id == "" {
		return nil, ErrInvalidSongID
	}

	return s.songDAO.GetSongByID(id)
}

// parseSpotifyResponseToSongModel takes an io.Reader (the body of the HTTP response) and populates the provided Song model with data parsed from the response.
func parseSpotifyResponseToSongModel(body io.Reader, song *model.Song) error {
	type SpotifyTrack struct {
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
	}

	var spotifyTrack SpotifyTrack
	if err := json.NewDecoder(body).Decode(&spotifyTrack); err != nil {
		return err
	}

	var artistNames []string
	for _, artist := range spotifyTrack.Artists {
		artistNames = append(artistNames, artist.Name)
	}

	var albumImageURL string
	if len(spotifyTrack.Album.Images) > 0 {
		albumImageURL = spotifyTrack.Album.Images[0].URL
	}

	// Populate the Song model with the data from the Spotify track
	song.SpotifyID = spotifyTrack.ID
	song.Name = spotifyTrack.Name
	song.Artist = strings.Join(artistNames, ", ")
	song.AlbumName = spotifyTrack.Album.Name
	song.AlbumImageURL = albumImageURL
	song.PreviewURL = spotifyTrack.PreviewURL
	song.ExternalURL = spotifyTrack.ExternalURLs.Spotify

	return nil
}

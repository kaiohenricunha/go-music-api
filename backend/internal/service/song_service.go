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
	ErrSongNameRequired  = errors.New("song name and artist are required")
	ErrSongAlreadyExists = errors.New("a song with the same name by the same artist already exists")
	ErrSongNotFound      = errors.New("song not found")
	// Adding Spotify API errors
	ErrFetchingFromSpotify = errors.New("error fetching data from Spotify")
)

type SongService interface {
	GetAllSongs() ([]model.Song, error)
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

func (s *songService) GetAllSongs() ([]model.Song, error) {
	return s.songDAO.GetAllSongs()
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

// SearchSongsFromSpotify searches for songs on Spotify based on track name and artist name.
func (s *songService) SearchSongsFromSpotify(trackName, artistName string) ([]model.Song, error) {
	// Check if the track name and artist name are provided
	log.Printf("Searching for song: %s by artist: %s", trackName, artistName)

	// Construct the search query with track name and artist name
	query := url.QueryEscape(fmt.Sprintf("track:%s artist:%s", trackName, artistName))

	// Construct the request URL with the encoded query
	requestURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track&include_external=audio", query)

	log.Printf("Searching Song Request URL: %s\n", requestURL)

	// Create and send the request to Spotify's search API
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, ErrFetchingFromSpotify
	}

	// Make the request using the httpClient
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, ErrFetchingFromSpotify
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err == nil {
			fmt.Printf("Spotify Error Response: %s\n", string(bodyBytes))
		}
		return nil, ErrFetchingFromSpotify
	}

	// Define a struct to unmarshal the search results
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

	var spotifyResponse SpotifyTrackResponse
	if err := json.NewDecoder(resp.Body).Decode(&spotifyResponse); err != nil {
		return nil, err
	}

	var songs []model.Song
	for _, item := range spotifyResponse.Tracks.Items {
		var artistNames []string
		for _, artist := range item.Artists {
			artistNames = append(artistNames, artist.Name)
		}

		var albumImageURL string
		if len(item.Album.Images) > 0 {
			albumImageURL = item.Album.Images[0].URL
		}

		song := model.Song{
			SpotifyID:     item.ID,
			Name:          item.Name,
			Artist:        strings.Join(artistNames, ", "),
			AlbumName:     item.Album.Name,
			AlbumImageURL: albumImageURL,
			PreviewURL:    item.PreviewURL,
			ExternalURL:   item.ExternalURLs.Spotify,
		}

		songs = append(songs, song)
	}

	return songs, nil
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

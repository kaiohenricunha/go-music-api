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

// parseSpotifyResponseToSongModel takes an io.Reader (the body of the HTTP response) and a pointer to a Song model,
// and populates the Song model with data parsed from the response.
func parseSpotifyResponseToSongModel(body io.Reader, song *model.Song) error {
	type SpotifyResponse struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Album struct {
			Name   string `json:"name"`
			Images []struct {
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
		} `json:"album"`
		PreviewURL   string `json:"preview_url"`
		ExternalURLs struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
	}

	var spotifyTrack SpotifyResponse
	if err := json.NewDecoder(body).Decode(&spotifyTrack); err != nil {
		return err
	}

	// Assuming you want to concatenate the artist names if there are multiple
	var artistNames []string
	for _, artist := range spotifyTrack.Artists {
		artistNames = append(artistNames, artist.Name)
	}

	// For the album image, you might want to choose one based on a specific size criteria.
	// Here, we'll simply take the first one.
	var albumImageURL string
	if len(spotifyTrack.Album.Images) > 0 {
		albumImageURL = spotifyTrack.Album.Images[0].URL // This could be adjusted based on your needs
	}

	// Populate your Song model with the data from the Spotify track
	song.SpotifyID = spotifyTrack.ID
	song.Name = spotifyTrack.Name
	song.Artist = strings.Join(artistNames, ", ")
	song.AlbumName = spotifyTrack.Album.Name
	song.AlbumImageURL = albumImageURL
	song.PreviewURL = spotifyTrack.PreviewURL
	song.ExternalURL = spotifyTrack.ExternalURLs.Spotify

	return nil
}

// SearchSongsFromSpotify searches for songs on Spotify based on track name and artist name.
func (s *songService) SearchSongsFromSpotify(trackName, artistName string) ([]model.Song, error) {
	log.Println("Entering SearchSongsFromSpotify") // Debug log

	// Early log to check parameter values
	log.Printf("Received trackName: %s, artistName: %s\n", trackName, artistName)

	// Construct the search query with track name and artist name
	query := fmt.Sprintf("track:%s artist:%s", url.QueryEscape(trackName), url.QueryEscape(artistName))

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
	var searchResults struct {
		Tracks struct {
			Items []model.Song `json:"items"`
		} `json:"tracks"`
	}

	// Decode the response body into the struct
	if err := json.NewDecoder(resp.Body).Decode(&searchResults); err != nil {
		return nil, err
	}

	// Extract the items from the search results and return them
	return searchResults.Tracks.Items, nil
}

package musicapp

import (
	"github.com/kaiohenricunha/go-music-k8s/backend/config"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"testing"
)

func TestPostSong(t *testing.T) {
	db := SetupTestDB()     // Initialize and migrate your test database
	defer clearDatabase(db) // Ensure the database is cleared after the test runs.
	cfg := &config.Config{DB: db}

	song := &Song{Name: "Test Song", Artist: "Test Artist"}

	err := PostSong(cfg, song)
	assert.NoError(t, err, "Failed to post a new song")
}

// TestGetAllSongs tests fetching all songs from the database
func TestGetAllSongs(t *testing.T) {
	db := SetupTestDB()
	defer clearDatabase(db) // Ensure the database is cleared after the test runs.
	cfg := &config.Config{DB: db}

	// Prepopulate the database with test data
	db.Create(&Song{Name: "Test Song 1", Artist: "Test Artist 1"})
	db.Create(&Song{Name: "Test Song 2", Artist: "Test Artist 2"})

	songs, err := GetAllSongs(cfg)
	assert.NoError(t, err)
	assert.NotEmpty(t, songs)
	assert.Equal(t, 2, len(songs), "Expected 2 songs in the database")
}

// TestUpdateSong tests updating a song's details in the database
func TestUpdateSong(t *testing.T) {
	db := SetupTestDB()
	defer clearDatabase(db) // Ensure the database is cleared after the test runs.
	cfg := &config.Config{DB: db}

	// Prepopulate the database with a test song
	db.Create(&Song{Model: gorm.Model{ID: 1}, Name: "Old Name", Artist: "Old Artist"})

	// Update the song
	updatedSong := &Song{Model: gorm.Model{ID: 1}, Name: "New Name", Artist: "New Artist"}
	err := UpdateSong(cfg, updatedSong)
	assert.NoError(t, err)

	// Verify the song was updated
	var song Song
	db.First(&song, 1)
	assert.Equal(t, "New Name", song.Name)
	assert.Equal(t, "New Artist", song.Artist)
}

// TestDeleteSong tests deleting a song from the database
func TestDeleteSong(t *testing.T) {
	db := SetupTestDB()
	defer clearDatabase(db) // Ensure the database is cleared after the test runs.
	cfg := &config.Config{DB: db}

	// Prepopulate the database with a test song
	db.Create(&Song{Name: "Test Song", Artist: "Test Artist"})

	// Delete the song
	err := DeleteSong(cfg, 1) // Assuming the song ID is 1
	assert.NoError(t, err)

	// Verify the song was deleted
	var songs []Song
	db.Find(&songs)
	assert.Empty(t, songs, "Expected no songs in the database after deletion")
}

// SetupTestDB initializes a new database connection for testing purposes,
// migrates necessary schemas, and returns the connection.
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}

	// Migrate the schema. Add all necessary migrations here.
	err = db.AutoMigrate(&Song{})
	if err != nil {
		log.Panicf("Failed to migrate database: %v", err)
	}

	return db
}

func clearDatabase(db *gorm.DB) {
	db.Exec("DELETE FROM songs")
	db.Exec("DELETE FROM sqlite_sequence WHERE name='songs'") // If using SQLite; resets auto-increment.
}

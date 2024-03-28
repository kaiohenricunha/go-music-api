package musicapp

import (
	"log"

	"github.com/kaiohenricunha/go-music-k8s/backend/config"
)

var musicOrm = NewMusicOrm()

// Buscar todas as músicas
func GetAllSongs(cfg *config.Config) ([]*Song, error) {
	songs, err := musicOrm.Get(cfg.DB)
	if err != nil {
		log.Print("error getting songs from DB", err.Error())
		return nil, err
	}

	return songs, nil
}

// Adicionar uma nova música
func PostSong(cfg *config.Config, newsong *Song) error {
	err := musicOrm.Post(cfg.DB, newsong)
	if err != nil {
		log.Print("error posting song to DB", err.Error())
		return err
	}

	return nil
}

// UpdateSong updates an existing song by ID
func UpdateSong(cfg *config.Config, id uint, updatedSong *Song) error {
	// No need to manually update fields if the Updates method is used correctly
	err := musicOrm.Update(cfg.DB, id, updatedSong)
	if err != nil {
		log.Printf("Error updating song: %v", err)
		return err
	}

	return nil
}

// Deletar uma música
func DeleteSong(cfg *config.Config, songID uint) error {
	err := musicOrm.Delete(cfg.DB, songID)
	if err != nil {
		log.Print("error deleting song from DB", err.Error())
		return err
	}

	return nil
}

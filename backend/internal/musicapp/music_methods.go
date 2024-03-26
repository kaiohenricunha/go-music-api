package musicapp

import (
	"github.com/kaiohenricunha/go-music-k8s/backend/config"
	"log"
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

// Atualizar uma música existente
func UpdateSong(cfg *config.Config, song *Song) error {
	err := musicOrm.Update(cfg.DB, song)
	if err != nil {
		log.Print("error updating song in DB", err.Error())
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

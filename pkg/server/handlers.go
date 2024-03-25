package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kaiohenricunha/go-music-k8s/pkg/musicapp"
	"log"
	"net/http"
	"strconv"
)

func GetSongHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("GET /api/v1/music")

	// get songs
	songs, err := musicapp.GetAllSongs(gcfg)
	if err != nil {
		log.Print("GET songs api failed", err.Error())
		w.Write([]byte("Error getting songs!"))
		return
	}

	// response
	json.NewEncoder(w).Encode(songs)
}

func PostSongHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("POST /api/v1/music")

	// post songs
	newsong := musicapp.Song{}
	json.NewDecoder(r.Body).Decode(&newsong)

	err := musicapp.PostSong(gcfg, &newsong)
	if err != nil {
		log.Print("POST songs api failed", err.Error())
		w.Write([]byte("Error posting song !"))
		return
	}

	// response
	w.Write([]byte("New song published !"))
}

func UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("PUT /api/v1/music")

	var song musicapp.Song
	json.NewDecoder(r.Body).Decode(&song)

	err := musicapp.UpdateSong(gcfg, &song)
	if err != nil {
		log.Print("PUT song api failed", err.Error())
		w.Write([]byte("Error updating song!"))
		return
	}

	w.Write([]byte("Song updated successfully!"))
}

func DeleteSongHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("DELETE /api/v1/music/{id}")

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		log.Print("Error parsing song ID", err.Error())
		w.Write([]byte("Invalid song ID"))
		return
	}

	err = musicapp.DeleteSong(gcfg, uint(id))
	if err != nil {
		log.Print("DELETE song api failed", err.Error())
		w.Write([]byte("Error deleting song!"))
		return
	}

	w.Write([]byte("Song deleted successfully!"))
}

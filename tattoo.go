package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gogotattoo/common/models"
	"github.com/gorilla/mux"
)

// Tattoo shows info on a single tattoo work by id
func Tattoo(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range tattoos {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(models.NewTattoo("", "brr", "", ""))
}

// TattooToml shows info of a single tattoo work by id in toml format
func TattooToml(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if len(params["artist"]) > 0 {
		if len(req.URL.Query().Get("status")) > 0 {
			renderTattoosToml(w, params["artist"], req.URL.Query().Get("status"))
		} else {
			for _, tat := range artistWorks[params["artist"]+"/tattoo"] {
				toml.NewEncoder(w).Encode(tat)
			}

		}
		return
	}
	toml.NewEncoder(w).Encode(tattoos)
}

// Tattoos returns the list of all tattoos
func Tattoos(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if len(params["artist"]) > 0 {
		if len(req.URL.Query().Get("status")) > 0 {
			renderTattoos(w, params["artist"], req.URL.Query().Get("status"))
		} else {
			json.NewEncoder(w).Encode(artistWorks[params["artist"]+"/tattoo"])
		}
		return
	}
	json.NewEncoder(w).Encode(tattoos)
}

// CreateTattoo adds a new tattoo to the memory
func CreateTattoo(w http.ResponseWriter, req *http.Request) {
	log.Println("POST /tattoo")
	params := mux.Vars(req)
	defer req.Body.Close()
	var tat models.Tattoo
	err := json.NewDecoder(req.Body).Decode(&tat)
	log.Println("TITLE\n", tat.Title)
	if err != nil {
		log.Println("ERROR\n", err)
		json.NewEncoder(w).Encode(err)
		return
	}
	tat.ID = params["id"]

	if len(params["artist"]) > 0 {
		artistTattoos[params["artist"]+"/tattoo/"+params["work_name"]+"?status=wip"] = tat
	} else {
		tattoos = append(tattoos, tat)
	}
	m, _ := json.Marshal(tat)
	log.Println("TATTOO\n", string(m)+"\n")
	json.NewEncoder(w).Encode(tat)
}

var artistTattoos = make(map[string]models.Tattoo)

func renderTattoos(w http.ResponseWriter, artistName, status string) {
	tts := make([]models.Tattoo, 0, 100)
	for key, tat := range artistTattoos {
		if strings.Contains(key, artistName) {
			tts = append(tts, tat)
		}
	}
	json.NewEncoder(w).Encode(tts)
}
func renderTattoosToml(w http.ResponseWriter, artistName, status string) {
	tts := make([]models.Tattoo, 0, 100)
	for key, tat := range artistTattoos {
		if strings.Contains(key, artistName) {
			tts = append(tts, tat)
		}
	}

	for _, tat := range tts {
		toml.NewEncoder(w).Encode(tat)
	}
}

// DeleteTattoo deletes a tattoo by id from the memory
func DeleteTattoo(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if len(params["artist"]) > 0 {
		delete(artistTattoos, params["artist"]+"/tattoo/"+params["work_name"]+"?status=wip")
		renderTattoos(w, params["artist"], "wip")
		return
	}
	for index, item := range tattoos {
		if item.ID == params["id"] {
			tattoos = append(tattoos[:index], tattoos[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(tattoos)
}

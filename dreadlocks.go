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

// Lock shows info on a single dreadlocks work by id
func Lock(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range locks {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(models.NewDreadlocks("", "brr", "", ""))
}

// LocksToml shows info of a single work by id in toml format
func LocksToml(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if len(params["artist"]) > 0 {
		if len(req.URL.Query().Get("status")) > 0 {
			renderDreadlocksToml(w, params["artist"], req.URL.Query().Get("status"))
		} else {
			for _, dr := range artistWorks[params["artist"]+"/dreadlocks"] {
				toml.NewEncoder(w).Encode(dr)
			}
		}
	}
}

// Locks returns the list of all dreadlocks
func Locks(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if len(params["artist"]) > 0 {
		if len(req.URL.Query().Get("status")) > 0 {
			renderDreadlocks(w, params["artist"], req.URL.Query().Get("status"))
		} else {
			json.NewEncoder(w).Encode(artistWorks[params["artist"]+"/dreadlocks"])
		}
	}
}

// CreateDreadlocks adds a new work to the memory
func CreateDreadlocks(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	defer req.Body.Close()
	var tat models.Dreadlocks
	err := json.NewDecoder(req.Body).Decode(&tat)
	log.Println("TITLE\n", tat.Title)
	if err != nil {
		log.Println("ERROR\n", err)
		json.NewEncoder(w).Encode(err)
		return
	}
	tat.ID = params["id"]

	if len(params["artist"]) > 0 {
		artistDreadlocks[params["artist"]+"/dreadlocks/"+params["work_name"]+"?status=wip"] = tat
	} else {
		locks = append(locks, tat)
	}
	m, _ := json.Marshal(tat)
	log.Println("TATTOO\n", string(m)+"\n")
	json.NewEncoder(w).Encode(tat)
}

// DeleteDreadlocks deletes a work by id from the memory
func DeleteDreadlocks(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if len(params["artist"]) > 0 {
		delete(artistDreadlocks, params["artist"]+"/dreadlocks/"+params["work_name"]+"?status=wip")
		renderDreadlocks(w, params["artist"], "wip")
		return
	}
	for index, item := range locks {
		if item.ID == params["id"] {
			locks = append(locks[:index], locks[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(locks)
}

var artistDreadlocks = make(map[string]models.Dreadlocks)

func renderDreadlocks(w http.ResponseWriter, artistName, status string) {
	drrlkts := make([]models.Dreadlocks, 0, 100)
	for key, dr := range artistDreadlocks {
		if strings.Contains(key, artistName) {
			drrlkts = append(drrlkts, dr)
		}
	}
	json.NewEncoder(w).Encode(drrlkts)
}

func renderDreadlocksToml(w http.ResponseWriter, artistName, status string) {
	drrlkts := make([]models.Dreadlocks, 0, 100)
	for key, dr := range artistDreadlocks {
		if strings.Contains(key, artistName) {
			drrlkts = append(drrlkts, dr)
		}
	}
	for _, d := range drrlkts {
		toml.NewEncoder(w).Encode(d)
	}
}

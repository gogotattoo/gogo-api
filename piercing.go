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

// CreatePiercing adds a new design artwork
func CreatePiercing(w http.ResponseWriter, req *http.Request) {
	log.Println("POST /piercing")
	params := mux.Vars(req)
	defer req.Body.Close()
	var per models.Piercing
	err := json.NewDecoder(req.Body).Decode(&per)
	log.Println("TITLE\n", per.Title)
	if err != nil {
		log.Println("ERROR\n", err)
		json.NewEncoder(w).Encode(err)
		return
	}
	per.ID = params["id"]
	if len(params["artist"]) > 0 {
		artistPiercing[params["artist"]+"/piercing/"+params["work_name"]+"?status=wip"] = per
	} else {
		piercing = append(piercing, per)
	}
	m, _ := json.Marshal(per)
	log.Println("PIERCING\n", string(m)+"\n")
	json.NewEncoder(w).Encode(per)
}

// Piercing returns the list of all piercing works
func Piercing(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if len(params["artist"]) > 0 {
		if len(req.URL.Query().Get("status")) > 0 {
			renderPiercing(w, params["artist"], req.URL.Query().Get("status"))
		} else {
			json.NewEncoder(w).Encode(artistWorks[params["artist"]+"/piercing"])
		}
		return
	}
	json.NewEncoder(w).Encode(piercing)
}

// PiercingToml shows info of a single piercing work by id in toml format
func PiercingToml(w http.ResponseWriter, req *http.Request) {
	toml.NewEncoder(w).Encode(piercing[len(piercing)-1])
}

// DeletePiercing deletes a piercing by id from the posted works in memory
func DeletePiercing(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if len(params["artist"]) > 0 {
		delete(artistPiercing, params["artist"]+"/piercing/"+params["work_name"]+"?status=wip")
		renderPiercing(w, params["artist"], "wip")
		return
	}
	for index, item := range piercing {
		if item.ID == params["id"] {
			piercing = append(piercing[:index], piercing[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(piercing)
}

var artistPiercing = make(map[string]models.Piercing)

func renderPiercing(w http.ResponseWriter, artistName, status string) {
	prsngs := make([]models.Piercing, 0, 100)
	for key, pier := range artistPiercing {
		if strings.Contains(key, artistName) {
			prsngs = append(prsngs, pier)
		}
	}
	json.NewEncoder(w).Encode(prsngs)
}

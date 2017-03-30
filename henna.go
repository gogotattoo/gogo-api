package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/gogotattoo/common/models"
	"github.com/gorilla/mux"
)

// Hennas returns the list of all hennas
func Hennas(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(hennas)
}

// HennaToml shows info of a single henna work by id in toml format
func HennaToml(w http.ResponseWriter, req *http.Request) {
	toml.NewEncoder(w).Encode(hennas[len(hennas)-1])
}

// CreateHenna adds a new henna object
func CreateHenna(w http.ResponseWriter, req *http.Request) {
	log.Println("POST /henna")
	params := mux.Vars(req)
	defer req.Body.Close()
	var hen models.Henna
	err := json.NewDecoder(req.Body).Decode(&hen)
	log.Println("TITLE\n", hen.Title)
	if err != nil {
		log.Println("ERROR\n", err)
		json.NewEncoder(w).Encode(err)
		return
	}
	hen.ID = params["id"]
	hennas = append(hennas, hen)
	m, _ := json.Marshal(hen)
	log.Println("HENNA\n", string(m)+"\n")
	json.NewEncoder(w).Encode(hen)
}

// DeleteHenna deletes a henna by id from the posted works in memory
func DeleteHenna(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range hennas {
		if item.ID == params["id"] {
			hennas = append(hennas[:index], hennas[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(hennas)
}

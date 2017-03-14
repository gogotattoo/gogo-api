package main

import (
	"encoding/json"
	"log"
	"net/http"

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
	piercing = append(piercing, per)
	m, _ := json.Marshal(per)
	log.Println("PIERCING\n", string(m)+"\n")
	json.NewEncoder(w).Encode(per)
}

// Piercing returns the list of all piercing works
func Piercing(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(piercing)
}

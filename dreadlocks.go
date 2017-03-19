package main

import (
	"encoding/json"
	"log"
	"net/http"

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
	toml.NewEncoder(w).Encode(locks[len(locks)-1])
}

// Locks returns the list of all dreadlocks
func Locks(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(locks)
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
	locks = append(locks, tat)
	m, _ := json.Marshal(tat)
	log.Println("TATTOO\n", string(m)+"\n")
	json.NewEncoder(w).Encode(tat)
}

// DeleteDreadlocks deletes a work by id from the memory
func DeleteDreadlocks(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range locks {
		if item.ID == params["id"] {
			locks = append(locks[:index], locks[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(locks)
}

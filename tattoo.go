package main

import (
	"encoding/json"
	"log"
	"net/http"

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
	toml.NewEncoder(w).Encode(tattoos[len(tattoos)-1])
}

// Tattoos returns the list of all tattoos
func Tattoos(w http.ResponseWriter, req *http.Request) {
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
	tattoos = append(tattoos, tat)
	m, _ := json.Marshal(tat)
	log.Println("TATTOO\n", string(m)+"\n")
	json.NewEncoder(w).Encode(tat)
}

// DeleteTattoo deletes a tattoo by id from the memory
func DeleteTattoo(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range tattoos {
		if item.ID == params["id"] {
			tattoos = append(tattoos[:index], tattoos[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(tattoos)
}

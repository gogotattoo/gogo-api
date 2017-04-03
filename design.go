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

// CreateDesign adds a new design object
func CreateDesign(w http.ResponseWriter, req *http.Request) {
	log.Println("POST /design")
	params := mux.Vars(req)
	defer req.Body.Close()
	var des models.Design
	err := json.NewDecoder(req.Body).Decode(&des)
	log.Println("TITLE\n", des.Title)
	if err != nil {
		log.Println("ERROR\n", err)
		json.NewEncoder(w).Encode(err)
		return
	}
	des.ID = params["id"]

	if len(params["artist"]) > 0 {
		artistDesigns[params["artist"]+"/design/"+params["work_name"]+"?status=wip"] = des
	} else {
		designs = append(designs, des)
	}
	m, _ := json.Marshal(des)
	log.Println("DESIGN\n", string(m)+"\n")
	json.NewEncoder(w).Encode(des)
}

// Designs returns the list of all designs
func Designs(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if len(params["artist"]) > 0 {
		if len(req.URL.Query().Get("status")) > 0 {
			renderDesigns(w, params["artist"], req.URL.Query().Get("status"))
		} else {
			json.NewEncoder(w).Encode(artistWorks[params["artist"]+"/design"])
		}
		return
	}
	json.NewEncoder(w).Encode(designs)
}

// DesignsToml shows info of a single design work by id in toml format
func DesignsToml(w http.ResponseWriter, req *http.Request) {
	toml.NewEncoder(w).Encode(designs[len(designs)-1])
}

// DeleteDesign deletes a design by id from the posted works in memory
func DeleteDesign(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if len(params["artist"]) > 0 {
		delete(artistDesigns, params["artist"]+"/design/"+params["work_name"]+"?status=wip")
		renderDesigns(w, params["artist"], "wip")
		return
	}
	for index, item := range designs {
		if item.ID == params["id"] {
			designs = append(designs[:index], designs[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(designs)
}

var artistDesigns = make(map[string]models.Design)

func renderDesigns(w http.ResponseWriter, artistName, status string) {
	dsgns := make([]models.Design, 0, 100)
	for key, des := range artistDesigns {
		if strings.Contains(key, artistName) {
			dsgns = append(dsgns, des)
		}
	}
	json.NewEncoder(w).Encode(dsgns)
}

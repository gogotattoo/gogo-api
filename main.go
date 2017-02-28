package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type tattoo struct {
	ID           string   `json:"id"`
	Link         string   `json:"link,omitempty"`
	Title        string   `json:"title,omitempty"`
	MadeDate     string   `json:"tattoodate,omitempty"`
	PublishDate  string   `json:"date,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	BodyParts    []string `json:"bodypart,omitempty"`
	ImageIpfs    string   `json:"image_ipfs"`
	ImagsIpfs    []string `json:"images_ipfs,omitempty"`
	MadeLocation Address  `json:"made_at"`
	DurationMin  int      `json:"duration_min"`
	Gender       string   `json:"gender"`
	Extra        string   `json:"extra"`
	Article      string   `json:"article"`
}

// Address stores the location information where the work was made
type Address struct {
	City    string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
	Shop    string `json:"shop,omitempty"`
}

var tattoos []tattoo

func NewTattoo(id, title, link, imageIpfs string) (t tattoo) {
	t.ID = id
	t.Link = link
	t.ImageIpfs = imageIpfs
	return
}

func GetTattooEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range tattoos {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(NewTattoo("", "brr", "", ""))
}

func GetTattoosEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(tattoos)
}

func CreateTattooEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var tattoo tattoo
	_ = json.NewDecoder(req.Body).Decode(&tattoo)
	tattoo.ID = params["id"]
	tattoos = append(tattoos, tattoo)
	json.NewEncoder(w).Encode(tattoos)
}

func DeleteTattooEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range tattoos {
		if item.ID == params["id"] {
			tattoos = append(tattoos[:index], tattoos[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(tattoos)
}

func main() {
	router := mux.NewRouter()
	tattoos = append(tattoos, NewTattoo("0", "Young forever", "gogo/tattoo/young_forever", "QmUgcdgXS7RGC837EzDkHEMaWtPgPAMN9ntNeMbXsy98fi"))
	router.HandleFunc("/tattoo", GetTattoosEndpoint).Methods("GET")
	router.HandleFunc("/tattoo/{id}", GetTattooEndpoint).Methods("GET")
	router.HandleFunc("/tattoo/{id}", CreateTattooEndpoint).Methods("POST")
	router.HandleFunc("/tattoo/{id}", DeleteTattooEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))
}

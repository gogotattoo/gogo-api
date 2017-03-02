package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type tattoo struct {
	ID    string `json:"id,omitempty"`
	Link  string `json:"link,omitempty"`
	Title string `json:"title,omitempty" toml:"title"`
	//MadeLocation Address  `json:"made_at"`
	DurationMin int    `json:"duration_min,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Extra       string `json:"extra,omitempty"`
	Article     string `json:"article,omitempty"`

	MadeDate        string   `json:"tattoodate,omitempty" toml:"tattoodate"`
	PublishDate     string   `json:"date,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	BodyParts       []string `json:"bodypart,omitempty"`
	ImageIpfs       string   `json:"image_ipfs" toml:"image_ipfs"`
	ImagesIpfs      []string `json:"images_ipfs,omitempty" toml:"images_ipfs"`
	LocationCity    string   `json:"made_at_city" toml:"location_city"`
	LocationCountry string   `json:"made_at_country" toml:"location_country"`
	MadeAtShop      string   `json:"made_at_shop,omitempty" toml:"made_at_shop"`
}

// Address stores the location information where the work was made
// type Address struct {
// 	City    string `json:"city,omitempty"`
// 	Country string `json:"country,omitempty"`
// 	Shop    string `json:"shop,omitempty"`
// }

var tattoos []tattoo

// NewTattoo returns a new tattoo , requires id, the unique title of the new work
// link, also unique and final image ipfs hash
func NewTattoo(id, title, link, hash string) (t tattoo) {
	t.ID = id
	t.Link = link
	t.ImageIpfs = hash
	return
}

// Tattoo shows info on a single tattoo work by id
func Tattoo(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range tattoos {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(NewTattoo("", "brr", "", ""))
}

// Tattoos returns the list of all tattoos
func Tattoos(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(tattoos)
}

var myClient = &http.Client{Timeout: 100 * time.Second}

func getJSON(url string, target interface{}) (io.ReadCloser, error) {
	r, err := myClient.Get(url)
	if err != nil {
		return nil, err
	}
	//defer r.Body.Close()
	t := r.Body
	return t, json.NewDecoder(r.Body).Decode(target)
}

// TattooRefresh returns the list of all tattoos
func TattooRefresh(w http.ResponseWriter, req *http.Request) {

	// foo2 := F[]struct {
	// 	Name string
	// }
	// r, err := getJSON("https://api.github.com/repos/gogotattoo/gogo/contents/content/tattoo?ref=master",
	// 	foo2)
	// if err != nil {
	// 	json.NewEncoder(w).Encode(err)
	// }
	m, _ := json.Marshal(req.URL)
	w.Write(m)
}

// CreateTattoo adds a new tattoo to the memory
func CreateTattoo(w http.ResponseWriter, req *http.Request) {
	log.Println("POST /tattoo")
	params := mux.Vars(req)
	defer req.Body.Close()
	var tat tattoo
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
	log.Println("TATTOO\n", string(m))
	json.NewEncoder(w).Encode(tattoos)
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

func main() {
	router := mux.NewRouter()
	tattoos = append(tattoos, NewTattoo("0", "Young forever", "gogo/tattoo/young_forever", "QmUgcdgXS7RGC837EzDkHEMaWtPgPAMN9ntNeMbXsy98fi"))
	router.HandleFunc("/tattoo", Tattoos).Methods("GET")
	router.HandleFunc("/tattoo/refresh", TattooRefresh).Methods("GET")
	router.HandleFunc("/tattoo/{id}", Tattoo).Methods("GET")
	router.HandleFunc("/tattoo/{id}", CreateTattoo).Methods("POST")
	router.HandleFunc("/tattoo/{id}", DeleteTattoo).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))
}

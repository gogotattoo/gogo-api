package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gogotattoo/common/models"
	"github.com/gogotattoo/gogo-api/artwork"
	"github.com/gogotattoo/gogo-upload/cli"
	"github.com/gogotattoo/gogo-upload/watermark"
	"github.com/gorilla/mux"
)

var tattoos []models.Tattoo
var hennas []models.Henna
var piercing []models.Piercing
var designs []models.Design

var artistWorks = make(map[string]models.Artworks)

// ArtistArtworkRefresh returns the list of all tattoos
// TODO: add a timer, allow only every 5-10 mins
func ArtistArtworkRefresh(artType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		artistName := mux.Vars(r)["name"]
		artistWorks[artistName+"/"+artType] = artwork.Refresh(artistName, artType)
		json.NewEncoder(w).Encode(artistWorks[artistName+"/"+artType])
	}
}

// ArtistArtwork returns the list of all artists' tattoos actually published to git repos
func ArtistArtwork(artType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(artistWorks[mux.Vars(r)["name"]+"/"+artType])
	}
}

var refreshLastCalled = time.Now()

// ArtistArtworkRefreshAll refreshes all known art types for given artistName
// returns {"tattoo: {...}, design: {}, etc..."} json
func ArtistArtworkRefreshAll(w http.ResponseWriter, r *http.Request) {
	artistName := mux.Vars(r)["name"]
	result := make(map[string]models.Artworks)
	for _, artType := range []string{"tattoo", "henna", "piercing", "design"} {
		if time.Now().After(refreshLastCalled.Add(time.Minute * 3)) {
			artistWorks[artistName+"/"+artType] = artwork.Refresh(artistName, artType)
		}
		result[artType] = artistWorks[artistName+"/"+artType]
	}
	refreshLastCalled = time.Now()
	json.NewEncoder(w).Encode(result)
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

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		type vars struct {
			Token string
			Date  string
		}
		t, _ := template.ParseFiles("template/upload.gtpl")
		t.Execute(w, vars{Token: token, Date: time.Now().Format("2006/01/02")})
	} else {
		r.ParseMultipartForm(32 << 20)
		artistName := r.Form.Get("artist_name")
		madeAt := r.Form.Get("made_at")
		madeDate := r.Form.Get("made_date")
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		//fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		watermark.NeedLabels = true
		watermark.WatermarkPath = os.Getenv("GOPATH") + "/src/github.com/gogotattoo/gogo-upload/watermarks/gogo-watermark.png"
		watermark.OutputDir = "./upload/"
		watermark.LabelMadeBy = artistName
		watermark.LabelMadeAt = madeAt
		watermark.LabelDate = madeDate

		hashes := cli.AddWatermarks("./upload/" + handler.Filename)

		if len(hashes) > 0 {
			http.Redirect(w, r, "https://ipfs.io/ipfs/"+hashes[0], 301)
		} else {
			fmt.Fprintf(w, "Cannot get hash, try again.")
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tattoo", Tattoos).Methods("GET")
	//router.HandleFunc("/tattoo/{id}", Tattoo).Methods("GET")
	router.HandleFunc("/tattoo.toml", TattooToml).Methods("GET")
	router.HandleFunc("/tattoo/{id}", CreateTattoo).Methods("POST")
	router.HandleFunc("/tattoo/{id}", DeleteTattoo).Methods("DELETE")

	router.HandleFunc("/henna", Hennas).Methods("GET")
	router.HandleFunc("/henna/{id}", CreateHenna).Methods("POST")

	router.HandleFunc("/design", Designs).Methods("GET")
	router.HandleFunc("/design/{id}", CreateDesign).Methods("POST")

	router.HandleFunc("/piercing", Piercing).Methods("GET")
	router.HandleFunc("/piercing/{id}", CreatePiercing).Methods("POST")

	router.HandleFunc("/upload", upload)
	// router.Handle("/uploaded/", http.StripPrefix("/uploaded/", http.FileServer(http.Dir("./upload/"))))
	router.PathPrefix("/uploaded/").Handler(http.StripPrefix("/uploaded/", http.FileServer(http.Dir("upload/"))))

	for _, t := range []string{"tattoo", "henna", "piercing", "design"} {
		router.HandleFunc("/"+t+"/{name}", ArtistArtwork(t)).Methods("GET")
		router.HandleFunc("/"+t+"/{name}/refresh", ArtistArtworkRefresh(t)).Methods("GET")
	}

	router.HandleFunc("/all/{name}/refresh", ArtistArtworkRefreshAll).Methods("GET")

	go func() {
		for _, artistName := range []string{"gogo", "aid", "xizi"} {
			for _, artType := range []string{"tattoo", "henna", "piercing", "design"} {
				artistWorks[artistName+"/"+artType] = artwork.Refresh(artistName, artType)
			}
		}
	}()

	log.Fatal(http.ListenAndServe(":12345", Log(router)))
}

// Log prints basic http request info
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

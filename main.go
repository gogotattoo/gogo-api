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
	"strings"
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
var locks []models.Dreadlocks

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
	for _, artType := range []string{"tattoo", "henna", "piercing", "design", "dreadlocks", "art", "food"} {
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
		token := r.Form.Get("token")
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		//fmt.Fprintf(w, "%v", handler.Header)
		dirName := "./upload/" + strings.Replace(time.Now().Format("2006/01/02"), "/", "_", -1) + "/"
		os.MkdirAll(dirName, os.ModePerm)
		f, err := os.OpenFile(dirName+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		// Initializing watermark maker. There definetely must be a better golang way for it
		watermark.NeedLabels = true
		watermark.OutputDir = dirName
		watermark.LabelMadeBy = artistName
		watermark.LabelMadeAt = madeAt
		watermark.LabelDate = madeDate

		// It's a bit dirty, but will work for now
		// by default, it's gogo's work
		watermark.WatermarkPath = os.Getenv("GOPATH") + "/src/github.com/gogotattoo/gogo-upload/watermarks/gogo-watermark.png"
		if artistName == "aid" {
			watermark.WatermarkPath = os.Getenv("GOPATH") + "/src/github.com/gogotattoo/gogo-upload/watermarks/aidlong.png"
		} else if artistName == "xizi" {
			watermark.WatermarkPath = os.Getenv("GOPATH") + "/src/github.com/gogotattoo/gogo-upload/watermarks/xizilong.png"
		} else if artistName == "klimin" {
			watermark.WatermarkPath = os.Getenv("GOPATH") + "/src/github.com/gogotattoo/gogo-upload/watermarks/klimin-watermark.png"
		} else if artistName == "jiaye" {
			watermark.WatermarkPath = os.Getenv("GOPATH") + "/src/github.com/gogotattoo/gogo-upload/watermarks/jiaye-watermark.png"
		}

		watermark.V3 = true
		watermark.WatermarkPath = os.Getenv("GOPATH") + "/src/github.com/gogotattoo/gogo-upload/watermarks/v3/" + artistName + ".png"

		hashes := cli.AddWatermarks(dirName + handler.Filename)

		if len(hashes) > 0 {
			if len(token) == 0 {
				type res struct {
					URL  string
					Hash string
				}
				json.NewEncoder(w).Encode(&res{URL: "https://ipfs.io/ipfs/" + hashes[0], Hash: hashes[0]})
			} else {
				http.Redirect(w, r, "https://ipfs.io/ipfs/"+hashes[0], 301)
			}
		} else {
			fmt.Fprintf(w, "Cannot get hash, try again.")
		}
	}
}

func main() {
	router := NewRouter()

	artists := make(models.Artists, 6)
	artists[0] = models.Artist{Name: "gogo", Services: []string{"tattoo", "henna", "piercing", "design", "dreadlocks"}}
	artists[1] = models.Artist{Name: "aid", Services: []string{"tattoo"}}
	artists[2] = models.Artist{Name: "xizi", Services: []string{"tattoo", "design", "henna", "piercing"}}
	artists[3] = models.Artist{Name: "kate", Services: []string{"tattoo", "design"}}
	artists[4] = models.Artist{Name: "klimin", Services: []string{"tattoo", "design"}}
	artists[5] = models.Artist{Name: "jiaye", Services: []string{"tattoo", "design"}}
	for _, artist := range artists {
		for _, service := range artist.Services {
			go func(name, service string) {
				artistWorks[name+"/"+service] = artwork.Refresh(name, service)
			}(artist.Name, service)
		}
	}

	log.Fatal(http.ListenAndServe(":12345", Log(router)))
}

// Log prints basic http request info
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

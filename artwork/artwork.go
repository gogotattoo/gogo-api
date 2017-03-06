package artwork

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gogotattoo/common/models"
	"github.com/gogotattoo/common/util"
)

func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	t, er := ioutil.ReadAll(r.Body)
	if er != nil {
		log.Println("Error", er)
	}
	log.Println(string(t))
	return json.Unmarshal(t, target)
}

type file struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}

const (
	gitURL = "https://api.github.com/repos/gogotattoo/%s/contents/content/%s?ref=master"
)

var myClient = &http.Client{Timeout: 100 * time.Second}

// Refresh since we are currntly not using a database, instead all the artwork data is
// on git services. Let's update our in memory database by utilizing github api.
func Refresh(artistName, artType string) models.Artworks {
	var files []file
	url := fmt.Sprintf(gitURL, artistName, artType)
	log.Println("Updating from " + url)
	of, err := myClient.Get(url)
	if err != nil {
		log.Panic(err)
	}
	defer of.Body.Close()
	// of, err := os.Open("gogo-2017-03-05-tattoo.json")
	// if err != nil {
	// 	log.Panic(err)
	// }
	// defer of.Close()
	err = json.NewDecoder(of.Body).Decode(&files)
	if err != nil {
		log.Panic(err)
	}
	var works models.Artworks
	for i, f := range files {
		//color.Green(":\t" + f.Name + "\t\t" + f.DownloadURL)
		of, err := myClient.Get(f.DownloadURL)
		if err != nil {
			log.Panic(err)
		}
		defer of.Body.Close()
		tomlStr, _ := util.ExtractTomlStr(of.Body)
		var work models.Artwork
		toml.Unmarshal([]byte(tomlStr), &work)
		work.ID = strconv.Itoa(i)
		works = append(works, work)
	}
	sort.Sort(works)
	return works
}

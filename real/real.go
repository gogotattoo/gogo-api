package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
)

var myClient = &http.Client{Timeout: 100 * time.Second}

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

func main() {
	date, _ := time.Parse(time.RFC822, "14 Mar 10 18:00 UTC")
	var config = map[string]interface{}{
		"date":   date,
		"counts": []int{1, 1, 2, 3, 5, 8},
		"hash": map[string]string{
			"key1": "val1",
			"key2": "val2",
		},
	}
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf.String())
	type file struct {
		Name string `json:"name"`
	}
	files := make(map[file]string, 0)
	var objmap map[string]*json.RawMessage
	err := getJSON("https://api.github.com/repos/gogotattoo/gogo/contents/content/tattoo?ref=master",
		&objmap)
	if err != nil {
		log.Panic(err)
	}
	//fmt.Println(r)
	fmt.Println(objmap)
	fmt.Println(files)
}

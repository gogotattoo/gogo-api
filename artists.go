package main

import (
	"encoding/json"
	"net/http"

	"github.com/gogotattoo/common/models"
)

func Artists(w http.ResponseWriter, r *http.Request) {
	artists := make(models.Artists, 4)
	artists[0] = models.Artist{Link: "gogo", Name: "Яна Gogo",
		Services: []string{"tattoo", "henna", "piercing", "design", "dreadlocks"}}
	artists[1] = models.Artist{Link: "aid", Name: "Valentin Aidov",
		Services: []string{"tattoo", "design"}}
	artists[2] = models.Artist{Link: "xizi", Name: "Xizi",
		Services: []string{"tattoo", "henna", "piercing", "design"}}
	artists[3] = models.Artist{Link: "kate", Name: "Екатерина",
		Services: []string{"tattoo", "henna", "design"}}
	json.NewEncoder(w).Encode(artists)
}

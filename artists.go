package main

import (
	"encoding/json"
	"net/http"

	"github.com/gogotattoo/common/models"
)

func Artists(w http.ResponseWriter, r *http.Request) {
	artists := make(models.Artists, 5)
	artists[0] = models.Artist{Link: "gogo", Name: "Яна Gogo",
		Services:   []string{"tattoo", "henna", "piercing", "design", "dreadlocks"},
		AvatarIpfs: "QmasSfXhWZgB1BT7Ytn7SCxxxnpZcMXofr2m93LpGTDGHh"}
	artists[1] = models.Artist{Link: "aid", Name: "Valentin Aidov",
		Services:   []string{"tattoo", "design"},
		AvatarIpfs: "Qmc32Lt6bh4ybvHcKMYr5bBU8aSSb4LUFrE7MKWuudNhs7"}
	artists[2] = models.Artist{Link: "xizi", Name: "Xizi",
		Services:   []string{"tattoo", "henna", "piercing", "design"},
		AvatarIpfs: "QmToh8L6B4ytHdFLErRhWJcUf8aqW6GQ3TMnwdHxt98GUR"}
	artists[3] = models.Artist{Link: "kate", Name: "Екатерина",
		Services:   []string{"tattoo", "henna", "design"},
		AvatarIpfs: "QmXRBH18LTZx3G29BtVVee266hPn1swherCYJ4Gz8oR5iN"}
	artists[4] = models.Artist{Link: "klimin", Name: "Vitaly Klimin",
		Services:   []string{"tattoo", "design"},
		AvatarIpfs: "QmThwsgPBLiPR33NU4toa2zE6ABAv2w2M6DK6fWTbQyN6b"}
	artists[5] = models.Artist{Link: "jiaye", Name: "Jia Ye",
		Services:   []string{"tattoo", "design"},
		AvatarIpfs: ""}
	json.NewEncoder(w).Encode(artists)
}

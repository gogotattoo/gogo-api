package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

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
}

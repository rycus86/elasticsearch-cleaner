package main

import (
	"fmt"
	"github.com/rycus86/elasticsearch-cleaner/elasticsearch"
	"github.com/rycus86/elasticsearch-cleaner/settings"
	"net/http"
	"time"
)

func main() {
	settings.Initialize()

	client := http.Client{Timeout: time.Duration(15) * time.Second}
	es := elasticsearch.NewClient(&client)

	keys, err := es.FetchIndices()
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch indices: %s", err))
	}

	if len(keys) <= settings.GetMaxIndices() {
		fmt.Printf("There are only %d items ( <= %d )\n", len(keys), settings.GetMaxIndices())
		return
	} else {
		fmt.Printf("There are %d items\n", len(keys))
	}

	for idx, key := range keys[0 : len(keys)-settings.GetMaxIndices()-1] {
		fmt.Printf("#%3d: Removing %s ... ", idx+1, key)

		err = es.DeleteIndex(key)

		if err != nil {
			fmt.Printf("Failed: %s\n", err)
		} else {
			fmt.Println("OK")
		}
	}
}

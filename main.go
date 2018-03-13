package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"time"
)

var baseUrl = "http://192.168.0.52:9200"
var pattern = "fluentd-.*"
var maxLength = 20

type EsStats struct {
	Indices map[string]interface{}
}

func main() {
	client := http.Client{Timeout: time.Duration(15) * time.Second}
	response, err := client.Get(fmt.Sprintf("%s/_stats", baseUrl))

	if err != nil {
		panic(fmt.Sprintf("Failed to fetch ES stats: %s", err))
	}

	var stats = EsStats{}
	err = json.NewDecoder(response.Body).Decode(&stats)

	if err != nil {
		panic(fmt.Sprintf("Failed to decode ES respose: %s", err))
	}

	keys := make([]string, 0, len(stats.Indices))
	for key := range stats.Indices {
		if matched, _ := regexp.MatchString(pattern, key); matched {
			keys = append(keys, key)
		}
	}

	if len(keys) <= maxLength {
		fmt.Printf("There are only %d items ( <= %d )\n", len(keys), maxLength)
		return
	} else {
		fmt.Printf("There are %d items\n", len(keys))
	}

	sort.Strings(keys)

	for idx, key := range keys[0 : len(keys)-maxLength-1] {
		fmt.Printf("#%3d: Removing %s ... ", idx+1, key)

		indexUrl, _ := url.Parse(fmt.Sprintf("%s/%s", baseUrl, key))
		response, err = client.Do(&http.Request{
			Method: "DELETE",
			URL:    indexUrl,
		})

		if err != nil {
			fmt.Printf("Failed: %s\n", err)
		} else {
			if response.StatusCode == http.StatusOK {
				fmt.Println("OK")
			} else {
				fmt.Printf("Failed: HTTP %s\n", response.Status)
			}
		}
	}
}

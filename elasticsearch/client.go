package elasticsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rycus86/elasticsearch-cleaner/settings"
	"net/http"
	"net/url"
	"regexp"
	"sort"
)

type esClient struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) Client {
	return &esClient{httpClient: httpClient}
}

func (es *esClient) FetchIndices() ([]string, error) {
	statsUrl, err := url.Parse(fmt.Sprintf("%s/_stats", settings.GetBaseUrl()))
	if err != nil {
		panic(fmt.Sprintf("Invalid URL: %s/_stats", settings.GetBaseUrl()))
	}

	response, err := es.httpClient.Do(&http.Request{
		Method: "GET",
		URL:    statsUrl,
	})

	if err != nil {
		fmt.Printf("Failed to fetch ES stats: %s\n", err)
		return []string{}, err
	}

	var stats = statistics{}
	err = json.NewDecoder(response.Body).Decode(&stats)

	if err != nil {
		fmt.Printf("Failed to decode ES respose: %s\n", err)
		return []string{}, err
	}

	keys := make([]string, 0, len(stats.Indices))
	for key := range stats.Indices {
		if matched, _ := regexp.MatchString(settings.GetPattern(), key); matched {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	return keys, nil
}

func (es *esClient) DeleteIndex(key string) error {
	indexUrl, _ := url.Parse(fmt.Sprintf("%s/%s", settings.GetBaseUrl(), key))
	response, err := es.httpClient.Do(&http.Request{
		Method: "DELETE",
		URL:    indexUrl,
	})

	if err != nil {
		return err
	} else {
		if response.StatusCode != http.StatusOK {
			return nil
		} else {
			return errors.New(fmt.Sprintf("HTTP %s", response.Status))
		}
	}
}

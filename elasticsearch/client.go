package elasticsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
)

type esClient struct {
	httpClient *http.Client
	baseUrl    string
}

func NewClient(httpClient *http.Client, baseUrl string) Client {
	return &esClient{httpClient: httpClient, baseUrl: baseUrl}
}

func (es *esClient) FetchIndices(pattern string) ([]string, error) {
	statsUrl, err := url.Parse(fmt.Sprintf("%s/_stats", es.baseUrl))
	if err != nil {
		panic(fmt.Sprintf("Invalid URL: %s/_stats", es.baseUrl))
	}

	response, err := es.httpClient.Do(&http.Request{
		Method: "GET",
		URL:    statsUrl,
	})

	if err != nil {
		fmt.Printf("Failed to fetch ES stats: %s\n", err)
		return []string{}, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Failed to fetch ES stats: %s\n", response.Status)
		return []string{}, errors.New(fmt.Sprintf("HTTP %s", response.Status))
	}

	stats := statistics{}
	err = json.NewDecoder(response.Body).Decode(&stats)

	if err != nil {
		fmt.Printf("Failed to decode ES respose: %s\n", err)
		return []string{}, err
	}

	keys := make([]string, 0, len(stats.Indices))
	for key := range stats.Indices {
		if matched, _ := regexp.MatchString(pattern, key); matched {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	return keys, nil
}

func (es *esClient) DeleteIndex(key string) error {
	indexUrl, _ := url.Parse(fmt.Sprintf("%s/%s", es.baseUrl, key))
	response, err := es.httpClient.Do(&http.Request{
		Method: "DELETE",
		URL:    indexUrl,
	})

	if err != nil {
		return err
	} else {
		if response.StatusCode == http.StatusOK {
			return nil
		} else {
			return errors.New(fmt.Sprintf("HTTP %s", response.Status))
		}
	}
}

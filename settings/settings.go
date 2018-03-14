package settings

import (
	"fmt"
	"github.com/rycus86/elasticsearch-cleaner/env"
	"net/url"
	"regexp"
	"strconv"
)

type applicationSettings struct {
	BaseUrl    string
	Pattern    string
	MaxIndices int
}

var settings = &applicationSettings{
	Pattern:    ".*",
	MaxIndices: 20,
}

func Initialize() {
	//baseUrl := env.Get("BASE_URL")
	baseUrl := env.GetOrDefault("BASE_URL", "http://192.168.0.52:9200")
	if _, err := url.Parse(baseUrl); baseUrl == "" || err != nil {
		panic(fmt.Sprintf("Missing or invalid BASE_URL: %s", baseUrl))
	}

	settings.BaseUrl = baseUrl

	//pattern := env.GetOrDefault("PATTERN", settings.Pattern)
	pattern := env.GetOrDefault("PATTERN", "fluentd-.*")
	if _, err := regexp.Compile(pattern); pattern == "" || err != nil {
		panic(fmt.Sprintf("Invalid regexp pattern: %s", pattern))
	}

	settings.Pattern = pattern

	maxIndices, err := strconv.Atoi(env.GetOrDefault("MAX_INDICES", strconv.Itoa(settings.MaxIndices)))
	if maxIndices <= 0 || err != nil {
		panic(fmt.Sprintf("Invalid MAX_INDICES value: %s - %s", env.Get("MAX_INDICES"), err))
	}

	settings.MaxIndices = maxIndices
}

func GetBaseUrl() string {
	return settings.BaseUrl
}

func GetPattern() string {
	return settings.Pattern
}

func GetMaxIndices() int {
	return settings.MaxIndices
}

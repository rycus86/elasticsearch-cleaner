package settings

import (
	"fmt"
	"github.com/rycus86/elasticsearch-cleaner/env"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type applicationSettings struct {
	BaseUrl    string
	Pattern    string
	MaxIndices int
	Interval   time.Duration
	Timeout    time.Duration
}

var settings = &applicationSettings{
	MaxIndices: 20,
	Interval:   12 * time.Hour,
	Timeout:    30 * time.Second,
}

func Initialize() {
	baseUrl := env.Get("BASE_URL")
	if _, err := url.Parse(baseUrl); baseUrl == "" || err != nil {
		panic(fmt.Sprintf("Missing or invalid BASE_URL: %s", baseUrl))
	}

	settings.BaseUrl = baseUrl

	pattern := env.Get("PATTERN")
	if _, err := regexp.Compile(pattern); pattern == "" || err != nil {
		panic(fmt.Sprintf("Invalid regexp pattern: '%s'", pattern))
	}

	settings.Pattern = pattern

	maxIndices, err := strconv.Atoi(env.GetOrDefault("MAX_INDICES", strconv.Itoa(settings.MaxIndices)))
	if maxIndices <= 0 || err != nil {
		panic(fmt.Sprintf("Invalid MAX_INDICES value: %s - %s", env.Get("MAX_INDICES"), err))
	}

	settings.MaxIndices = maxIndices

	interval, err := time.ParseDuration(env.GetOrDefault("INTERVAL", settings.Interval.String()))
	if interval <= 0 || err != nil {
		panic(fmt.Sprintf("Invalid INTERVAL value: %s - %s", env.Get("INTERVAL"), err))
	}

	settings.Interval = interval

	timeout, err := time.ParseDuration(env.GetOrDefault("TIMEOUT", settings.Timeout.String()))
	if timeout <= 0 || err != nil {
		panic(fmt.Sprintf("Invalid TIMEOUT value: %s - %s", env.Get("TIMEOUT"), err))
	}

	settings.Timeout = timeout
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

func GetInterval() time.Duration {
	return settings.Interval
}

func GetTimeout() time.Duration {
	return settings.Timeout
}

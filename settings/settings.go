package settings

import (
	"fmt"
	"github.com/rycus86/elasticsearch-cleaner/env"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type ApplicationSettings struct {
	baseUrl    string
	pattern    string
	maxIndices int
	interval   time.Duration
	timeout    time.Duration
}

func Initialize() *ApplicationSettings {
	settings := &ApplicationSettings{
		maxIndices: 20,
		interval:   12 * time.Hour,
		timeout:    30 * time.Second,
	}

	baseUrl := env.Get("BASE_URL")
	if _, err := url.Parse(baseUrl); baseUrl == "" || err != nil {
		panic(fmt.Sprintf("Missing or invalid BASE_URL: %s", baseUrl))
	}

	settings.baseUrl = baseUrl

	pattern := env.Get("PATTERN")
	if _, err := regexp.Compile(pattern); pattern == "" || err != nil {
		panic(fmt.Sprintf("Invalid regexp pattern: '%s'", pattern))
	}

	settings.pattern = pattern

	maxIndices, err := strconv.Atoi(env.GetOrDefault("MAX_INDICES", strconv.Itoa(settings.maxIndices)))
	if maxIndices <= 0 || err != nil {
		panic(fmt.Sprintf("Invalid MAX_INDICES value: %s - %s", env.Get("MAX_INDICES"), err))
	}

	settings.maxIndices = maxIndices

	interval, err := time.ParseDuration(env.GetOrDefault("INTERVAL", settings.interval.String()))
	if interval <= 0 || err != nil {
		panic(fmt.Sprintf("Invalid INTERVAL value: %s - %s", env.Get("INTERVAL"), err))
	}

	settings.interval = interval

	timeout, err := time.ParseDuration(env.GetOrDefault("TIMEOUT", settings.timeout.String()))
	if timeout <= 0 || err != nil {
		panic(fmt.Sprintf("Invalid TIMEOUT value: %s - %s", env.Get("TIMEOUT"), err))
	}

	settings.timeout = timeout

	return settings
}

func (settings *ApplicationSettings) GetBaseUrl() string {
	return settings.baseUrl
}

func (settings *ApplicationSettings) GetPattern() string {
	return settings.pattern
}

func (settings *ApplicationSettings) GetMaxIndices() int {
	return settings.maxIndices
}

func (settings *ApplicationSettings) GetInterval() time.Duration {
	return settings.interval
}

func (settings *ApplicationSettings) GetTimeout() time.Duration {
	return settings.timeout
}

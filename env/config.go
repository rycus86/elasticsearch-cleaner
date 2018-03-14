package env

import "os"

func Get(key string) string {
	return os.Getenv(key)
}

func GetOrDefault(key string, defaultValue string) string {
	value, found := os.LookupEnv(key)
	if found {
		return value
	} else {
		return defaultValue
	}
}

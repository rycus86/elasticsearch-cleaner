package env

import (
	"os"
	"testing"
)

func TestLookup(t *testing.T) {
	os.Setenv("FOR_TEST", "test_value")
	defer os.Unsetenv("FOR_TEST")

	if "test_value" != GetOrDefault("FOR_TEST", "not-found") {
		t.Error("Environment variable not read properly")
	}
}

func TestDefault(t *testing.T) {
	if "default-value" != GetOrDefault("NOT_CONFIGURED", "default-value") {
		t.Error("Default value expected")
	}
}

func TestGetFound(t *testing.T) {
	os.Setenv("CONFIGURED", "found")
	defer os.Unsetenv("CONFIGURED")

	if "found" != Get("CONFIGURED") {
		t.Error("Unexpected value")
	}
}

func TestGetNotFound(t *testing.T) {
	if "" != Get("NOT_CONFIGURED") {
		t.Error("Unexpected value")
	}
}

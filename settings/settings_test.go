package settings

import (
	"os"
	"testing"
)

func TestDefaults(t *testing.T) {
	os.Setenv("BASE_URL", "base-url")
	defer os.Unsetenv("BASE_URL")

	Initialize()

	if baseUrl := GetBaseUrl(); baseUrl != "base-url" {
		t.Errorf("Base URL not set properly: %s", baseUrl)
	}

	if pattern := GetPattern(); pattern != ".*" {
		t.Errorf("Unexpected pattern: %s", pattern)
	}

	if maxIndices := GetMaxIndices(); maxIndices != 20 {
		t.Errorf("Unexpected max indices: %d", maxIndices)
	}
}

func TestNonDefaults(t *testing.T) {
	os.Setenv("BASE_URL", "sample-url")
	defer os.Unsetenv("BASE_URL")
	os.Setenv("PATTERN", "changed")
	defer os.Unsetenv("PATTERN")
	os.Setenv("MAX_INDICES", "42")
	defer os.Unsetenv("MAX_INDICES")

	Initialize()

	if baseUrl := GetBaseUrl(); baseUrl != "sample-url" {
		t.Errorf("Base URL not set properly: %s", baseUrl)
	}

	if pattern := GetPattern(); pattern != "changed" {
		t.Errorf("Unexpected pattern: %s", pattern)
	}

	if maxIndices := GetMaxIndices(); maxIndices != 42 {
		t.Errorf("Unexpected max indices: %d", maxIndices)
	}
}

func TestNoBaseUrl(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected to panic")
		}
	}()

	Initialize()
}

func TestInvalidBaseUrl(t *testing.T) {
	os.Setenv("BASE_URL", "@:@")
	defer os.Unsetenv("BASE_URL")

	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected to panic")
		}
	}()

	Initialize()
}

func TestInvalidPattern(t *testing.T) {
	os.Setenv("BASE_URL", "valid-url")
	defer os.Unsetenv("BASE_URL")
	os.Setenv("PATTERN", "?][")
	defer os.Unsetenv("PATTERN")

	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected to panic")
		}
	}()

	Initialize()
}

func TestInvalidEmptyPattern(t *testing.T) {
	os.Setenv("BASE_URL", "valid-url")
	defer os.Unsetenv("BASE_URL")
	os.Setenv("PATTERN", "")
	defer os.Unsetenv("PATTERN")

	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected to panic")
		}
	}()

	Initialize()
}

func TestInvalidMaxIndices(t *testing.T) {
	os.Setenv("BASE_URL", "valid-url")
	defer os.Unsetenv("BASE_URL")
	os.Setenv("MAX_INDICES", "invalid-int")
	defer os.Unsetenv("MAX_INDICES")

	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected to panic")
		}
	}()

	Initialize()
}

func TestMaxIndicesOutOfRange(t *testing.T) {
	os.Setenv("BASE_URL", "valid-url")
	defer os.Unsetenv("BASE_URL")
	os.Setenv("MAX_INDICES", "-12")
	defer os.Unsetenv("MAX_INDICES")

	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected to panic")
		}
	}()

	Initialize()
}

func TestMaxIndicesNonInteger(t *testing.T) {
	os.Setenv("BASE_URL", "valid-url")
	defer os.Unsetenv("BASE_URL")
	os.Setenv("MAX_INDICES", "3.14")
	defer os.Unsetenv("MAX_INDICES")

	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected to panic")
		}
	}()

	Initialize()
}

package settings

import (
	"os"
	"testing"
)

func TestDefaults(t *testing.T) {
	os.Setenv("BASE_URL", "base-url")
	defer os.Unsetenv("BASE_URL")
	os.Setenv("PATTERN", "pattern")
	defer os.Unsetenv("PATTERN")

	config := Initialize()

	if baseUrl := config.GetBaseUrl(); baseUrl != "base-url" {
		t.Errorf("Base URL not set properly: %s", baseUrl)
	}

	if pattern := config.GetPattern(); pattern != "pattern" {
		t.Errorf("Unexpected pattern: %s", pattern)
	}

	if maxIndices := config.GetMaxIndices(); maxIndices != 20 {
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

	config := Initialize()

	if baseUrl := config.GetBaseUrl(); baseUrl != "sample-url" {
		t.Errorf("Base URL not set properly: %s", baseUrl)
	}

	if pattern := config.GetPattern(); pattern != "changed" {
		t.Errorf("Unexpected pattern: %s", pattern)
	}

	if maxIndices := config.GetMaxIndices(); maxIndices != 42 {
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

func TestNoPattern(t *testing.T) {
	os.Setenv("BASE_URL", "base-set")
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
	os.Setenv("PATTERN", "pattern")
	defer os.Unsetenv("PATTERN")
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
	os.Setenv("PATTERN", "pattern")
	defer os.Unsetenv("PATTERN")
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
	os.Setenv("PATTERN", "pattern")
	defer os.Unsetenv("PATTERN")
	os.Setenv("MAX_INDICES", "3.14")
	defer os.Unsetenv("MAX_INDICES")

	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected to panic")
		}
	}()

	Initialize()
}

func TestInterval(t *testing.T) {
	os.Setenv("BASE_URL", "sample-url")
	defer os.Unsetenv("BASE_URL")
	os.Setenv("PATTERN", "changed")
	defer os.Unsetenv("PATTERN")
	os.Setenv("INTERVAL", "3m12s")
	defer os.Unsetenv("INTERVAL")

	config := Initialize()

	if interval := config.GetInterval(); interval.Seconds() != 192 {
		t.Errorf("Unexpected interval: %s", interval)
	}
}

func TestInvalidInterval(t *testing.T) {
	os.Setenv("BASE_URL", "valid-url")
	defer os.Unsetenv("BASE_URL")
	os.Setenv("PATTERN", "pattern")
	defer os.Unsetenv("PATTERN")
	os.Setenv("INTERVAL", "x-1")
	defer os.Unsetenv("INTERVAL")

	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected to panic")
		}
	}()

	Initialize()
}

func TestTimeout(t *testing.T) {
	os.Setenv("BASE_URL", "sample-url")
	defer os.Unsetenv("BASE_URL")
	os.Setenv("PATTERN", "changed")
	defer os.Unsetenv("PATTERN")
	os.Setenv("TIMEOUT", "42s")
	defer os.Unsetenv("TIMEOUT")

	config := Initialize()

	if timeout := config.GetTimeout(); timeout.Seconds() != 42 {
		t.Errorf("Unexpected timeout: %s", timeout)
	}
}

func TestInvalidTimeout(t *testing.T) {
	os.Setenv("BASE_URL", "valid-url")
	defer os.Unsetenv("BASE_URL")
	os.Setenv("PATTERN", "pattern")
	defer os.Unsetenv("PATTERN")
	os.Setenv("TIMEOUT", "x-1")
	defer os.Unsetenv("TIMEOUT")

	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected to panic")
		}
	}()

	Initialize()
}

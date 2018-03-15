package elasticsearch

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"fmt"
	"github.com/rycus86/elasticsearch-cleaner/settings"
	"os"
)

var server *httptest.Server

func TestFetchIndices(t *testing.T) {
	setup := mocks()
	setup.Pattern = "sample-.*"
	setup.Response = `
	{
		"indices": {
			"sample-003": 1,
			"sample-001": 2,
			"sample-002": 3,
			"not-matching-001": 4,
			".nomatch": 5
		}
	}
	`

	initialize(setup)

	client := NewClient(&http.Client{})
	response, err := client.FetchIndices()

	if err != nil {
		t.Errorf("Invalid response: %s", err)
	}

	if len(response) != 3 {
		t.Errorf("Wrong number of results: %d", len(response))
	}

	for idx, key := range response {
		if key != fmt.Sprintf("sample-%03d", idx+1) {
			t.Errorf("Invalid result item: %s", key)
		}
	}
}

func TestInvalidIndicesResponse(t *testing.T) {
	setup := mocks()
	setup.Response = "not-a-json"

	initialize(setup)

	client := NewClient(&http.Client{})
	_, err := client.FetchIndices()

	if err == nil {
		t.Error("Expected to fail")
	}
}

func TestFailingFetchIndicesRequest(t *testing.T) {
	setup := mocks()
	setup.Status = 500
	setup.Response = "{}"

	initialize(setup)

	client := NewClient(&http.Client{})
	_, err := client.FetchIndices()

	if err == nil {
		t.Error("Expected to fail")
	}
}

func TestDeleteIndex(t *testing.T) {
	initialize(mocks())

	client := NewClient(&http.Client{})
	err := client.DeleteIndex("x")

	if err != nil {
		t.Errorf("Unexpected failure: %s", err)
	}
}

func TestFailingDeleteIndexRequest(t *testing.T) {
	setup := mocks()
	setup.Status = 500
	setup.Response = "{}"

	initialize(setup)

	client := NewClient(&http.Client{})
	err := client.DeleteIndex("failing")

	if err == nil {
		t.Error("Expected to fail")
	}
}

type mockSettings struct {
	Pattern		string
	MaxIndices	int
	Status		int
	Response	string
	ExpectedPath	string
}

func mocks() *mockSettings {
	return &mockSettings{
		Status:		http.StatusOK,
		Pattern:	"x",
		MaxIndices:	10,
	}
}

func initialize(m *mockSettings) *httptest.Server {
	server = httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(m.Status)
		w.Write([]byte(m.Response))
	}))

	os.Setenv("BASE_URL", server.URL)
	os.Setenv("PATTERN", fmt.Sprintf("%s.*", m.Pattern))
	os.Setenv("MAX_INDICES", fmt.Sprintf("%d", m.MaxIndices))

	settings.Initialize()

	return server
}

func TestMain(m *testing.M) {
	defer os.Unsetenv("BASE_URL")
	defer os.Unsetenv("PATTERN")
	defer os.Unsetenv("MAX_INDICES")

	defer server.Close()

	os.Exit(m.Run())
}

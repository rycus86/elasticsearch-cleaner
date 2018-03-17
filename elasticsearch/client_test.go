package elasticsearch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var server *httptest.Server

func TestFetchIndices(t *testing.T) {
	setup := mocks()
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

	baseUrl := initialize(setup)

	client := NewClient(&http.Client{}, baseUrl)
	response, err := client.FetchIndices("sample-.*")

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

	baseUrl := initialize(setup)

	client := NewClient(&http.Client{}, baseUrl)
	_, err := client.FetchIndices("fail")

	if err == nil {
		t.Error("Expected to fail")
	}
}

func TestFailingFetchIndicesRequest(t *testing.T) {
	setup := mocks()
	setup.Status = 500
	setup.Response = "{}"

	baseUrl := initialize(setup)

	client := NewClient(&http.Client{}, baseUrl)
	_, err := client.FetchIndices("fail")

	if err == nil {
		t.Error("Expected to fail")
	}
}

func TestDeleteIndex(t *testing.T) {
	baseUrl := initialize(mocks())

	client := NewClient(&http.Client{}, baseUrl)
	err := client.DeleteIndex("x")

	if err != nil {
		t.Errorf("Unexpected failure: %s", err)
	}
}

func TestFailingDeleteIndexRequest(t *testing.T) {
	setup := mocks()
	setup.Status = 500
	setup.Response = "{}"

	baseUrl := initialize(setup)

	client := NewClient(&http.Client{}, baseUrl)
	err := client.DeleteIndex("failing")

	if err == nil {
		t.Error("Expected to fail")
	}
}

type mockSettings struct {
	Status       int
	Response     string
	ExpectedPath string
}

func mocks() *mockSettings {
	return &mockSettings{
		Status: http.StatusOK,
	}
}

func initialize(m *mockSettings) string {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(m.Status)
		w.Write([]byte(m.Response))
	}))

	return server.URL
}

func TestMain(m *testing.M) {
	defer server.Close()

	os.Exit(m.Run())
}

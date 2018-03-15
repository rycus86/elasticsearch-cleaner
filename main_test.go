package main

import (
	"errors"
	"github.com/rycus86/elasticsearch-cleaner/settings"
	"os"
	"strconv"
	"testing"
	"time"
)

func ExampleStartCommunication() {
	os.Setenv("MAX_INDICES", "1")
	defer os.Unsetenv("MAX_INDICES")

	settings.Initialize()

	StartCommunication(make(chan bool), &mockClient{
		Keys: []string{"abc", "def", "ghi"},
	})

	time.Sleep(30 * time.Millisecond)

	//go WaitForDeletes(deleteChannel, &mockClient{})
	//
	//WaitForEvents(eventsChannel, deleteChannel, &mockClient{
	//	Keys: []string{"abc", "def", "ghi"},
	//})

	// Output:
	// There are 3 items
	// Removing abc ... OK
	// Removing def ... OK
}

func TestCheckIndices(t *testing.T) {
	os.Setenv("MAX_INDICES", "2")
	defer os.Unsetenv("MAX_INDICES")

	settings.Initialize()

	ch := make(chan string, 5)
	defer close(ch)

	CheckIndices(ch, &mockClient{
		Keys: []string{"1", "2", "3", "4", "5"},
	})

	for idx := 0; idx < 3; idx++ {
		key := <-ch

		if key != strconv.Itoa(idx+1) {
			t.Errorf("Unexpected key received: %s (#%d)", key, idx+1)
		}
	}
}

func ExampleCheckIndices() {
	os.Setenv("MAX_INDICES", "2")
	defer os.Unsetenv("MAX_INDICES")

	settings.Initialize()

	CheckIndices(make(chan string, 5), &mockClient{
		Keys: []string{"1", "2", "3", "4", "5"},
	})

	// Output:
	// There are 5 items
}

func TestCheckIndicesFails(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected to panic")
		}
	}()

	CheckIndices(make(chan string, 1), &mockClient{
		FetchError: errors.New("failed to fetch"),
	})
}

func TestLessThanMaximumIndices(t *testing.T) {
	os.Setenv("MAX_INDICES", "10")
	defer os.Unsetenv("MAX_INDICES")

	settings.Initialize()

	ch := make(chan string, 5)
	defer close(ch)

	CheckIndices(ch, &mockClient{
		Keys: []string{},
	})

	ch <- "done"

	first := <-ch
	if first != "done" {
		t.Error("Unexpected item:", first)
	}
}

func ExampleLessThanMaximumIndices() {
	os.Setenv("MAX_INDICES", "10")
	defer os.Unsetenv("MAX_INDICES")

	settings.Initialize()

	CheckIndices(make(chan string), &mockClient{
		Keys: []string{"1", "2"},
	})

	// Output:
	// There are only 2 items ( <= 10 )
}

func TestDeleteIndex(t *testing.T) {
	var client = &mockClient{}
	DeleteIndex("test", client)

	if client.DeletedKey != "test" {
		t.Error("Unexpected key:", client.DeletedKey)
	}
}

func ExampleDeleteIndexError() {
	DeleteIndex("test", &mockClient{
		DeleteError: errors.New("failed to delete"),
	})

	// Output:
	// Failed to remove key test : failed to delete
}

type mockClient struct {
	Keys        []string
	FetchError  error
	DeleteError error
	DeletedKey  string
}

func (m *mockClient) FetchIndices() ([]string, error) {
	var keys []string
	keys, m.Keys = m.Keys, m.Keys[:0]
	return keys, m.FetchError
}

func (m *mockClient) DeleteIndex(key string) error {
	m.DeletedKey = key
	return m.DeleteError
}

func TestMain(m *testing.M) {
	os.Setenv("BASE_URL", "base-url")
	defer os.Unsetenv("BASE_URL")
	os.Setenv("PATTERN", "[0-9]")
	defer os.Unsetenv("PATTERN")
	os.Setenv("MAX_INDICES", "1")
	defer os.Unsetenv("MAX_INDICES")

	os.Exit(m.Run())
}

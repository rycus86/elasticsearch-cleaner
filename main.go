package main

import (
	"fmt"
	"github.com/rycus86/elasticsearch-cleaner/elasticsearch"
	"github.com/rycus86/elasticsearch-cleaner/settings"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var config *settings.ApplicationSettings

func StartCommunication(eventChannel chan bool, es elasticsearch.Client) {
	deleteChannel := make(chan string)

	go WaitForEvents(eventChannel, deleteChannel, es)
	go WaitForDeletes(deleteChannel, es)

	eventChannel <- true
}

func WaitForEvents(eventChannel <-chan bool, deleteChannel chan string, es elasticsearch.Client) {
	for range eventChannel {
		CheckIndices(deleteChannel, es)
	}

	close(deleteChannel)
}

func WaitForDeletes(deleteChannel chan string, es elasticsearch.Client) {
	for key := range deleteChannel {
		DeleteIndex(key, es)
	}
}

func CheckIndices(delete chan<- string, es elasticsearch.Client) {
	keys, err := es.FetchIndices(config.GetPattern())
	if err != nil {
		log.Panic("Failed to fetch indices: ", err)
	}

	numKeys, maxIndices := len(keys), config.GetMaxIndices()
	if numKeys <= maxIndices {
		fmt.Println("There are only", numKeys, "items ( <=", maxIndices, ")")
		return
	} else {
		fmt.Println("There are", numKeys, "items")
	}

	lastKey := numKeys - maxIndices
	for _, key := range keys[0:lastKey] {
		delete <- key
	}
}

func DeleteIndex(key string, es elasticsearch.Client) {
	err := es.DeleteIndex(key)
	if err != nil {
		fmt.Println("Failed to remove key", key, ":", err)
	} else {
		fmt.Println("Removing", key, "... OK")
	}
}

func RunMain(es elasticsearch.Client) {
	settings.Initialize()

	// setup the trigger channel
	eventChannel := make(chan bool)

	// start taking events from the channels
	StartCommunication(eventChannel, es)

	// schedule the repeated runs
	ticker := time.NewTicker(config.GetInterval())

	// setup signal handlers
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	// shorthand function to kick off an execution
	SendEvent := func() { eventChannel <- true }

	// wait for repeat runs and signals
	for {
		select {

		case <-ticker.C: // on repeat runs

			SendEvent()

		case s := <-signals: // on signals

			if s == syscall.SIGHUP {
				SendEvent()
			} else {
				close(eventChannel)
				return
			}

		}
	}
}

func main() {
	config = settings.Initialize()

	client := http.Client{Timeout: config.GetTimeout()}
	es := elasticsearch.NewClient(&client, config.GetBaseUrl())

	RunMain(es)
}

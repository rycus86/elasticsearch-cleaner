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
	keys, err := es.FetchIndices()
	if err != nil {
		log.Panic("Failed to fetch indices: ", err)
	}

	if len(keys) <= settings.GetMaxIndices() {
		fmt.Println("There are only", len(keys), "items ( <=", settings.GetMaxIndices(), ")")
		return
	} else {
		fmt.Println("There are", len(keys), "items")
	}

	lastKey := len(keys) - settings.GetMaxIndices()
	for _, key := range keys[0:lastKey] {
		delete <- key
		//go DeleteIndex(key, es)
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
	ticker := time.NewTicker(settings.GetInterval())

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
	client := http.Client{Timeout: settings.GetTimeout()}
	es := elasticsearch.NewClient(&client)

	RunMain(es)
}

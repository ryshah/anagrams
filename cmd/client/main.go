// Package main provides a test client for the Anagram Finder service.
// It simulates concurrent requests to the anagram server to test performance,
// load handling, and verify the API functionality under concurrent load.
package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"

	"github.com/ryshah/anagrams/pkg/config"
)

// cfg holds the application configuration loaded from config.yaml
var cfg config.Config

// worker is a goroutine function that makes HTTP requests to the anagram server.
// It sends a GET request to find anagrams for a given word and logs the results.
// This function is designed to be run concurrently to simulate load on the server.
//
// Parameters:
//   - id: A unique identifier for this worker (used for tracking/debugging)
//   - word: The word to query for anagrams
//   - wg: A WaitGroup to signal completion of this worker's task
//
// Behavior:
//   - Constructs a URL with the word as a query parameter
//   - Makes an HTTP GET request to the local anagram server
//   - Logs success or error messages based on the response
//   - In debug mode, logs the full response body
//   - Always calls wg.Done() when finished (via defer)
//
// Example:
//
//	var wg sync.WaitGroup
//	wg.Add(1)
//	go worker(1, "listen", &wg)
func worker(id int, word string, wg *sync.WaitGroup) {

	defer wg.Done()

	url := fmt.Sprintf(
		"http://localhost%s/v1/anagrams?word=%s",
		cfg.Server.Port, word,
	)

	resp, err := http.Get(url)
	if err == nil {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != 200 {
			slog.Error("Request error: " + string(body))
		} else {
			slog.Info("Anagrams found for " + word)
			if cfg.Log.Debug {
				slog.Info(word + " => " + string(body))
			}
		}
	}

}

// main is the entry point for the Anagram Finder test client.
// It simulates concurrent load on the anagram server by spawning multiple
// goroutines that make simultaneous requests.
//
// Workflow:
//  1. Loads configuration from config.yaml
//  2. Retrieves the list of test words from configuration
//  3. Creates a specified number of concurrent workers (from config)
//  4. Each worker queries for anagrams of words from the test list (round-robin)
//  5. Waits for all workers to complete before exiting
//
// Configuration parameters used:
//   - Client.TestWords: Array of words to query
//   - Client.ConcurrentRequests: Number of concurrent requests to make
//   - Server.Port: Port where the anagram server is running
//   - Log.Debug: Whether to log detailed response information
//
// This client is useful for:
//   - Load testing the anagram server
//   - Verifying concurrent request handling
//   - Testing cache performance under load
//   - Benchmarking response times
func main() {
	cfgPtr, usedDefaults := config.Load()
	if usedDefaults {
		slog.Warn("Client using default configuration")
	}
	cfg = *cfgPtr
	words := cfg.Client.TestWords
	var wg sync.WaitGroup
	for i := 0; i < cfg.Client.ConcurrentRequests; i++ {

		wg.Add(1)

		go worker(i, words[i%len(words)], &wg)
	}

	wg.Wait()
}

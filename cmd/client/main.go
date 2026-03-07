package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"

	"github.com/ryshah/anagrams/pkg/config"
)

var cfg config.Config

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

// A simple client that simulates multiple requests to the running server
// it iterates through fixed set of words specified in dictionary
func main() {
	cfg = *config.Load()
	words := cfg.Client.TestWords
	var wg sync.WaitGroup
	for i := 0; i < cfg.Client.ConcurrentRequests; i++ {

		wg.Add(1)

		go worker(i, words[i%len(words)], &wg)
	}

	wg.Wait()
}

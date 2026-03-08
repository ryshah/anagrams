package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ryshah/anagrams/pkg/cache"
	"github.com/ryshah/anagrams/pkg/config"
	"github.com/ryshah/anagrams/pkg/middleware"
	"github.com/ryshah/anagrams/pkg/service"
)

type response struct {
	Word     string   `json:"word"`
	Anagrams []string `json:"anagrams"`
}

func anagramHandler(finder *service.AnagramFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var result []string
		word := r.URL.Query().Get("word")
		if finder.Ready() {
			cache, found := lru.Get(word)
			if !found {
				result = finder.Find(word)
				// Add to LRU cache so found next time
				lru.Put(word, result)
			} else {
				result = cache
			}
			json.NewEncoder(w).Encode(response{
				Word:     word,
				Anagrams: result,
			})

		} else {
			http.Error(w, "Error loading the dictionary", http.StatusServiceUnavailable)
		}
	}
}

var cfg config.Config
var lru *cache.LRU

func setupLogging() {
	if cfg.Log.Debug == true {
		options := &slog.HandlerOptions{
			Level: slog.LevelDebug, // Set the minimum level to Debug
		}
		// Create a new logger with the options
		logger := slog.New(slog.NewJSONHandler(os.Stdout, options))

		// Set this new logger as the default for top-level slog functions
		slog.SetDefault(logger)
	}
}

func main() {
	setupLogging()
	lru = cache.New()
	cfg = *config.Load()
	setupLogging()
	finder := service.NewAnagramFinder()
	finder.LoadDictionary()

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/anagrams", anagramHandler(finder))
	mux.Handle("/metrics", promhttp.Handler())
	authMiddleware := middleware.LocalHostAuth()

	handler := middleware.Metrics(authMiddleware(mux))

	slog.Info(cfg.Service.Name + " Server running on " + cfg.Server.Port)
	err := http.ListenAndServe(cfg.Server.Port, handler)
	if err != nil {
		slog.Error("Fatal error starting server" + err.Error())
	}
}

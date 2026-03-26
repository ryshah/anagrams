// Package main provides the HTTP server for the Anagram Finder service.
// It exposes a REST API endpoint for finding anagrams of given words,
// includes Prometheus metrics, and implements LRU caching for performance.
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

// response represents the JSON structure returned by the anagram API endpoint.
// It contains the original word and a list of its anagrams found in the dictionary.
type response struct {
	Word     string   `json:"word"`     // The original word queried
	Anagrams []string `json:"anagrams"` // List of anagrams found for the word
}

// anagramHandler creates an HTTP handler function for the /v1/anagrams endpoint.
// It processes GET requests with a "word" query parameter and returns JSON responses
// containing the word and its anagrams. Results are cached in an LRU cache for performance.
//
// Parameters:
//   - finder: An AnagramFinder instance used to search for anagrams in the loaded dictionary
//
// Returns:
//   - http.HandlerFunc: A handler function that processes anagram lookup requests
//
// Response format:
//   - 200 OK: JSON with word and anagrams array
//   - 503 Service Unavailable: When dictionary is not loaded
//
// Example usage:
//
//	GET /v1/anagrams?word=listen
//	Response: {"word":"listen","anagrams":["silent","enlist"]}
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

// cfg holds the application configuration loaded from config.yaml
var cfg config.Config

// lru is the LRU (Least Recently Used) cache for storing anagram lookup results
var lru *cache.LRU

// setupLogging configures the application's logging system based on the configuration.
// If debug mode is enabled in the config, it sets the log level to Debug and uses
// JSON formatting for structured logging output.
//
// The function reads from the global cfg variable and sets up slog as the default logger.
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

// main is the entry point for the Anagram Finder HTTP server.
// It performs the following initialization steps:
//  1. Sets up logging configuration
//  2. Initializes the LRU cache
//  3. Loads application configuration from config.yaml
//  4. Creates and initializes the AnagramFinder service
//  5. Loads the dictionary into memory
//  6. Sets up HTTP routes and middleware
//  7. Starts the HTTP server
//
// The server exposes the following endpoints:
//   - /v1/anagrams: Main API endpoint for finding anagrams (requires localhost authentication)
//   - /metrics: Prometheus metrics endpoint for monitoring
//
// Middleware chain:
//   - LocalHostAuth: Restricts access to localhost only
//   - Metrics: Collects Prometheus metrics for all requests
//
// The server runs until terminated or encounters a fatal error.
func main() {
	setupLogging()

	// Initialize LRU cache with error handling
	var err error
	lru, err = cache.New()
	if err != nil {
		slog.Error("Fatal error initializing cache: " + err.Error())
		os.Exit(1)
	}

	cfgPtr, usedDefaults := config.Load()
	if usedDefaults {
		slog.Warn("Server using default configuration - config file not found or invalid")
	}
	cfg = *cfgPtr
	setupLogging()

	// Initialize anagram finder and load dictionary
	finder := service.NewAnagramFinder()
	err = finder.LoadDictionary()
	if err != nil {
		slog.Error("Fatal error loading dictionary: " + err.Error())
		os.Exit(1)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/anagrams", anagramHandler(finder))
	mux.Handle("/metrics", promhttp.Handler())
	authMiddleware := middleware.LocalHostAuth()

	handler := middleware.Metrics(authMiddleware(mux))

	slog.Info(cfg.Service.Name + " Server running on " + cfg.Server.Port)
	err = http.ListenAndServe(cfg.Server.Port, handler)
	if err != nil {
		slog.Error("Fatal error starting server: " + err.Error())
		os.Exit(1)
	}
}

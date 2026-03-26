// Package service provides the core anagram finding functionality.
// It loads dictionary files into memory and provides fast anagram lookups
// using a sorted-string key approach.
//
// The service supports:
// - Multiple dictionary files (e.g., English, French)
// - Unicode normalization for international characters
// - Thread-safe dictionary loading
// - Fast O(1) anagram lookups using hash maps
//
// Algorithm:
// Words are indexed by their sorted character representation. For example,
// "listen", "silent", and "enlist" all sort to "eilnst" and are stored
// together in the map under that key.
package service

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ryshah/anagrams/pkg/config"
	"golang.org/x/text/unicode/norm"
)

var (
	// cfg holds the application configuration
	cfg config.Config
	// once ensures dictionary is loaded only once
	once sync.Once
	// loadErr stores any error that occurred during dictionary loading
	loadErr error
)

// AnagramFinder provides anagram lookup functionality.
// It maintains an in-memory map of words indexed by their sorted characters.
//
// Structure:
//   - anagramMap: Maps sorted character keys to sets of words
//   - loaded: Indicates whether the dictionary has been loaded
//
// Example:
//
//	Key "ader" -> {"read", "dear", "dare"}
type AnagramFinder struct {
	anagramMap map[string]map[string]struct{} // Map of sorted keys to word sets
	loaded     bool                           // Whether dictionary is loaded
}

// runeEncoding returns a deterministic encoding of rune frequencies.
// This is an alternative to sortedString for anagram detection.
//
// Algorithm:
//  1. Normalizes the input string
//  2. Counts frequency of each rune
//  3. Sorts runes alphabetically
//  4. Returns concatenated sorted runes
//
// Parameters:
//   - s: The input string to encode
//
// Returns:
//   - string: Deterministic encoding of the string's characters
//
// Example:
//
//	runeEncoding("read") -> "ader"
//	runeEncoding("dear") -> "ader"
//
// Note: Currently not used in favor of sortedString, but kept for reference.
func runeEncoding(s string) string {
	s = normalizeInput(s)

	// rune frequency map
	freq := make(map[rune]int)
	for _, r := range s {
		freq[r]++
	}
	// build deterministic encoding
	var builder strings.Builder

	keys := make([]rune, 0, len(freq))
	for k := range freq {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	for _, r := range keys {
		builder.WriteRune(r)
	}
	return builder.String()
}

// NewAnagramFinder creates and initializes a new AnagramFinder instance.
// It loads the application configuration and prepares the anagram map.
//
// Returns:
//   - *AnagramFinder: A new AnagramFinder ready to load dictionaries
//
// Example:
//
//	finder := service.NewAnagramFinder()
//	finder.LoadDictionary()
//	anagrams := finder.Find("listen")
func NewAnagramFinder() *AnagramFinder {
	cfgPtr, usedDefaults := config.Load()
	if usedDefaults {
		slog.Warn("AnagramFinder using default configuration")
	}
	cfg = *cfgPtr

	return &AnagramFinder{
		anagramMap: make(map[string]map[string]struct{}),
		loaded:     false,
	}
}

// normalizeInput normalizes a string for consistent anagram matching.
// It converts to lowercase and applies Unicode NFC normalization.
//
// Unicode normalization ensures that composed and decomposed characters
// are treated consistently (e.g., "é" vs "e" + accent).
//
// Parameters:
//   - s: The input string to normalize
//
// Returns:
//   - string: Normalized lowercase string
//
// Example:
//
//	normalizeInput("ÉCOLE") -> "école"
//	normalizeInput("e\u0301cole") -> "école" (same as above)
func normalizeInput(s string) string {
	return norm.NFC.String(strings.ToLower(s))
}

// sortedString returns a string with its characters sorted alphabetically.
// This is the key function for anagram detection - anagrams produce the same sorted string.
//
// Algorithm:
//  1. Normalizes the input (lowercase + Unicode normalization)
//  2. Converts to rune slice for proper Unicode handling
//  3. Sorts runes alphabetically
//  4. Returns the sorted string
//
// Parameters:
//   - s: The input string to sort
//
// Returns:
//   - string: String with characters sorted alphabetically
//
// Example:
//
//	sortedString("listen") -> "eilnst"
//	sortedString("silent") -> "eilnst"
//	sortedString("enlist") -> "eilnst"
func sortedString(s string) string {
	s = normalizeInput(s)
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}

// Ready indicates whether the AnagramFinder has loaded its dictionary
// and is ready to process anagram queries.
//
// Returns:
//   - bool: true if dictionary is loaded, false otherwise
//
// Example:
//
//	if finder.Ready() {
//	    anagrams := finder.Find("listen")
//	}
func (a *AnagramFinder) Ready() bool {
	return a.loaded
}

// LoadDictionary loads all configured dictionary files into memory.
// This method is thread-safe and ensures dictionaries are loaded only once,
// even if called concurrently from multiple goroutines.
//
// Process:
//  1. Validates configuration has dictionary files specified
//  2. Reads each dictionary file specified in configuration
//  3. Scans words line by line
//  4. Normalizes and indexes each word by its sorted character key
//  5. Logs loading time for each dictionary
//  6. Marks the service as ready
//
// Dictionary format:
//   - One word per line
//   - Any encoding (UTF-8 recommended for international characters)
//   - Empty lines are skipped
//
// Performance:
//   - Uses sync.Once to ensure single execution
//   - Logs loading time in milliseconds
//   - Efficient O(1) lookup after loading
//
// Returns:
//   - error: Error if dictionary loading fails (file not found, read error, etc.)
//
// Example:
//
//	finder := service.NewAnagramFinder()
//	if err := finder.LoadDictionary(); err != nil {
//	    log.Fatal("Failed to load dictionary:", err)
//	}
//	// Service is now ready for queries
func (a *AnagramFinder) LoadDictionary() error {
	// Use sync.Once to ensure dictionary is loaded only once
	// even if multiple goroutines call LoadDictionary concurrently
	once.Do(func() {
		// Validate configuration
		if len(cfg.Dictionary.Files) == 0 {
			loadErr = errors.New("no dictionary files specified in configuration")
			slog.Error(loadErr.Error())
			return
		}

		start := time.Now()
		for _, filepath := range cfg.Dictionary.Files {
			file, err := os.Open(filepath)
			if err != nil {
				loadErr = fmt.Errorf("unable to open dictionary file '%s': %w", filepath, err)
				slog.Error(loadErr.Error())
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanWords)
			wordCount := 0

			for scanner.Scan() {
				word := scanner.Text()

				// word := normalizeInput(strings.TrimSpace(tempword))
				if word == "" {
					continue
				}
				// All the keys are sorted alphabetically and saved to map
				// each map value will hold all the words that are anagrams
				// for that sorted key
				// Another approach would be to use rune encodings instead
				key := sortedString(word)
				// key := runeEncoding(word)
				if _, ok := a.anagramMap[key]; !ok {
					a.anagramMap[key] = make(map[string]struct{})
				}
				a.anagramMap[key][word] = struct{}{}
				wordCount++
			}

			// Check for scanner errors
			if err := scanner.Err(); err != nil {
				loadErr = fmt.Errorf("error reading dictionary file '%s': %w", filepath, err)
				slog.Error(loadErr.Error())
				return
			}

			elapsed := float64(time.Since(start).Microseconds()) / 1000
			slog.Info(fmt.Sprintf("%s dictionary loaded %d words in %.2fms", filepath, wordCount, elapsed))
		}

		// Mark as loaded only if all dictionaries loaded successfully
		a.loaded = true
		slog.Info("All dictionaries loaded successfully")
	})

	return loadErr
}

// Find searches for anagrams of the given word in the loaded dictionary.
//
// Algorithm:
//  1. Checks if dictionary is loaded
//  2. Sorts the input word to create a lookup key
//  3. Retrieves all words with the same sorted key
//  4. Returns the list of anagrams
//
// Performance:
//   - O(n log n) for sorting the input word (where n is word length)
//   - O(1) for map lookup
//   - O(m) for building result slice (where m is number of anagrams)
//
// Parameters:
//   - word: The word to find anagrams for
//
// Returns:
//   - []string: List of anagrams found (empty if none found or service not ready)
//
// Example:
//
//	anagrams := finder.Find("listen")
//	// Returns: ["silent", "enlist", "inlets", "listen"]
//
// Note: The input word itself is included in the results if it exists in the dictionary.
func (a *AnagramFinder) Find(word string) []string {
	if !a.loaded {
		slog.Error("Service not ready")
		return []string{}
	}

	start := time.Now()

	// preword := normalizeInput(strings.TrimSpace(word))

	key := sortedString(word)
	// key := runeEncoding(word)

	words, ok := a.anagramMap[key]
	results := make([]string, 0, len(words))
	elapsed := float64(time.Since(start).Microseconds()) / 1000

	if !ok || len(word) == 0 {
		slog.Debug("No anagrams found", "word", word, "elapsed_ms", elapsed)
		return results
	}

	for w := range words {
		results = append(results, w)
	}
	slog.Debug("Anagrams found", "count", len(results), "word", word, "elapsed_ms", elapsed)
	return results
}

// a simple service that finds anagrams from
// specified dictionaries
package service

import (
	"bufio"
	"log/slog"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ryshah/anagrams/pkg/config"
	"golang.org/x/text/unicode/norm"
)

var (
	cfg  config.Config
	once sync.Once
)

type AnagramFinder struct {
	anagramMap map[string]map[string]struct{}
	loaded     bool
}

// runeEncoding returns a deterministic encoding of rune frequencies
// Example: "read" -> a1d1e1r1
func runeEncoding(s string) string {

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
		// fmt.Fprintf(&builder, "%d", freq[r])
	}
	return builder.String()
}

func NewAnagramFinder() *AnagramFinder {
	// fmt.Println("\nWelcome to the Anagram Finder")
	// fmt.Println("-----------------------------")
	cfg = *config.Load()

	return &AnagramFinder{
		anagramMap: make(map[string]map[string]struct{}),
		loaded:     false,
	}
}

func normalizeInput(s string) string {
	return norm.NFC.String(strings.ToLower(s))
}

func sortedString(s string) string {
	// s = normalizeInput(s)
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}

// Indicates if the service can be used
func (a *AnagramFinder) Ready() bool {
	return a.loaded
}

// Load all the dictionary files only once
// service will not be marked as ready till all the dictionaries
// are loaded
func (a *AnagramFinder) LoadDictionary() {
	// only do once even if multiple requests waiting
	once.Do(func() {
		start := time.Now()
		for _, filepath := range cfg.Dictionary.Files {
			file, err := os.Open(filepath)
			if err != nil {
				slog.Error("Unable to load dictionary: " + filepath)
				return
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanWords)

			for scanner.Scan() {
				tempword := scanner.Text()

				word := normalizeInput(strings.TrimSpace(tempword))
				if word == "" {
					continue
				}
				// All the keys are sorted alphabetically and saved to map
				// each map value will hold all the words that are anagrams
				// for that sorted key
				// Another approach would be to use rune encodings instead
				key := sortedString(word)
				if _, ok := a.anagramMap[key]; !ok {
					a.anagramMap[key] = make(map[string]struct{})
				}
				a.anagramMap[key][word] = struct{}{}
			}
			slog.Info(filepath + "dictionary loaded in " +
				strconv.FormatFloat(float64(time.Since(start).Microseconds())/1000, 'f', 2, 64) + "ms")
			a.loaded = true
		}
	})
}

// Sort the input word and then search for entry in the map
func (a *AnagramFinder) Find(word string) []string {
	if !a.loaded {
		slog.Error("Service not ready")
		return []string{}
	}

	start := time.Now()
	key := sortedString(word)

	words, ok := a.anagramMap[key]
	results := make([]string, 0, len(words))
	elapsed := float64(time.Since(start).Microseconds()) / 1000

	if !ok || len(word) == 0 {
		slog.Debug("No anagrams found for %s in %.2fms\n", word, elapsed)
		return results
	}

	for w := range words {
		results = append(results, w)
	}
	return results
}

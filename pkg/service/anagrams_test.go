// Package service provides unit tests for the anagram finding service.
// Tests cover string sorting, Unicode normalization, dictionary loading,
// and anagram lookup functionality.
package service

import (
	"os"
	"testing"

	"github.com/ryshah/anagrams/pkg/config"
)

// TestSortedString verifies that the sortedString function correctly
// sorts characters in a word alphabetically.
//
// Test scenario:
//   - Input: "read"
//   - Expected output: "ader"
//
// This is the foundation of the anagram detection algorithm.
func TestSortedString(t *testing.T) {

	result := sortedString("read")
	expected := "ader"

	if result != expected {
		t.Errorf("expected %s but got %s", expected, result)
	}
}

// TestSortedStringAnagramsMatch verifies that anagrams produce the same
// sorted string, which is the key to anagram detection.
//
// Test scenario:
//   - Sorts "read" and "dear"
//   - Verifies both produce the same sorted result
//
// This confirms that anagrams can be identified by comparing sorted strings.
func TestSortedStringAnagramsMatch(t *testing.T) {

	a := sortedString("read")
	b := sortedString("dear")

	if a != b {
		t.Errorf("expected anagrams to match but got %s and %s", a, b)
	}
}

// TestRuneEncoding verifies the alternative rune encoding approach
// for anagram detection.
//
// Test scenario:
//   - Encodes "read" and "dear"
//   - Verifies both produce the same encoding
//
// This tests an alternative algorithm (not currently used in production).
func TestRuneEncoding(t *testing.T) {

	a := runeEncoding("read")
	b := runeEncoding("dear")

	if a != b {
		t.Errorf("expected rune encoding to match but got %s and %s", a, b)
	}
}

// TestRuneEncodingDifferentWords verifies that non-anagrams produce
// different rune encodings.
//
// Test scenario:
//   - Encodes "read" and "book"
//   - Verifies they produce different encodings
//
// This ensures the encoding algorithm can distinguish non-anagrams.
func TestRuneEncodingDifferentWords(t *testing.T) {

	a := runeEncoding("read")
	b := runeEncoding("book")

	if a == b {
		t.Errorf("expected different encoding but got %s", a)
	}
}

// TestLoadDictionaryAndLookup verifies the complete dictionary loading
// and anagram lookup workflow.
//
// Test scenario:
//  1. Creates a temporary dictionary file with test words
//  2. Loads the dictionary into an AnagramFinder
//  3. Verifies that anagrams are correctly grouped together
//
// This is an integration test covering the full anagram finding process.
func TestLoadDictionaryAndLookup(t *testing.T) {

	// Create temporary dictionary file with test data
	cfg, _ := config.Load()
	content := "read\ndear\ndare\nhello\n"
	file, err := os.CreateTemp("", "dict")

	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	// Replace the dictionary files list with only our test file
	cfg.Dictionary.Files = []string{file.Name()}

	file.WriteString(content)
	file.Close()
	finder := NewAnagramFinder()

	err = finder.LoadDictionary()
	if err != nil {
		t.Fatalf("failed to load dictionary: %v", err)
	}

	key := sortedString("read")

	words := finder.anagramMap[key]

	if len(words) != 3 {
		t.Errorf("expected 3 anagrams but got %d", len(words))
	}
}

// TestNoAnagramsFound verifies that the service correctly handles
// queries for words that have no anagrams in the dictionary.
//
// Test scenario:
//  1. Creates a dictionary with limited words
//  2. Searches for a word not in the dictionary
//  3. Verifies that no anagrams are found
//
// This ensures the service handles "not found" cases gracefully.
func TestNoAnagramsFound(t *testing.T) {
	cfg, _ := config.Load()
	content := "read\ndear\n"
	file, err := os.CreateTemp("", "dict")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	file.WriteString(content)
	file.Close()
	// Replace the dictionary files list with only our test file
	cfg.Dictionary.Files = []string{file.Name()}
	finder := NewAnagramFinder()

	err = finder.LoadDictionary()
	if err != nil {
		t.Fatalf("failed to load dictionary: %v", err)
	}

	key := sortedString("xyz")

	words := finder.anagramMap[key]

	if len(words) != 0 {
		t.Errorf("expected no anagrams but found some")
	}
}

// TestNormalizeInput_ComposedVsDecomposed verifies that Unicode normalization
// correctly handles composed vs decomposed character representations.
//
// Test scenario:
//   - "é" as a single composed character
//   - "é" as "e" + combining accent (decomposed)
//   - Verifies both normalize to the same result
//
// This is critical for international character support (French, Spanish, etc.).
func TestNormalizeInput_ComposedVsDecomposed(t *testing.T) {

	// "é" as a single composed character
	composed := "écart"

	// "é" as decomposed form (e + combining accent)
	decomposed := "e\u0301cart"

	n1 := normalizeInput(composed)
	n2 := normalizeInput(decomposed)

	if n1 != n2 {
		t.Errorf("expected normalized strings to match: %s vs %s", n1, n2)
	}
}

// TestNormalizeInput_LowercaseUnicode verifies that normalization
// correctly converts Unicode characters to lowercase.
//
// Test scenario:
//   - Input: "ÉCOLE" (uppercase with accent)
//   - Expected: "école" (lowercase with accent)
//
// This ensures case-insensitive anagram matching for international characters.
func TestNormalizeInput_LowercaseUnicode(t *testing.T) {

	input := "ÉCOLE"
	expected := "école"

	result := normalizeInput(input)

	if result != expected {
		t.Errorf("expected %s but got %s", expected, result)
	}
}

// TestSortedString_FrenchAnagram verifies that the sorting algorithm
// works correctly with French words containing accented characters.
//
// Test scenario:
//   - Sorts "écart" and "tracé" (French anagrams)
//   - Verifies they produce the same sorted result
//
// This confirms international character support in the core algorithm.
func TestSortedString_FrenchAnagram(t *testing.T) {

	a := sortedString("écart")
	b := sortedString("tracé")

	if a != b {
		t.Errorf("expected French anagrams to match but got %s and %s", a, b)
	}
}

// TestSortedString_RuneSorting verifies that runes (Unicode characters)
// are sorted correctly, not just ASCII bytes.
//
// Test scenario:
//   - Input: "çba" (contains cedilla)
//   - Expected: "abç" (sorted by Unicode code point)
//
// This ensures proper Unicode handling in the sorting algorithm.
func TestSortedString_RuneSorting(t *testing.T) {

	input := "çba"
	expected := "abç"

	result := sortedString(input)

	if result != expected {
		t.Errorf("expected %s but got %s", expected, result)
	}
}

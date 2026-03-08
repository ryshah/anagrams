package service

import (
	"os"
	"testing"

	"github.com/ryshah/anagrams/pkg/config"
)

func TestSortedString(t *testing.T) {

	result := sortedString("read")
	expected := "ader"

	if result != expected {
		t.Errorf("expected %s but got %s", expected, result)
	}
}

func TestSortedStringAnagramsMatch(t *testing.T) {

	a := sortedString("read")
	b := sortedString("dear")

	if a != b {
		t.Errorf("expected anagrams to match but got %s and %s", a, b)
	}
}

func TestRuneEncoding(t *testing.T) {

	a := runeEncoding("read")
	b := runeEncoding("dear")

	if a != b {
		t.Errorf("expected rune encoding to match but got %s and %s", a, b)
	}
}

func TestRuneEncodingDifferentWords(t *testing.T) {

	a := runeEncoding("read")
	b := runeEncoding("book")

	if a == b {
		t.Errorf("expected different encoding but got %s", a)
	}
}

func TestLoadDictionaryAndLookup(t *testing.T) {

	// create temporary dictionary file
	cfg := config.Load()
	content := "read\ndear\ndare\nhello\n"
	file, err := os.CreateTemp("", "dict")

	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	cfg.Dictionary.Files = append(cfg.Dictionary.Files, file.Name())

	file.WriteString(content)
	file.Close()
	finder := NewAnagramFinder()

	finder.LoadDictionary()

	key := sortedString("read")

	words := finder.anagramMap[key]

	if len(words) != 3 {
		t.Errorf("expected 3 anagrams but got %d", len(words))
	}
}

func TestNoAnagramsFound(t *testing.T) {
	cfg := config.Load()
	content := "read\ndear\n"
	file, err := os.CreateTemp("", "dict")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	file.WriteString(content)
	file.Close()
	cfg.Dictionary.Files = append(cfg.Dictionary.Files, file.Name())
	finder := NewAnagramFinder()

	finder.LoadDictionary()

	key := sortedString("xyz")

	words := finder.anagramMap[key]

	if len(words) != 0 {
		t.Errorf("expected no anagrams but found some")
	}
}

func TestNormalizeInput_ComposedVsDecomposed(t *testing.T) {

	// "é" composed
	composed := "écart"

	// "é" decomposed (e + accent)
	decomposed := "e\u0301cart"

	n1 := normalizeInput(composed)
	n2 := normalizeInput(decomposed)

	if n1 != n2 {
		t.Errorf("expected normalized strings to match: %s vs %s", n1, n2)
	}
}

func TestNormalizeInput_LowercaseUnicode(t *testing.T) {

	input := "ÉCOLE"
	expected := "école"

	result := normalizeInput(input)

	if result != expected {
		t.Errorf("expected %s but got %s", expected, result)
	}
}

// func TestSortedString_FrenchAnagram(t *testing.T) {

// 	a := sortedString("écart")
// 	b := sortedString("trace")

// 	if a != b {
// 		t.Errorf("expected French anagrams to match but got %s and %s", a, b)
// 	}
// }

func TestSortedString_RuneSorting(t *testing.T) {

	input := "çba"
	expected := "abç"

	result := sortedString(input)

	if result != expected {
		t.Errorf("expected %s but got %s", expected, result)
	}
}

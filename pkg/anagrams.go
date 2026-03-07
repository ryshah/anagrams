package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"golang.org/x/text/unicode/norm"
)

type AnagramFinder struct {
	anagramMap map[string]map[string]struct{}
	loaded     bool
}

// runeEncoding returns a deterministic encoding of rune frequencies
// Example: "read" -> a1d1e1r1
func runeEncoding(s string) string {

	s = strings.ToLower(s)

	// rune frequency map
	freq := make(map[rune]int)

	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			freq[r]++
		}
	}

	// build deterministic encoding
	var builder strings.Builder

	keys := make([]rune, 0, len(freq))
	for k := range freq {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, r := range keys {
		builder.WriteRune(r)
		builder.WriteString(fmt.Sprintf("%d", freq[r]))
	}

	return builder.String()
}

func NewAnagramFinder() *AnagramFinder {
	fmt.Println("\nWelcome to the Anagram Finder")
	fmt.Println("-----------------------------")

	return &AnagramFinder{
		anagramMap: make(map[string]map[string]struct{}),
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

// func sortedString(s string) string {
// 	s = strings.ToLower(s)
// 	chars := strings.Split(s, "")
// 	sort.Strings(chars)
// 	return strings.Join(chars, "")
// }

func (a *AnagramFinder) loadDictionary(path string) error {
	start := time.Now()

	file, err := os.Open(path)
	if err != nil {
		return err
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

		key := sortedString(word)
		// fmt.Println(tempword, " => ", key)
		// key := runeEncoding(word)

		if _, ok := a.anagramMap[key]; !ok {
			a.anagramMap[key] = make(map[string]struct{})
		}

		a.anagramMap[key][word] = struct{}{}
	}

	fmt.Printf("Dictionary loaded in %.2fms\n",
		float64(time.Since(start).Microseconds())/1000)
	fmt.Print(len(a.anagramMap))
	return scanner.Err()
}

func (a *AnagramFinder) find() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nAnagramFinder> ")

		if !scanner.Scan() {
			return
		}

		input := strings.TrimSpace(scanner.Text())

		if strings.EqualFold(input, "exit") {
			fmt.Println()
			return
		}

		start := time.Now()
		key := sortedString(input)

		words, ok := a.anagramMap[key]

		elapsed := float64(time.Since(start).Microseconds()) / 1000

		if !ok || len(words) == 0 {
			fmt.Printf("No anagrams found for %s in %.2fms\n", input, elapsed)
			continue
		}

		results := make([]string, 0, len(words))
		for w := range words {
			results = append(results, w)
		}

		fmt.Printf("%d anagrams found for %s in %.2fms\n",
			len(results), input, elapsed)

		fmt.Println(strings.Join(results, ", "))
	}
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("\nUsage: go run anagrams.go <path-to-dictionary-file>")
		return
	}

	finder := NewAnagramFinder()

	err := finder.loadDictionary(os.Args[1])
	if err != nil {
		fmt.Println("\nError loading dictionary file")
		os.Exit(-1)
	}

	finder.find()
}

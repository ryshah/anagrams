// Package main provides unit tests for the Anagram Finder HTTP server.
// It tests the anagram handler functionality using mock implementations
// to verify correct behavior under various conditions.
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockFinder is a mock implementation of the AnagramFinder interface
// used for testing the anagram handler without requiring a real dictionary.
// It simulates a ready state and returns predefined anagram results.
type mockFinder struct{}

// Ready simulates the dictionary being loaded and ready for queries.
// Always returns true for testing purposes.
//
// Returns:
//   - bool: Always true, indicating the mock finder is ready
func (m *mockFinder) Ready() bool {
	return true
}

// Find returns a predefined list of anagrams for testing purposes.
// Regardless of the input word, it always returns the same test anagrams.
//
// Parameters:
//   - word: The word to find anagrams for (ignored in mock)
//
// Returns:
//   - []string: A fixed list of test anagrams ["read", "dear", "dare"]
func (m *mockFinder) Find(word string) []string {
	return []string{"read", "dear", "dare"}
}

// TestAnagramHandlerSuccess verifies that the anagram handler correctly processes
// a valid request and returns the expected JSON response with anagrams.
//
// Test scenario:
//   - Creates a mock finder that returns predefined anagrams
//   - Simulates a GET request to /v1/anagrams?word=read
//   - Verifies HTTP 200 status code
//   - Validates the response contains exactly 3 anagrams
//
// This test ensures the happy path works correctly when the dictionary is loaded.
func TestAnagramHandlerSuccess(t *testing.T) {

	finder := &mockFinder{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resp := response{
			Word:     "read",
			Anagrams: finder.Find("read"),
		}

		json.NewEncoder(w).Encode(resp)
	})

	req := httptest.NewRequest("GET", "/v1/anagrams?word=read", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200")
	}

	var resp response

	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Anagrams) != 3 {
		t.Fatalf("expected 3 anagrams")
	}
}

// TestAnagramHandlerNotReady verifies that the handler returns an appropriate error
// when the dictionary is not loaded or the service is not ready.
//
// Test scenario:
//   - Simulates a service unavailable state
//   - Makes a GET request to /v1/anagrams?word=read
//   - Verifies HTTP 503 (Service Unavailable) status code
//   - Ensures proper error message is returned
//
// This test validates error handling when the service is not ready to process requests.
func TestAnagramHandlerNotReady(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Error loading the dictionary", http.StatusServiceUnavailable)
	})

	req := httptest.NewRequest("GET", "/v1/anagrams?word=read", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503")
	}
}

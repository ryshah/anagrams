package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockFinder struct{}

func (m *mockFinder) Ready() bool {
	return true
}

func (m *mockFinder) Find(word string) []string {
	return []string{"read", "dear", "dare"}
}

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

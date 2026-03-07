package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Word     string   `json:"word"`
	Anagrams []string `json:"anagrams"`
}

func main() {

	// dict := []string{
	// 	"read", "dear", "dare",
	// 	"écart", "trace", "crate", "react",
	// }

	// finder := anagram.New(dict)

	http.HandleFunc("/v1/anagrams", func(w http.ResponseWriter, r *http.Request) {

		word := r.URL.Query().Get("word")
		result := []string{}
		// result := finder.Find(word)

		resp := response{
			Word:     word,
			Anagrams: result,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("Server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

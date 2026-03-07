package cache

import (
	"fmt"
	"testing"

	"github.com/ryshah/anagrams/pkg/config"
)

func TestLRU(t *testing.T) {
	cfg := config.Load()
	cfg.LRUCache.Capacity = 2
	fmt.Printf("%v+", cfg)
	lru := New()

	lru.Put("a", []string{"1"})
	lru.Put("b", []string{"2"})

	if _, ok := lru.Get("a"); !ok {
		t.Fatal("expected cache hit")
	}

	lru.Put("c", []string{"3"}) // evict b

	if _, ok := lru.Get("b"); ok {
		t.Fatal("expected b to be evicted")
	}
}

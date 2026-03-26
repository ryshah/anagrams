// Package cache provides unit tests for the LRU cache implementation.
// Tests verify cache operations including insertion, retrieval, and eviction.
package cache

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ryshah/anagrams/pkg/config"
)

// TestLRU verifies the basic LRU cache functionality including:
// - Insertion of entries
// - Retrieval of cached entries (cache hits)
// - Eviction of least recently used entries when capacity is reached
//
// Test scenario:
//  1. Creates a cache with capacity of 2
//  2. Inserts entries "a" and "b"
//  3. Verifies "a" can be retrieved (cache hit)
//  4. Inserts entry "c" which should evict "b" (least recently used)
//  5. Verifies "b" is no longer in cache (cache miss)
func TestLRU(t *testing.T) {
	cfg, _ := config.Load()
	cfg.LRUCache.Capacity = 2
	fmt.Printf("%v+", cfg)
	lru, err := New()
	if err != nil {
		t.Fatalf("failed to initialize cache: %v", err)
	}

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

// TestLRUInvalidCapacity verifies that the cache initialization fails
// when an invalid capacity (zero or negative) is provided.
//
// Test scenario:
//  1. Sets cache capacity to 0
//  2. Attempts to create a new cache
//  3. Verifies that an error is returned
//  4. Verifies that the cache pointer is nil
func TestLRUInvalidCapacity(t *testing.T) {
	// Reset the singleton for this test
	once = sync.Once{}
	cache = nil
	initErr = nil

	cfg, _ := config.Load()
	cfg.LRUCache.Capacity = 0

	lru, err := New()
	if err == nil {
		t.Fatal("expected error when capacity is 0, got nil")
	}
	if lru != nil {
		t.Fatal("expected nil cache when initialization fails")
	}

	// Reset for other tests
	once = sync.Once{}
	cache = nil
	initErr = nil
}

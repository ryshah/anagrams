// Package cache implements an LRU (Least Recently Used) caching mechanism
// for storing anagram lookup results in memory.
//
// In production environments, this should be replaced with a distributed
// caching solution like Redis or Memcached for better scalability and
// persistence across server restarts.
//
// The cache is thread-safe and uses a mutex to protect concurrent access.
package cache

import (
	"container/list"
	"errors"
	"log/slog"
	"strconv"
	"sync"

	"github.com/ryshah/anagrams/pkg/config"
)

// entry represents a single cache entry containing a key-value pair.
// The key is the word being queried, and the value is the list of anagrams.
type entry struct {
	key string   // The word that was queried
	val []string // The list of anagrams found for the word
}

// LRU implements a Least Recently Used cache with a fixed capacity.
// When the cache reaches capacity, the least recently accessed item is evicted.
//
// Structure:
//   - cap: Maximum number of entries the cache can hold
//   - ll: Doubly-linked list to maintain access order (most recent at front)
//   - m: Hash map for O(1) lookup of cache entries
//   - mu: Mutex for thread-safe operations
type LRU struct {
	cap int                      // Maximum capacity of the cache
	ll  *list.List               // Doubly-linked list for LRU ordering
	m   map[string]*list.Element // Map for fast lookups
	mu  sync.Mutex               // Mutex for thread safety
}

var (
	// cache is the singleton LRU cache instance
	cache *LRU
	// once ensures the cache is initialized only once
	once sync.Once
	// initErr stores any error that occurred during cache initialization
	initErr error
)

// New creates and returns a singleton LRU cache instance.
// The cache capacity is loaded from the application configuration.
// This function is thread-safe and ensures only one cache instance exists.
//
// Returns:
//   - *LRU: Pointer to the singleton cache instance (nil if initialization failed)
//   - error: Error if cache initialization failed (e.g., invalid capacity)
//
// Errors:
//   - Returns error if cache capacity is zero or negative
//   - Returns error if list initialization fails
//
// Example:
//
//	cache, err := cache.New()
//	if err != nil {
//	    log.Fatal("Failed to initialize cache:", err)
//	}
//	cache.Put("listen", []string{"silent", "enlist"})
func New() (*LRU, error) {
	// Only initialize cache once
	once.Do(func() {
		cfg, usedDefaults := config.Load()
		if usedDefaults {
			slog.Warn("Cache initialization using default configuration")
		}

		// Validate cache capacity
		if cfg.LRUCache.Capacity <= 0 {
			initErr = errors.New("cache capacity must be greater than zero, got: " + strconv.Itoa(cfg.LRUCache.Capacity))
			slog.Error("Failed to initialize cache: " + initErr.Error())
			return
		}

		// Initialize the doubly-linked list
		ll := list.New()
		if ll == nil {
			initErr = errors.New("failed to initialize cache list")
			slog.Error("Failed to initialize cache: " + initErr.Error())
			return
		}

		// Initialize the map
		m := make(map[string]*list.Element)
		if m == nil {
			initErr = errors.New("failed to initialize cache map")
			slog.Error("Failed to initialize cache: " + initErr.Error())
			return
		}

		slog.Info("Initialized cache with " + strconv.Itoa(cfg.LRUCache.Capacity) + " capacity")
		cache = &LRU{
			cap: cfg.LRUCache.Capacity,
			ll:  ll,
			m:   m,
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return cache, nil
}

// Get retrieves a value from the cache if it exists.
// When a cache hit occurs, the accessed entry is moved to the front of the
// LRU list to mark it as recently used.
//
// This method is thread-safe and logs cache hits/misses at debug level.
//
// Parameters:
//   - k: The key (word) to look up in the cache
//
// Returns:
//   - []string: The cached list of anagrams (nil if not found)
//   - bool: true if the key was found (cache hit), false otherwise (cache miss)
//
// Example:
//
//	anagrams, found := cache.Get("listen")
//	if found {
//	    fmt.Println("Cache hit:", anagrams)
//	}
func (l *LRU) Get(k string) ([]string, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if e, ok := l.m[k]; ok {
		l.ll.MoveToFront(e)
		value := e.Value.(*entry)
		slog.Debug("Cache hit for " + k)
		return value.val, true
	}
	slog.Debug("Cache miss for " + k)
	return nil, false
}

// Put adds a new entry to the LRU cache.
// The new entry is placed at the front of the list (most recently used position).
// If the cache is at capacity, the least recently used entry (at the back) is evicted.
//
// This method is thread-safe and logs cache operations at debug/warn level.
//
// Parameters:
//   - k: The key (word) to store in the cache
//   - v: The value (list of anagrams) to associate with the key
//
// Behavior:
//   - New entries are added to the front of the LRU list
//   - If cache is full, the least recently used entry is removed
//   - Logs a warning when eviction occurs
//
// Example:
//
//	cache.Put("listen", []string{"silent", "enlist", "inlets"})
func (l *LRU) Put(k string, v []string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := l.ll.PushFront(&entry{k, v})
	l.m[k] = e
	slog.Debug("Adding to cache " + k)

	if l.ll.Len() > l.cap {
		slog.Warn("Maximum capacity reached, removing an entry")
		b := l.ll.Back()
		l.ll.Remove(b)
		delete(l.m, b.Value.(*entry).key)
	}
}

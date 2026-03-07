// Implements caching functionality
// In real environments, redis or similar databases should be used
package cache

import (
	"container/list"
	"log/slog"
	"strconv"
	"sync"

	"github.com/ryshah/anagrams/pkg/config"
)

type entry struct {
	key string
	val []string
}

type LRU struct {
	cap int
	ll  *list.List
	m   map[string]*list.Element
	mu  sync.Mutex
}

var (
	cache LRU
	once  sync.Once
)

func New() *LRU {
	cfg := config.Load()
	// Only initialize cache once
	once.Do(func() {
		slog.Info("Initialized cache with " + strconv.Itoa(cfg.LRUCache.Capacity) + " capacity")
		cache = LRU{
			cap: cfg.LRUCache.Capacity,
			ll:  list.New(),
			m:   map[string]*list.Element{},
		}
	})
	return &cache
}

// Gets the entry from cache if present, else return nil
// Element accessed is moved to front
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

// Adds the elment to the LRU cache.
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

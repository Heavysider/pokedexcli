package pokecache

import (
	"sync"
	"time"
)

type Cache map[string]cacheEntry

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

var mu *sync.Mutex

func NewCache(interval time.Duration) Cache {
	mu = &sync.Mutex{}
	cache := Cache{}
	go cache.reapLoop(interval)
	return cache
}

func (cache Cache) Add(key string, val []byte) {
	mu.Lock()
	defer mu.Unlock()
	cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (cache Cache) Get(key string) ([]byte, bool) {
	mu.Lock()
	defer mu.Unlock()
	res, ok := cache[key]
	if !ok {
		return nil, ok
	}
	return res.val, ok
}

func (cache Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		t := <-ticker.C
		for k, v := range cache {
			if v.createdAt.Add(interval).Unix() <= t.Unix() {
				mu.Lock()
				delete(cache, k)
				mu.Unlock()
			}
		}
	}
}

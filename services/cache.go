package services

import (
	"sync"
	"time"
)

type CacheItem struct {
	Data      map[string]string
	ExpiresAt time.Time
}

var prayerCache = make(map[string]CacheItem)
var mu sync.RWMutex

func GetCache(key string) (map[string]string, bool) {
	mu.RLock()
	defer mu.RUnlock()

	item, exists := prayerCache[key]
	if !exists || time.Now().After(item.ExpiresAt) {
		return nil, false
	}
	return item.Data, true
}

func SetCache(key string, data map[string]string) {
	mu.Lock()
	defer mu.Unlock()

	prayerCache[key] = CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour),
	}
}

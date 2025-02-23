package in_memory

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"sync"
	"time"
)

type ExpirationType int

const (
	ExpireAfterWrite ExpirationType = iota
	ExpireAfterAccess
)

type CacheItem struct {
	value interface{}
}

type CacheInstance struct {
	lruCache       *lru.Cache[string, CacheItem]
	expirationTime map[string]time.Time
	mutex          sync.RWMutex
	expiration     time.Duration
	expirationType ExpirationType
}

func NewCacheInstance(size int, expiration time.Duration, expType ExpirationType) *CacheInstance {
	lruCache, _ := lru.New[string, CacheItem](size)
	cache := &CacheInstance{
		lruCache:       lruCache,
		expirationTime: make(map[string]time.Time),
		expiration:     expiration,
		expirationType: expType,
	}

	// Start cleanup goroutine
	go cache.cleanup()
	return cache
}

func (c *CacheInstance) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	c.lruCache.Add(key, CacheItem{value: value})

	// Store expiration timestamp based on the mode
	if c.expirationType == ExpireAfterWrite {
		c.expirationTime[key] = now.Add(c.expiration)
	} else if c.expirationType == ExpireAfterAccess {
		c.expirationTime[key] = now.Add(c.expiration)
	}
}

// Get retrieves a value and updates expiration if using "expire after access"
func (c *CacheInstance) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if the key exists
	item, found := c.lruCache.Get(key)
	if !found {
		return nil, false
	}

	now := time.Now()

	if expiry, exists := c.expirationTime[key]; exists && now.After(expiry) {
		c.lruCache.Remove(key)
		delete(c.expirationTime, key)
		return nil, false
	}

	if c.expirationType == ExpireAfterAccess {
		c.expirationTime[key] = now.Add(c.expiration)
	}

	return item.value, true
}

// Delete removes a key
func (c *CacheInstance) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.lruCache.Remove(key)
	delete(c.expirationTime, key)
}

func (c *CacheInstance) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.lruCache.Purge()
	c.expirationTime = make(map[string]time.Time)
}

func (c *CacheInstance) cleanup() {
	for {
		time.Sleep(1 * time.Minute)
		now := time.Now()

		c.mutex.Lock()
		for key, expiry := range c.expirationTime {
			if now.After(expiry) {
				c.lruCache.Remove(key)
				delete(c.expirationTime, key)
			}
		}
		c.mutex.Unlock()
	}
}

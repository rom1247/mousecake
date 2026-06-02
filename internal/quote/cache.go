package quote

import (
	"sync"
	"time"

	"github.com/mousecake-go/mousecake-go/internal/quote/domain"
)

// cacheEntry 缓存条目。
type cacheEntry struct {
	value    *domain.QuoteResult
	expireAt time.Time
}

// MemoryCache 基于 sync.RWMutex + map 的内存缓存，仅缓存 QuoteResult。
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]*cacheEntry
	ttl   time.Duration
}

// NewMemoryCache 创建内存缓存，使用指定 TTL。
func NewMemoryCache(ttl time.Duration) *MemoryCache {
	return &MemoryCache{
		items: make(map[string]*cacheEntry),
		ttl:   ttl,
	}
}

// Get 读取缓存，过期条目自动删除。
func (c *MemoryCache) Get(key string) (*domain.QuoteResult, bool) {
	c.mu.RLock()
	entry, ok := c.items[key]
	if ok && !time.Now().After(entry.expireAt) {
		defer c.mu.RUnlock()
		return entry.value, true
	}
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	// 过期条目：获取写锁后 double-check 再删除
	c.mu.Lock()
	entry, ok = c.items[key]
	if ok && time.Now().After(entry.expireAt) {
		delete(c.items, key)
	}
	c.mu.Unlock()
	return nil, false
}

// Set 写入缓存。
func (c *MemoryCache) Set(key string, value *domain.QuoteResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &cacheEntry{
		value:    value,
		expireAt: time.Now().Add(c.ttl),
	}
}

package quote

import (
	"sync"
	"testing"
	"time"

	"github.com/mousecake-go/mousecake-go/internal/quote/domain"
)

func TestCache_GetMiss(t *testing.T) {
	cache := NewMemoryCache(10 * time.Second)
	result, ok := cache.Get("quote:okx:1:0xA:0xB:1000")
	if ok {
		t.Error("期望缓存未命中, 得到命中")
	}
	if result != nil {
		t.Error("期望 nil, 得到非 nil")
	}
}

func TestCache_GetHit(t *testing.T) {
	cache := NewMemoryCache(10 * time.Second)
	quoteResult := &domain.QuoteResult{
		Provider:   "okx",
		FromAmount: "1000",
		ToAmount:   "2000",
	}
	key := "quote:okx:1:0xA:0xB:1000"
	cache.Set(key, quoteResult)

	result, ok := cache.Get(key)
	if !ok {
		t.Fatal("期望缓存命中, 得到未命中")
	}
	if result.Provider != "okx" {
		t.Errorf("期望 provider=okx, 得到 %s", result.Provider)
	}
}

func TestCache_GetExpired(t *testing.T) {
	cache := NewMemoryCache(50 * time.Millisecond)
	quoteResult := &domain.QuoteResult{Provider: "okx"}
	key := "quote:okx:1:0xA:0xB:1000"
	cache.Set(key, quoteResult)

	// 等待过期
	time.Sleep(100 * time.Millisecond)

	result, ok := cache.Get(key)
	if ok {
		t.Error("期望缓存已过期返回未命中, 得到命中")
	}
	if result != nil {
		t.Error("期望 nil result, 得到非 nil")
	}
}

func TestCache_SetAndGet(t *testing.T) {
	cache := NewMemoryCache(10 * time.Second)
	quoteResult := &domain.QuoteResult{
		Provider:   "okx",
		FromAmount: "1000",
		ToAmount:   "2000",
	}
	cache.Set("quote:okx:1:0xA:0xB:1000", quoteResult)

	result, ok := cache.Get("quote:okx:1:0xA:0xB:1000")
	if !ok {
		t.Fatal("期望缓存命中")
	}
	if result.FromAmount != "1000" {
		t.Errorf("期望 from_amount=1000, 得到 %s", result.FromAmount)
	}
}

func TestCache_ConcurrentReadWrite(t *testing.T) {
	cache := NewMemoryCache(10 * time.Second)
	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key := "quote:okx:1:0xA:0xB:1000"
			cache.Set(key, &domain.QuoteResult{Provider: "okx", FromAmount: "1000"})
			cache.Get(key)
			cache.Get("nonexistent")
		}()
		if i%10 == 0 {
			_ = i // 使用变量避免 lint 警告
		}
	}
	wg.Wait()
}

func TestCache_DifferentKeys(t *testing.T) {
	cache := NewMemoryCache(10 * time.Second)
	cache.Set("quote:okx:1:0xA:0xB:1000", &domain.QuoteResult{Provider: "okx"})
	cache.Set("quote:okx:1:0xA:0xB:2000", &domain.QuoteResult{Provider: "okx", FromAmount: "2000"})

	result1, ok1 := cache.Get("quote:okx:1:0xA:0xB:1000")
	result2, ok2 := cache.Get("quote:okx:1:0xA:0xB:2000")

	if !ok1 || !ok2 {
		t.Fatal("两个 key 都应该命中")
	}
	if result1.FromAmount == result2.FromAmount {
		t.Error("不同 key 应返回不同的结果")
	}
}

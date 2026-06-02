// Package quote 实现代币报价聚合功能，支持多供应商报价比较和交易数据获取。
package quote

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/mousecake-go/mousecake-go/internal/quote/domain"
)

// QuoteProvider 定义报价供应商的统一接口，所有供应商适配器必须实现该接口。
type QuoteProvider interface {
	// Name 返回供应商名称（如 "okx"、"zerox"）。
	Name() string
	// GetQuote 获取报价数据。
	GetQuote(ctx context.Context, params domain.QuoteParams) (*domain.QuoteResult, error)
	// GetSwap 获取交易数据。
	GetSwap(ctx context.Context, params domain.SwapParams) (*domain.SwapResult, error)
}

// ProviderRegistry 供应商注册表，支持按名称注册和查找 QuoteProvider 实例。
type ProviderRegistry struct {
	mu        sync.RWMutex
	providers map[string]QuoteProvider
}

// NewProviderRegistry 创建空的供应商注册表。
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]QuoteProvider),
	}
}

// Register 注册供应商。
func (r *ProviderRegistry) Register(provider QuoteProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[provider.Name()] = provider
}

// Get 按名称查找供应商，未找到时返回 ErrProviderNotFound。
func (r *ProviderRegistry) Get(name string) (QuoteProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("查找供应商 %s: %w", name, domain.ErrProviderNotFound)
	}
	return p, nil
}

// List 返回所有已注册供应商的名称列表（按字母序排列）。
func (r *ProviderRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

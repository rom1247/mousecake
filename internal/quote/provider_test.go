package quote

import (
	"context"
	"errors"
	"testing"

	"github.com/mousecake-go/mousecake-go/internal/quote/domain"
)

// mockProvider 用于测试的模拟供应商。
type mockProvider struct {
	name string
}

func (m *mockProvider) Name() string { return m.name }
func (m *mockProvider) GetQuote(_ context.Context, _ domain.QuoteParams) (*domain.QuoteResult, error) {
	return nil, nil
}
func (m *mockProvider) GetSwap(_ context.Context, _ domain.SwapParams) (*domain.SwapResult, error) {
	return nil, nil
}

func TestProviderRegistry_RegisterAndGet(t *testing.T) {
	registry := NewProviderRegistry()
	p := &mockProvider{name: "okx"}
	registry.Register(p)

	got, err := registry.Get("okx")
	if err != nil {
		t.Fatalf("不期望错误, 得到 %v", err)
	}
	if got.Name() != "okx" {
		t.Errorf("期望 name=okx, 得到 %s", got.Name())
	}
}

func TestProviderRegistry_GetNotFound(t *testing.T) {
	registry := NewProviderRegistry()
	_, err := registry.Get("unknown")
	if err == nil {
		t.Fatal("期望错误, 得到 nil")
	}
	if !errors.Is(err, domain.ErrProviderNotFound) {
		t.Fatalf("期望 ErrProviderNotFound, 得到 %v", err)
	}
}

func TestProviderRegistry_List(t *testing.T) {
	registry := NewProviderRegistry()
	registry.Register(&mockProvider{name: "zerox"})
	registry.Register(&mockProvider{name: "okx"})

	names := registry.List()
	if len(names) != 2 {
		t.Fatalf("期望 2 个供应商, 得到 %d", len(names))
	}
	// 按字母序排列
	if names[0] != "okx" || names[1] != "zerox" {
		t.Errorf("期望 [okx, zerox], 得到 %v", names)
	}
}

func TestProviderRegistry_ListEmpty(t *testing.T) {
	registry := NewProviderRegistry()
	names := registry.List()
	if len(names) != 0 {
		t.Errorf("期望空列表, 得到 %v", names)
	}
}

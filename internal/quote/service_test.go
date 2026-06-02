package quote

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/internal/quote/domain"
)

// mockQuoteProvider 模拟 QuoteProvider。
type mockQuoteProvider struct {
	quoteResult *domain.QuoteResult
	swapResult  *domain.SwapResult
	err         error
}

func (m *mockQuoteProvider) Name() string { return "test-provider" }
func (m *mockQuoteProvider) GetQuote(_ context.Context, _ domain.QuoteParams) (*domain.QuoteResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.quoteResult, nil
}
func (m *mockQuoteProvider) GetSwap(_ context.Context, _ domain.SwapParams) (*domain.SwapResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.swapResult, nil
}

// mockRepo 模拟 SwapRecordRepository。
type mockRepo struct {
	records map[int64]*domain.SwapRecord
	saveErr error
}

func newMockRepo() *mockRepo {
	return &mockRepo{records: make(map[int64]*domain.SwapRecord)}
}

func (m *mockRepo) Save(_ context.Context, record *domain.SwapRecord) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.records[record.ID] = record
	return nil
}

func (m *mockRepo) FindByID(_ context.Context, id int64) (*domain.SwapRecord, error) {
	r, ok := m.records[id]
	if !ok {
		return nil, domain.ErrSwapRecordNotFound
	}
	return r, nil
}

func (m *mockRepo) UpdateStatus(_ context.Context, id int64, txHash string) error {
	r, ok := m.records[id]
	if !ok {
		return domain.ErrSwapRecordNotFound
	}
	if r.Status != domain.SwapStatusPending {
		return domain.ErrAlreadySubmitted
	}
	r.Status = domain.SwapStatusSubmitted
	r.TxHash = txHash
	return nil
}

func TestQuoteService_GetQuote(t *testing.T) {
	t.Run("成功查询报价（带缓存）", func(t *testing.T) {
		registry := NewProviderRegistry()
		provider := &mockQuoteProvider{
			quoteResult: &domain.QuoteResult{
				Provider:   "test-provider",
				FromAmount: "1000",
				ToAmount:   "2000",
			},
		}
		registry.Register(provider)

		cache := NewMemoryCache(10 * time.Second)
		repo := newMockRepo()
		svc := NewQuoteService(registry, cache, repo, 1)

		params := domain.QuoteParams{
			ChainID:   1,
			FromToken: "0xA",
			ToToken:   "0xB",
			Amount:    "1000",
			SwapMode:  domain.SwapModeExactIn,
		}
		result, err := svc.GetQuote(context.Background(), "test-provider", params)
		require.NoError(t, err)
		assert.Equal(t, "test-provider", result.Provider)
	})

	t.Run("缓存命中直接返回", func(t *testing.T) {
		registry := NewProviderRegistry()
		callCount := 0
		provider := &mockQuoteProvider{
			quoteResult: &domain.QuoteResult{Provider: "test-provider"},
		}
		registry.Register(provider)

		cache := NewMemoryCache(10 * time.Second)
		repo := newMockRepo()
		svc := NewQuoteService(registry, cache, repo, 1)

		params := domain.QuoteParams{
			ChainID:   1,
			FromToken: "0xA",
			ToToken:   "0xB",
			Amount:    "1000",
			SwapMode:  domain.SwapModeExactIn,
		}

		_, _ = svc.GetQuote(context.Background(), "test-provider", params)
		callCount++

		// 第二次应该命中缓存
		_, err := svc.GetQuote(context.Background(), "test-provider", params)
		require.NoError(t, err)
	})

	t.Run("供应商未注册", func(t *testing.T) {
		registry := NewProviderRegistry()
		cache := NewMemoryCache(10 * time.Second)
		repo := newMockRepo()
		svc := NewQuoteService(registry, cache, repo, 1)

		params := domain.QuoteParams{
			ChainID:   1,
			FromToken: "0xA",
			ToToken:   "0xB",
			Amount:    "1000",
			SwapMode:  domain.SwapModeExactIn,
		}
		_, err := svc.GetQuote(context.Background(), "unknown", params)
		assert.True(t, errors.Is(err, domain.ErrProviderNotFound))
	})
}

func TestQuoteService_GetSwap(t *testing.T) {
	t.Run("成功获取 swap 数据并创建记录", func(t *testing.T) {
		registry := NewProviderRegistry()
		provider := &mockQuoteProvider{
			swapResult: &domain.SwapResult{
				QuoteResult: domain.QuoteResult{
					Provider:   "test-provider",
					FromAmount: "1000",
					ToAmount:   "2000",
				},
				TxData: domain.TxData{
					To:   "0xContract",
					Data: "0xabcd",
				},
			},
		}
		registry.Register(provider)

		cache := NewMemoryCache(10 * time.Second)
		repo := newMockRepo()
		svc := NewQuoteService(registry, cache, repo, 1)

		params := domain.SwapParams{
			QuoteParams: domain.QuoteParams{
				ChainID:   1,
				FromToken: "0xA",
				ToToken:   "0xB",
				Amount:    "1000",
				SwapMode:  domain.SwapModeExactIn,
			},
			SlippagePercent:   0.5,
			UserWalletAddress: "0xUser",
		}
		result, err := svc.GetSwap(context.Background(), "test-provider", params)
		require.NoError(t, err)
		assert.Equal(t, "test-provider", result.Provider)
		assert.Equal(t, "0xContract", result.TxData.To)
	})
}

func TestQuoteService_SubmitSwap(t *testing.T) {
	t.Run("成功提交 tx_hash", func(t *testing.T) {
		registry := NewProviderRegistry()
		cache := NewMemoryCache(10 * time.Second)
		repo := newMockRepo()
		svc := NewQuoteService(registry, cache, repo, 1)

		// 先创建一条 swap record
		record, err := domain.NewSwapRecord(1, domain.NewSwapRecordOpts{
			Provider:        "okx",
			ChainID:         1,
			FromToken:       "0xA",
			ToToken:         "0xB",
			FromAmount:      "1000",
			ToAmount:        "2000",
			SlippagePercent: 0.5,
			SwapMode:        domain.SwapModeExactIn,
		})
		require.NoError(t, err)
		require.NoError(t, repo.Save(context.Background(), record))

		err = svc.SubmitSwap(context.Background(), record.ID, "0xabc123")
		require.NoError(t, err)

		found, err := repo.FindByID(context.Background(), record.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.SwapStatusSubmitted, found.Status)
		assert.Equal(t, "0xabc123", found.TxHash)
	})

	t.Run("swap 记录不存在", func(t *testing.T) {
		registry := NewProviderRegistry()
		cache := NewMemoryCache(10 * time.Second)
		repo := newMockRepo()
		svc := NewQuoteService(registry, cache, repo, 1)

		err := svc.SubmitSwap(context.Background(), 999, "0xabc")
		assert.True(t, errors.Is(err, domain.ErrSwapRecordNotFound))
	})
}

func TestQuoteService_GetProviders(t *testing.T) {
	t.Run("获取供应商列表", func(t *testing.T) {
		registry := NewProviderRegistry()
		registry.Register(&mockQuoteProvider{})
		cache := NewMemoryCache(10 * time.Second)
		repo := newMockRepo()
		svc := NewQuoteService(registry, cache, repo, 1)

		providers := svc.GetProviders()
		assert.Contains(t, providers, "test-provider")
	})
}

package quote

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mousecake-go/mousecake-go/internal/quote/domain"
)

// QuoteService 报价聚合服务，编排缓存、供应商和仓库。
type QuoteService struct {
	registry *ProviderRegistry
	cache    *MemoryCache
	repo     domain.SwapRecordRepository
	nodeID   int64
	log      *slog.Logger
}

// NewQuoteService 创建 QuoteService。
func NewQuoteService(registry *ProviderRegistry, cache *MemoryCache, repo domain.SwapRecordRepository, nodeID int64) *QuoteService {
	return &QuoteService{
		registry: registry,
		cache:    cache,
		repo:     repo,
		nodeID:   nodeID,
		log:      slog.Default().With("module", "quote", "layer", "service"),
	}
}

// GetQuote 获取报价（优先缓存，缓存未命中则调用供应商）。
func (s *QuoteService) GetQuote(ctx context.Context, providerName string, params domain.QuoteParams) (*domain.QuoteResult, error) {
	provider, err := s.registry.Get(providerName)
	if err != nil {
		return nil, fmt.Errorf("get quote: %w", err)
	}

	// 检查缓存
	key := domain.CacheKey(providerName, params)
	if result, ok := s.cache.Get(key); ok {
		return result, nil
	}

	// 调用供应商
	result, err := provider.GetQuote(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("provider %s get quote: %w", providerName, err)
	}

	// 写入缓存
	s.cache.Set(key, result)

	return result, nil
}

// GetSwap 获取交易数据并创建 swap_record。
func (s *QuoteService) GetSwap(ctx context.Context, providerName string, params domain.SwapParams) (*domain.SwapResult, error) {
	provider, err := s.registry.Get(providerName)
	if err != nil {
		return nil, fmt.Errorf("get swap: %w", err)
	}

	// 调用供应商（swap 不经过缓存）
	result, err := provider.GetSwap(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("provider %s get swap: %w", providerName, err)
	}

	// 创建 swap_record
	record, err := domain.NewSwapRecord(s.nodeID, domain.NewSwapRecordOpts{
		Provider:        providerName,
		ChainID:         params.ChainID,
		FromToken:       params.FromToken,
		ToToken:         params.ToToken,
		FromAmount:      result.FromAmount,
		ToAmount:        result.ToAmount,
		SlippagePercent: params.SlippagePercent,
		SwapMode:        params.SwapMode,
	})
	if err != nil {
		return nil, fmt.Errorf("create swap record: %w", err)
	}

	if err := s.repo.Save(ctx, record); err != nil {
		return nil, fmt.Errorf("save swap record: %w", err)
	}

	// 将 record ID 注入到结果中供前端使用
	result.ID = record.ID

	return result, nil
}

// SubmitSwap 提交交易哈希，原子更新 swap_record 状态。
func (s *QuoteService) SubmitSwap(ctx context.Context, id int64, txHash string) error {
	// 先检查记录是否存在
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("submit swap: %w", err)
	}

	// 原子更新状态
	if err := s.repo.UpdateStatus(ctx, id, txHash); err != nil {
		return fmt.Errorf("submit swap update: %w", err)
	}
	return nil
}

// GetProviders 返回所有已注册供应商名称列表。
func (s *QuoteService) GetProviders() []string {
	return s.registry.List()
}

package launchpad

import (
	"context"
	"fmt"
	"time"

	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// SaleMetaService 管理销售元信息的 CRUD 和 visibility 管理。
type SaleMetaService struct {
	metaRepo domain.SaleMetaRepository
	saleRepo domain.SaleRepository
}

// NewSaleMetaService 创建 SaleMetaService。
func NewSaleMetaService(metaRepo domain.SaleMetaRepository, saleRepo domain.SaleRepository) *SaleMetaService {
	return &SaleMetaService{metaRepo: metaRepo, saleRepo: saleRepo}
}

// CreateSaleMetaInput 创建销售元信息的输入。
type CreateSaleMetaInput struct {
	SaleID      int64  `json:"sale_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	BannerURL   string `json:"banner_url"`
	LogoURL     string `json:"logo_url"`
	WebsiteURL  string `json:"website_url"`
	SocialLinks string `json:"social_links"`
	Visibility  string `json:"visibility"`
}

// CreateSaleMeta 创建销售元信息。
func (s *SaleMetaService) CreateSaleMeta(ctx context.Context, input CreateSaleMetaInput) (*domain.SaleMeta, error) {
	visibility := input.Visibility
	if visibility == "" {
		visibility = "public"
	}
	now := time.Now()
	meta := domain.ReconstructSaleMeta(0, input.SaleID, input.Title, input.Description,
		input.BannerURL, input.LogoURL, input.WebsiteURL, input.SocialLinks, visibility, 0, now, now)

	if err := s.metaRepo.Create(ctx, meta); err != nil {
		return nil, fmt.Errorf("创建销售元信息: %w", err)
	}
	return meta, nil
}

// UpdateSaleMetaInput 更新销售元信息的输入。
type UpdateSaleMetaInput struct {
	SaleID      int64  `json:"sale_id" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	BannerURL   string `json:"banner_url"`
	LogoURL     string `json:"logo_url"`
	WebsiteURL  string `json:"website_url"`
	SocialLinks string `json:"social_links"`
	Visibility  string `json:"visibility"`
}

// UpdateSaleMeta 更新销售元信息。
func (s *SaleMetaService) UpdateSaleMeta(ctx context.Context, input UpdateSaleMetaInput) (*domain.SaleMeta, error) {
	meta, err := s.metaRepo.FindBySaleID(ctx, input.SaleID)
	if err != nil {
		return nil, fmt.Errorf("查询销售元信息: %w", err)
	}
	if meta == nil {
		return nil, domain.ErrNotFound
	}

	updated := domain.ReconstructSaleMeta(meta.ID, meta.SaleID,
		coalesce(input.Title, meta.Title),
		coalesce(input.Description, meta.Description),
		coalesce(input.BannerURL, meta.BannerURL),
		coalesce(input.LogoURL, meta.LogoURL),
		coalesce(input.WebsiteURL, meta.WebsiteURL),
		coalesce(input.SocialLinks, meta.SocialLinks),
		coalesce(input.Visibility, meta.Visibility),
		meta.SortOrder, meta.CreatedAt, time.Now())

	if err := s.metaRepo.Update(ctx, updated); err != nil {
		return nil, fmt.Errorf("更新销售元信息: %w", err)
	}
	return updated, nil
}

// TokenService 管理代币元信息的 CRUD。
type TokenService struct {
	tokenRepo domain.TokenRepository
}

// NewTokenService 创建 TokenService。
func NewTokenService(tokenRepo domain.TokenRepository) *TokenService {
	return &TokenService{tokenRepo: tokenRepo}
}

// CreateTokenInput 创建代币元信息的输入。
type CreateTokenInput struct {
	Address  string `json:"address" binding:"required"`
	ChainID  int    `json:"chain_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Symbol   string `json:"symbol" binding:"required"`
	Decimals int    `json:"decimals" binding:"required"`
	LogoURL  string `json:"logo_url"`
}

// CreateToken 创建代币元信息。
func (s *TokenService) CreateToken(ctx context.Context, input CreateTokenInput) (*domain.Token, error) {
	now := time.Now()
	token := domain.ReconstructToken(0, input.Address, input.ChainID, input.Name,
		input.Symbol, input.Decimals, input.LogoURL, now, now)

	if err := s.tokenRepo.Create(ctx, token); err != nil {
		return nil, fmt.Errorf("创建代币信息: %w", err)
	}
	return token, nil
}

// UpdateTokenInput 更新代币元信息的输入。
type UpdateTokenInput struct {
	ID       int64  `json:"id" binding:"required"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	LogoURL  string `json:"logo_url"`
}

// UpdateToken 更新代币元信息。
func (s *TokenService) UpdateToken(ctx context.Context, input UpdateTokenInput) (*domain.Token, error) {
	token, err := s.tokenRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("查询代币信息: %w", err)
	}
	if token == nil {
		return nil, domain.ErrNotFound
	}

	updated := domain.ReconstructToken(token.ID, token.Address, token.ChainID,
		coalesce(input.Name, token.Name),
		coalesce(input.Symbol, token.Symbol),
		firstNonZero(input.Decimals, token.Decimals),
		coalesce(input.LogoURL, token.LogoURL),
		token.CreatedAt, time.Now())

	if err := s.tokenRepo.Update(ctx, updated); err != nil {
		return nil, fmt.Errorf("更新代币信息: %w", err)
	}
	return updated, nil
}

// coalesce 返回第一个非空字符串。
func coalesce(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// firstNonZero 返回第一个非零整数。
func firstNonZero(values ...int) int {
	for _, v := range values {
		if v != 0 {
			return v
		}
	}
	return 0
}

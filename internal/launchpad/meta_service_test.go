package launchpad

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// --- Mock 仓库 ---

type mockSaleMetaRepo struct {
	metas map[int64]*domain.SaleMeta
	err   error
}

func newMockSaleMetaRepo() *mockSaleMetaRepo {
	return &mockSaleMetaRepo{metas: make(map[int64]*domain.SaleMeta)}
}

func (m *mockSaleMetaRepo) FindBySaleID(_ context.Context, saleID int64) (*domain.SaleMeta, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, meta := range m.metas {
		if meta.SaleID == saleID {
			return meta, nil
		}
	}
	return nil, nil
}

func (m *mockSaleMetaRepo) Create(_ context.Context, meta *domain.SaleMeta) error {
	if m.err != nil {
		return m.err
	}
	meta.ID = int64(len(m.metas) + 1)
	m.metas[meta.ID] = meta
	return nil
}

func (m *mockSaleMetaRepo) Update(_ context.Context, meta *domain.SaleMeta) error {
	if m.err != nil {
		return m.err
	}
	m.metas[meta.ID] = meta
	return nil
}

type mockTokenRepo struct {
	tokens map[int64]*domain.Token
	nextID int64
	err    error
}

func newMockTokenRepo() *mockTokenRepo {
	return &mockTokenRepo{tokens: make(map[int64]*domain.Token), nextID: 1}
}

func (m *mockTokenRepo) FindByID(_ context.Context, id int64) (*domain.Token, error) {
	if m.err != nil {
		return nil, m.err
	}
	t, ok := m.tokens[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (m *mockTokenRepo) FindByAddress(_ context.Context, _ string, _ int) (*domain.Token, error) {
	return nil, nil
}

func (m *mockTokenRepo) Create(_ context.Context, token *domain.Token) error {
	if m.err != nil {
		return m.err
	}
	token.ID = m.nextID
	m.nextID++
	m.tokens[token.ID] = token
	return nil
}

func (m *mockTokenRepo) Update(_ context.Context, token *domain.Token) error {
	if m.err != nil {
		return m.err
	}
	m.tokens[token.ID] = token
	return nil
}

type mockSaleRepo struct {
	sales map[int64]*domain.Sale
}

func newMockSaleRepoForMeta() *mockSaleRepo {
	return &mockSaleRepo{sales: make(map[int64]*domain.Sale)}
}

func (m *mockSaleRepo) FindByID(_ context.Context, id int64) (*domain.Sale, error) {
	s, ok := m.sales[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return s, nil
}

func (m *mockSaleRepo) FindByContractAddress(_ context.Context, _ string) (*domain.Sale, error) {
	return nil, nil
}

func (m *mockSaleRepo) FindPublicList(_ context.Context, _, _ int) ([]*domain.Sale, int64, error) {
	return nil, 0, nil
}

func (m *mockSaleRepo) Create(_ context.Context, _ *domain.Sale) error { return nil }

func (m *mockSaleRepo) Update(_ context.Context, _ *domain.Sale) error { return nil }

// --- SaleMeta CRUD 测试（Task 7.3）---

func TestSaleMetaService_CreateSaleMeta(t *testing.T) {
	metaRepo := newMockSaleMetaRepo()
	saleRepo := newMockSaleRepoForMeta()
	svc := NewSaleMetaService(metaRepo, saleRepo)

	meta, err := svc.CreateSaleMeta(context.Background(), CreateSaleMetaInput{
		SaleID:      1,
		Title:       "测试 IDO",
		Description: "描述",
		Visibility:  "public",
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), meta.SaleID)
	assert.Equal(t, "测试 IDO", meta.Title)
	assert.Equal(t, "public", meta.Visibility)
	assert.False(t, meta.CreatedAt.IsZero())
}

func TestSaleMetaService_CreateSaleMeta_DefaultVisibility(t *testing.T) {
	metaRepo := newMockSaleMetaRepo()
	saleRepo := newMockSaleRepoForMeta()
	svc := NewSaleMetaService(metaRepo, saleRepo)

	meta, err := svc.CreateSaleMeta(context.Background(), CreateSaleMetaInput{
		SaleID: 1,
		Title:  "测试",
	})
	require.NoError(t, err)
	assert.Equal(t, "public", meta.Visibility)
}

func TestSaleMetaService_UpdateSaleMeta(t *testing.T) {
	metaRepo := newMockSaleMetaRepo()
	saleRepo := newMockSaleRepoForMeta()
	svc := NewSaleMetaService(metaRepo, saleRepo)

	// 先创建
	created, err := svc.CreateSaleMeta(context.Background(), CreateSaleMetaInput{
		SaleID: 1,
		Title:  "原始标题",
	})
	require.NoError(t, err)

	// 更新
	updated, err := svc.UpdateSaleMeta(context.Background(), UpdateSaleMetaInput{
		SaleID: 1,
		Title:  "新标题",
	})
	require.NoError(t, err)
	assert.Equal(t, "新标题", updated.Title)
	assert.Equal(t, created.ID, updated.ID)
}

func TestSaleMetaService_UpdateSaleMeta_NotFound(t *testing.T) {
	metaRepo := newMockSaleMetaRepo()
	saleRepo := newMockSaleRepoForMeta()
	svc := NewSaleMetaService(metaRepo, saleRepo)

	_, err := svc.UpdateSaleMeta(context.Background(), UpdateSaleMetaInput{
		SaleID: 999,
		Title:  "不存在",
	})
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestSaleMetaService_CreateSaleMeta_RepoError(t *testing.T) {
	metaRepo := newMockSaleMetaRepo()
	metaRepo.err = assert.AnError
	saleRepo := newMockSaleRepoForMeta()
	svc := NewSaleMetaService(metaRepo, saleRepo)

	_, err := svc.CreateSaleMeta(context.Background(), CreateSaleMetaInput{
		SaleID: 1,
		Title:  "测试",
	})
	assert.Error(t, err)
}

// --- Token CRUD 测试（Task 7.5）---

func TestTokenService_CreateToken(t *testing.T) {
	tokenRepo := newMockTokenRepo()
	svc := NewTokenService(tokenRepo)

	token, err := svc.CreateToken(context.Background(), CreateTokenInput{
		Address:  "0xToken",
		ChainID:  1,
		Name:     "Test Token",
		Symbol:   "TT",
		Decimals: 18,
		LogoURL:  "https://example.com/logo.png",
	})
	require.NoError(t, err)
	assert.Equal(t, "0xToken", token.Address)
	assert.Equal(t, 1, token.ChainID)
	assert.Equal(t, "Test Token", token.Name)
	assert.Equal(t, "TT", token.Symbol)
	assert.Equal(t, 18, token.Decimals)
	assert.False(t, token.CreatedAt.IsZero())
}

func TestTokenService_UpdateToken(t *testing.T) {
	tokenRepo := newMockTokenRepo()
	svc := NewTokenService(tokenRepo)

	// 先创建
	created, err := svc.CreateToken(context.Background(), CreateTokenInput{
		Address:  "0xToken",
		ChainID:  1,
		Name:     "Old Name",
		Symbol:   "TT",
		Decimals: 18,
	})
	require.NoError(t, err)

	// 更新
	updated, err := svc.UpdateToken(context.Background(), UpdateTokenInput{
		ID:      created.ID,
		Name:    "New Name",
		Symbol:  "NT",
		LogoURL: "https://example.com/new.png",
	})
	require.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)
	assert.Equal(t, "NT", updated.Symbol)
	assert.Equal(t, "https://example.com/new.png", updated.LogoURL)
	assert.Equal(t, 18, updated.Decimals) // 未修改字段保持不变
}

func TestTokenService_UpdateToken_NotFound(t *testing.T) {
	tokenRepo := newMockTokenRepo()
	svc := NewTokenService(tokenRepo)

	_, err := svc.UpdateToken(context.Background(), UpdateTokenInput{
		ID:   999,
		Name: "不存在",
	})
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestTokenService_CreateToken_RepoError(t *testing.T) {
	tokenRepo := newMockTokenRepo()
	tokenRepo.err = assert.AnError
	svc := NewTokenService(tokenRepo)

	_, err := svc.CreateToken(context.Background(), CreateTokenInput{
		Address:  "0xToken",
		ChainID:  1,
		Name:     "Test",
		Symbol:   "T",
		Decimals: 18,
	})
	assert.Error(t, err)
}

// --- 辅助测试 ---

func TestCoalesce(t *testing.T) {
	assert.Equal(t, "first", coalesce("first", "second"))
	assert.Equal(t, "second", coalesce("", "second"))
	assert.Equal(t, "", coalesce("", ""))
}

func TestFirstNonZero(t *testing.T) {
	assert.Equal(t, 1, firstNonZero(1, 2))
	assert.Equal(t, 2, firstNonZero(0, 2))
	assert.Equal(t, 0, firstNonZero(0, 0))
}

// 确保 mock 类型实现接口
var _ domain.SaleMetaRepository = (*mockSaleMetaRepo)(nil)
var _ domain.TokenRepository = (*mockTokenRepo)(nil)
var _ domain.SaleRepository = (*mockSaleRepo)(nil)

// 用于修复 mockSaleRepo 的重复声明——让 admin_test 使用不同的 mock 名称
// prepare_service_test.go 中没有 saleRepo mock，不会冲突

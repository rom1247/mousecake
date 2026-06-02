package launchpad

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// --- 用户操作 Mock ---

type mockUserSaleRepo struct {
	sale *domain.Sale
	err  error
}

func (m *mockUserSaleRepo) FindByID(_ context.Context, _ int64) (*domain.Sale, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.sale, nil
}

func (m *mockUserSaleRepo) FindByContractAddress(_ context.Context, _ string) (*domain.Sale, error) {
	return nil, nil
}

func (m *mockUserSaleRepo) FindPublicList(_ context.Context, _, _ int) ([]*domain.Sale, int64, error) {
	return nil, 0, nil
}

func (m *mockUserSaleRepo) Create(_ context.Context, _ *domain.Sale) error { return nil }
func (m *mockUserSaleRepo) Update(_ context.Context, _ *domain.Sale) error { return nil }

type mockUserPoolRepo struct {
	pool *domain.Pool
	err  error
}

func (m *mockUserPoolRepo) FindBySaleID(_ context.Context, _ int64) ([]*domain.Pool, error) {
	return nil, nil
}

func (m *mockUserPoolRepo) FindBySaleAndPool(_ context.Context, _ int64, _ int) (*domain.Pool, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.pool, nil
}

func (m *mockUserPoolRepo) Create(_ context.Context, _ *domain.Pool) error { return nil }
func (m *mockUserPoolRepo) Update(_ context.Context, _ *domain.Pool) error { return nil }

type mockUserTierLimitRepo struct {
	limit *domain.TierLimit
	err   error
}

func (m *mockUserTierLimitRepo) FindBySaleID(_ context.Context, _ int64) ([]*domain.TierLimit, error) {
	return nil, nil
}

func (m *mockUserTierLimitRepo) FindBySaleAndTier(_ context.Context, _ int64, _ int) (*domain.TierLimit, error) {
	return m.limit, m.err
}

func (m *mockUserTierLimitRepo) Save(_ context.Context, _ *domain.TierLimit) error { return nil }

type mockUserWhitelistRepo struct {
	result bool
	err    error
}

func (m *mockUserWhitelistRepo) IsWhitelisted(_ context.Context, _ int64, _ string) (bool, error) {
	return m.result, m.err
}

type mockUserUserPoolStateRepo struct {
	state *domain.UserPoolState
	err   error
}

func (m *mockUserUserPoolStateRepo) FindByUserAndPool(_ context.Context, _ string, _ int64, _ int) (*domain.UserPoolState, error) {
	return m.state, m.err
}

type mockUserCreditRepo struct {
	credit *domain.UserCredit
	err    error
}

func (m *mockUserCreditRepo) FindByUserAndSale(_ context.Context, _ string, _ int64) (*domain.UserCredit, error) {
	return m.credit, m.err
}

// 辅助：创建测试用的普通池
func makeNormalPool(saleID int64) *domain.Pool {
	return domain.ReconstructPool(saleID, saleID, 0,
		big.NewInt(10000), big.NewInt(50000), big.NewInt(5000),
		nil, nil, nil,
		false, false, 0, 0, 0, 0, now, now)
}

// 辅助：创建测试用的特殊池
func makeSpecialPool(saleID int64) *domain.Pool {
	return domain.ReconstructPool(saleID, saleID, 0,
		big.NewInt(10000), big.NewInt(50000), big.NewInt(5000),
		nil, nil, nil,
		true, false, 0, 0, 0, 0, now, now)
}

// 辅助：创建测试用的 UserService
func newTestUserService(
	saleRepo domain.SaleRepository,
	poolRepo domain.PoolRepository,
	tierLimitRepo domain.TierLimitRepository,
	whitelistRepo domain.WhitelistRepository,
	userPoolStateRepo domain.UserPoolStateRepository,
	creditRepo domain.UserCreditRepository,
) *UserService {
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("test_calldata"), hash: "test_hash"})
	chain := &mockChainReader{}
	prepareSvc := NewPrepareService(newMockPrepareTxRepo(), chain, encoder, 30*time.Minute)
	querySvc := newQuerySvc(saleRepo, poolRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	return NewUserService(prepareSvc, querySvc, encoder, "MousePadByTier", chain,
		saleRepo, poolRepo, tierLimitRepo, whitelistRepo, userPoolStateRepo, creditRepo)
}

// --- Deposit 测试（Task 9.1）---

func TestUserService_Deposit_NormalPool(t *testing.T) {
	sale := domain.ReconstructSale(1, "0x1", 1, "", "", "", "", "", 0, 0, 0, false, 0, now, now, nil)
	saleRepo := &mockUserSaleRepo{sale: sale}
	poolRepo := &mockUserPoolRepo{pool: makeNormalPool(1)}
	tierLimitRepo := &mockUserTierLimitRepo{}
	userPoolStateRepo := &mockUserUserPoolStateRepo{}

	svc := newTestUserService(saleRepo, poolRepo, tierLimitRepo, nil, userPoolStateRepo, nil)

	tx, err := svc.Deposit(context.Background(), DepositInput{
		CallerAddress: "0xUser",
		SaleID:        1,
		PoolIndex:     0,
		Amount:        "1000",
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpDeposit, tx.OperationType)
}

func TestUserService_Deposit_ExceedTierLimit(t *testing.T) {
	sale := domain.ReconstructSale(1, "0x1", 1, "", "", "", "", "", 0, 0, 0, false, 0, now, now, nil)
	saleRepo := &mockUserSaleRepo{sale: sale}
	poolRepo := &mockUserPoolRepo{pool: makeNormalPool(1)}
	tierLimitRepo := &mockUserTierLimitRepo{
		limit: domain.ReconstructTierLimit(1, 1, 1, big.NewInt(500), now, now),
	}
	// 用户已投入 400
	userPoolStateRepo := &mockUserUserPoolStateRepo{
		state: domain.ReconstructUserPoolState(1, 1, 0, "0xUser", big.NewInt(400), false, now, now),
	}

	svc := newTestUserService(saleRepo, poolRepo, tierLimitRepo, nil, userPoolStateRepo, nil)

	// 尝试投入 200（400+200=600 > 500 限制）
	_, err := svc.Deposit(context.Background(), DepositInput{
		CallerAddress: "0xUser",
		SaleID:        1,
		PoolIndex:     0,
		Amount:        "200",
	})
	assert.Error(t, err)
}

func TestUserService_Deposit_SpecialPool_NotWhitelisted(t *testing.T) {
	sale := domain.ReconstructSale(1, "0x1", 1, "", "", "", "", "", 0, 0, 0, false, 0, now, now, nil)
	saleRepo := &mockUserSaleRepo{sale: sale}
	poolRepo := &mockUserPoolRepo{pool: makeSpecialPool(1)}
	wlRepo := &mockUserWhitelistRepo{result: false}
	tierLimitRepo := &mockUserTierLimitRepo{}

	svc := newTestUserService(saleRepo, poolRepo, tierLimitRepo, wlRepo, nil, nil)

	_, err := svc.Deposit(context.Background(), DepositInput{
		CallerAddress: "0xUser",
		SaleID:        1,
		PoolIndex:     0,
		Amount:        "1000",
	})
	assert.Error(t, err)
}

func TestUserService_Deposit_SpecialPool_Whitelisted(t *testing.T) {
	sale := domain.ReconstructSale(1, "0x1", 1, "", "", "", "", "", 0, 0, 0, false, 0, now, now, nil)
	saleRepo := &mockUserSaleRepo{sale: sale}
	poolRepo := &mockUserPoolRepo{pool: makeSpecialPool(1)}
	wlRepo := &mockUserWhitelistRepo{result: true}
	tierLimitRepo := &mockUserTierLimitRepo{}

	svc := newTestUserService(saleRepo, poolRepo, tierLimitRepo, wlRepo, nil, nil)

	tx, err := svc.Deposit(context.Background(), DepositInput{
		CallerAddress: "0xUser",
		SaleID:        1,
		PoolIndex:     0,
		Amount:        "1000",
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpDeposit, tx.OperationType)
}

func TestUserService_Deposit_InvalidAmount(t *testing.T) {
	sale := domain.ReconstructSale(1, "0x1", 1, "", "", "", "", "", 0, 0, 0, false, 0, now, now, nil)
	saleRepo := &mockUserSaleRepo{sale: sale}
	poolRepo := &mockUserPoolRepo{pool: makeNormalPool(1)}
	tierLimitRepo := &mockUserTierLimitRepo{}

	svc := newTestUserService(saleRepo, poolRepo, tierLimitRepo, nil, nil, nil)

	_, err := svc.Deposit(context.Background(), DepositInput{
		CallerAddress: "0xUser",
		SaleID:        1,
		PoolIndex:     0,
		Amount:        "invalid",
	})
	assert.Error(t, err)
}

// --- Harvest 测试（Task 9.3）---

func TestUserService_Harvest(t *testing.T) {
	sale := domain.ReconstructSale(1, "0x1", 1, "", "", "", "", "", 0, 0, 0, false, 0, now, now, nil)
	saleRepo := &mockUserSaleRepo{sale: sale}

	svc := newTestUserService(saleRepo, nil, nil, nil, nil, nil)

	tx, err := svc.Harvest(context.Background(), HarvestInput{
		CallerAddress: "0xUser",
		SaleID:        1,
		PoolIndex:     0,
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpHarvest, tx.OperationType)
}

func TestUserService_Harvest_SaleNotFound(t *testing.T) {
	saleRepo := &mockUserSaleRepo{err: domain.ErrNotFound}

	svc := newTestUserService(saleRepo, nil, nil, nil, nil, nil)

	_, err := svc.Harvest(context.Background(), HarvestInput{
		CallerAddress: "0xUser",
		SaleID:        999,
		PoolIndex:     0,
	})
	assert.Error(t, err)
}

// --- Release 测试（Task 9.5）---

func TestUserService_Release(t *testing.T) {
	svc := newTestUserService(nil, nil, nil, nil, nil, nil)

	tx, err := svc.Release(context.Background(), ReleaseInput{
		CallerAddress: "0xUser",
		ScheduleID:    1,
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpRelease, tx.OperationType)
}

// 确保接口实现
var _ domain.SaleRepository = (*mockUserSaleRepo)(nil)
var _ domain.PoolRepository = (*mockUserPoolRepo)(nil)
var _ domain.TierLimitRepository = (*mockUserTierLimitRepo)(nil)
var _ domain.WhitelistRepository = (*mockUserWhitelistRepo)(nil)
var _ domain.UserPoolStateRepository = (*mockUserUserPoolStateRepo)(nil)
var _ domain.UserCreditRepository = (*mockUserCreditRepo)(nil)

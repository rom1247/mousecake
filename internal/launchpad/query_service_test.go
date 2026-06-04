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

// --- 查询服务 Mock 仓库 ---

type mockQuerySaleRepo struct {
	sale  *domain.Sale
	sales []*domain.Sale
	total int64
	err   error
}

func (m *mockQuerySaleRepo) FindByID(_ context.Context, _ int64) (*domain.Sale, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.sale, nil
}

func (m *mockQuerySaleRepo) FindByContractAddress(_ context.Context, _ string) (*domain.Sale, error) {
	return nil, nil
}

func (m *mockQuerySaleRepo) FindPublicList(_ context.Context, _, _ int) ([]*domain.Sale, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.sales, m.total, nil
}

func (m *mockQuerySaleRepo) Create(_ context.Context, _ *domain.Sale) error { return nil }
func (m *mockQuerySaleRepo) Update(_ context.Context, _ *domain.Sale) error { return nil }

type mockQueryPoolRepo struct {
	pools []*domain.Pool
	pool  *domain.Pool
	err   error
}

func (m *mockQueryPoolRepo) FindBySaleID(_ context.Context, _ int64) ([]*domain.Pool, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.pools, nil
}

func (m *mockQueryPoolRepo) FindBySaleAndPool(_ context.Context, _ int64, _ int) (*domain.Pool, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.pool, nil
}

func (m *mockQueryPoolRepo) Create(_ context.Context, _ *domain.Pool) error { return nil }
func (m *mockQueryPoolRepo) Update(_ context.Context, _ *domain.Pool) error { return nil }

type mockQueryWhitelistRepo struct {
	result bool
	err    error
}

func (m *mockQueryWhitelistRepo) IsWhitelisted(_ context.Context, _ int64, _ string) (bool, error) {
	return m.result, m.err
}

type mockQueryDepositRepo struct {
	deposits []*domain.Deposit
	err      error
}

func (m *mockQueryDepositRepo) FindByUserAndSale(_ context.Context, _ string, _ int64) ([]*domain.Deposit, error) {
	return m.deposits, m.err
}

func (m *mockQueryDepositRepo) FindByUserAndPool(_ context.Context, _ string, _ int64, _ int) ([]*domain.Deposit, error) {
	return nil, nil
}

type mockQueryHarvestRepo struct {
	harvests []*domain.Harvest
	err      error
}

func (m *mockQueryHarvestRepo) FindByUserAndSale(_ context.Context, _ string, _ int64) ([]*domain.Harvest, error) {
	return m.harvests, m.err
}

type mockQueryVestingRepo struct {
	schedules []*domain.VestingSchedule
	schedule  *domain.VestingSchedule
	err       error
}

func (m *mockQueryVestingRepo) FindByBeneficiary(_ context.Context, _ string) ([]*domain.VestingSchedule, error) {
	return m.schedules, m.err
}

func (m *mockQueryVestingRepo) FindByID(_ context.Context, _ int64) (*domain.VestingSchedule, error) {
	return m.schedule, m.err
}

type mockQueryReleaseRepo struct {
	releases []*domain.VestingRelease
	err      error
}

func (m *mockQueryReleaseRepo) FindByScheduleID(_ context.Context, _ int64) ([]*domain.VestingRelease, error) {
	return m.releases, m.err
}

type mockQueryUserPoolRepo struct {
	state *domain.UserPoolState
	err   error
}

func (m *mockQueryUserPoolRepo) FindByUserAndPool(_ context.Context, _ string, _ int64, _ int) (*domain.UserPoolState, error) {
	return m.state, m.err
}

type mockQueryCreditRepo struct {
	credit *domain.UserCredit
	err    error
}

func (m *mockQueryCreditRepo) FindByUserAndSale(_ context.Context, _ string, _ int64) (*domain.UserCredit, error) {
	return m.credit, m.err
}

type mockQueryTierLimitRepo struct {
	limits []*domain.TierLimit
	limit  *domain.TierLimit
	err    error
}

func (m *mockQueryTierLimitRepo) FindBySaleID(_ context.Context, _ int64) ([]*domain.TierLimit, error) {
	return m.limits, m.err
}

func (m *mockQueryTierLimitRepo) FindBySaleAndTier(_ context.Context, _ int64, _ int) (*domain.TierLimit, error) {
	return m.limit, m.err
}

func (m *mockQueryTierLimitRepo) Save(_ context.Context, _ *domain.TierLimit) error { return nil }

// now 用于测试中的时间戳。
var now = time.Now()

// --- 用户查询测试（Task 8.1）---

func TestUserQueryService_ListPublicSales(t *testing.T) {
	saleRepo := &mockQuerySaleRepo{
		sales: []*domain.Sale{
			domain.ReconstructSale(1, "0x1", domain.SaleDeployed, 1, "", "", "", "", "", 0, 0, 0, false, 0, now, now, nil),
		},
		total: 1,
	}
	svc := newQuerySvc(saleRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	result, err := svc.ListPublicSales(context.Background(), 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)
	assert.Len(t, result.Sales, 1)
	assert.Equal(t, int64(1), result.Sales[0].Sale.ID)
}

func TestUserQueryService_GetSaleDetail(t *testing.T) {
	sale := domain.ReconstructSale(1, "0x1", domain.SaleDeployed, 1, "", "", "", "", "", 0, 0, 0, false, 0, now, now, nil)
	saleRepo := &mockQuerySaleRepo{sale: sale}
	metaRepo := newMockSaleMetaRepo()
	_ = metaRepo.Create(context.Background(), domain.ReconstructSaleMeta(1, 1, "标题", "", "", "", "", "", "public", 0, now, now))
	poolRepo := &mockQueryPoolRepo{pools: []*domain.Pool{}}

	svc := newQuerySvc(saleRepo, poolRepo, metaRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	detail, err := svc.GetSaleDetail(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, int64(1), detail.Sale.ID)
	assert.NotNil(t, detail.Meta)
}

func TestUserQueryService_GetSaleDetail_NotFound(t *testing.T) {
	saleRepo := &mockQuerySaleRepo{err: domain.ErrNotFound}
	svc := newQuerySvc(saleRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	_, err := svc.GetSaleDetail(context.Background(), 999)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

// --- 用户 Tier 查询测试（Task 8.3）---

func TestUserQueryService_GetUserTier(t *testing.T) {
	chain := &mockChainReader{tier: 3}
	svc := newQuerySvc(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, chain)

	tier, err := svc.GetUserTier(context.Background(), "0xTier", "0xUser")
	require.NoError(t, err)
	assert.Equal(t, 3, tier)
}

func TestUserQueryService_GetUserTier_Error(t *testing.T) {
	chain := &mockChainReader{err: assert.AnError}
	svc := newQuerySvc(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, chain)

	_, err := svc.GetUserTier(context.Background(), "0xTier", "0xUser")
	assert.Error(t, err)
}

// --- 白名单状态查询测试（Task 8.5）---

func TestUserQueryService_CheckWhitelist_Yes(t *testing.T) {
	wlRepo := &mockQueryWhitelistRepo{result: true}
	svc := newQuerySvc(nil, nil, nil, nil, wlRepo, nil, nil, nil, nil, nil, nil, nil)

	ok, err := svc.CheckWhitelist(context.Background(), 1, "0xUser")
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestUserQueryService_CheckWhitelist_No(t *testing.T) {
	wlRepo := &mockQueryWhitelistRepo{result: false}
	svc := newQuerySvc(nil, nil, nil, nil, wlRepo, nil, nil, nil, nil, nil, nil, nil)

	ok, err := svc.CheckWhitelist(context.Background(), 1, "0xUser")
	require.NoError(t, err)
	assert.False(t, ok)
}

// --- 用户申购记录查询测试（Task 8.7）---

func TestUserQueryService_GetUserDeposits(t *testing.T) {
	depositRepo := &mockQueryDepositRepo{
		deposits: []*domain.Deposit{
			domain.ReconstructDeposit(1, 1, 0, "0xUser", "0xTx", big.NewInt(1000), 100, now),
		},
	}
	svc := newQuerySvc(nil, nil, nil, nil, nil, depositRepo, nil, nil, nil, nil, nil, nil)

	deposits, err := svc.GetUserDeposits(context.Background(), "0xUser", 1)
	require.NoError(t, err)
	assert.Len(t, deposits, 1)
}

// --- 用户结算信息查询测试（Task 8.9）---

func TestUserQueryService_GetUserHarvest(t *testing.T) {
	harvestRepo := &mockQueryHarvestRepo{
		harvests: []*domain.Harvest{
			domain.ReconstructHarvest(1, 1, 0, "0xUser", false,
				big.NewInt(500), big.NewInt(500), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0),
				"0xTx", 200, now),
		},
	}
	svc := newQuerySvc(nil, nil, nil, nil, nil, nil, nil, harvestRepo, nil, nil, nil, nil)

	harvests, err := svc.GetUserHarvest(context.Background(), "0xUser", 1)
	require.NoError(t, err)
	assert.Len(t, harvests, 1)
}

// --- Vesting 可释放量测试（Task 8.11）---

func TestCalculateReleasable(t *testing.T) {
	schedule := domain.ReconstructVestingSchedule(
		1, 1, 0, 1, "0xBeneficiary",
		big.NewInt(1000), big.NewInt(300),
		now, now,
	)

	releasable := CalculateReleasable(schedule)
	assert.Equal(t, big.NewInt(700), releasable)
}

func TestCalculateReleasable_FullyReleased(t *testing.T) {
	schedule := domain.ReconstructVestingSchedule(
		1, 1, 0, 1, "0xBeneficiary",
		big.NewInt(1000), big.NewInt(1000),
		now, now,
	)

	releasable := CalculateReleasable(schedule)
	assert.Zero(t, releasable.Sign())
}

// --- 预估配售计算测试（Task 8.13）---

func makeTestPool(saleID int64, raising, offering *big.Int) *domain.Pool {
	return domain.ReconstructPool(saleID, saleID, 0, raising, offering, nil, nil, nil, nil, false, false, 0, 0, 0, 0, now, now)
}

func TestUserQueryService_EstimateAllocation(t *testing.T) {
	pool := makeTestPool(1, big.NewInt(1000), big.NewInt(5000))
	poolRepo := &mockQueryPoolRepo{pool: pool}
	svc := newQuerySvc(nil, poolRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	result, err := svc.EstimateAllocation(context.Background(), EstimateAllocationInput{
		SaleID:    1,
		PoolIndex: 0,
		UserTotal: "100",
	})
	require.NoError(t, err)
	// 用户投入 100 * (发售量 5000 / 募资量 1000) = 500
	assert.Equal(t, big.NewInt(500), result.Allocation)
}

func TestUserQueryService_EstimateAllocation_PoolNotConfigured(t *testing.T) {
	pool := makeTestPool(1, nil, nil)
	poolRepo := &mockQueryPoolRepo{pool: pool}
	svc := newQuerySvc(nil, poolRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	_, err := svc.EstimateAllocation(context.Background(), EstimateAllocationInput{
		SaleID:    1,
		PoolIndex: 0,
		UserTotal: "100",
	})
	assert.Error(t, err)
}

// --- 辅助函数 ---

// newQuerySvc 用 mock 仓库创建 UserQueryService。
func newQuerySvc(
	saleRepo domain.SaleRepository,
	poolRepo domain.PoolRepository,
	metaRepo domain.SaleMetaRepository,
	tierLimitRepo domain.TierLimitRepository,
	whitelistRepo domain.WhitelistRepository,
	depositRepo domain.DepositRepository,
	userPoolRepo domain.UserPoolStateRepository,
	harvestRepo domain.HarvestRepository,
	vestingRepo domain.VestingScheduleRepository,
	releaseRepo domain.VestingReleaseRepository,
	creditRepo domain.UserCreditRepository,
	chain domain.ChainReader,
) *UserQueryService {
	return NewUserQueryService(
		saleRepo, poolRepo, metaRepo, tierLimitRepo, whitelistRepo,
		depositRepo, userPoolRepo, creditRepo, harvestRepo, vestingRepo,
		releaseRepo, chain,
	)
}

// 确保接口实现
var _ domain.SaleRepository = (*mockQuerySaleRepo)(nil)
var _ domain.PoolRepository = (*mockQueryPoolRepo)(nil)
var _ domain.WhitelistRepository = (*mockQueryWhitelistRepo)(nil)
var _ domain.DepositRepository = (*mockQueryDepositRepo)(nil)
var _ domain.HarvestRepository = (*mockQueryHarvestRepo)(nil)
var _ domain.VestingScheduleRepository = (*mockQueryVestingRepo)(nil)
var _ domain.VestingReleaseRepository = (*mockQueryReleaseRepo)(nil)
var _ domain.UserPoolStateRepository = (*mockQueryUserPoolRepo)(nil)
var _ domain.UserCreditRepository = (*mockQueryCreditRepo)(nil)
var _ domain.TierLimitRepository = (*mockQueryTierLimitRepo)(nil)

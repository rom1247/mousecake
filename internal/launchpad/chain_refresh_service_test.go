package launchpad

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// --- Mock 类型 ---

// mockChainRefreshReader mock chainRefreshReader 接口。
type mockChainRefreshReader struct {
	batchCallResult []CallResult
	batchCallErr    error
	poolInfo        *PoolInfo
	poolInfoErr     error
	mouseTierAddr   common.Address
	mouseTierErr    error
	userPoolInfo    *UserPoolInfo
	userPoolInfoErr error

	batchTierLimitsResult map[int]*big.Int
	batchTierLimitsErr    error
	batchUserPoolResult   map[int]*UserPoolInfo
	batchUserPoolErr      error
	batchVestingResult    map[int]*VestingScheduleInfo
	batchVestingErr       error
	batchReleasableResult map[int]*big.Int
	batchReleasableErr    error
}

func (m *mockChainRefreshReader) BatchCall(_ context.Context, _ common.Address, _ []Call) ([]CallResult, error) {
	return m.batchCallResult, m.batchCallErr
}

func (m *mockChainRefreshReader) GetPoolInfo(_ context.Context, _ common.Address, _ *big.Int) (*PoolInfo, error) {
	return m.poolInfo, m.poolInfoErr
}

func (m *mockChainRefreshReader) GetMouseTier(_ context.Context, _ common.Address) (common.Address, error) {
	return m.mouseTierAddr, m.mouseTierErr
}

func (m *mockChainRefreshReader) GetUserPoolInfo(_ context.Context, _ common.Address, _ common.Address, _ *big.Int) (*UserPoolInfo, error) {
	return m.userPoolInfo, m.userPoolInfoErr
}

func (m *mockChainRefreshReader) BatchGetTierLimits(_ context.Context, _ common.Address, _ []*big.Int) (map[int]*big.Int, error) {
	return m.batchTierLimitsResult, m.batchTierLimitsErr
}

func (m *mockChainRefreshReader) BatchGetUserPoolInfos(_ context.Context, _ common.Address, _ []common.Address, _ *big.Int) (map[int]*UserPoolInfo, error) {
	return m.batchUserPoolResult, m.batchUserPoolErr
}

func (m *mockChainRefreshReader) BatchGetVestingSchedules(_ context.Context, _ common.Address, _ []*big.Int) (map[int]*VestingScheduleInfo, error) {
	return m.batchVestingResult, m.batchVestingErr
}

func (m *mockChainRefreshReader) BatchGetReleasableAmounts(_ context.Context, _ common.Address, _ []*big.Int) (map[int]*big.Int, error) {
	return m.batchReleasableResult, m.batchReleasableErr
}

// mockChainRefreshWriter mock chainRefreshWriter 接口。
type mockChainRefreshWriter struct {
	saleConfigErr      error
	poolConfigErr      error
	tierParamsErr      error
	tierLimitsErr      error
	userPoolStateErr   error
	vestingScheduleErr error

	addresses      []string
	addressesErr   error
	scheduleIDs    []int64
	scheduleIDsErr error
}

func (m *mockChainRefreshWriter) UpdateSaleConfig(_ context.Context, _ int64, _ map[string]any) error {
	return m.saleConfigErr
}

func (m *mockChainRefreshWriter) UpdatePoolConfig(_ context.Context, _ int64, _ int, _ map[string]any) error {
	return m.poolConfigErr
}

func (m *mockChainRefreshWriter) UpdateTierParams(_ context.Context, _ int, _ map[string]any) error {
	return m.tierParamsErr
}

func (m *mockChainRefreshWriter) UpdateTierLimits(_ context.Context, _ int64, _ map[int]string) error {
	return m.tierLimitsErr
}

func (m *mockChainRefreshWriter) UpdateUserPoolState(_ context.Context, _ int64, _ int, _ string, _ map[string]any) error {
	return m.userPoolStateErr
}

func (m *mockChainRefreshWriter) ListUserAddressesByPool(_ context.Context, _ int64, _ int) ([]string, error) {
	return m.addresses, m.addressesErr
}

func (m *mockChainRefreshWriter) UpdateVestingSchedule(_ context.Context, _ int64, _ map[string]any) error {
	return m.vestingScheduleErr
}

func (m *mockChainRefreshWriter) ListVestingScheduleIDsByUser(_ context.Context, _ string) ([]int64, error) {
	return m.scheduleIDs, m.scheduleIDsErr
}

// mockSaleLookup mock saleLookup 接口。
type mockSaleLookup struct {
	sale *domain.Sale
	err  error
}

func (m *mockSaleLookup) FindByContractAddress(_ context.Context, _ string) (*domain.Sale, error) {
	return m.sale, m.err
}

// 确保 mock 类型满足接口。
var _ chainRefreshReader = (*mockChainRefreshReader)(nil)
var _ chainRefreshWriter = (*mockChainRefreshWriter)(nil)
var _ saleLookup = (*mockSaleLookup)(nil)

// --- 测试辅助 ---

// testSale 创建测试用 Sale。
func testSale() *domain.Sale {
	return domain.ReconstructSale(1, "0x1234567890123456789012345678901234567890", domain.SaleDeployed, 1, "0xDeployer",
		"0xOwner", "0xRaise", "0xOffer", "0xTier", 1000, 2000, 0, false, 10,
		time.Now(), time.Now(), nil)
}

// --- 参数校验测试 ---

// TestChainStateRefresh_ScopesEmpty 验证 scopes 为空时返回错误。
func TestChainStateRefresh_ScopesEmpty(t *testing.T) {
	t.Parallel()

	svc := NewChainRefreshService(&mockChainRefreshReader{}, &mockChainRefreshWriter{}, &mockSaleLookup{})
	_, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "scopes 不能为空")
}

// TestChainStateRefresh_InvalidScope 验证无效 scope 值返回错误。
func TestChainStateRefresh_InvalidScope(t *testing.T) {
	t.Parallel()

	svc := NewChainRefreshService(&mockChainRefreshReader{}, &mockChainRefreshWriter{}, &mockSaleLookup{})
	_, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"sale", "invalid_scope"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "scope 值无效")
}

// TestChainStateRefresh_SaleNotFound 验证 sale_address 不存在时返回错误。
func TestChainStateRefresh_SaleNotFound(t *testing.T) {
	t.Parallel()

	svc := NewChainRefreshService(
		&mockChainRefreshReader{},
		&mockChainRefreshWriter{},
		&mockSaleLookup{err: fmt.Errorf("not found")},
	)
	_, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"sale"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "sale 记录不存在")
}

// TestChainStateRefresh_PoolMissingIndex 验证 scope=pool 缺少 pool_index 时返回错误。
func TestChainStateRefresh_PoolMissingIndex(t *testing.T) {
	t.Parallel()

	svc := NewChainRefreshService(
		&mockChainRefreshReader{},
		&mockChainRefreshWriter{},
		&mockSaleLookup{sale: testSale()},
	)
	_, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"pool"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pool_index 必填")
}

// --- 单 scope 刷新测试 ---

// TestChainStateRefresh_SingleScopeSale 验证 scope=sale 刷新成功。
func TestChainStateRefresh_SingleScopeSale(t *testing.T) {
	t.Parallel()

	sale := testSale()
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			batchCallResult: []CallResult{
				{Method: "startBlock", Data: big.NewInt(1000).Bytes()},
				{Method: "endBlock", Data: big.NewInt(2000).Bytes()},
				{Method: "vestingRevoked", Data: []byte{0}},
				{Method: "vestingStartTime", Data: big.NewInt(0).Bytes()},
				{Method: "owner", Data: common.HexToAddress("0xOwner").Bytes()},
				{Method: "raiseToken", Data: common.HexToAddress("0xRaise").Bytes()},
				{Method: "offeringToken", Data: common.HexToAddress("0xOffer").Bytes()},
				{Method: "mouseTier", Data: common.HexToAddress("0xTier").Bytes()},
				{Method: "nextScheduleId", Data: big.NewInt(5).Bytes()},
			},
		},
		&mockChainRefreshWriter{},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"sale"},
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), resp.SaleID)
	assert.NotNil(t, resp.Results["sale"])
	assert.Contains(t, resp.Results["sale"].UpdatedFields, "start_block")
	assert.Contains(t, resp.Results["sale"].UpdatedFields, "end_block")
	assert.Empty(t, resp.Results["sale"].Error)
}

// --- 多 scope 部分成功测试 ---

// TestChainStateRefresh_PartialSuccess 验证多个 scope 部分失败时不影响其他 scope。
func TestChainStateRefresh_PartialSuccess(t *testing.T) {
	t.Parallel()

	sale := testSale()
	poolIndex := 0
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			// BatchCall 用于 sale scope，返回成功数据
			batchCallResult: []CallResult{
				{Method: "startBlock", Data: big.NewInt(1000).Bytes()},
				{Method: "endBlock", Data: big.NewInt(2000).Bytes()},
				{Method: "vestingRevoked", Data: []byte{0}},
				{Method: "vestingStartTime", Data: big.NewInt(0).Bytes()},
				{Method: "owner", Data: common.HexToAddress("0xOwner").Bytes()},
				{Method: "raiseToken", Data: common.HexToAddress("0xRaise").Bytes()},
				{Method: "offeringToken", Data: common.HexToAddress("0xOffer").Bytes()},
				{Method: "mouseTier", Data: common.HexToAddress("0xTier").Bytes()},
				{Method: "nextScheduleId", Data: big.NewInt(5).Bytes()},
			},
			// poolInfo 用于 pool scope，返回错误
			poolInfoErr: fmt.Errorf("eth_call timeout"),
		},
		&mockChainRefreshWriter{},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"sale", "pool"},
		PoolIndex:   &poolIndex,
	})
	require.NoError(t, err)

	// sale scope 成功
	assert.NotNil(t, resp.Results["sale"])
	assert.Empty(t, resp.Results["sale"].Error)

	// pool scope 失败但不影响整体
	assert.NotNil(t, resp.Results["pool"])
	assert.Contains(t, resp.Results["pool"].Error, "读取 poolInfo 失败")
}

// --- tier_params scope 测试 ---

func TestChainStateRefresh_TierParamsSuccess(t *testing.T) {
	t.Parallel()

	sale := testSale()
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			mouseTierAddr: common.HexToAddress("0xTierContract"),
			batchCallResult: []CallResult{
				{Method: "ceiling", Data: big.NewInt(100).Bytes()},
				{Method: "multiplier", Data: big.NewInt(200).Bytes()},
				{Method: "tierBaseAmount", Data: big.NewInt(300).Bytes()},
			},
		},
		&mockChainRefreshWriter{},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"tier_params"},
	})
	require.NoError(t, err)
	assert.Empty(t, resp.Results["tier_params"].Error)
	assert.Contains(t, resp.Results["tier_params"].UpdatedFields, "ceiling")
	assert.Contains(t, resp.Results["tier_params"].UpdatedFields, "multiplier")
	assert.Contains(t, resp.Results["tier_params"].UpdatedFields, "tier_base_amount")
}

func TestChainStateRefresh_TierParamsMouseTierFail(t *testing.T) {
	t.Parallel()

	sale := testSale()
	svc := NewChainRefreshService(
		&mockChainRefreshReader{mouseTierErr: fmt.Errorf("rpc timeout")},
		&mockChainRefreshWriter{},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"tier_params"},
	})
	require.NoError(t, err)
	assert.Contains(t, resp.Results["tier_params"].Error, "读取 mouseTier 地址失败")
}

func TestChainStateRefresh_TierParamsDBFail(t *testing.T) {
	t.Parallel()

	sale := testSale()
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			mouseTierAddr: common.HexToAddress("0xTierContract"),
			batchCallResult: []CallResult{
				{Method: "ceiling", Data: big.NewInt(100).Bytes()},
				{Method: "multiplier", Data: big.NewInt(200).Bytes()},
				{Method: "tierBaseAmount", Data: big.NewInt(300).Bytes()},
			},
		},
		&mockChainRefreshWriter{tierParamsErr: fmt.Errorf("db connection lost")},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"tier_params"},
	})
	require.NoError(t, err)
	assert.Contains(t, resp.Results["tier_params"].Error, "更新 DB 失败")
}

// --- tier_limits scope 测试 ---

func TestChainStateRefresh_TierLimitsSpecificTier(t *testing.T) {
	t.Parallel()

	sale := testSale()
	tier := 3
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			batchTierLimitsResult: map[int]*big.Int{0: big.NewInt(500)},
		},
		&mockChainRefreshWriter{},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"tier_limits"},
		Tier:        &tier,
	})
	require.NoError(t, err)
	assert.Empty(t, resp.Results["tier_limits"].Error)
	assert.Contains(t, resp.Results["tier_limits"].UpdatedFields, "tier_3")
}

func TestChainStateRefresh_TierLimitsFullScan(t *testing.T) {
	t.Parallel()

	sale := testSale()
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			batchTierLimitsResult: map[int]*big.Int{
				1: big.NewInt(100),
				2: big.NewInt(200),
				3: big.NewInt(0), // 零值应被跳过
			},
		},
		&mockChainRefreshWriter{},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"tier_limits"},
	})
	require.NoError(t, err)
	assert.Empty(t, resp.Results["tier_limits"].Error)
	assert.Contains(t, resp.Results["tier_limits"].UpdatedFields, "tier_1")
	assert.Contains(t, resp.Results["tier_limits"].UpdatedFields, "tier_2")
	assert.NotContains(t, resp.Results["tier_limits"].UpdatedFields, "tier_3")
}

func TestChainStateRefresh_TierLimitsBatchFail(t *testing.T) {
	t.Parallel()

	sale := testSale()
	svc := NewChainRefreshService(
		&mockChainRefreshReader{batchTierLimitsErr: fmt.Errorf("rpc error")},
		&mockChainRefreshWriter{},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"tier_limits"},
	})
	require.NoError(t, err)
	assert.Contains(t, resp.Results["tier_limits"].Error, "批量读取 Tier 额度失败")
}

func TestChainStateRefresh_TierLimitsDBFail(t *testing.T) {
	t.Parallel()

	sale := testSale()
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			batchTierLimitsResult: map[int]*big.Int{0: big.NewInt(100)},
		},
		&mockChainRefreshWriter{tierLimitsErr: fmt.Errorf("db error")},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"tier_limits"},
	})
	require.NoError(t, err)
	assert.Contains(t, resp.Results["tier_limits"].Error, "更新 DB 失败")
}

// --- user_pool scope 测试 ---

func TestChainStateRefresh_UserPoolSingleUser(t *testing.T) {
	t.Parallel()

	sale := testSale()
	poolIndex := 0
	userAddr := "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			userPoolInfo: &UserPoolInfo{
				AmountPool:  big.NewInt(1000),
				ClaimedPool: false,
			},
		},
		&mockChainRefreshWriter{},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"user_pool"},
		PoolIndex:   &poolIndex,
		UserAddress: &userAddr,
	})
	require.NoError(t, err)
	assert.Empty(t, resp.Results["user_pool"].Error)
	assert.Contains(t, resp.Results["user_pool"].UpdatedFields, "total_deposited")
	assert.Contains(t, resp.Results["user_pool"].UpdatedFields, "claimed")
}

func TestChainStateRefresh_UserPoolBatch(t *testing.T) {
	t.Parallel()

	sale := testSale()
	poolIndex := 0
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			batchUserPoolResult: map[int]*UserPoolInfo{
				0: {AmountPool: big.NewInt(100), ClaimedPool: false},
				1: {AmountPool: big.NewInt(200), ClaimedPool: true},
			},
		},
		&mockChainRefreshWriter{
			addresses: []string{"0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B", "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"},
		},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"user_pool"},
		PoolIndex:   &poolIndex,
	})
	require.NoError(t, err)
	assert.Empty(t, resp.Results["user_pool"].Error)
	assert.Len(t, resp.Results["user_pool"].UpdatedFields, 2)
}

func TestChainStateRefresh_UserPoolBatchPartialFail(t *testing.T) {
	t.Parallel()

	sale := testSale()
	poolIndex := 0
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			batchUserPoolResult: map[int]*UserPoolInfo{
				0: {AmountPool: big.NewInt(100), ClaimedPool: false},
				// index 1 缺失，模拟链上调用失败
			},
		},
		&mockChainRefreshWriter{
			addresses: []string{"0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B", "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"},
		},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"user_pool"},
		PoolIndex:   &poolIndex,
	})
	require.NoError(t, err)
	assert.Empty(t, resp.Results["user_pool"].Error)
	assert.Len(t, resp.Results["user_pool"].UpdatedFields, 1)
}

func TestChainStateRefresh_UserPoolBatchRPCFail(t *testing.T) {
	t.Parallel()

	sale := testSale()
	poolIndex := 0
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			batchUserPoolErr: fmt.Errorf("batch rpc timeout"),
		},
		&mockChainRefreshWriter{
			addresses: []string{"0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"},
		},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"user_pool"},
		PoolIndex:   &poolIndex,
	})
	require.NoError(t, err)
	assert.Contains(t, resp.Results["user_pool"].Error, "批量读取用户 Pool 状态失败")
}

// --- vesting scope 测试 ---

func TestChainStateRefresh_VestingByScheduleIDs(t *testing.T) {
	t.Parallel()

	sale := testSale()
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			batchVestingResult: map[int]*VestingScheduleInfo{
				0: {
					Beneficiary: common.HexToAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"),
					Pid:         big.NewInt(0),
					AmountTotal: big.NewInt(1000),
					Released:    big.NewInt(200),
				},
				1: {
					Beneficiary: common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"),
					Pid:         big.NewInt(1),
					AmountTotal: big.NewInt(2000),
					Released:    big.NewInt(500),
				},
			},
			batchReleasableResult: map[int]*big.Int{
				0: big.NewInt(300),
				1: big.NewInt(600),
			},
		},
		&mockChainRefreshWriter{},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"vesting"},
		ScheduleIDs: []int64{10, 20},
	})
	require.NoError(t, err)
	assert.Empty(t, resp.Results["vesting"].Error)
	assert.Len(t, resp.Results["vesting"].UpdatedFields, 2)
	assert.Contains(t, resp.Results["vesting"].UpdatedFields, "schedule_10")
	assert.Contains(t, resp.Results["vesting"].UpdatedFields, "schedule_20")
}

func TestChainStateRefresh_VestingByUser(t *testing.T) {
	t.Parallel()

	sale := testSale()
	userAddr := "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			batchVestingResult: map[int]*VestingScheduleInfo{
				0: {
					Beneficiary: common.HexToAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"),
					Pid:         big.NewInt(0),
					AmountTotal: big.NewInt(1000),
					Released:    big.NewInt(100),
				},
			},
			batchReleasableResult: map[int]*big.Int{0: big.NewInt(50)},
		},
		&mockChainRefreshWriter{
			scheduleIDs: []int64{5},
		},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"vesting"},
		UserAddress: &userAddr,
	})
	require.NoError(t, err)
	assert.Empty(t, resp.Results["vesting"].Error)
	assert.Contains(t, resp.Results["vesting"].UpdatedFields, "schedule_5")
}

func TestChainStateRefresh_VestingScheduleNotFound(t *testing.T) {
	t.Parallel()

	sale := testSale()
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			// index 0 缺失，模拟 schedule 不存在
			batchReleasableResult: map[int]*big.Int{},
		},
		&mockChainRefreshWriter{},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"vesting"},
		ScheduleIDs: []int64{99},
	})
	require.NoError(t, err)
	assert.Empty(t, resp.Results["vesting"].Error)
	assert.Empty(t, resp.Results["vesting"].UpdatedFields)
}

func TestChainStateRefresh_VestingBatchFail(t *testing.T) {
	t.Parallel()

	sale := testSale()
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			batchVestingErr: fmt.Errorf("batch rpc error"),
		},
		&mockChainRefreshWriter{},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"vesting"},
		ScheduleIDs: []int64{1},
	})
	require.NoError(t, err)
	assert.Contains(t, resp.Results["vesting"].Error, "批量读取 vesting schedule 失败")
}

func TestChainStateRefresh_VestingReleasableFail(t *testing.T) {
	t.Parallel()

	sale := testSale()
	svc := NewChainRefreshService(
		&mockChainRefreshReader{
			batchVestingResult: map[int]*VestingScheduleInfo{
				0: {
					Beneficiary: common.HexToAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"),
					Pid:         big.NewInt(0),
					AmountTotal: big.NewInt(1000),
					Released:    big.NewInt(0),
				},
			},
			batchReleasableErr: fmt.Errorf("rpc timeout"),
		},
		&mockChainRefreshWriter{},
		&mockSaleLookup{sale: sale},
	)

	resp, err := svc.ChainStateRefresh(context.Background(), ChainStateRefreshRequest{
		SaleAddress: "0x1234567890123456789012345678901234567890",
		Scopes:      []string{"vesting"},
		ScheduleIDs: []int64{1},
	})
	require.NoError(t, err)
	// releasable 失败不阻塞，schedule 仍应更新成功
	assert.Empty(t, resp.Results["vesting"].Error)
	assert.Contains(t, resp.Results["vesting"].UpdatedFields, "schedule_1")
}

package launchpad

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// ChainStateRefreshRequest 链上状态刷新请求。
type ChainStateRefreshRequest struct {
	// SaleAddress 销售合约地址。
	SaleAddress string `json:"sale_address" binding:"required"`
	// Scopes 刷新范围列表，可选值：sale/pool/tier_params/tier_limits/user_pool/vesting。
	Scopes []string `json:"scopes" binding:"required"`
	// PoolIndex 可选，scope=pool 和 scope=user_pool 时必填。
	PoolIndex *int `json:"pool_index"`
	// UserAddress 可选，scope=user_pool 时指定单用户。
	UserAddress *string `json:"user_address"`
	// ScheduleIDs 可选，scope=vesting 时指定 schedule IDs。
	ScheduleIDs []int64 `json:"schedule_ids"`
	// Tier 可选，scope=tier_limits 时指定单个 tier。
	Tier *int `json:"tier"`
}

// ScopeResult 单个 scope 的刷新结果。
type ScopeResult struct {
	// UpdatedFields 更新的字段列表。
	UpdatedFields []string `json:"updated_fields"`
	// Error 刷新错误信息。
	Error string `json:"error,omitempty"`
}

// ChainStateRefreshResponse 链上状态刷新响应。
type ChainStateRefreshResponse struct {
	// SaleID 销售记录 ID。
	SaleID int64 `json:"sale_id"`
	// Results 各 scope 的刷新结果，key 为 scope 名。
	Results map[string]*ScopeResult `json:"results"`
}

// validScopes 合法的 scope 值集合。
var validScopes = map[string]bool{
	"sale":        true,
	"pool":        true,
	"tier_params": true,
	"tier_limits": true,
	"user_pool":   true,
	"vesting":     true,
}

// chainRefreshReader 定义 ChainRefreshService 需要的链上读取能力。
type chainRefreshReader interface {
	BatchCall(ctx context.Context, saleAddr common.Address, calls []Call) ([]CallResult, error)
	GetPoolInfo(ctx context.Context, saleAddr common.Address, pid *big.Int) (*PoolInfo, error)
	GetMouseTier(ctx context.Context, saleAddr common.Address) (common.Address, error)
	GetUserPoolInfo(ctx context.Context, saleAddr common.Address, user common.Address, pid *big.Int) (*UserPoolInfo, error)
	BatchGetTierLimits(ctx context.Context, saleAddr common.Address, tiers []*big.Int) (map[int]*big.Int, error)
	BatchGetUserPoolInfos(ctx context.Context, saleAddr common.Address, userAddresses []common.Address, pid *big.Int) (map[int]*UserPoolInfo, error)
	BatchGetVestingSchedules(ctx context.Context, saleAddr common.Address, scheduleIDs []*big.Int) (map[int]*VestingScheduleInfo, error)
	BatchGetReleasableAmounts(ctx context.Context, saleAddr common.Address, scheduleIDs []*big.Int) (map[int]*big.Int, error)
}

// chainRefreshWriter 定义 ChainRefreshService 需要的数据写入能力。
type chainRefreshWriter interface {
	UpdateSaleConfig(ctx context.Context, saleID int64, fields map[string]any) error
	UpdatePoolConfig(ctx context.Context, saleID int64, poolIndex int, fields map[string]any) error
	UpdateTierParams(ctx context.Context, chainID int, fields map[string]any) error
	UpdateTierLimits(ctx context.Context, saleID int64, limits map[int]string) error
	UpdateUserPoolState(ctx context.Context, saleID int64, poolIndex int, userAddress string, fields map[string]any) error
	ListUserAddressesByPool(ctx context.Context, saleID int64, poolIndex int) ([]string, error)
	UpdateVestingSchedule(ctx context.Context, scheduleID int64, fields map[string]any) error
	ListVestingScheduleIDsByUser(ctx context.Context, userAddress string) ([]int64, error)
}

// saleLookup 定义 ChainRefreshService 需要的销售查询能力。
type saleLookup interface {
	FindByContractAddress(ctx context.Context, address string) (*domain.Sale, error)
}

// ChainRefreshService 链上状态刷新用例，按 scope 编排 ChainReader 和 Repository。
type ChainRefreshService struct {
	reader chainRefreshReader
	repo   chainRefreshWriter
	saleDB saleLookup
}

// NewChainRefreshService 创建 ChainRefreshService 实例。
func NewChainRefreshService(reader chainRefreshReader, repo chainRefreshWriter, saleDB saleLookup) *ChainRefreshService {
	return &ChainRefreshService{
		reader: reader,
		repo:   repo,
		saleDB: saleDB,
	}
}

// ChainStateRefresh 执行链上状态刷新。
func (s *ChainRefreshService) ChainStateRefresh(ctx context.Context, req ChainStateRefreshRequest) (*ChainStateRefreshResponse, error) {
	// 参数校验
	if len(req.Scopes) == 0 {
		return nil, fmt.Errorf("scopes 不能为空")
	}
	for _, scope := range req.Scopes {
		if !validScopes[scope] {
			return nil, fmt.Errorf("scope 值无效: %s", scope)
		}
	}

	sale, err := s.saleDB.FindByContractAddress(ctx, req.SaleAddress)
	if err != nil {
		return nil, fmt.Errorf("sale 记录不存在: %w", err)
	}

	// scope=pool 必须提供 pool_index
	for _, scope := range req.Scopes {
		if scope == "pool" && req.PoolIndex == nil {
			return nil, fmt.Errorf("scope=pool 时 pool_index 必填")
		}
	}

	saleAddr := common.HexToAddress(sale.ContractAddress)
	resp := &ChainStateRefreshResponse{
		SaleID:  sale.ID,
		Results: make(map[string]*ScopeResult),
	}

	// 逐 scope 执行刷新
	for _, scope := range req.Scopes {
		result := &ScopeResult{}
		switch scope {
		case "sale":
			result = s.refreshSale(ctx, saleAddr, sale)
		case "pool":
			result = s.refreshPool(ctx, saleAddr, sale.ID, *req.PoolIndex)
		case "tier_params":
			result = s.refreshTierParams(ctx, saleAddr, sale)
		case "tier_limits":
			result = s.refreshTierLimits(ctx, saleAddr, sale.ID, req.Tier)
		case "user_pool":
			result = s.refreshUserPool(ctx, saleAddr, sale.ID, req.PoolIndex, req.UserAddress)
		case "vesting":
			result = s.refreshVesting(ctx, saleAddr, sale.ID, req.ScheduleIDs, req.UserAddress)
		}
		resp.Results[scope] = result
	}

	return resp, nil
}

// refreshSale 刷新 Sale 配置（9 个字段一次 BatchCall）。
func (s *ChainRefreshService) refreshSale(ctx context.Context, saleAddr common.Address, sale *domain.Sale) *ScopeResult {
	result := &ScopeResult{}

	calls := []Call{
		{Method: "startBlock"},
		{Method: "endBlock"},
		{Method: "vestingRevoked"},
		{Method: "vestingStartTime"},
		{Method: "owner"},
		{Method: "raiseToken"},
		{Method: "offeringToken"},
		{Method: "mouseTier"},
		{Method: "nextScheduleId"},
	}

	results, err := s.reader.BatchCall(ctx, saleAddr, calls)
	if err != nil {
		result.Error = fmt.Sprintf("BatchCall 失败: %v", err)
		return result
	}

	fields := make(map[string]any)
	if results[0].Err == nil {
		startBlock := new(big.Int).SetBytes(results[0].Data)
		fields["start_block"] = startBlock.Int64()
		result.UpdatedFields = append(result.UpdatedFields, "start_block")
	}
	if results[1].Err == nil {
		endBlock := new(big.Int).SetBytes(results[1].Data)
		fields["end_block"] = endBlock.Int64()
		result.UpdatedFields = append(result.UpdatedFields, "end_block")
	}
	if results[2].Err == nil {
		revoked := len(results[2].Data) > 0 && results[2].Data[len(results[2].Data)-1] == 1
		fields["vesting_revoked"] = revoked
		result.UpdatedFields = append(result.UpdatedFields, "vesting_revoked")
	}
	if results[3].Err == nil {
		vestingStartTime := new(big.Int).SetBytes(results[3].Data)
		fields["vesting_start_time"] = vestingStartTime.Int64()
		result.UpdatedFields = append(result.UpdatedFields, "vesting_start_time")
	}
	if results[4].Err == nil && len(results[4].Data) >= 32 {
		owner := common.BytesToAddress(results[4].Data[12:])
		fields["owner_address"] = owner.Hex()
		result.UpdatedFields = append(result.UpdatedFields, "owner_address")
	}
	if results[5].Err == nil && len(results[5].Data) >= 32 {
		raiseToken := common.BytesToAddress(results[5].Data[12:])
		fields["raise_token_address"] = raiseToken.Hex()
		result.UpdatedFields = append(result.UpdatedFields, "raise_token_address")
	}
	if results[6].Err == nil && len(results[6].Data) >= 32 {
		offeringToken := common.BytesToAddress(results[6].Data[12:])
		fields["offering_token_address"] = offeringToken.Hex()
		result.UpdatedFields = append(result.UpdatedFields, "offering_token_address")
	}
	if results[7].Err == nil && len(results[7].Data) >= 32 {
		mouseTier := common.BytesToAddress(results[7].Data[12:])
		fields["mouse_tier_address"] = mouseTier.Hex()
		result.UpdatedFields = append(result.UpdatedFields, "mouse_tier_address")
	}
	if results[8].Err == nil {
		nextScheduleId := new(big.Int).SetBytes(results[8].Data)
		_ = nextScheduleId // 不更新到 DB，仅用于参考
	}

	if len(fields) > 0 {
		if err := s.repo.UpdateSaleConfig(ctx, sale.ID, fields); err != nil {
			result.Error = fmt.Sprintf("更新 DB 失败: %v", err)
			return result
		}
	}

	return result
}

// refreshPool 刷新 Pool 配置（poolInfo 11 个字段）。
func (s *ChainRefreshService) refreshPool(ctx context.Context, saleAddr common.Address, saleID int64, poolIndex int) *ScopeResult {
	result := &ScopeResult{}

	poolInfo, err := s.reader.GetPoolInfo(ctx, saleAddr, big.NewInt(int64(poolIndex)))
	if err != nil {
		result.Error = fmt.Sprintf("读取 poolInfo 失败: %v", err)
		return result
	}

	fields := map[string]any{
		"offering_amount":      poolInfo.OfferingAmountPool.String(),
		"raising_amount":       poolInfo.RaisingAmountPool.String(),
		"total_amount":         poolInfo.TotalAmountPool.String(),
		"limit_per_user":       poolInfo.LimitPerUser.String(),
		"is_special_sale":      poolInfo.IsSpecialSale,
		"has_tax":              poolInfo.HasTax,
		"vesting_percentage":   poolInfo.VestingPercentage.Int64(),
		"vesting_cliff":        poolInfo.VestingCliff.Int64(),
		"vesting_slice_period": poolInfo.VestingSlicePeriod.Int64(),
		"vesting_duration":     poolInfo.VestingDuration.Int64(),
	}

	result.UpdatedFields = []string{
		"offering_amount", "raising_amount", "total_amount", "limit_per_user",
		"is_special_sale", "has_tax", "vesting_percentage", "vesting_cliff",
		"vesting_slice_period", "vesting_duration",
	}

	if err := s.repo.UpdatePoolConfig(ctx, saleID, poolIndex, fields); err != nil {
		result.Error = fmt.Sprintf("更新 DB 失败: %v", err)
		return result
	}

	return result
}

// refreshTierParams 刷新全局 Tier 参数（两步：先读 mouseTier 地址，再读参数）。
func (s *ChainRefreshService) refreshTierParams(ctx context.Context, saleAddr common.Address, sale *domain.Sale) *ScopeResult {
	result := &ScopeResult{}

	// 第一步：读取 mouseTier 地址
	mouseTierAddr, err := s.reader.GetMouseTier(ctx, saleAddr)
	if err != nil {
		result.Error = fmt.Sprintf("读取 mouseTier 地址失败: %v", err)
		return result
	}

	// 第二步：BatchCall 读取 ceiling/multiplier/tierBaseAmount
	calls := []Call{
		{Method: "ceiling", To: &mouseTierAddr},
		{Method: "multiplier", To: &mouseTierAddr},
		{Method: "tierBaseAmount", To: &mouseTierAddr},
	}

	results, err := s.reader.BatchCall(ctx, saleAddr, calls)
	if err != nil {
		result.Error = fmt.Sprintf("BatchCall 失败: %v", err)
		return result
	}

	fields := make(map[string]any)
	if results[0].Err == nil {
		ceiling := new(big.Int).SetBytes(results[0].Data)
		fields["ceiling"] = ceiling.String()
		result.UpdatedFields = append(result.UpdatedFields, "ceiling")
	}
	if results[1].Err == nil {
		multiplier := new(big.Int).SetBytes(results[1].Data)
		fields["multiplier"] = multiplier.String()
		result.UpdatedFields = append(result.UpdatedFields, "multiplier")
	}
	if results[2].Err == nil {
		tierBaseAmount := new(big.Int).SetBytes(results[2].Data)
		fields["tier_base_amount"] = tierBaseAmount.String()
		result.UpdatedFields = append(result.UpdatedFields, "tier_base_amount")
	}

	if len(fields) > 0 {
		if err := s.repo.UpdateTierParams(ctx, sale.ChainID, fields); err != nil {
			result.Error = fmt.Sprintf("更新 DB 失败: %v", err)
			return result
		}
	}

	return result
}

// refreshTierLimits 刷新 Tier 额度限制（支持全量遍历和指定 tier）。
func (s *ChainRefreshService) refreshTierLimits(ctx context.Context, saleAddr common.Address, saleID int64, tier *int) *ScopeResult {
	result := &ScopeResult{}

	var tierBigInts []*big.Int
	if tier != nil {
		tierBigInts = []*big.Int{big.NewInt(int64(*tier))}
	} else {
		for i := 0; i < 10; i++ {
			tierBigInts = append(tierBigInts, big.NewInt(int64(i)))
		}
	}

	limits, err := s.reader.BatchGetTierLimits(ctx, saleAddr, tierBigInts)
	if err != nil {
		result.Error = fmt.Sprintf("批量读取 Tier 额度失败: %v", err)
		return result
	}

	validLimits := make(map[int]string)
	for i, t := range tierBigInts {
		limit, ok := limits[i]
		if !ok {
			continue
		}
		if limit.Sign() == 0 {
			continue
		}
		tIdx := int(t.Int64())
		validLimits[tIdx] = limit.String()
		result.UpdatedFields = append(result.UpdatedFields, fmt.Sprintf("tier_%d", tIdx))
	}

	if len(validLimits) > 0 {
		if err := s.repo.UpdateTierLimits(ctx, saleID, validLimits); err != nil {
			result.Error = fmt.Sprintf("更新 DB 失败: %v", err)
			return result
		}
	}

	return result
}

// refreshUserPool 刷新用户 Pool 状态（支持单用户和批量）。
func (s *ChainRefreshService) refreshUserPool(ctx context.Context, saleAddr common.Address, saleID int64, poolIndex *int, userAddress *string) *ScopeResult {
	result := &ScopeResult{}

	pid := 0
	if poolIndex != nil {
		pid = *poolIndex
	}

	if userAddress != nil {
		// 单用户模式
		userAddr := common.HexToAddress(*userAddress)
		info, err := s.reader.GetUserPoolInfo(ctx, saleAddr, userAddr, big.NewInt(int64(pid)))
		if err != nil {
			result.Error = fmt.Sprintf("读取用户 Pool 状态失败: %v", err)
			return result
		}

		fields := map[string]any{
			"total_deposited": info.AmountPool.String(),
			"claimed":         info.ClaimedPool,
		}
		if err := s.repo.UpdateUserPoolState(ctx, saleID, pid, *userAddress, fields); err != nil {
			result.Error = fmt.Sprintf("更新 DB 失败: %v", err)
			return result
		}
		result.UpdatedFields = []string{"total_deposited", "claimed"}
		return result
	}

	// 批量模式：从 DB 查询所有用户
	addresses, err := s.repo.ListUserAddressesByPool(ctx, saleID, pid)
	if err != nil {
		result.Error = fmt.Sprintf("查询用户列表失败: %v", err)
		return result
	}

	if len(addresses) == 0 {
		return result
	}

	userAddrs := make([]common.Address, len(addresses))
	for i, addr := range addresses {
		userAddrs[i] = common.HexToAddress(addr)
	}

	infos, err := s.reader.BatchGetUserPoolInfos(ctx, saleAddr, userAddrs, big.NewInt(int64(pid)))
	if err != nil {
		result.Error = fmt.Sprintf("批量读取用户 Pool 状态失败: %v", err)
		return result
	}

	for i, addr := range addresses {
		info, ok := infos[i]
		if !ok {
			continue
		}
		fields := map[string]any{
			"total_deposited": info.AmountPool.String(),
			"claimed":         info.ClaimedPool,
		}
		if err := s.repo.UpdateUserPoolState(ctx, saleID, pid, addr, fields); err != nil {
			slog.WarnContext(ctx, "更新用户 Pool 状态失败", "addr", addr, "error", err)
			continue
		}
		result.UpdatedFields = append(result.UpdatedFields, addr)
	}

	return result
}

// refreshVesting 刷新锁仓计划（支持 schedule_ids 和 user_address 两种模式）。
func (s *ChainRefreshService) refreshVesting(ctx context.Context, saleAddr common.Address, saleID int64, scheduleIDs []int64, userAddress *string) *ScopeResult {
	result := &ScopeResult{}

	var ids []int64
	if len(scheduleIDs) > 0 {
		ids = scheduleIDs
	} else if userAddress != nil {
		dbIDs, err := s.repo.ListVestingScheduleIDsByUser(ctx, *userAddress)
		if err != nil {
			result.Error = fmt.Sprintf("查询用户 schedule IDs 失败: %v", err)
			return result
		}
		ids = dbIDs
	}

	if len(ids) == 0 {
		return result
	}

	scheduleBigInts := make([]*big.Int, len(ids))
	for i, id := range ids {
		scheduleBigInts[i] = big.NewInt(id)
	}

	schedules, err := s.reader.BatchGetVestingSchedules(ctx, saleAddr, scheduleBigInts)
	if err != nil {
		result.Error = fmt.Sprintf("批量读取 vesting schedule 失败: %v", err)
		return result
	}

	releasables, relErr := s.reader.BatchGetReleasableAmounts(ctx, saleAddr, scheduleBigInts)
	if relErr != nil {
		slog.WarnContext(ctx, "批量读取可释放量失败", "error", relErr)
	}

	for i, scheduleID := range ids {
		schedule, ok := schedules[i]
		if !ok {
			continue
		}
		fields := map[string]any{
			"beneficiary":  schedule.Beneficiary.Hex(),
			"pool_index":   schedule.Pid.Int64(),
			"amount_total": schedule.AmountTotal.String(),
			"released":     schedule.Released.String(),
		}
		if releasable, ok := releasables[i]; ok && releasable != nil {
			fields["releasable_amount"] = releasable.String()
		}
		if err := s.repo.UpdateVestingSchedule(ctx, scheduleID, fields); err != nil {
			slog.WarnContext(ctx, "更新 vesting schedule 失败", "schedule_id", scheduleID, "error", err)
			continue
		}
		result.UpdatedFields = append(result.UpdatedFields, fmt.Sprintf("schedule_%d", scheduleID))
	}

	return result
}

// Package launchpad 实现 Launchpad IDO 募资模块的业务逻辑和数据访问。
package launchpad

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/mousecake-go/mousecake-go/internal/chain"
	"github.com/mousecake-go/mousecake-go/internal/chain/contract/launchpad"
	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// PoolInfo 表示链上 poolInfo(pid) 返回的完整 Pool 配置。
type PoolInfo struct {
	// OfferingAmountPool 发售代币数量。
	OfferingAmountPool *big.Int
	// RaisingAmountPool 募资目标数量。
	RaisingAmountPool *big.Int
	// TotalAmountPool 实际募资总额。
	TotalAmountPool *big.Int
	// LimitPerUser 单用户申购上限。
	LimitPerUser *big.Int
	// IsSpecialSale 是否为特殊池。
	IsSpecialSale bool
	// HasTax 是否收取税费。
	HasTax bool
	// VestingPercentage TGE 释放比例。
	VestingPercentage *big.Int
	// VestingCliff 悬崖期（秒）。
	VestingCliff *big.Int
	// VestingSlicePeriod 线性释放间隔（秒）。
	VestingSlicePeriod *big.Int
	// VestingDuration 总锁仓时长（秒）。
	VestingDuration *big.Int
	// Configured 池子是否已配置。
	Configured bool
}

// VestingScheduleInfo 表示链上 vestingSchedules(id) 返回的锁仓计划。
type VestingScheduleInfo struct {
	// Beneficiary 受益人地址。
	Beneficiary common.Address
	// Pid 池子索引。
	Pid *big.Int
	// AmountTotal 锁仓总量。
	AmountTotal *big.Int
	// Released 已释放量。
	Released *big.Int
}

// UserPoolInfo 表示链上 viewUserPoolInfo 返回的用户 Pool 状态。
type UserPoolInfo struct {
	// AmountPool 用户申购金额。
	AmountPool *big.Int
	// ClaimedPool 是否已结算。
	ClaimedPool bool
}

// Call 表示一次链上合约调用请求。
type Call struct {
	// Method 合约方法名。
	Method string
	// Args 编码参数（传递给 Pack 方法）。
	Args []any
	// To 可选的目标合约地址，为 nil 时使用 saleAddr。
	To *common.Address
}

// CallResult 表示一次链上合约调用的结果。
type CallResult struct {
	// Method 合约方法名。
	Method string
	// Data 解码后的返回值（nil 表示调用失败）。
	Data []byte
	// Err 调用错误。
	Err error
}

// ChainReader 实现 domain.ChainReader 接口，通过 NodePool 执行只读链上调用。
type ChainReader struct {
	pool    chain.NodePool
	tierAbi *launchpad.MouseTier
	padAbi  *launchpad.MousePadByTier
}

// NewChainReader 创建 ChainReader 实例。
// pool 参数提供多节点容错和熔断能力，由调用方管理生命周期。
func NewChainReader(pool chain.NodePool) *ChainReader {
	return &ChainReader{
		pool:    pool,
		tierAbi: launchpad.NewMouseTier(),
		padAbi:  launchpad.NewMousePadByTier(),
	}
}

// GetTransactionReceipt 查询交易 Receipt。
func (r *ChainReader) GetTransactionReceipt(ctx context.Context, txHash string) (*domain.ReceiptInfo, error) {
	hash := common.HexToHash(txHash)
	receipt, err := r.pool.TransactionReceipt(ctx, hash)
	if err != nil {
		if err == ethereum.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询 receipt %s: %w", txHash, err)
	}

	return &domain.ReceiptInfo{
		Status:      receipt.Status,
		BlockNumber: receipt.BlockNumber.Uint64(),
		GasUsed:     receipt.GasUsed,
	}, nil
}

// GetUserTier 查询用户在 MouseTier 合约的实时 Tier。
func (r *ChainReader) GetUserTier(ctx context.Context, mouseTierAddress string, userAddress string) (int, error) {
	contract := common.HexToAddress(mouseTierAddress)
	user := common.HexToAddress(userAddress)

	data := r.tierAbi.PackGetUserTier(user)

	result, err := r.callContract(ctx, contract, data)
	if err != nil {
		return 0, fmt.Errorf("调用 getUserTier(%s): %w", userAddress, err)
	}

	tier, err := r.tierAbi.UnpackGetUserTier(result)
	if err != nil {
		return 0, fmt.Errorf("解析 getUserTier 返回值: %w", err)
	}

	return int(tier.Int64()), nil
}

// GetUserCredit 查询用户在 MouseTier 合约的实时 Credit。
func (r *ChainReader) GetUserCredit(ctx context.Context, mouseTierAddress string, userAddress string) (*big.Int, error) {
	contract := common.HexToAddress(mouseTierAddress)
	user := common.HexToAddress(userAddress)

	data := r.tierAbi.PackGetUserCredit(user)

	result, err := r.callContract(ctx, contract, data)
	if err != nil {
		return nil, fmt.Errorf("调用 getUserCredit(%s): %w", userAddress, err)
	}

	credit, err := r.tierAbi.UnpackGetUserCredit(result)
	if err != nil {
		return nil, fmt.Errorf("解析 getUserCredit 返回值: %w", err)
	}

	return credit, nil
}

// GetStartBlock 读取 Sale 合约的 startBlock。
func (r *ChainReader) GetStartBlock(ctx context.Context, saleAddr common.Address) (*big.Int, error) {
	data := r.padAbi.PackStartBlock()
	result, err := r.callContract(ctx, saleAddr, data)
	if err != nil {
		return nil, fmt.Errorf("调用 startBlock: %w", err)
	}
	return r.padAbi.UnpackStartBlock(result)
}

// GetEndBlock 读取 Sale 合约的 endBlock。
func (r *ChainReader) GetEndBlock(ctx context.Context, saleAddr common.Address) (*big.Int, error) {
	data := r.padAbi.PackEndBlock()
	result, err := r.callContract(ctx, saleAddr, data)
	if err != nil {
		return nil, fmt.Errorf("调用 endBlock: %w", err)
	}
	return r.padAbi.UnpackEndBlock(result)
}

// GetVestingRevoked 读取 Sale 合约的 vestingRevoked。
func (r *ChainReader) GetVestingRevoked(ctx context.Context, saleAddr common.Address) (bool, error) {
	data := r.padAbi.PackVestingRevoked()
	result, err := r.callContract(ctx, saleAddr, data)
	if err != nil {
		return false, fmt.Errorf("调用 vestingRevoked: %w", err)
	}
	return r.padAbi.UnpackVestingRevoked(result)
}

// GetVestingStartTime 读取 Sale 合约的 vestingStartTime。
func (r *ChainReader) GetVestingStartTime(ctx context.Context, saleAddr common.Address) (*big.Int, error) {
	data := r.padAbi.PackVestingStartTime()
	result, err := r.callContract(ctx, saleAddr, data)
	if err != nil {
		return nil, fmt.Errorf("调用 vestingStartTime: %w", err)
	}
	return r.padAbi.UnpackVestingStartTime(result)
}

// GetOwner 读取 Sale 合约的 owner。
func (r *ChainReader) GetOwner(ctx context.Context, saleAddr common.Address) (common.Address, error) {
	data := r.padAbi.PackOwner()
	result, err := r.callContract(ctx, saleAddr, data)
	if err != nil {
		return common.Address{}, fmt.Errorf("调用 owner: %w", err)
	}
	return r.padAbi.UnpackOwner(result)
}

// GetRaiseToken 读取 Sale 合约的 raiseToken。
func (r *ChainReader) GetRaiseToken(ctx context.Context, saleAddr common.Address) (common.Address, error) {
	data := r.padAbi.PackRaiseToken()
	result, err := r.callContract(ctx, saleAddr, data)
	if err != nil {
		return common.Address{}, fmt.Errorf("调用 raiseToken: %w", err)
	}
	return r.padAbi.UnpackRaiseToken(result)
}

// GetOfferingToken 读取 Sale 合约的 offeringToken。
func (r *ChainReader) GetOfferingToken(ctx context.Context, saleAddr common.Address) (common.Address, error) {
	data := r.padAbi.PackOfferingToken()
	result, err := r.callContract(ctx, saleAddr, data)
	if err != nil {
		return common.Address{}, fmt.Errorf("调用 offeringToken: %w", err)
	}
	return r.padAbi.UnpackOfferingToken(result)
}

// GetMouseTier 读取 Sale 合约的 mouseTier 地址。
func (r *ChainReader) GetMouseTier(ctx context.Context, saleAddr common.Address) (common.Address, error) {
	data := r.padAbi.PackMouseTier()
	result, err := r.callContract(ctx, saleAddr, data)
	if err != nil {
		return common.Address{}, fmt.Errorf("调用 mouseTier: %w", err)
	}
	return r.padAbi.UnpackMouseTier(result)
}

// GetNextScheduleId 读取 Sale 合约的 nextScheduleId。
func (r *ChainReader) GetNextScheduleId(ctx context.Context, saleAddr common.Address) (*big.Int, error) {
	data := r.padAbi.PackNextScheduleId()
	result, err := r.callContract(ctx, saleAddr, data)
	if err != nil {
		return nil, fmt.Errorf("调用 nextScheduleId: %w", err)
	}
	return r.padAbi.UnpackNextScheduleId(result)
}

// GetPoolInfo 读取链上 poolInfo(pid) 返回的完整 Pool 配置。
func (r *ChainReader) GetPoolInfo(ctx context.Context, saleAddr common.Address, pid *big.Int) (*PoolInfo, error) {
	data := r.padAbi.PackPoolInfo(pid)
	result, err := r.callContract(ctx, saleAddr, data)
	if err != nil {
		return nil, fmt.Errorf("调用 poolInfo(%s): %w", pid.String(), err)
	}
	output, err := r.padAbi.UnpackPoolInfo(result)
	if err != nil {
		return nil, fmt.Errorf("解析 poolInfo 返回值: %w", err)
	}
	return &PoolInfo{
		OfferingAmountPool: output.OfferingAmountPool,
		RaisingAmountPool:  output.RaisingAmountPool,
		TotalAmountPool:    output.TotalAmountPool,
		LimitPerUser:       output.LimitPerUserInLP,
		IsSpecialSale:      output.IsSpecialSale,
		HasTax:             output.HasTax,
		VestingPercentage:  output.VestingPercentage,
		VestingCliff:       output.VestingCliff,
		VestingSlicePeriod: output.VestingSlicePeriodSeconds,
		VestingDuration:    output.VestingDuration,
		Configured:         output.Configured,
	}, nil
}

// GetTierLimit 读取 Sale 合约的 tierLimits(tier)。
func (r *ChainReader) GetTierLimit(ctx context.Context, saleAddr common.Address, tier *big.Int) (*big.Int, error) {
	data := r.padAbi.PackTierLimits(tier)
	result, err := r.callContract(ctx, saleAddr, data)
	if err != nil {
		return nil, fmt.Errorf("调用 tierLimits(%s): %w", tier.String(), err)
	}
	return r.padAbi.UnpackTierLimits(result)
}

// GetCeiling 读取 MouseTier 合约的 ceiling。
func (r *ChainReader) GetCeiling(ctx context.Context, mouseTierAddr common.Address) (*big.Int, error) {
	data := r.tierAbi.PackCeiling()
	result, err := r.callContract(ctx, mouseTierAddr, data)
	if err != nil {
		return nil, fmt.Errorf("调用 ceiling: %w", err)
	}
	return r.tierAbi.UnpackCeiling(result)
}

// GetMultiplier 读取 MouseTier 合约的 multiplier。
func (r *ChainReader) GetMultiplier(ctx context.Context, mouseTierAddr common.Address) (*big.Int, error) {
	data := r.tierAbi.PackMultiplier()
	result, err := r.callContract(ctx, mouseTierAddr, data)
	if err != nil {
		return nil, fmt.Errorf("调用 multiplier: %w", err)
	}
	return r.tierAbi.UnpackMultiplier(result)
}

// GetTierBaseAmount 读取 MouseTier 合约的 tierBaseAmount。
func (r *ChainReader) GetTierBaseAmount(ctx context.Context, mouseTierAddr common.Address) (*big.Int, error) {
	data := r.tierAbi.PackTierBaseAmount()
	result, err := r.callContract(ctx, mouseTierAddr, data)
	if err != nil {
		return nil, fmt.Errorf("调用 tierBaseAmount: %w", err)
	}
	return r.tierAbi.UnpackTierBaseAmount(result)
}

// GetUserPoolInfo 读取链上 viewUserPoolInfo(user, pid) 返回的用户 Pool 状态。
func (r *ChainReader) GetUserPoolInfo(ctx context.Context, saleAddr common.Address, user common.Address, pid *big.Int) (*UserPoolInfo, error) {
	data := r.padAbi.PackViewUserPoolInfo(user, pid)
	result, err := r.callContract(ctx, saleAddr, data)
	if err != nil {
		return nil, fmt.Errorf("调用 viewUserPoolInfo(%s, %s): %w", user.Hex(), pid.String(), err)
	}
	output, err := r.padAbi.UnpackViewUserPoolInfo(result)
	if err != nil {
		return nil, fmt.Errorf("解析 viewUserPoolInfo 返回值: %w", err)
	}
	return &UserPoolInfo{
		AmountPool:  output.AmountPool,
		ClaimedPool: output.ClaimedPool,
	}, nil
}

// GetVestingSchedule 读取链上 vestingSchedules(scheduleId) 返回的锁仓计划。
func (r *ChainReader) GetVestingSchedule(ctx context.Context, saleAddr common.Address, scheduleId *big.Int) (*VestingScheduleInfo, error) {
	data := r.padAbi.PackVestingSchedules(scheduleId)
	result, err := r.callContract(ctx, saleAddr, data)
	if err != nil {
		return nil, fmt.Errorf("调用 vestingSchedules(%s): %w", scheduleId.String(), err)
	}
	output, err := r.padAbi.UnpackVestingSchedules(result)
	if err != nil {
		return nil, fmt.Errorf("解析 vestingSchedules 返回值: %w", err)
	}
	return &VestingScheduleInfo{
		Beneficiary: output.Beneficiary,
		Pid:         output.Pid,
		AmountTotal: output.AmountTotal,
		Released:    output.Released,
	}, nil
}

// GetReleasableAmount 读取链上 computeReleasableAmount(scheduleId) 返回的可释放量。
func (r *ChainReader) GetReleasableAmount(ctx context.Context, saleAddr common.Address, scheduleId *big.Int) (*big.Int, error) {
	data := r.padAbi.PackComputeReleasableAmount(scheduleId)
	result, err := r.callContract(ctx, saleAddr, data)
	if err != nil {
		return nil, fmt.Errorf("调用 computeReleasableAmount(%s): %w", scheduleId.String(), err)
	}
	return r.padAbi.UnpackComputeReleasableAmount(result)
}

// BatchGetTierLimits 批量读取指定 tier 的额度限制。
func (r *ChainReader) BatchGetTierLimits(ctx context.Context, saleAddr common.Address, tiers []*big.Int) (map[int]*big.Int, error) {
	calls := make([]Call, len(tiers))
	for i, tier := range tiers {
		calls[i] = Call{Method: "tierLimits", Args: []any{tier}}
	}
	results, err := r.BatchCall(ctx, saleAddr, calls)
	if err != nil {
		return nil, fmt.Errorf("批量调用 tierLimits: %w", err)
	}
	limits := make(map[int]*big.Int, len(tiers))
	for i, result := range results {
		if result.Err != nil || result.Data == nil {
			continue
		}
		limit, err := r.padAbi.UnpackTierLimits(result.Data)
		if err != nil {
			continue
		}
		limits[i] = limit
	}
	return limits, nil
}

// BatchGetUserPoolInfos 批量读取用户的 Pool 状态。
func (r *ChainReader) BatchGetUserPoolInfos(ctx context.Context, saleAddr common.Address, userAddresses []common.Address, pid *big.Int) (map[int]*UserPoolInfo, error) {
	calls := make([]Call, len(userAddresses))
	for i, addr := range userAddresses {
		calls[i] = Call{Method: "viewUserPoolInfo", Args: []any{addr, pid}}
	}
	results, err := r.BatchCall(ctx, saleAddr, calls)
	if err != nil {
		return nil, fmt.Errorf("批量调用 viewUserPoolInfo: %w", err)
	}
	infos := make(map[int]*UserPoolInfo, len(userAddresses))
	for i, result := range results {
		if result.Err != nil || result.Data == nil {
			continue
		}
		output, err := r.padAbi.UnpackViewUserPoolInfo(result.Data)
		if err != nil {
			continue
		}
		infos[i] = &UserPoolInfo{
			AmountPool:  output.AmountPool,
			ClaimedPool: output.ClaimedPool,
		}
	}
	return infos, nil
}

// BatchGetVestingSchedules 批量读取锁仓计划。
func (r *ChainReader) BatchGetVestingSchedules(ctx context.Context, saleAddr common.Address, scheduleIDs []*big.Int) (map[int]*VestingScheduleInfo, error) {
	calls := make([]Call, len(scheduleIDs))
	for i, id := range scheduleIDs {
		calls[i] = Call{Method: "vestingSchedules", Args: []any{id}}
	}
	results, err := r.BatchCall(ctx, saleAddr, calls)
	if err != nil {
		return nil, fmt.Errorf("批量调用 vestingSchedules: %w", err)
	}
	schedules := make(map[int]*VestingScheduleInfo, len(scheduleIDs))
	for i, result := range results {
		if result.Err != nil || result.Data == nil {
			continue
		}
		output, err := r.padAbi.UnpackVestingSchedules(result.Data)
		if err != nil {
			continue
		}
		schedules[i] = &VestingScheduleInfo{
			Beneficiary: output.Beneficiary,
			Pid:         output.Pid,
			AmountTotal: output.AmountTotal,
			Released:    output.Released,
		}
	}
	return schedules, nil
}

// BatchGetReleasableAmounts 批量读取可释放量。
func (r *ChainReader) BatchGetReleasableAmounts(ctx context.Context, saleAddr common.Address, scheduleIDs []*big.Int) (map[int]*big.Int, error) {
	calls := make([]Call, len(scheduleIDs))
	for i, id := range scheduleIDs {
		calls[i] = Call{Method: "computeReleasableAmount", Args: []any{id}}
	}
	results, err := r.BatchCall(ctx, saleAddr, calls)
	if err != nil {
		return nil, fmt.Errorf("批量调用 computeReleasableAmount: %w", err)
	}
	amounts := make(map[int]*big.Int, len(scheduleIDs))
	for i, result := range results {
		if result.Err != nil || result.Data == nil {
			continue
		}
		amount, err := r.padAbi.UnpackComputeReleasableAmount(result.Data)
		if err != nil {
			continue
		}
		amounts[i] = amount
	}
	return amounts, nil
}

// BatchCall 批量执行链上合约调用，将多个 Call 合并为一次 NodePool.BatchCallContract 请求。
func (r *ChainReader) BatchCall(ctx context.Context, saleAddr common.Address, calls []Call) ([]CallResult, error) {
	if len(calls) == 0 {
		return nil, nil
	}

	// 编码每个调用为 CallMsg
	msgs := make([]ethereum.CallMsg, len(calls))
	// 保存 Pack 方法名到解码信息的映射
	methodNames := make([]string, len(calls))

	for i, call := range calls {
		target := saleAddr
		if call.To != nil {
			target = *call.To
		}
		data, err := r.packCall(call)
		if err != nil {
			return nil, fmt.Errorf("编码调用 %s: %w", call.Method, err)
		}
		msgs[i] = ethereum.CallMsg{
			To:   &target,
			Data: data,
		}
		methodNames[i] = call.Method
	}

	results, err := r.pool.BatchCallContract(ctx, msgs)
	if err != nil {
		return nil, fmt.Errorf("批量调用合约: %w", err)
	}

	callResults := make([]CallResult, len(calls))
	for i, data := range results {
		callResults[i] = CallResult{
			Method: methodNames[i],
			Data:   data,
		}
		if data == nil {
			callResults[i].Err = fmt.Errorf("合约调用 %s 返回空数据", methodNames[i])
		}
	}

	return callResults, nil
}

// packCall 将 Call 编码为合约调用数据。
func (r *ChainReader) packCall(call Call) ([]byte, error) {
	switch call.Method {
	case "startBlock":
		return r.padAbi.PackStartBlock(), nil
	case "endBlock":
		return r.padAbi.PackEndBlock(), nil
	case "vestingRevoked":
		return r.padAbi.PackVestingRevoked(), nil
	case "vestingStartTime":
		return r.padAbi.PackVestingStartTime(), nil
	case "owner":
		return r.padAbi.PackOwner(), nil
	case "raiseToken":
		return r.padAbi.PackRaiseToken(), nil
	case "offeringToken":
		return r.padAbi.PackOfferingToken(), nil
	case "mouseTier":
		return r.padAbi.PackMouseTier(), nil
	case "nextScheduleId":
		return r.padAbi.PackNextScheduleId(), nil
	case "poolInfo":
		arg, err := argBigInt(call.Args, 0)
		if err != nil {
			return nil, fmt.Errorf("poolInfo: %w", err)
		}
		return r.padAbi.PackPoolInfo(arg), nil
	case "tierLimits":
		arg, err := argBigInt(call.Args, 0)
		if err != nil {
			return nil, fmt.Errorf("tierLimits: %w", err)
		}
		return r.padAbi.PackTierLimits(arg), nil
	case "viewUserPoolInfo":
		addr, err := argAddress(call.Args, 0)
		if err != nil {
			return nil, fmt.Errorf("viewUserPoolInfo: %w", err)
		}
		pid, err := argBigInt(call.Args, 1)
		if err != nil {
			return nil, fmt.Errorf("viewUserPoolInfo: %w", err)
		}
		return r.padAbi.PackViewUserPoolInfo(addr, pid), nil
	case "vestingSchedules":
		arg, err := argBigInt(call.Args, 0)
		if err != nil {
			return nil, fmt.Errorf("vestingSchedules: %w", err)
		}
		return r.padAbi.PackVestingSchedules(arg), nil
	case "computeReleasableAmount":
		arg, err := argBigInt(call.Args, 0)
		if err != nil {
			return nil, fmt.Errorf("computeReleasableAmount: %w", err)
		}
		return r.padAbi.PackComputeReleasableAmount(arg), nil
	case "ceiling":
		return r.tierAbi.PackCeiling(), nil
	case "multiplier":
		return r.tierAbi.PackMultiplier(), nil
	case "tierBaseAmount":
		return r.tierAbi.PackTierBaseAmount(), nil
	default:
		return nil, fmt.Errorf("未知方法: %s", call.Method)
	}
}

// argBigInt 安全地从 args 中提取 *big.Int。
func argBigInt(args []any, index int) (*big.Int, error) {
	if len(args) <= index {
		return nil, fmt.Errorf("参数索引 %d 越界（长度 %d）", index, len(args))
	}
	v, ok := args[index].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("参数 %d 类型错误: 期望 *big.Int，实际 %T", index, args[index])
	}
	return v, nil
}

// argAddress 安全地从 args 中提取 common.Address。
func argAddress(args []any, index int) (common.Address, error) {
	if len(args) <= index {
		return common.Address{}, fmt.Errorf("参数索引 %d 越界（长度 %d）", index, len(args))
	}
	v, ok := args[index].(common.Address)
	if !ok {
		return common.Address{}, fmt.Errorf("参数 %d 类型错误: 期望 common.Address，实际 %T", index, args[index])
	}
	return v, nil
}

// callContract 执行合约 view 调用。
func (r *ChainReader) callContract(ctx context.Context, to common.Address, data []byte) ([]byte, error) {
	msg := ethereum.CallMsg{
		To:   &to,
		Data: data,
	}
	return r.pool.CallContract(ctx, msg, nil)
}

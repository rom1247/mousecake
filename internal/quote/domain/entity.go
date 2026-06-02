// Package domain 定义 quote 模块的领域模型，包含实体、值对象、状态机和仓库接口。
// 本包不能导入任何外层包（编译期隔离）。
package domain

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// SwapStatus 表示 Swap 记录的状态。
type SwapStatus string

const (
	// SwapStatusPending 待提交交易哈希。
	SwapStatusPending SwapStatus = "pending"
	// SwapStatusSubmitted 已提交交易哈希。
	SwapStatusSubmitted SwapStatus = "submitted"
)

// SwapMode 表示兑换模式。
type SwapMode string

const (
	// SwapModeExactIn 精确输入模式。
	SwapModeExactIn SwapMode = "exactIn"
	// SwapModeExactOut 精确输出模式。
	SwapModeExactOut SwapMode = "exactOut"
)

// validSwapModes 合法的兑换模式集合。
var validSwapModes = map[SwapMode]bool{
	SwapModeExactIn:  true,
	SwapModeExactOut: true,
}

// 领域哨兵错误
var (
	// ErrProviderEmpty provider 不能为空。
	ErrProviderEmpty = errors.New("quote: provider 不能为空")
	// ErrFromTokenEmpty from_token 不能为空。
	ErrFromTokenEmpty = errors.New("quote: from_token 不能为空")
	// ErrToTokenEmpty to_token 不能为空。
	ErrToTokenEmpty = errors.New("quote: to_token 不能为空")
	// ErrFromAmountEmpty from_amount 不能为空。
	ErrFromAmountEmpty = errors.New("quote: from_amount 不能为空")
	// ErrToAmountEmpty to_amount 不能为空。
	ErrToAmountEmpty = errors.New("quote: to_amount 不能为空")
	// ErrInvalidSwapMode 无效的兑换模式。
	ErrInvalidSwapMode = errors.New("quote: 无效的兑换模式")
	// ErrAlreadySubmitted swap 记录已提交。
	ErrAlreadySubmitted = errors.New("quote: swap 记录已提交")
	// ErrSwapRecordNotFound swap 记录不存在。
	ErrSwapRecordNotFound = errors.New("quote: swap 记录不存在")
	// ErrProviderNotFound 供应商不存在。
	ErrProviderNotFound = errors.New("quote: 供应商不存在")
)

// SwapRecord 是 Swap 记录实体，表示一次代币兑换请求。
type SwapRecord struct {
	// ID Snowflake 风格主键。
	ID int64
	// Provider 供应商名称。
	Provider string
	// ChainID 链 ID。
	ChainID int
	// FromToken 源代币地址。
	FromToken string
	// ToToken 目标代币地址。
	ToToken string
	// FromAmount 源代币数量（wei）。
	FromAmount string
	// ToAmount 目标代币数量（wei）。
	ToAmount string
	// SlippagePercent 滑点百分比。
	SlippagePercent float64
	// SwapMode 兑换模式。
	SwapMode SwapMode
	// Status 当前状态。
	Status SwapStatus
	// TxHash 交易哈希。
	TxHash string
	// CreatedAt 创建时间。
	CreatedAt time.Time
	// UpdatedAt 更新时间。
	UpdatedAt time.Time
}

// NewSwapRecordOpts 创建 SwapRecord 的输入参数。
type NewSwapRecordOpts struct {
	Provider        string
	ChainID         int
	FromToken       string
	ToToken         string
	FromAmount      string
	ToAmount        string
	SlippagePercent float64
	SwapMode        SwapMode
}

// NewSwapRecord 创建 SwapRecord 实体，校验必填字段后设置 status=pending。
func NewSwapRecord(nodeID int64, opts NewSwapRecordOpts) (*SwapRecord, error) {
	if opts.Provider == "" {
		return nil, ErrProviderEmpty
	}
	if opts.FromToken == "" {
		return nil, ErrFromTokenEmpty
	}
	if opts.ToToken == "" {
		return nil, ErrToTokenEmpty
	}
	if opts.FromAmount == "" {
		return nil, ErrFromAmountEmpty
	}
	if opts.ToAmount == "" {
		return nil, ErrToAmountEmpty
	}
	if !validSwapModes[opts.SwapMode] {
		return nil, ErrInvalidSwapMode
	}

	now := time.Now()
	return &SwapRecord{
		ID:              defaultIDGenerator.Generate(nodeID),
		Provider:        opts.Provider,
		ChainID:         opts.ChainID,
		FromToken:       opts.FromToken,
		ToToken:         opts.ToToken,
		FromAmount:      opts.FromAmount,
		ToAmount:        opts.ToAmount,
		SlippagePercent: opts.SlippagePercent,
		SwapMode:        opts.SwapMode,
		Status:          SwapStatusPending,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// Submit 提交交易哈希，将状态从 pending 转换为 submitted。
func (r *SwapRecord) Submit(txHash string) error {
	if r.Status == SwapStatusSubmitted {
		return ErrAlreadySubmitted
	}
	r.Status = SwapStatusSubmitted
	r.TxHash = txHash
	r.UpdatedAt = time.Now()
	return nil
}

// SwapRecordSnapshot 从数据库重建 SwapRecord 实体所需的全部字段快照。
type SwapRecordSnapshot struct {
	ID              int64
	Provider        string
	ChainID         int
	FromToken       string
	ToToken         string
	FromAmount      string
	ToAmount        string
	SlippagePercent float64
	SwapMode        SwapMode
	Status          SwapStatus
	TxHash          string
}

// ReconstructSwapRecord 从数据库重建 SwapRecord 实体（跳过业务规则校验）。
func ReconstructSwapRecord(snapshot SwapRecordSnapshot) *SwapRecord {
	return &SwapRecord{
		ID:              snapshot.ID,
		Provider:        snapshot.Provider,
		ChainID:         snapshot.ChainID,
		FromToken:       snapshot.FromToken,
		ToToken:         snapshot.ToToken,
		FromAmount:      snapshot.FromAmount,
		ToAmount:        snapshot.ToAmount,
		SlippagePercent: snapshot.SlippagePercent,
		SwapMode:        snapshot.SwapMode,
		Status:          snapshot.Status,
		TxHash:          snapshot.TxHash,
	}
}

// --- Snowflake ID 生成器 ---

// idGenerator 基于 Snowflake 风格的 ID 生成器。
type idGenerator struct {
	mu        sync.Mutex
	timestamp int64
	sequence  int64
}

// defaultIDGenerator 全局默认 ID 生成器实例。
var defaultIDGenerator = &idGenerator{}

const (
	epoch          int64 = 1700000000000 // 自定义纪元（毫秒）
	nodeIDBits     int64 = 10
	sequenceBits   int64 = 12
	maxNodeID      int64 = (1 << nodeIDBits) - 1     // 1023
	maxSequence    int64 = (1 << sequenceBits) - 1   // 4095
	nodeIDShift    int64 = sequenceBits              // 12
	timestampShift int64 = nodeIDBits + sequenceBits // 22
)

// Generate 生成一个 Snowflake 风格 ID。
func (g *idGenerator) Generate(nodeID int64) int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	if nodeID < 0 || nodeID > maxNodeID {
		nodeID = 1
	}

	now := time.Now().UnixMilli()
	if now == g.timestamp {
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			for now <= g.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		g.sequence = 0
	}
	g.timestamp = now

	id := ((now - epoch) << timestampShift) | (nodeID << nodeIDShift) | g.sequence
	return id
}

// GenerateID 生成一个全局唯一的 Snowflake 风格 ID（使用默认生成器）。
func GenerateID(nodeID int64) int64 {
	return defaultIDGenerator.Generate(nodeID)
}

// ValidateIDRange 校验 ID 是否在有效范围内（17-19 位数字）。
func ValidateIDRange(id int64) error {
	if id < 10000000000000000 {
		return fmt.Errorf("quote: ID %d 不在有效范围 (17-19 位数字)", id)
	}
	return nil
}

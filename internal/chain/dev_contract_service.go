package chain

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ContractQueryRequest 通用合约调用请求。
type ContractQueryRequest struct { // @name ContractQueryRequest
	// Contract 合约名称（注册表中的名称）。
	Contract string `json:"contract" binding:"required"`
	// Method 方法名称。
	Method string `json:"method" binding:"required"`
	// Args 方法参数列表（JSON 值）。
	Args []any `json:"args"`
	// Address 可选，指定合约地址（覆盖注册表中的默认地址）。
	Address string `json:"address"`
	// Block 可选，查询的区块号（十六进制字符串或 "latest"/"pending"）。
	Block string `json:"block"`
}

// ContractQueryResult 通用合约调用结果。
type ContractQueryResult struct { // @name ContractQueryResult
	// Method 方法名称。
	Method string `json:"method"`
	// Mutability 状态可变性。
	Mutability string `json:"mutability"`
	// Outputs 只读调用的返回值列表（JSON 安全格式）。
	Outputs []any `json:"outputs,omitempty"`
	// TxHash 写入调用的交易哈希。
	TxHash string `json:"tx_hash,omitempty"`
	// BlockNumber 交易所在区块号。
	BlockNumber uint64 `json:"block_number,omitempty"`
	// GasUsed 交易消耗的 gas。
	GasUsed uint64 `json:"gas_used,omitempty"`
	// Events 写入调用触发的事件。
	Events []ParsedEvent `json:"events,omitempty"`
}

// writeResult 写入调用的结果。
type writeResult struct {
	txHash      string
	blockNumber uint64
	gasUsed     uint64
	events      []ParsedEvent
}

// DevContractService 通用合约查询/执行服务，仅供开发环境使用。
// 支持通过 ABI 动态调用任意已注册合约的方法。
type DevContractService struct {
	registry *ABIRegistry
	pool     NodePool
	signer   *Signer
	parser   *EventParser
}

// NewDevContractService 创建通用合约查询服务。
func NewDevContractService(registry *ABIRegistry, pool NodePool, signer *Signer, parser *EventParser) *DevContractService {
	return &DevContractService{
		registry: registry,
		pool:     pool,
		signer:   signer,
		parser:   parser,
	}
}

// Query 执行通用合约调用，根据方法的 StateMutability 自动选择只读查询或写入交易。
func (s *DevContractService) Query(ctx context.Context, req ContractQueryRequest) (*ContractQueryResult, error) {
	// 1. 查找合约
	entry, err := s.registry.Get(req.Contract)
	if err != nil {
		return nil, err
	}

	// 2. 查找方法
	method, ok := entry.ABI.Methods[req.Method]
	if !ok {
		return nil, fmt.Errorf("合约 %s 没有方法 %s", req.Contract, req.Method)
	}

	// 3. 确定合约地址
	addr := entry.Address
	if req.Address != "" {
		if !common.IsHexAddress(req.Address) {
			return nil, fmt.Errorf("无效的合约地址: %s", req.Address)
		}
		addr = common.HexToAddress(req.Address)
	}
	if addr == (common.Address{}) {
		return nil, fmt.Errorf("合约 %s 没有默认地址，请在请求中指定 address 字段", req.Contract)
	}

	// 4. 转换参数
	convertedArgs, err := ConvertParams(req.Args, method.Inputs)
	if err != nil {
		return nil, fmt.Errorf("参数转换: %w", err)
	}

	// 5. 编码调用数据（方法选择器 + 参数）
	callData, err := packCallData(method, convertedArgs)
	if err != nil {
		return nil, fmt.Errorf("编码调用数据: %w", err)
	}

	// 6. 根据读写类型执行
	result := &ContractQueryResult{
		Method:     method.Name,
		Mutability: method.StateMutability,
	}

	if isReadMethod(method.StateMutability) {
		outputs, err := s.executeRead(ctx, addr, callData, method.Outputs, req.Block)
		if err != nil {
			return nil, err
		}
		result.Outputs = outputs
	} else {
		if s.signer == nil {
			return nil, fmt.Errorf("写入方法 %s 需要 admin 私钥，请在配置中设置 launchpad.admin_private_key", method.Name)
		}
		txResult, err := s.executeWrite(ctx, addr, callData)
		if err != nil {
			return nil, err
		}
		result.TxHash = txResult.txHash
		result.BlockNumber = txResult.blockNumber
		result.GasUsed = txResult.gasUsed
		result.Events = txResult.events
	}

	return result, nil
}

// executeRead 执行只读合约调用（eth_call），解码返回值。
func (s *DevContractService) executeRead(ctx context.Context, addr common.Address, data []byte, outputs abi.Arguments, blockTag string) ([]any, error) {
	blockNum, err := parseBlockNumber(blockTag)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &addr,
		Data: data,
	}

	returnData, err := s.pool.CallContract(ctx, msg, blockNum)
	if err != nil {
		return nil, fmt.Errorf("eth_call: %w", err)
	}

	// 解码返回值
	if len(outputs) == 0 {
		return nil, nil
	}

	values, err := outputs.Unpack(returnData)
	if err != nil {
		return nil, fmt.Errorf("解码返回值: %w", err)
	}

	// 转换为 JSON 安全格式
	result := make([]any, len(values))
	for i, v := range values {
		result[i] = convertValue(v)
	}
	return result, nil
}

// executeWrite 执行写入调用（签名广播交易 → 等待 receipt → 解析事件）。
func (s *DevContractService) executeWrite(ctx context.Context, addr common.Address, data []byte) (*writeResult, error) {
	txHash, err := s.signer.SignAndBroadcast(ctx, addr, data, nil)
	if err != nil {
		return nil, fmt.Errorf("签名广播: %w", err)
	}

	receipt, err := s.signer.WaitForReceipt(ctx, txHash)
	if err != nil {
		return &writeResult{txHash: txHash.Hex()}, fmt.Errorf("等待 receipt: %w", err)
	}

	if receipt.Status == types.ReceiptStatusFailed {
		return &writeResult{txHash: txHash.Hex()}, fmt.Errorf("链上交易执行失败 (tx=%s)", txHash.Hex())
	}

	var events []ParsedEvent
	if s.parser != nil {
		events = s.parser.ParseLogs(receipt.Logs)
	}

	return &writeResult{
		txHash:      txHash.Hex(),
		blockNumber: receipt.BlockNumber.Uint64(),
		gasUsed:     receipt.GasUsed,
		events:      events,
	}, nil
}

// packCallData 将方法选择器和参数打包为调用数据。
func packCallData(method abi.Method, args []any) ([]byte, error) {
	packed, err := method.Inputs.Pack(args...)
	if err != nil {
		return nil, err
	}
	callData := make([]byte, len(method.ID)+len(packed))
	copy(callData, method.ID)
	copy(callData[len(method.ID):], packed)
	return callData, nil
}

// parseBlockNumber 解析区块号标签或十六进制字符串。
func parseBlockNumber(blockTag string) (*big.Int, error) {
	switch {
	case blockTag == "" || blockTag == "latest":
		return nil, nil
	case blockTag == "pending":
		return big.NewInt(-1), nil
	default:
		hexStr := strings.TrimPrefix(blockTag, "0x")
		n, ok := new(big.Int).SetString(hexStr, 16)
		if !ok {
			return nil, fmt.Errorf("无效的区块号: %s", blockTag)
		}
		return n, nil
	}
}

// isReadMethod 判断方法是否为只读（view/pure）。
func isReadMethod(mutability string) bool {
	return mutability == "view" || mutability == "pure"
}

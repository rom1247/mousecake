package chain

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
)

// ContractEntry 注册表中的合约条目，包含名称、ABI 和默认部署地址。
type ContractEntry struct {
	// Name 合约名称（如 "MousePadByTier"）。
	Name string
	// ABI 解析后的合约 ABI。
	ABI abi.ABI
	// Address 默认部署地址，可为零值（调用时必须显式指定）。
	Address common.Address
}

// ABIRegistry 合约 ABI 注册中心，管理所有已知合约的 ABI 和默认地址。
// 仅供开发环境使用，支持按合约名查找 ABI 和方法。
type ABIRegistry struct {
	mu        sync.RWMutex
	contracts map[string]*ContractEntry
}

// NewABIRegistry 创建空的 ABI 注册中心。
func NewABIRegistry() *ABIRegistry {
	return &ABIRegistry{
		contracts: make(map[string]*ContractEntry),
	}
}

// Register 从 bind.MetaData 注册合约 ABI，可同时指定默认部署地址。
// 如果同名合约已注册，返回错误。
func (r *ABIRegistry) Register(name string, metaData *bind.MetaData, address ...common.Address) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.contracts[name]; exists {
		return fmt.Errorf("合约 %s 已注册", name)
	}

	parsed, err := metaData.ParseABI()
	if err != nil {
		return fmt.Errorf("解析 %s ABI: %w", name, err)
	}

	entry := &ContractEntry{
		Name: name,
		ABI:  *parsed,
	}
	if len(address) > 0 {
		entry.Address = address[0]
	}

	r.contracts[name] = entry
	return nil
}

// Get 根据合约名获取注册条目，未找到返回错误。
func (r *ABIRegistry) Get(name string) (*ContractEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, ok := r.contracts[name]
	if !ok {
		return nil, fmt.Errorf("合约 %s 未注册", name)
	}
	return entry, nil
}

// ListContracts 返回所有已注册合约的名称列表。
func (r *ABIRegistry) ListContracts() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.contracts))
	for name := range r.contracts {
		names = append(names, name)
	}
	return names
}

// ListMethods 返回指定合约的所有方法签名（名称 + StateMutability）。
func (r *ABIRegistry) ListMethods(contractName string) ([]MethodInfo, error) {
	entry, err := r.Get(contractName)
	if err != nil {
		return nil, err
	}

	methods := make([]MethodInfo, 0, len(entry.ABI.Methods))
	for _, m := range entry.ABI.Methods {
		methods = append(methods, MethodInfo{
			Name:        m.Name,
			Selector:    "0x" + common.Bytes2Hex(m.ID),
			Mutability:  m.StateMutability,
			InputTypes:  formatArguments(m.Inputs),
			OutputTypes: formatArguments(m.Outputs),
		})
	}
	return methods, nil
}

// MethodInfo 合约方法的摘要信息。
type MethodInfo struct { // @name MethodInfo
	// Name 方法名。
	Name string `json:"name"`
	// Selector 函数选择器（4 字节十六进制，如 "0xedac985b"）。
	Selector string `json:"selector"`
	// Mutability 状态可变性（view/pure/nonpayable/payable）。
	Mutability string `json:"mutability"`
	// InputTypes 输入参数类型描述列表。
	InputTypes []ArgInfo `json:"input_types"`
	// OutputTypes 输出参数类型描述列表。
	OutputTypes []ArgInfo `json:"output_types"`
}

// ArgInfo ABI 参数的类型描述。
type ArgInfo struct { // @name ArgInfo
	// Name 参数名。
	Name string `json:"name"`
	// Type Solidity 类型字符串（如 uint256、address）。
	Type string `json:"type"`
}

// formatArguments 将 abi.Arguments 转换为 []ArgInfo。
func formatArguments(args abi.Arguments) []ArgInfo {
	result := make([]ArgInfo, len(args))
	for i, a := range args {
		result[i] = ArgInfo{
			Name: a.Name,
			Type: a.Type.String(),
		}
	}
	return result
}

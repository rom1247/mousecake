package launchpad

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mousecake-go/mousecake-go/internal/chain/contract/launchpad"
)

// ContractCodeEncoder 使用生成的合约绑定代码编码 calldata，无需加载 ABI JSON 文件。
type ContractCodeEncoder struct {
	padABI      *launchpad.MousePadByTier
	deployerABI *launchpad.MousePadByTierDeployer
}

// NewContractCodeEncoder 创建使用生成合约代码的编码器。
func NewContractCodeEncoder() *ContractCodeEncoder {
	return &ContractCodeEncoder{
		padABI:      launchpad.NewMousePadByTier(),
		deployerABI: launchpad.NewMousePadByTierDeployer(),
	}
}

// EncodeCall 编码合约函数调用，返回 calldata。
// contractName 参数指定合约："MousePadByTierDeployer" 或 "MousePadByTier"。
func (e *ContractCodeEncoder) EncodeCall(contractName, methodName string, args ...any) ([]byte, error) {
	if contractName == "MousePadByTierDeployer" {
		return e.encodeDeployer(methodName, args...)
	}
	// 默认使用 MousePadByTier 合约
	return e.encodeMousePadByTier(methodName, args...)
}

// encodeDeployer 编码 MousePadByTierDeployer 合约调用。
func (e *ContractCodeEncoder) encodeDeployer(methodName string, args ...any) ([]byte, error) {
	switch methodName {
	case "createSale":
		// createSale(address raiseToken, address offeringToken, address admin, address mouseTier, uint256 startBlock, uint256 endBlock)
		return e.deployerABI.TryPackCreateSale(
			toAddress(args[0]), // raiseToken
			toAddress(args[1]), // offeringToken
			toAddress(args[2]), // admin
			toAddress(args[3]), // mouseTier
			bigInt(args[4]),    // startBlock
			bigInt(args[5]),    // endBlock
		)
	default:
		return nil, &ErrUnsupportedMethod{Method: methodName, Contract: "MousePadByTierDeployer"}
	}
}

// encodeMousePadByTier 编码 MousePadByTier 合约调用。
func (e *ContractCodeEncoder) encodeMousePadByTier(methodName string, args ...any) ([]byte, error) {
	switch methodName {
	case "setPool":
		// setPool(uint256 pid, uint256 offeringAmount, uint256 raisingAmount, uint256 limitPerUser, bool isSpecialSale, bool hasTax, uint256 vestingPercentage, uint256 vestingCliff, uint256 vestingSlicePeriodSeconds, uint256 vestingDuration)
		return e.padABI.TryPackSetPool(
			bigInt(args[0]), // pid
			bigInt(args[1]), // offeringAmount
			bigInt(args[2]), // raisingAmount
			bigInt(args[3]), // limitPerUser
			toBool(args[4]), // isSpecialSale
			toBool(args[5]), // hasTax
			bigInt(args[6]), // vestingPercentage
			bigInt(args[7]), // vestingCliff
			bigInt(args[8]), // vestingSlicePeriodSeconds
			bigInt(args[9]), // vestingDuration
		)
	case "setTierLimits":
		// setTierLimits(uint256 tier, uint256 limit)
		return e.padABI.TryPackSetTierLimits(
			bigInt(args[0]), // tier
			bigInt(args[1]), // limit
		)
	case "addWhitelist":
		// addWhitelist(address[] users)
		addresses := toAddressArray(args[0])
		return e.padABI.TryPackAddWhitelist(addresses)
	case "removeWhitelist":
		// removeWhitelist(address[] users)
		addresses := toAddressArray(args[0])
		return e.padABI.TryPackRemoveWhitelist(addresses)
	case "setStartEndBlock":
		// setStartEndBlock(uint256 startBlock, uint256 endBlock)
		return e.padABI.TryPackSetStartEndBlock(
			bigInt(args[0]), // startBlock
			bigInt(args[1]), // endBlock
		)
	case "revoke":
		// revoke()
		return e.padABI.TryPackRevoke()
	case "finalWithdraw":
		// finalWithdraw(uint256 raiseAmount, uint256 offeringAmount)
		return e.padABI.TryPackFinalWithdraw(
			bigInt(args[0]), // raiseAmount
			bigInt(args[1]), // offeringAmount
		)
	case "recoverToken":
		// recoverToken(address token, address to, uint256 amount)
		return e.padABI.TryPackRecoverToken(
			toAddress(args[0]), // token
			toAddress(args[1]), // to
			bigInt(args[2]),    // amount
		)
	case "deposit":
		// deposit(uint256 amount, uint256 pid)
		return e.padABI.TryPackDeposit(
			bigInt(args[0]), // amount
			bigInt(args[1]), // pid
		)
	case "harvest":
		// harvest(uint256 pid)
		return e.padABI.TryPackHarvest(
			bigInt(args[0]), // pid
		)
	case "release":
		// release(uint256 scheduleId)
		return e.padABI.TryPackRelease(
			bigInt(args[0]), // scheduleId
		)
	default:
		return nil, &ErrUnsupportedMethod{Method: methodName, Contract: "MousePadByTier"}
	}
}

// CalldataHash 计算 calldata 的 SHA-256 指纹，用于去重。
func (e *ContractCodeEncoder) CalldataHash(calldata []byte) string {
	hash := sha256.Sum256(calldata)
	return "0x" + hex.EncodeToString(hash[:])
}

// toAddress 将任意类型转换为 common.Address。
func toAddress(v any) common.Address {
	switch val := v.(type) {
	case common.Address:
		return val
	case string:
		return common.HexToAddress(val)
	default:
		return common.Address{}
	}
}

// toAddressArray 将字符串数组转换为 common.Address 数组。
func toAddressArray(v any) []common.Address {
	switch val := v.(type) {
	case []common.Address:
		return val
	case []string:
		addresses := make([]common.Address, len(val))
		for i, addr := range val {
			addresses[i] = common.HexToAddress(addr)
		}
		return addresses
	default:
		return nil
	}
}

// toBool 将任意类型转换为 bool。
func toBool(v any) bool {
	switch val := v.(type) {
	case bool:
		return val
	default:
		return false
	}
}

// bigInt 将任意类型转换为 *big.Int。
func bigInt(v any) *big.Int {
	switch val := v.(type) {
	case *big.Int:
		return val
	case int64:
		return big.NewInt(val)
	case int:
		return big.NewInt(int64(val))
	case uint64:
		return new(big.Int).SetUint64(val)
	case string:
		result := new(big.Int)
		result.SetString(val, 10)
		return result
	default:
		return big.NewInt(0)
	}
}

// ErrUnsupportedMethod 表示不支持的方法错误。
type ErrUnsupportedMethod struct {
	Method   string
	Contract string
}

func (e *ErrUnsupportedMethod) Error() string {
	if e.Contract != "" {
		return "不支持的方法: " + e.Contract + "." + e.Method
	}
	return "不支持的方法: " + e.Method
}

// EncoderInterface 定义 PrepareService 需要的 ABI 编码器接口。
type EncoderInterface interface {
	// EncodeCall 编码合约函数调用，返回 calldata。
	EncodeCall(contractName, methodName string, args ...any) ([]byte, error)
	// CalldataHash 计算 calldata 的 SHA-256 指纹。
	CalldataHash(calldata []byte) string
}

// 确保实现了 EncoderInterface 接口
var _ EncoderInterface = (*ContractCodeEncoder)(nil)

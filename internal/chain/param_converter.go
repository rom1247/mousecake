package chain

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// ConvertParams 根据 ABI 参数类型定义，将 JSON 值（[]any）转换为 Go 类型值（[]any）。
// JSON 输入的 number 类型在 Go 中为 float64，string 保持 string，bool 保持 bool。
func ConvertParams(jsonArgs []any, abiArgs abi.Arguments) ([]any, error) {
	if len(jsonArgs) != len(abiArgs) {
		return nil, fmt.Errorf("参数数量不匹配: 传入 %d 个, ABI 需要 %d 个", len(jsonArgs), len(abiArgs))
	}

	result := make([]any, len(jsonArgs))
	for i, arg := range abiArgs {
		converted, err := convertABIParam(jsonArgs[i], arg.Type)
		if err != nil {
			return nil, fmt.Errorf("转换参数 %d (%s, 类型 %s): %w", i, arg.Name, arg.Type.String(), err)
		}
		result[i] = converted
	}
	return result, nil
}

// convertABIParam 根据 abi.Type 将单个 JSON 值转换为 Go 类型。
func convertABIParam(val any, typ abi.Type) (any, error) {
	if val == nil {
		return nil, nil
	}

	switch typ.T {
	case abi.IntTy, abi.UintTy:
		return convertInt(val, typ.Size)
	case abi.AddressTy:
		return convertAddress(val)
	case abi.BoolTy:
		return convertBool(val)
	case abi.StringTy:
		return convertString(val)
	case abi.BytesTy:
		return convertBytes(val)
	case abi.FixedBytesTy:
		return convertFixedBytes(val, typ.Size)
	case abi.SliceTy:
		return convertSlice(val, *typ.Elem)
	case abi.ArrayTy:
		return convertArrayParam(val, *typ.Elem, typ.Size)
	case abi.TupleTy:
		return convertTuple(val, typ.TupleElems, typ.TupleRawNames)
	default:
		return nil, fmt.Errorf("不支持的 ABI 类型: %s", typ.String())
	}
}

// convertInt 将 JSON 值转换为 *big.Int。
// 支持 float64、string、int64。
func convertInt(val any, bitSize int) (*big.Int, error) {
	switch v := val.(type) {
	case float64:
		if v != float64(int64(v)) {
			return nil, fmt.Errorf("浮点数 %v 不能无损转换为整数", v)
		}
		return big.NewInt(int64(v)), nil
	case string:
		s := strings.TrimSpace(v)
		// 尝试十进制
		n := new(big.Int)
		if _, ok := n.SetString(s, 10); ok {
			return n, nil
		}
		// 尝试十六进制
		hexStr := strings.TrimPrefix(s, "0x")
		hexStr = strings.TrimPrefix(hexStr, "0X")
		if _, ok := n.SetString(hexStr, 16); ok {
			return n, nil
		}
		return nil, fmt.Errorf("无法解析为整数: %s", s)
	case int64:
		return big.NewInt(v), nil
	case int:
		return big.NewInt(int64(v)), nil
	default:
		return nil, fmt.Errorf("期望数字或字符串，实际类型: %T", val)
	}
}

// convertAddress 将 JSON 值转换为 common.Address。
func convertAddress(val any) (common.Address, error) {
	s, ok := val.(string)
	if !ok {
		return common.Address{}, fmt.Errorf("期望字符串，实际类型: %T", val)
	}
	if !common.IsHexAddress(s) {
		return common.Address{}, fmt.Errorf("无效的以太坊地址: %s", s)
	}
	return common.HexToAddress(s), nil
}

// convertBool 将 JSON 值转换为 bool。
func convertBool(val any) (bool, error) {
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false, fmt.Errorf("无法解析为布尔值: %s", v)
		}
		return b, nil
	default:
		return false, fmt.Errorf("期望布尔值，实际类型: %T", val)
	}
}

// convertString 将 JSON 值转换为 string。
func convertString(val any) (string, error) {
	s, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("期望字符串，实际类型: %T", val)
	}
	return s, nil
}

// convertBytes 将 JSON 值转换为 []byte。
func convertBytes(val any) ([]byte, error) {
	s, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("期望十六进制字符串，实际类型: %T", val)
	}
	return hex.DecodeString(strings.TrimPrefix(s, "0x"))
}

// convertFixedBytes 将 JSON 值转换为固定长度字节数组。
func convertFixedBytes(val any, size int) (any, error) {
	s, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("期望十六进制字符串，实际类型: %T", val)
	}
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		return nil, fmt.Errorf("解析十六进制: %w", err)
	}
	if len(b) > size {
		return nil, fmt.Errorf("字节长度 %d 超过固定大小 %d", len(b), size)
	}

	// 创建固定大小数组并右对齐（与 EVM 一致）
	switch size {
	case 32:
		var arr [32]byte
		copy(arr[32-len(b):], b)
		return arr, nil
	default:
		// 对于其他固定大小，返回 []byte（abi.Arguments.Pack 可以接受）
		padded := make([]byte, size)
		copy(padded[size-len(b):], b)
		return padded, nil
	}
}

// convertSlice 将 JSON 数组转换为 Go 切片。
func convertSlice(val any, elemType abi.Type) (any, error) {
	arr, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("期望数组，实际类型: %T", val)
	}

	result := make([]any, len(arr))
	for i, elem := range arr {
		converted, err := convertABIParam(elem, elemType)
		if err != nil {
			return nil, fmt.Errorf("数组元素 %d: %w", i, err)
		}
		result[i] = converted
	}

	// 将 []any 转为具体类型的切片，以便 abi.Pack 正确处理
	return typedSlice(result, elemType)
}

// convertArrayParam 将 JSON 数组转换为固定长度的 Go 切片。
func convertArrayParam(val any, elemType abi.Type, size int) (any, error) {
	arr, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("期望数组，实际类型: %T", val)
	}
	if len(arr) != size {
		return nil, fmt.Errorf("数组长度不匹配: 传入 %d, 期望 %d", len(arr), size)
	}
	return convertSlice(val, elemType)
}

// convertTuple 将 JSON 数组/对象转换为 Go 结构体（以 []any 形式传递给 abi.Pack）。
func convertTuple(val any, elems []*abi.Type, names []string) ([]any, error) {
	// 先尝试按数组解析（按位置）
	if arr, ok := val.([]any); ok {
		if len(arr) != len(elems) {
			return nil, fmt.Errorf("元组长度不匹配: 传入 %d, 期望 %d", len(arr), len(elems))
		}
		result := make([]any, len(arr))
		for i, elem := range arr {
			converted, err := convertABIParam(elem, *elems[i])
			if err != nil {
				return nil, fmt.Errorf("元组字段 %d (%s): %w", i, names[i], err)
			}
			result[i] = converted
		}
		return result, nil
	}

	// 尝试按对象解析（按字段名）
	if m, ok := val.(map[string]any); ok {
		result := make([]any, len(elems))
		for i, name := range names {
			fieldVal, exists := m[name]
			if !exists {
				return nil, fmt.Errorf("元组缺少字段: %s", name)
			}
			converted, err := convertABIParam(fieldVal, *elems[i])
			if err != nil {
				return nil, fmt.Errorf("元组字段 %s: %w", name, err)
			}
			result[i] = converted
		}
		return result, nil
	}

	return nil, fmt.Errorf("期望数组或对象，实际类型: %T", val)
}

// typedSlice 将 []any 转为具体的类型化切片，使 abi.Arguments.Pack 能正确处理。
func typedSlice(items []any, elemType abi.Type) (any, error) {
	switch elemType.T {
	case abi.AddressTy:
		result := make([]common.Address, len(items))
		for i, item := range items {
			addr, ok := item.(common.Address)
			if !ok {
				return nil, fmt.Errorf("数组元素 %d 不是 Address 类型", i)
			}
			result[i] = addr
		}
		return result, nil
	case abi.UintTy, abi.IntTy:
		result := make([]*big.Int, len(items))
		for i, item := range items {
			bi, ok := item.(*big.Int)
			if !ok {
				return nil, fmt.Errorf("数组元素 %d 不是 *big.Int 类型", i)
			}
			result[i] = bi
		}
		return result, nil
	case abi.BoolTy:
		result := make([]bool, len(items))
		for i, item := range items {
			b, ok := item.(bool)
			if !ok {
				return nil, fmt.Errorf("数组元素 %d 不是 bool 类型", i)
			}
			result[i] = b
		}
		return result, nil
	case abi.StringTy:
		result := make([]string, len(items))
		for i, item := range items {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("数组元素 %d 不是 string 类型", i)
			}
			result[i] = s
		}
		return result, nil
	case abi.BytesTy:
		result := make([][]byte, len(items))
		for i, item := range items {
			b, ok := item.([]byte)
			if !ok {
				return nil, fmt.Errorf("数组元素 %d 不是 []byte 类型", i)
			}
			result[i] = b
		}
		return result, nil
	default:
		// 复杂元素类型保持 []any
		return items, nil
	}
}

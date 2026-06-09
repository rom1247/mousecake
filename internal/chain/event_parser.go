// Package chain 提供 RPC 节点管理基础设施，包括多节点池、熔断器、限流、健康检查。
package chain

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/mousecake-go/mousecake-go/internal/chain/contract/launchpad"
)

// ParsedEvent 解析后的合约事件，包含事件名和字段值映射。
type ParsedEvent struct {
	// Name 事件名称（如 "Deposited"、"Harvested"）。
	Name string
	// Fields 事件字段的键值映射，big.Int 已转为十进制字符串，Address 已转为十六进制字符串。
	Fields map[string]any
	// Address 触发事件的合约地址。
	Address common.Address
	// LogIndex 日志在区块中的索引。
	LogIndex uint
}

// EventParser 使用多个合约 ABI 解析 receipt logs 中的事件。
// 仅供开发环境使用，通过注册的合约 ABI 遍历匹配事件签名。
type EventParser struct {
	abis []*abiSource
}

// abiSource 封装合约名称和解析后的 ABI。
type abiSource struct {
	name string
	abi  abi.ABI
}

// NewEventParser 创建事件解析器，注册 MousePadByTier、MousePadByTierDeployer、MouseTier 三个合约 ABI。
func NewEventParser() *EventParser {
	padABI, err := launchpad.MousePadByTierMetaData.ParseABI()
	if err != nil {
		panic(fmt.Sprintf("解析 MousePadByTier ABI: %v", err))
	}
	deployerABI, err := launchpad.MousePadByTierDeployerMetaData.ParseABI()
	if err != nil {
		panic(fmt.Sprintf("解析 MousePadByTierDeployer ABI: %v", err))
	}
	tierABI, err := launchpad.MouseTierMetaData.ParseABI()
	if err != nil {
		panic(fmt.Sprintf("解析 MouseTier ABI: %v", err))
	}

	return &EventParser{
		abis: []*abiSource{
			{name: "MousePadByTier", abi: *padABI},
			{name: "MousePadByTierDeployer", abi: *deployerABI},
			{name: "MouseTier", abi: *tierABI},
		},
	}
}

// ParseLogs 解析 receipt logs 中的所有事件，返回解析后的事件列表。
// 无法识别的日志（topic 不匹配任何已注册合约的事件签名）会被静默跳过。
func (p *EventParser) ParseLogs(logs []*types.Log) []ParsedEvent {
	var events []ParsedEvent
	for _, log := range logs {
		if len(log.Topics) == 0 {
			continue
		}
		topic := log.Topics[0]
		for _, src := range p.abis {
			event, ok := src.abi.Events[topic.Hex()]
			if !ok {
				// 尝试按 ID 匹配
				for _, e := range src.abi.Events {
					if e.ID == topic {
						event = e
						ok = true
						break
					}
				}
				if !ok {
					continue
				}
			}

			fields, err := p.decodeEventFields(event, log)
			if err != nil {
				continue
			}

			events = append(events, ParsedEvent{
				Name:     event.Name,
				Fields:   fields,
				Address:  log.Address,
				LogIndex: log.Index,
			})
			break
		}
	}
	return events
}

// decodeEventFields 解码单条日志中的事件字段，返回字段名到值的映射。
func (p *EventParser) decodeEventFields(event abi.Event, log *types.Log) (map[string]any, error) {
	fields := make(map[string]any)

	// 分离 indexed 和 non-indexed 参数
	var nonIndexedArgs abi.Arguments
	var indexedArgs abi.Arguments
	for _, input := range event.Inputs {
		if input.Indexed {
			indexedArgs = append(indexedArgs, input)
		} else {
			nonIndexedArgs = append(nonIndexedArgs, input)
		}
	}

	// 解码 non-indexed 字段（从 log.Data）
	if len(log.Data) > 0 && len(nonIndexedArgs) > 0 {
		values, err := nonIndexedArgs.Unpack(log.Data)
		if err != nil {
			return nil, fmt.Errorf("解码 non-indexed 字段: %w", err)
		}
		idx := 0
		for _, input := range event.Inputs {
			if !input.Indexed && idx < len(values) {
				fields[input.Name] = convertValue(values[idx])
				idx++
			}
		}
	}

	// 解码 indexed 字段（从 log.Topics[1:]）
	topicIdx := 1
	for _, input := range event.Inputs {
		if input.Indexed && topicIdx < len(log.Topics) {
			val, err := decodeIndexedTopic(input.Type, log.Topics[topicIdx])
			if err != nil {
				return nil, fmt.Errorf("解码 indexed 字段 %s: %w", input.Name, err)
			}
			fields[input.Name] = convertValue(val)
			topicIdx++
		}
	}

	return fields, nil
}

// decodeIndexedTopic 从 topic 中解码 indexed 参数值。
func decodeIndexedTopic(typ abi.Type, topic common.Hash) (any, error) {
	switch typ.T {
	case abi.AddressTy:
		return common.BytesToAddress(topic[12:]), nil
	case abi.UintTy, abi.IntTy:
		return new(big.Int).SetBytes(topic[:]), nil
	case abi.BoolTy:
		return topic[31] == 1, nil
	case abi.BytesTy, abi.FixedBytesTy:
		return topic[:], nil
	default:
		return topic[:], nil
	}
}

// convertValue 将 ABI 解码值转换为 JSON 安全的表示。
// *big.Int → 十进制字符串，common.Address → 十六进制字符串，[]common.Address → []string。
func convertValue(v any) any {
	if v == nil {
		return nil
	}

	// 处理指针类型
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil
		}
		v = rv.Elem().Interface()
	}

	switch val := v.(type) {
	case *big.Int:
		return val.String()
	case big.Int:
		return val.String()
	case common.Address:
		return val.Hex()
	case []common.Address:
		result := make([]string, len(val))
		for i, addr := range val {
			result[i] = addr.Hex()
		}
		return result
	case bool:
		return val
	case string:
		return val
	case []byte:
		return common.Bytes2Hex(val)
	case common.Hash:
		return val.Hex()
	default:
		// 尝试处理切片类型
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Slice {
			result := make([]any, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				result[i] = convertValue(rv.Index(i).Interface())
			}
			return result
		}
		return fmt.Sprintf("%v", val)
	}
}

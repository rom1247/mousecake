package sync

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

// TestExtractEventNameFromLog 测试从链上日志中提取事件名称（topic[0] 前 10 位）。
func TestExtractEventNameFromLog(t *testing.T) {
	t.Parallel()

	// 构造一个完整的 topic hash，对应 Transfer(address,address,uint256) 的签名
	transferTopic := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	// common.Hash 始终为 32 字节，Hex() 返回 "0x" + 64 字符，长度恒为 66 >= 10。
	// 因此不存在 "topic hex 长度不足 10" 的情况，所有有效的 common.Hash 都会走截取逻辑。
	// 零值 Hash 的前 10 字符为 "0x00000000"。
	zeroHash := common.Hash{}

	tests := []struct {
		name string
		log  types.Log
		want string
	}{
		{
			name: "正常事件 — 提取 topic[0] 前 10 字符",
			log: types.Log{
				Topics: []common.Hash{transferTopic},
			},
			want: "0xddf252ad",
		},
		{
			name: "无 topic — 返回 Unknown",
			log: types.Log{
				Topics: []common.Hash{},
			},
			want: "Unknown",
		},
		{
			name: "零值 topic — 仍然截取前 10 字符",
			log: types.Log{
				Topics: []common.Hash{zeroHash},
			},
			want: "0x00000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := extractEventNameFromLog(tt.log)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestTopicsToHex 测试将 topics 转换为十六进制字符串列表。
func TestTopicsToHex(t *testing.T) {
	t.Parallel()

	topic1 := common.HexToHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	topic2 := common.HexToHash("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	tests := []struct {
		name   string
		topics []common.Hash
		want   []string
	}{
		{
			name:   "空 topics",
			topics: nil,
			want:   []string{},
		},
		{
			name:   "单个 topic",
			topics: []common.Hash{topic1},
			want:   []string{topic1.Hex()},
		},
		{
			name:   "多个 topics",
			topics: []common.Hash{topic1, topic2},
			want:   []string{topic1.Hex(), topic2.Hex()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := topicsToHex(tt.topics)
			assert.Equal(t, tt.want, got)
		})
	}
}

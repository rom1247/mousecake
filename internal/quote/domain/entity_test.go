package domain

import (
	"testing"
)

func TestNewSwapRecord(t *testing.T) {
	tests := []struct {
		name    string
		opts    NewSwapRecordOpts
		wantErr error
	}{
		{
			name: "正常创建 SwapRecord",
			opts: NewSwapRecordOpts{
				Provider:        "okx",
				ChainID:         1,
				FromToken:       "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
				ToToken:         "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				FromAmount:      "1000000",
				ToAmount:        "2000000",
				SlippagePercent: 0.5,
				SwapMode:        SwapModeExactIn,
			},
			wantErr: nil,
		},
		{
			name: "provider 为空",
			opts: NewSwapRecordOpts{
				Provider:   "",
				ChainID:    1,
				FromToken:  "0xA",
				ToToken:    "0xB",
				FromAmount: "1000",
				ToAmount:   "2000",
				SwapMode:   SwapModeExactIn,
			},
			wantErr: ErrProviderEmpty,
		},
		{
			name: "from_token 为空",
			opts: NewSwapRecordOpts{
				Provider:   "okx",
				ChainID:    1,
				FromToken:  "",
				ToToken:    "0xB",
				FromAmount: "1000",
				ToAmount:   "2000",
				SwapMode:   SwapModeExactIn,
			},
			wantErr: ErrFromTokenEmpty,
		},
		{
			name: "to_token 为空",
			opts: NewSwapRecordOpts{
				Provider:   "okx",
				ChainID:    1,
				FromToken:  "0xA",
				ToToken:    "",
				FromAmount: "1000",
				ToAmount:   "2000",
				SwapMode:   SwapModeExactIn,
			},
			wantErr: ErrToTokenEmpty,
		},
		{
			name: "from_amount 为空",
			opts: NewSwapRecordOpts{
				Provider:   "okx",
				ChainID:    1,
				FromToken:  "0xA",
				ToToken:    "0xB",
				FromAmount: "",
				ToAmount:   "2000",
				SwapMode:   SwapModeExactIn,
			},
			wantErr: ErrFromAmountEmpty,
		},
		{
			name: "to_amount 为空",
			opts: NewSwapRecordOpts{
				Provider:   "okx",
				ChainID:    1,
				FromToken:  "0xA",
				ToToken:    "0xB",
				FromAmount: "1000",
				ToAmount:   "",
				SwapMode:   SwapModeExactIn,
			},
			wantErr: ErrToAmountEmpty,
		},
		{
			name: "无效 swap_mode",
			opts: NewSwapRecordOpts{
				Provider:   "okx",
				ChainID:    1,
				FromToken:  "0xA",
				ToToken:    "0xB",
				FromAmount: "1000",
				ToAmount:   "2000",
				SwapMode:   "invalid",
			},
			wantErr: ErrInvalidSwapMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record, err := NewSwapRecord(1, tt.opts)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("期望错误 %v, 得到 nil", tt.wantErr)
				}
				if err != tt.wantErr {
					t.Fatalf("期望错误 %v, 得到 %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("不期望错误, 得到 %v", err)
			}
			if record.Status != SwapStatusPending {
				t.Errorf("期望 status=%s, 得到 %s", SwapStatusPending, record.Status)
			}
			if record.Provider != tt.opts.Provider {
				t.Errorf("期望 provider=%s, 得到 %s", tt.opts.Provider, record.Provider)
			}
			if record.ID == 0 {
				t.Error("ID 不应为 0")
			}
		})
	}
}

func TestSwapRecord_Submit(t *testing.T) {
	record := mustNewSwapRecord(t)

	t.Run("正常状态转换 pending → submitted", func(t *testing.T) {
		txHash := "0xabc123def456abc123def456abc123def456abc123def456abc123def456abcd"
		err := record.Submit(txHash)
		if err != nil {
			t.Fatalf("不期望错误, 得到 %v", err)
		}
		if record.Status != SwapStatusSubmitted {
			t.Errorf("期望 status=%s, 得到 %s", SwapStatusSubmitted, record.Status)
		}
		if record.TxHash != txHash {
			t.Errorf("期望 tx_hash=%s, 得到 %s", txHash, record.TxHash)
		}
	})

	t.Run("非法状态转换 submitted → submitted", func(t *testing.T) {
		err := record.Submit("0xdef")
		if err != ErrAlreadySubmitted {
			t.Fatalf("期望 ErrAlreadySubmitted, 得到 %v", err)
		}
	})
}

func TestReconstructSwapRecord(t *testing.T) {
	record := ReconstructSwapRecord(SwapRecordSnapshot{
		ID:              17369650207862784,
		Provider:        "okx",
		ChainID:         1,
		FromToken:       "0xA",
		ToToken:         "0xB",
		FromAmount:      "1000",
		ToAmount:        "2000",
		SlippagePercent: 0.5,
		SwapMode:        SwapModeExactIn,
		Status:          SwapStatusSubmitted,
		TxHash:          "0xabc...",
	})
	if record.ID != 17369650207862784 {
		t.Errorf("期望 ID=17369650207862784, 得到 %d", record.ID)
	}
	if record.Status != SwapStatusSubmitted {
		t.Errorf("期望 status=%s, 得到 %s", SwapStatusSubmitted, record.Status)
	}
	if record.TxHash != "0xabc..." {
		t.Errorf("期望 tx_hash=0xabc..., 得到 %s", record.TxHash)
	}
}

func mustNewSwapRecord(t *testing.T) *SwapRecord {
	t.Helper()
	record, err := NewSwapRecord(1, NewSwapRecordOpts{
		Provider:        "okx",
		ChainID:         1,
		FromToken:       "0xA",
		ToToken:         "0xB",
		FromAmount:      "1000",
		ToAmount:        "2000",
		SlippagePercent: 0.5,
		SwapMode:        SwapModeExactIn,
	})
	if err != nil {
		t.Fatalf("创建 SwapRecord 失败: %v", err)
	}
	return record
}

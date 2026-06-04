package domain

import (
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"
)

func TestSale_Phase(t *testing.T) {
	now := time.Now()
	sale := ReconstructSale(1, "0xsale", SaleDeployed, 1, "0xdep", "0xowner", "0xraise", "0xoffer", "0xtier", 100, 200, 0, false, 1000, now, now, nil)

	tests := []struct {
		name          string
		currentBlock  int64
		expectedPhase SalePhase
	}{
		{"配置阶段", 50, SalePhaseConfiguring},
		{"募资阶段-起始块", 100, SalePhaseLive},
		{"募资阶段-中间块", 150, SalePhaseLive},
		{"募资阶段-结束块", 200, SalePhaseLive},
		{"结算阶段", 201, SalePhaseEnded},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sale.Phase(tt.currentBlock)
			if got != tt.expectedPhase {
				t.Errorf("Phase(%d) = %s, want %s", tt.currentBlock, got, tt.expectedPhase)
			}
		})
	}
}

func TestPool_IsConfigured(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		raising  *big.Int
		offering *big.Int
		expected bool
	}{
		{"配置完整", big.NewInt(1000), big.NewInt(5000), true},
		{"募资为零", big.NewInt(0), big.NewInt(5000), false},
		{"发售为零", big.NewInt(1000), big.NewInt(0), false},
		{"都为零", big.NewInt(0), big.NewInt(0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := ReconstructPool(1, 1, 0, tt.raising, tt.offering, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), false, false, 0, 0, 0, 0, now, now)
			if got := pool.IsConfigured(); got != tt.expected {
				t.Errorf("IsConfigured() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestVestingSchedule_Remaining(t *testing.T) {
	now := time.Now()
	schedule := ReconstructVestingSchedule(1, 1, 0, 1, "0xabc", big.NewInt(1000), big.NewInt(300), now, now)

	remaining := schedule.Remaining()
	if remaining.Cmp(big.NewInt(700)) != 0 {
		t.Errorf("Remaining() = %s, want 700", remaining.String())
	}
}

func TestPrepareTxStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		from    PrepareTxStatus
		to      PrepareTxStatus
		isValid bool
	}{
		{PrepareTxPending, PrepareTxBroadcast, true},
		{PrepareTxPending, PrepareTxSigned, true},
		{PrepareTxPending, PrepareTxExpired, true},
		{PrepareTxBroadcast, PrepareTxConfirmed, true},
		{PrepareTxBroadcast, PrepareTxReverted, true},
		{PrepareTxBroadcast, PrepareTxFailed, true},
		{PrepareTxConfirmed, PrepareTxBroadcast, false},
		{PrepareTxConfirmed, PrepareTxReverted, false},
		{PrepareTxReverted, PrepareTxConfirmed, false},
		{PrepareTxExpired, PrepareTxPending, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.from)+"→"+string(tt.to), func(t *testing.T) {
			got := tt.from.CanTransitionTo(tt.to)
			if got != tt.isValid {
				t.Errorf("CanTransitionTo(%s→%s) = %v, want %v", tt.from, tt.to, got, tt.isValid)
			}
		})
	}
}

func TestPrepareTx_Transition(t *testing.T) {
	tx := &PrepareTx{
		Status:    PrepareTxPending,
		UpdatedAt: time.Now(),
	}

	if err := tx.Transition(PrepareTxBroadcast); err != nil {
		t.Errorf("合法转换不应报错: %v", err)
	}
	if tx.Status != PrepareTxBroadcast {
		t.Errorf("状态应为 broadcast, got %s", tx.Status)
	}

	if err := tx.Transition(PrepareTxPending); err == nil {
		t.Error("非法转换应报错")
	}
}

func TestPrepareTxOperationType(t *testing.T) {
	tests := []struct {
		op      PrepareTxOperationType
		isAdmin bool
		isUser  bool
		isValid bool
	}{
		{OpCreateSale, true, false, true},
		{OpSetPool, true, false, true},
		{OpDeposit, false, true, true},
		{OpHarvest, false, true, true},
		{OpRelease, false, true, true},
		{PrepareTxOperationType("invalid"), false, false, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.op), func(t *testing.T) {
			if got := tt.op.IsAdminOperation(); got != tt.isAdmin {
				t.Errorf("IsAdminOperation() = %v, want %v", got, tt.isAdmin)
			}
			if got := tt.op.IsUserOperation(); got != tt.isUser {
				t.Errorf("IsUserOperation() = %v, want %v", got, tt.isUser)
			}
			if got := IsValidOperationType(tt.op); got != tt.isValid {
				t.Errorf("IsValidOperationType() = %v, want %v", got, tt.isValid)
			}
		})
	}
}

func TestSaleMeta_IsPublic(t *testing.T) {
	now := time.Now()
	pub := ReconstructSaleMeta(1, 1, "Title", "", "", "", "", "", "public", 0, now, now)
	if !pub.IsPublic() {
		t.Error("public 的 IsPublic() 应为 true")
	}

	hidden := ReconstructSaleMeta(1, 1, "Title", "", "", "", "", "", "hidden", 0, now, now)
	if hidden.IsPublic() {
		t.Error("hidden 的 IsPublic() 应为 false")
	}
}

func TestPrepareTx_IsExpired(t *testing.T) {
	tx := &PrepareTx{
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if !tx.IsExpired() {
		t.Error("已过期的 PrepareTx 应返回 true")
	}

	tx.ExpiresAt = time.Now().Add(1 * time.Hour)
	if tx.IsExpired() {
		t.Error("未过期的 PrepareTx 应返回 false")
	}
}

func TestErrSaleNotDeployed_可通过ErrorsIs判断(t *testing.T) {
	err := ErrSaleNotDeployed
	if !errors.Is(err, ErrSaleNotDeployed) {
		t.Error("errors.Is(ErrSaleNotDeployed, ErrSaleNotDeployed) 应返回 true")
	}

	// 包装后仍然能通过 errors.Is 判断
	wrapped := fmt.Errorf("业务校验: %w", err)
	if !errors.Is(wrapped, ErrSaleNotDeployed) {
		t.Error("包装后的错误仍应通过 errors.Is 匹配 ErrSaleNotDeployed")
	}
}

func TestReconstructSale_支持Status参数(t *testing.T) {
	now := time.Now()
	sale := ReconstructSale(
		1, "0xsale", SaleDeploying, 1,
		"0xdep", "0xowner", "0xraise", "0xoffer", "0xtier",
		100, 200, 0, false, 1000,
		now, now, nil,
	)

	if sale.Status != SaleDeploying {
		t.Errorf("Status = %s, want %s", sale.Status, SaleDeploying)
	}
	if sale.ID != 1 {
		t.Errorf("ID = %d, want 1", sale.ID)
	}
	if sale.ContractAddress != "0xsale" {
		t.Errorf("ContractAddress = %s, want 0xsale", sale.ContractAddress)
	}
}

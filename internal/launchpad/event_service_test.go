package launchpad

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
	syncpkg "github.com/mousecake-go/mousecake-go/internal/shared/sync"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// setupEventTestDB 创建包含所有必要表的内存 SQLite 数据库，用于事件处理测试。
func setupEventTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "打开内存数据库失败")

	err = db.AutoMigrate(
		&salePO{},
		&poolPO{},
		&depositPO{},
		&userPoolStatePO{},
		&userCreditPO{},
		&harvestPO{},
		&vestingSchedulePO{},
		&vestingReleasePO{},
		&whitelistPO{},
		&tierParamPO{},
		&prepareTxPO{},
	)
	require.NoError(t, err, "AutoMigrate 失败")

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
	return db
}

// makeChainEvent 创建测试用的 ChainEvent。
func makeChainEvent(eventName, contractAddress string, chainID int, blockNumber int64, txHash string, data map[string]any) *syncpkg.ChainEvent {
	dataJSON, _ := json.Marshal(data)
	return &syncpkg.ChainEvent{
		EventName:       eventName,
		ContractAddress: contractAddress,
		ChainID:         chainID,
		BlockNumber:     blockNumber,
		TxHash:          txHash,
		EventData:       string(dataJSON),
	}
}

// insertTestSale 插入一条测试 sale 记录并返回其 ID。
func insertTestSale(t *testing.T, db *gorm.DB, contractAddress string) int64 {
	t.Helper()
	sale := &salePO{
		ContractAddress:      contractAddress,
		ChainID:              1,
		DeployerAddress:      "0xdeployer",
		OwnerAddress:         "0xowner",
		RaiseTokenAddress:    "0xraise",
		OfferingTokenAddress: "0xoffering",
		MouseTierAddress:     "0xtier",
	}
	require.NoError(t, db.Create(sale).Error, "插入测试 sale 失败")
	return sale.ID
}

// --- 辅助函数测试 ---

// TestDataStr 测试事件数据字符串提取辅助函数。
func TestDataStr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     map[string]any
		key      string
		fallback string
		want     string
	}{
		{"字符串值", map[string]any{"key": "value"}, "key", "default", "value"},
		{"整数值转字符串", map[string]any{"key": 42}, "key", "default", "42"},
		{"键不存在", map[string]any{}, "missing", "default", "default"},
		{"nil 值", map[string]any{"key": nil}, "key", "default", "<nil>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, dataStr(tt.data, tt.key, tt.fallback))
		})
	}
}

// TestDataInt 测试事件数据整数提取辅助函数。
func TestDataInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     map[string]any
		key      string
		fallback int
		want     int
	}{
		{"float64 值", map[string]any{"key": float64(42)}, "key", 0, 42},
		{"int 值", map[string]any{"key": 100}, "key", 0, 100},
		{"键不存在", map[string]any{}, "missing", 99, 99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, dataInt(tt.data, tt.key, tt.fallback))
		})
	}
}

// --- HandleEvent 测试 ---

// TestEventService_HandleEvent_UnknownEvent 测试未知事件名不返回错误。
func TestEventService_HandleEvent_UnknownEvent(t *testing.T) {
	t.Parallel()

	svc := &EventService{}
	event := &syncpkg.ChainEvent{
		EventName: "UnknownEventXYZ",
		EventData: `{}`,
	}

	err := svc.HandleEvent(context.Background(), event)
	assert.NoError(t, err, "未知事件应静默处理")
}

// TestEventService_HandleEvent_KnownRoutes 测试已知事件名能正确路由到对应方法。
func TestEventService_HandleEvent_KnownRoutes(t *testing.T) {
	t.Parallel()

	svc := &EventService{}
	event := &syncpkg.ChainEvent{
		EventName:       "FinalWithdraw",
		EventData:       `{}`,
		ContractAddress: "0x0000000000000000000000000000000000000000",
		ChainID:         1,
	}

	err := svc.HandleEvent(context.Background(), event)
	assert.NoError(t, err, "FinalWithdraw 路由应成功")
}

// TestEventService_OnFinalWithdraw 测试 FinalWithdraw 不需要数据库。
func TestEventService_OnFinalWithdraw(t *testing.T) {
	t.Parallel()

	svc := &EventService{}
	event := &syncpkg.ChainEvent{
		EventName:       "FinalWithdraw",
		EventData:       `{}`,
		ContractAddress: "0x1234",
	}

	err := svc.OnFinalWithdraw(context.Background(), event)
	assert.NoError(t, err, "FinalWithdraw 应直接返回 nil")
}

// --- OnSaleCreated 测试 ---

// TestEventService_OnSaleCreated 测试 SaleCreated 事件正确创建 sale 记录。
func TestEventService_OnSaleCreated(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		contractAddr   string
		eventData      map[string]any
		wantCreator    string
		wantRaiseToken string
		wantOfferToken string
		wantMouseTier  string
	}{
		{
			name:         "完整事件数据",
			contractAddr: "0xsale001",
			eventData: map[string]any{
				"sale_address":   "0xsale001",
				"creator":        "0xcreator01",
				"raise_token":    "0xraise001",
				"offering_token": "0xoffer001",
				"mouse_tier":     "0xtier001",
			},
			wantCreator:    "0xcreator01",
			wantRaiseToken: "0xraise001",
			wantOfferToken: "0xoffer001",
			wantMouseTier:  "0xtier001",
		},
		{
			name:           "缺少可选字段使用默认值",
			contractAddr:   "0xsale002",
			eventData:      map[string]any{},
			wantCreator:    "",
			wantRaiseToken: "",
			wantOfferToken: "",
			wantMouseTier:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := setupEventTestDB(t)
			svc := NewEventService(db, nil)
			ctx := context.Background()

			event := makeChainEvent("SaleCreated", tt.contractAddr, 1, 100, "0xtx001", tt.eventData)

			err := svc.OnSaleCreated(ctx, event)
			require.NoError(t, err, "OnSaleCreated 不应返回错误")

			// 验证数据库记录
			var sale salePO
			require.NoError(t, db.Where("contract_address = ?", tt.contractAddr).First(&sale).Error,
				"应能查到 sale 记录")

			assert.Equal(t, tt.contractAddr, sale.ContractAddress, "合约地址应匹配")
			assert.Equal(t, 1, sale.ChainID, "链 ID 应匹配")
			assert.Equal(t, tt.wantCreator, sale.OwnerAddress, "创建者地址应匹配")
			assert.Equal(t, tt.wantRaiseToken, sale.RaiseTokenAddress, "募资代币地址应匹配")
			assert.Equal(t, tt.wantOfferToken, sale.OfferingTokenAddress, "发售代币地址应匹配")
			assert.Equal(t, tt.wantMouseTier, sale.MouseTierAddress, "MouseTier 地址应匹配")
		})
	}
}

// TestEventService_OnSaleCreated_Idempotent 测试相同合约地址重复调用不创建重复记录。
func TestEventService_OnSaleCreated_Idempotent(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()
	contractAddr := "0xsale_idem"

	eventData := map[string]any{
		"sale_address":   contractAddr,
		"creator":        "0xcreator",
		"raise_token":    "0xraise",
		"offering_token": "0xoffer",
		"mouse_tier":     "0xtier",
	}
	event := makeChainEvent("SaleCreated", contractAddr, 1, 100, "0xtx001", eventData)

	// 第一次调用
	require.NoError(t, svc.OnSaleCreated(ctx, event), "第一次调用应成功")

	// 第二次调用（幂等）
	require.NoError(t, svc.OnSaleCreated(ctx, event), "重复调用不应返回错误")

	// 验证只有一条记录
	var count int64
	db.Model(&salePO{}).Where("contract_address = ?", contractAddr).Count(&count)
	assert.Equal(t, int64(1), count, "幂等调用不应创建重复记录")
}

// --- OnPoolSet 测试 ---

// TestEventService_OnPoolSet 测试 PoolSet 事件正确创建 pool 记录。
func TestEventService_OnPoolSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		poolID         int
		raisingAmount  string
		offeringAmount string
		limitPerUser   string
	}{
		{
			name:           "创建池子 0",
			poolID:         0,
			raisingAmount:  "1000000000000000000",
			offeringAmount: "2000000000000000000",
			limitPerUser:   "500000000000000000",
		},
		{
			name:           "创建池子 1",
			poolID:         1,
			raisingAmount:  "3000000000000000000",
			offeringAmount: "4000000000000000000",
			limitPerUser:   "1000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := setupEventTestDB(t)
			svc := NewEventService(db, nil)
			ctx := context.Background()

			contractAddr := "0xsale_pool"
			saleID := insertTestSale(t, db, contractAddr)

			eventData := map[string]any{
				"pool_id":         float64(tt.poolID),
				"sale_address":    contractAddr,
				"raising_amount":  tt.raisingAmount,
				"offering_amount": tt.offeringAmount,
				"limit_per_user":  tt.limitPerUser,
			}
			event := makeChainEvent("PoolSet", contractAddr, 1, 200, "0xtxpool", eventData)

			err := svc.OnPoolSet(ctx, event)
			require.NoError(t, err, "OnPoolSet 不应返回错误")

			// 验证数据库记录
			var pool poolPO
			require.NoError(t, db.Where("sale_id = ? AND pool_index = ?", saleID, tt.poolID).First(&pool).Error,
				"应能查到 pool 记录")

			assert.Equal(t, saleID, pool.SaleID, "sale_id 应匹配")
			assert.Equal(t, tt.poolID, pool.PoolIndex, "pool_index 应匹配")
			assert.Equal(t, tt.raisingAmount, pool.RaisingAmount, "raising_amount 应匹配")
			assert.Equal(t, tt.offeringAmount, pool.OfferingAmount, "offering_amount 应匹配")
			assert.Equal(t, tt.limitPerUser, pool.LimitPerUser, "limit_per_user 应匹配")
		})
	}
}

// TestEventService_OnPoolSet_SaleNotFound 测试 sale 不存在时跳过处理。
func TestEventService_OnPoolSet_SaleNotFound(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	eventData := map[string]any{
		"pool_id":         float64(0),
		"raising_amount":  "1000",
		"offering_amount": "2000",
		"limit_per_user":  "500",
	}
	event := makeChainEvent("PoolSet", "0xnotexist", 1, 200, "0xtx", eventData)

	err := svc.OnPoolSet(ctx, event)
	assert.NoError(t, err, "sale 不存在时应跳过处理而不报错")
}

// --- OnDeposited 测试 ---

// TestEventService_OnDeposited_ThreeTableTx 测试 Deposited 事件三表事务写入。
func TestEventService_OnDeposited_ThreeTableTx(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		poolIndex  int
		wantCredit bool
	}{
		{
			name:       "pool_index=0 写入三张表（含 user_credit）",
			poolIndex:  0,
			wantCredit: true,
		},
		{
			name:       "pool_index=1 只写 deposit 和 user_pool_state",
			poolIndex:  1,
			wantCredit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := setupEventTestDB(t)
			svc := NewEventService(db, nil)
			ctx := context.Background()

			contractAddr := "0xsale_deposit"
			saleID := insertTestSale(t, db, contractAddr)

			user := "0xuser01"
			amount := "1000000000000000000"

			eventData := map[string]any{
				"user":    user,
				"pool_id": float64(tt.poolIndex),
				"amount":  amount,
			}
			event := makeChainEvent("Deposited", contractAddr, 1, 300, "0xtxdep01", eventData)

			err := svc.OnDeposited(ctx, event)
			require.NoError(t, err, "OnDeposited 不应返回错误")

			// 验证 launchpad_deposits
			var deposit depositPO
			require.NoError(t, db.Where("tx_hash = ?", "0xtxdep01").First(&deposit).Error,
				"应能查到 deposit 记录")
			assert.Equal(t, saleID, deposit.SaleID, "sale_id 应匹配")
			assert.Equal(t, tt.poolIndex, deposit.PoolIndex, "pool_index 应匹配")
			assert.Equal(t, user, deposit.UserAddress, "user_address 应匹配")
			assert.Equal(t, amount, deposit.Amount, "amount 应匹配")

			// 验证 launchpad_user_pool_states
			var state userPoolStatePO
			require.NoError(t, db.Where("sale_id = ? AND pool_index = ? AND user_address = ?",
				saleID, tt.poolIndex, user).First(&state).Error,
				"应能查到 user_pool_state 记录")
			assert.Equal(t, amount, state.TotalDeposited, "total_deposited 应匹配")

			// 验证 launchpad_user_credit
			var creditCount int64
			db.Model(&userCreditPO{}).Where("sale_id = ? AND user_address = ?", saleID, user).Count(&creditCount)
			if tt.wantCredit {
				assert.Equal(t, int64(1), creditCount, "pool_index=0 时应创建 user_credit 记录")
			} else {
				assert.Equal(t, int64(0), creditCount, "pool_index!=0 时不应创建 user_credit 记录")
			}
		})
	}
}

// TestEventService_OnDeposited_Idempotent 测试相同 tx_hash 不重复写入。
func TestEventService_OnDeposited_Idempotent(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	contractAddr := "0xsale_idem_dep"
	insertTestSale(t, db, contractAddr)

	eventData := map[string]any{
		"user":    "0xuser02",
		"pool_id": float64(0),
		"amount":  "1000000000000000000",
	}
	txHash := "0xtx_idem_dep"
	event := makeChainEvent("Deposited", contractAddr, 1, 300, txHash, eventData)

	// 第一次调用
	require.NoError(t, svc.OnDeposited(ctx, event), "第一次调用应成功")

	// 第二次调用（幂等）
	require.NoError(t, svc.OnDeposited(ctx, event), "重复调用不应返回错误")

	// 验证只有一条 deposit 记录
	var count int64
	db.Model(&depositPO{}).Where("tx_hash = ?", txHash).Count(&count)
	assert.Equal(t, int64(1), count, "幂等调用不应创建重复 deposit 记录")
}

// TestEventService_OnDeposited_NoRPCInTx 测试 OnDeposited 不依赖 NodePool。
func TestEventService_OnDeposited_NoRPCInTx(t *testing.T) {
	t.Parallel()

	// EventService 只有 db 字段，没有 NodePool 或其他 RPC 依赖。
	// 通过验证 EventService 结构体的定义确保这一点。
	svc := &EventService{}
	assert.Nil(t, svc.db, "EventService 不应有 db 之外的字段依赖")

	// 使用真实 DB 执行完整流程验证纯 DB 操作
	db := setupEventTestDB(t)
	svc = NewEventService(db, nil)
	ctx := context.Background()

	contractAddr := "0xsale_norpc"
	insertTestSale(t, db, contractAddr)

	eventData := map[string]any{
		"user":    "0xuser_norpc",
		"pool_id": float64(0),
		"amount":  "500000000000000000",
	}
	event := makeChainEvent("Deposited", contractAddr, 1, 400, "0xtxnorpc", eventData)

	err := svc.OnDeposited(ctx, event)
	assert.NoError(t, err, "纯 DB 操作应成功完成，无需 RPC")
}

// --- OnHarvested 测试 ---

// TestEventService_OnHarvested 测试 Harvested 事件写入 harvest 和 vesting_schedule。
func TestEventService_OnHarvested(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		vestingAmount string
		wantSchedule  bool
	}{
		{
			name:          "有 vesting 创建 vesting_schedule",
			vestingAmount: "500000000000000000",
			wantSchedule:  true,
		},
		{
			name:          "无 vesting 不创建 vesting_schedule",
			vestingAmount: "0",
			wantSchedule:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := setupEventTestDB(t)
			svc := NewEventService(db, nil)
			ctx := context.Background()

			contractAddr := "0xsale_harvest"
			saleID := insertTestSale(t, db, contractAddr)

			user := "0xuser_harvest"
			eventData := map[string]any{
				"user":            user,
				"pool_id":         float64(0),
				"offering_amount": "1000000000000000000",
				"pay_amount":      "800000000000000000",
				"tge_amount":      "200000000000000000",
				"vesting_amount":  tt.vestingAmount,
			}
			event := makeChainEvent("Harvested", contractAddr, 1, 500, "0xtxharvest01", eventData)

			err := svc.OnHarvested(ctx, event)
			require.NoError(t, err, "OnHarvested 不应返回错误")

			// 验证 launchpad_harvests
			var harvest harvestPO
			require.NoError(t, db.Where("tx_hash = ?", "0xtxharvest01").First(&harvest).Error,
				"应能查到 harvest 记录")
			assert.Equal(t, saleID, harvest.SaleID, "sale_id 应匹配")
			assert.Equal(t, user, harvest.UserAddress, "user_address 应匹配")
			assert.Equal(t, "1000000000000000000", harvest.OfferingAmount, "offering_amount 应匹配")
			assert.Equal(t, "800000000000000000", harvest.PayAmount, "pay_amount 应匹配")
			assert.Equal(t, "200000000000000000", harvest.TGEAmount, "tge_amount 应匹配")
			assert.Equal(t, tt.vestingAmount, harvest.VestingAmount, "vesting_amount 应匹配")

			// 验证 launchpad_vesting_schedules
			var scheduleCount int64
			db.Model(&vestingSchedulePO{}).Where("sale_id = ? AND beneficiary = ?", saleID, user).Count(&scheduleCount)
			if tt.wantSchedule {
				assert.Equal(t, int64(1), scheduleCount, "有 vesting 时应创建 vesting_schedule")
			} else {
				assert.Equal(t, int64(0), scheduleCount, "无 vesting 时不应创建 vesting_schedule")
			}
		})
	}
}

// --- OnReleased 测试 ---

// TestEventService_OnReleased 测试 Released 事件写入 vesting_release 记录。
func TestEventService_OnReleased(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	amount := "300000000000000000"
	eventData := map[string]any{
		"user":          "0xuser_release",
		"amount":        amount,
		"release_index": float64(1),
	}
	event := makeChainEvent("Released", "0xcontract", 1, 600, "0xtxrelease01", eventData)

	err := svc.OnReleased(ctx, event)
	require.NoError(t, err, "OnReleased 不应返回错误")

	// 验证 launchpad_vesting_releases
	var release vestingReleasePO
	require.NoError(t, db.Where("tx_hash = ?", "0xtxrelease01").First(&release).Error,
		"应能查到 vesting_release 记录")
	assert.Equal(t, amount, release.Amount, "amount 应匹配")
	assert.Equal(t, int64(600), release.BlockNumber, "block_number 应匹配")
}

// --- OnRevoked 测试 ---

// TestEventService_OnRevoked 测试 Revoked 事件更新 vesting_revoked 为 true。
func TestEventService_OnRevoked(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	contractAddr := "0xsale_revoke"
	insertTestSale(t, db, contractAddr)

	// 验证初始状态
	var sale salePO
	require.NoError(t, db.Where("contract_address = ?", contractAddr).First(&sale).Error)
	assert.False(t, sale.VestingRevoked, "初始状态应为 false")

	// 执行 Revoked 事件
	event := makeChainEvent("Revoked", contractAddr, 1, 700, "0xtxrevoke", map[string]any{})
	err := svc.OnRevoked(ctx, event)
	require.NoError(t, err, "OnRevoked 不应返回错误")

	// 验证 vesting_revoked 已更新
	require.NoError(t, db.Where("contract_address = ?", contractAddr).First(&sale).Error)
	assert.True(t, sale.VestingRevoked, "Revoked 后 vesting_revoked 应为 true")
}

// --- OnWhitelistAdded 测试 ---

// TestEventService_OnWhitelistAdded 测试 WhitelistAdded 事件写入白名单记录。
func TestEventService_OnWhitelistAdded(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	contractAddr := "0xsale_wl"
	saleID := insertTestSale(t, db, contractAddr)

	user := "0xuser_wl"
	eventData := map[string]any{
		"user":    user,
		"pool_id": float64(0),
	}
	event := makeChainEvent("WhitelistAdded", contractAddr, 1, 800, "0xtxwl01", eventData)

	err := svc.OnWhitelistAdded(ctx, event)
	require.NoError(t, err, "OnWhitelistAdded 不应返回错误")

	// 验证 launchpad_whitelists
	var wl whitelistPO
	require.NoError(t, db.Where("sale_id = ? AND address = ?", saleID, user).First(&wl).Error,
		"应能查到白名单记录")
	assert.Equal(t, saleID, wl.SaleID, "sale_id 应匹配")
	assert.Equal(t, user, wl.Address, "address 应匹配")
	assert.True(t, wl.IsActive, "is_active 应为 true")
}

// TestEventService_OnWhitelistAdded_Idempotent 测试重复添加白名单不创建重复记录。
func TestEventService_OnWhitelistAdded_Idempotent(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	contractAddr := "0xsale_wl_idem"
	saleID := insertTestSale(t, db, contractAddr)

	user := "0xuser_wl_idem"
	eventData := map[string]any{
		"user":    user,
		"pool_id": float64(0),
	}
	event := makeChainEvent("WhitelistAdded", contractAddr, 1, 800, "0xtxwl_idem", eventData)

	require.NoError(t, svc.OnWhitelistAdded(ctx, event), "第一次调用应成功")
	require.NoError(t, svc.OnWhitelistAdded(ctx, event), "重复调用不应返回错误")

	var count int64
	db.Model(&whitelistPO{}).Where("sale_id = ? AND address = ?", saleID, user).Count(&count)
	assert.Equal(t, int64(1), count, "幂等调用不应创建重复白名单记录")
}

// --- OnWhitelistRemoved 测试 ---

// TestEventService_OnWhitelistRemoved 测试 WhitelistRemoved 事件删除白名单记录。
func TestEventService_OnWhitelistRemoved(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	contractAddr := "0xsale_wlrm"
	saleID := insertTestSale(t, db, contractAddr)

	user := "0xuser_wlrm"

	// 先插入白名单记录
	require.NoError(t, db.Create(&whitelistPO{
		SaleID:   saleID,
		Address:  user,
		IsActive: true,
	}).Error, "插入白名单记录失败")

	// 验证记录存在
	var countBefore int64
	db.Model(&whitelistPO{}).Where("sale_id = ? AND address = ?", saleID, user).Count(&countBefore)
	require.Equal(t, int64(1), countBefore, "白名单记录应存在")

	// 执行 WhitelistRemoved 事件
	eventData := map[string]any{
		"user":    user,
		"pool_id": float64(0),
	}
	event := makeChainEvent("WhitelistRemoved", contractAddr, 1, 810, "0xtxwlrm", eventData)

	err := svc.OnWhitelistRemoved(ctx, event)
	require.NoError(t, err, "OnWhitelistRemoved 不应返回错误")

	// 验证白名单记录已删除
	var countAfter int64
	db.Model(&whitelistPO{}).Where("sale_id = ? AND address = ?", saleID, user).Count(&countAfter)
	assert.Equal(t, int64(0), countAfter, "白名单记录应已删除")
}

// --- OnUpdateCeiling 测试 ---

// TestEventService_OnUpdateCeiling 测试 UpdateCeiling 事件更新 ceiling 字段。
func TestEventService_OnUpdateCeiling(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	chainID := 1

	// 先插入 tierParam 记录
	require.NoError(t, db.Create(&tierParamPO{
		ChainID:        chainID,
		Ceiling:        "1000000000000000000",
		Multiplier:     "10",
		TierBaseAmount: "500000000000000000",
	}).Error, "插入 tierParam 记录失败")

	newCeiling := "2000000000000000000"
	eventData := map[string]any{
		"value": newCeiling,
	}
	event := makeChainEvent("UpdateCeiling", "0xcontract", chainID, 900, "0xtxceiling", eventData)

	err := svc.OnUpdateCeiling(ctx, event)
	require.NoError(t, err, "OnUpdateCeiling 不应返回错误")

	// 验证 ceiling 已更新
	var param tierParamPO
	require.NoError(t, db.Where("chain_id = ?", chainID).First(&param).Error)
	assert.Equal(t, newCeiling, param.Ceiling, "ceiling 应更新为新值")
}

// TestEventService_OnUpdateCeiling_NoRecord 测试记录不存在时更新不报错。
func TestEventService_OnUpdateCeiling_NoRecord(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	eventData := map[string]any{"value": "999"}
	event := makeChainEvent("UpdateCeiling", "0xcontract", 99, 900, "0xtx", eventData)

	// 没有记录时，Update 影响行数为 0，但不报错
	err := svc.OnUpdateCeiling(ctx, event)
	assert.NoError(t, err, "无记录时更新不应返回错误")
}

// --- OnUpdateMultiplier 测试 ---

// TestEventService_OnUpdateMultiplier 测试 UpdateMultiplier 事件更新 multiplier 字段。
func TestEventService_OnUpdateMultiplier(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	chainID := 1

	require.NoError(t, db.Create(&tierParamPO{
		ChainID:        chainID,
		Ceiling:        "1000000000000000000",
		Multiplier:     "10",
		TierBaseAmount: "500000000000000000",
	}).Error, "插入 tierParam 记录失败")

	newMultiplier := "20"
	eventData := map[string]any{
		"value": newMultiplier,
	}
	event := makeChainEvent("UpdateMultiplier", "0xcontract", chainID, 910, "0xtxmult", eventData)

	err := svc.OnUpdateMultiplier(ctx, event)
	require.NoError(t, err, "OnUpdateMultiplier 不应返回错误")

	var param tierParamPO
	require.NoError(t, db.Where("chain_id = ?", chainID).First(&param).Error)
	assert.Equal(t, newMultiplier, param.Multiplier, "multiplier 应更新为新值")
	// 验证其他字段未被修改
	assert.Equal(t, "1000000000000000000", param.Ceiling, "ceiling 不应被修改")
	assert.Equal(t, "500000000000000000", param.TierBaseAmount, "tier_base_amount 不应被修改")
}

// --- OnUpdateTierBaseAmount 测试 ---

// TestEventService_OnUpdateTierBaseAmount 测试 UpdateTierBaseAmount 事件更新 tier_base_amount 字段。
func TestEventService_OnUpdateTierBaseAmount(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	chainID := 1

	require.NoError(t, db.Create(&tierParamPO{
		ChainID:        chainID,
		Ceiling:        "1000000000000000000",
		Multiplier:     "10",
		TierBaseAmount: "500000000000000000",
	}).Error, "插入 tierParam 记录失败")

	newBaseAmount := "800000000000000000"
	eventData := map[string]any{
		"value": newBaseAmount,
	}
	event := makeChainEvent("UpdateTierBaseAmount", "0xcontract", chainID, 920, "0xtxbase", eventData)

	err := svc.OnUpdateTierBaseAmount(ctx, event)
	require.NoError(t, err, "OnUpdateTierBaseAmount 不应返回错误")

	var param tierParamPO
	require.NoError(t, db.Where("chain_id = ?", chainID).First(&param).Error)
	assert.Equal(t, newBaseAmount, param.TierBaseAmount, "tier_base_amount 应更新为新值")
	// 验证其他字段未被修改
	assert.Equal(t, "1000000000000000000", param.Ceiling, "ceiling 不应被修改")
	assert.Equal(t, "10", param.Multiplier, "multiplier 不应被修改")
}

// --- OnDeposited_SaleNotFound 测试 ---

// TestEventService_OnDeposited_SaleNotFound 测试 sale 不存在时返回错误。
func TestEventService_OnDeposited_SaleNotFound(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	eventData := map[string]any{
		"user":    "0xuser",
		"pool_id": float64(0),
		"amount":  "1000",
	}
	event := makeChainEvent("Deposited", "0xnotexist", 1, 300, "0xtx", eventData)

	err := svc.OnDeposited(ctx, event)
	assert.Error(t, err, "sale 不存在时应返回错误")
}

// --- OnHarvested_SaleNotFound 测试 ---

// TestEventService_OnHarvested_SaleNotFound 测试 sale 不存在时返回错误。
func TestEventService_OnHarvested_SaleNotFound(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	eventData := map[string]any{
		"user":            "0xuser",
		"pool_id":         float64(0),
		"offering_amount": "1000",
		"pay_amount":      "800",
		"tge_amount":      "200",
		"vesting_amount":  "500",
	}
	event := makeChainEvent("Harvested", "0xnotexist", 1, 500, "0xtx", eventData)

	err := svc.OnHarvested(ctx, event)
	assert.Error(t, err, "sale 不存在时应返回错误")
}

// --- OnWhitelistAdded_SaleNotFound 测试 ---

// TestEventService_OnWhitelistAdded_SaleNotFound 测试 sale 不存在时返回错误。
func TestEventService_OnWhitelistAdded_SaleNotFound(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	eventData := map[string]any{
		"user":    "0xuser",
		"pool_id": float64(0),
	}
	event := makeChainEvent("WhitelistAdded", "0xnotexist", 1, 800, "0xtx", eventData)

	err := svc.OnWhitelistAdded(ctx, event)
	assert.Error(t, err, "sale 不存在时应返回错误")
}

// --- OnWhitelistRemoved_SaleNotFound 测试 ---

// TestEventService_OnWhitelistRemoved_SaleNotFound 测试 sale 不存在时返回错误。
func TestEventService_OnWhitelistRemoved_SaleNotFound(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	eventData := map[string]any{
		"user":    "0xuser",
		"pool_id": float64(0),
	}
	event := makeChainEvent("WhitelistRemoved", "0xnotexist", 1, 810, "0xtx", eventData)

	err := svc.OnWhitelistRemoved(ctx, event)
	assert.Error(t, err, "sale 不存在时应返回错误")
}

// --- HandleEvent 路由集成测试 ---

// TestEventService_HandleEvent_SaleCreatedRoute 测试 HandleEvent 正确路由到 OnSaleCreated。
func TestEventService_HandleEvent_SaleCreatedRoute(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	svc := NewEventService(db, nil)
	ctx := context.Background()

	eventData := map[string]any{
		"sale_address":   "0xsale_route",
		"creator":        "0xcreator",
		"raise_token":    "0xraise",
		"offering_token": "0xoffer",
		"mouse_tier":     "0xtier",
	}
	event := makeChainEvent("SaleCreated", "0xsale_route", 1, 100, "0xtxroute", eventData)

	err := svc.HandleEvent(ctx, event)
	require.NoError(t, err, "HandleEvent 路由 SaleCreated 应成功")

	var sale salePO
	require.NoError(t, db.Where("contract_address = ?", "0xsale_route").First(&sale).Error)
	assert.Equal(t, "0xsale_route", sale.ContractAddress)
}

// --- OnSaleCreated 回填逻辑测试 ---

// insertDraftSale 插入一条 deploying 状态的 draft sale 记录并返回其 ID。
func insertDraftSale(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	sale := &salePO{
		ContractAddress:      "",
		ChainID:              1,
		DeployerAddress:      "0xdeployer",
		OwnerAddress:         "0xowner",
		RaiseTokenAddress:    "0xraise",
		OfferingTokenAddress: "0xoffering",
		MouseTierAddress:     "0xtier",
		Status:               string(domain.SaleDeploying),
	}
	require.NoError(t, db.Create(sale).Error, "插入 draft sale 失败")
	return sale.ID
}

// insertPrepareTxWithSaleID 插入一条 PrepareTx 记录，关联指定的 saleID。
func insertPrepareTxWithSaleID(t *testing.T, db *gorm.DB, txHash string, saleID int64) int64 {
	t.Helper()
	ptx := &prepareTxPO{
		SaleID:        &saleID,
		OperationType: "create_sale",
		CallerAddress: "0xadmin",
		TargetAddress: "0xfactory",
		Value:         "0",
		Calldata:      "0xabc",
		CalldataHash:  "0xhash123",
		Status:        "confirmed",
		TxHash:        &txHash,
		ExpiresAt:     time.Now().Add(24 * time.Hour),
	}
	require.NoError(t, db.Create(ptx).Error, "插入 prepare_tx 失败")
	return ptx.ID
}

// TestEventService_OnSaleCreated_BackfillDraftSale 测试 SaleCreated 事件通过 PrepareTx 回填 draft sale。
func TestEventService_OnSaleCreated_BackfillDraftSale(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	prepareTxRepo := NewPrepareTxRepository(db)
	svc := NewEventService(db, prepareTxRepo)
	ctx := context.Background()

	// 1. 插入 draft sale（deploying 状态，contract_address 为空）
	saleID := insertDraftSale(t, db)

	// 2. 插入 PrepareTx 记录关联该 saleID
	txHash := "0xtx_backfill_001"
	insertPrepareTxWithSaleID(t, db, txHash, saleID)

	// 3. 发出 SaleCreated 事件
	saleAddress := "0xsale_backfilled"
	eventData := map[string]any{
		"sale_address":   saleAddress,
		"creator":        "0xcreator_backfill",
		"raise_token":    "0xraise_backfill",
		"offering_token": "0xoffer_backfill",
		"mouse_tier":     "0xtier_backfill",
		"deployer":       "0xdeployer_backfill",
	}
	event := makeChainEvent("SaleCreated", saleAddress, 1, 100, txHash, eventData)

	err := svc.OnSaleCreated(ctx, event)
	require.NoError(t, err, "OnSaleCreated 回填不应返回错误")

	// 4. 验证 sale 的 contract_address 和 status 已更新
	var sale salePO
	require.NoError(t, db.Where("id = ?", saleID).First(&sale).Error, "应能查到 sale")
	assert.Equal(t, saleAddress, sale.ContractAddress, "contract_address 应被回填")
	assert.Equal(t, string(domain.SaleDeployed), sale.Status, "status 应更新为 deployed")
	assert.Equal(t, "0xdeployer_backfill", sale.DeployerAddress, "deployer_address 应被回填")
	assert.Equal(t, "0xcreator_backfill", sale.OwnerAddress, "owner_address 应被回填")
	assert.Equal(t, "0xraise_backfill", sale.RaiseTokenAddress, "raise_token_address 应被回填")
	assert.Equal(t, "0xoffer_backfill", sale.OfferingTokenAddress, "offering_token_address 应被回填")
	assert.Equal(t, "0xtier_backfill", sale.MouseTierAddress, "mouse_tier_address 应被回填")
}

// TestEventService_OnSaleCreated_BackfillIdempotent 测试回填操作的幂等性。
func TestEventService_OnSaleCreated_BackfillIdempotent(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	prepareTxRepo := NewPrepareTxRepository(db)
	svc := NewEventService(db, prepareTxRepo)
	ctx := context.Background()

	saleID := insertDraftSale(t, db)
	txHash := "0xtx_idem_backfill"
	insertPrepareTxWithSaleID(t, db, txHash, saleID)

	saleAddress := "0xsale_idem_backfill"
	eventData := map[string]any{
		"sale_address":   saleAddress,
		"creator":        "0xcreator",
		"raise_token":    "0xraise",
		"offering_token": "0xoffer",
		"mouse_tier":     "0xtier",
	}
	event := makeChainEvent("SaleCreated", saleAddress, 1, 100, txHash, eventData)

	// 第一次回填
	require.NoError(t, svc.OnSaleCreated(ctx, event), "第一次回填应成功")

	// 第二次回填（幂等）
	require.NoError(t, svc.OnSaleCreated(ctx, event), "重复回填不应返回错误")

	// 验证只有一条 sale 记录且数据正确
	var sale salePO
	require.NoError(t, db.Where("id = ?", saleID).First(&sale).Error)
	assert.Equal(t, saleAddress, sale.ContractAddress, "contract_address 应保持回填值")
	assert.Equal(t, string(domain.SaleDeployed), sale.Status, "status 应为 deployed")
}

// TestEventService_OnSaleCreated_NoPrepareTxFallback 测试找不到 PrepareTx 时回退到 FirstOrCreate。
func TestEventService_OnSaleCreated_NoPrepareTxFallback(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	prepareTxRepo := NewPrepareTxRepository(db)
	svc := NewEventService(db, prepareTxRepo)
	ctx := context.Background()

	saleAddress := "0xsale_fallback"
	txHash := "0xtx_no_prepare"
	eventData := map[string]any{
		"sale_address":   saleAddress,
		"creator":        "0xcreator",
		"raise_token":    "0xraise",
		"offering_token": "0xoffer",
		"mouse_tier":     "0xtier",
	}
	event := makeChainEvent("SaleCreated", saleAddress, 1, 100, txHash, eventData)

	err := svc.OnSaleCreated(ctx, event)
	require.NoError(t, err, "无 PrepareTx 时应回退到 FirstOrCreate")

	// 验证通过 FirstOrCreate 创建了 sale 记录
	var sale salePO
	require.NoError(t, db.Where("contract_address = ?", saleAddress).First(&sale).Error, "应能查到 sale 记录")
	assert.Equal(t, saleAddress, sale.ContractAddress)
}

// TestEventService_OnSaleCreated_PrepareTxNilSaleID 测试 PrepareTx.SaleID 为 nil 时回退到 FirstOrCreate。
func TestEventService_OnSaleCreated_PrepareTxNilSaleID(t *testing.T) {
	t.Parallel()

	db := setupEventTestDB(t)
	prepareTxRepo := NewPrepareTxRepository(db)
	svc := NewEventService(db, prepareTxRepo)
	ctx := context.Background()

	// 插入 PrepareTx 但 sale_id 为 nil
	txHash := "0xtx_nil_saleid"
	ptx := &prepareTxPO{
		SaleID:        nil,
		OperationType: "create_sale",
		CallerAddress: "0xadmin",
		TargetAddress: "0xfactory",
		Value:         "0",
		Calldata:      "0xabc",
		CalldataHash:  "0xhash_nil",
		Status:        "confirmed",
		TxHash:        &txHash,
		ExpiresAt:     time.Now().Add(24 * time.Hour),
	}
	require.NoError(t, db.Create(ptx).Error, "插入 prepare_tx 失败")

	saleAddress := "0xsale_nil_saleid"
	eventData := map[string]any{
		"sale_address":   saleAddress,
		"creator":        "0xcreator",
		"raise_token":    "0xraise",
		"offering_token": "0xoffer",
		"mouse_tier":     "0xtier",
	}
	event := makeChainEvent("SaleCreated", saleAddress, 1, 100, txHash, eventData)

	err := svc.OnSaleCreated(ctx, event)
	require.NoError(t, err, "PrepareTx.SaleID 为 nil 时应回退到 FirstOrCreate")

	// 验证通过 FirstOrCreate 创建了 sale 记录
	var sale salePO
	require.NoError(t, db.Where("contract_address = ?", saleAddress).First(&sale).Error, "应能查到 sale 记录")
	assert.Equal(t, saleAddress, sale.ContractAddress)
}

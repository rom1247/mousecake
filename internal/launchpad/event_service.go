package launchpad

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mousecake-go/mousecake-go/internal/shared/sync"
	"gorm.io/gorm"
)

// EventService 处理 Launchpad 相关的链上事件投影。
type EventService struct {
	db *gorm.DB
}

// NewEventService 创建 EventService 实例。
func NewEventService(db *gorm.DB) *EventService {
	return &EventService{db: db}
}

// HandleEvent 根据事件名称路由到对应的处理方法。
func (s *EventService) HandleEvent(ctx context.Context, event *sync.ChainEvent) error {
	switch event.EventName {
	case "SaleCreated":
		return s.OnSaleCreated(ctx, event)
	case "PoolSet":
		return s.OnPoolSet(ctx, event)
	case "Deposited":
		return s.OnDeposited(ctx, event)
	case "Harvested":
		return s.OnHarvested(ctx, event)
	case "Released":
		return s.OnReleased(ctx, event)
	case "Revoked":
		return s.OnRevoked(ctx, event)
	case "FinalWithdraw":
		return s.OnFinalWithdraw(ctx, event)
	case "WhitelistAdded":
		return s.OnWhitelistAdded(ctx, event)
	case "WhitelistRemoved":
		return s.OnWhitelistRemoved(ctx, event)
	case "UpdateCeiling":
		return s.OnUpdateCeiling(ctx, event)
	case "UpdateMultiplier":
		return s.OnUpdateMultiplier(ctx, event)
	case "UpdateTierBaseAmount":
		return s.OnUpdateTierBaseAmount(ctx, event)
	default:
		slog.Warn("未知事件名", "event_name", event.EventName, "event_id", event.ID)
		return nil
	}
}

// OnSaleCreated 处理 SaleCreated 事件，写入 launchpad_sales。
func (s *EventService) OnSaleCreated(ctx context.Context, event *sync.ChainEvent) error {
	var data map[string]any
	if err := json.Unmarshal([]byte(event.EventData), &data); err != nil {
		return fmt.Errorf("解析 SaleCreated 事件数据: %w", err)
	}

	return s.db.WithContext(ctx).Where("contract_address = ?", event.ContractAddress).
		FirstOrCreate(&salePO{
			ContractAddress:      dataStr(data, "sale_address", event.ContractAddress),
			ChainID:              event.ChainID,
			DeployerAddress:      dataStr(data, "deployer", ""),
			OwnerAddress:         dataStr(data, "creator", ""),
			RaiseTokenAddress:    dataStr(data, "raise_token", ""),
			OfferingTokenAddress: dataStr(data, "offering_token", ""),
			MouseTierAddress:     dataStr(data, "mouse_tier", ""),
		}).Error
}

// OnPoolSet 处理 PoolSet 事件，写入 launchpad_pools。
func (s *EventService) OnPoolSet(ctx context.Context, event *sync.ChainEvent) error {
	var data map[string]any
	if err := json.Unmarshal([]byte(event.EventData), &data); err != nil {
		return fmt.Errorf("解析 PoolSet 事件数据: %w", err)
	}

	var sale salePO
	if err := s.db.WithContext(ctx).Where("contract_address = ?", event.ContractAddress).First(&sale).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("PoolSet: sale 记录不存在，跳过处理", "contract", event.ContractAddress, "event_id", event.ID)
			return nil
		}
		return fmt.Errorf("查找 sale: %w", err)
	}

	poolIndex := dataInt(data, "pool_id", 0)
	return s.db.WithContext(ctx).
		Where("sale_id = ? AND pool_index = ?", sale.ID, poolIndex).
		FirstOrCreate(&poolPO{
			SaleID:         sale.ID,
			PoolIndex:      poolIndex,
			RaisingAmount:  dataStr(data, "raising_amount", "0"),
			OfferingAmount: dataStr(data, "offering_amount", "0"),
			LimitPerUser:   dataStr(data, "limit_per_user", "0"),
		}).Error
}

// OnDeposited 处理 Deposited 事件，三表事务写入。
func (s *EventService) OnDeposited(ctx context.Context, event *sync.ChainEvent) error {
	var data map[string]any
	if err := json.Unmarshal([]byte(event.EventData), &data); err != nil {
		return fmt.Errorf("解析 Deposited 事件数据: %w", err)
	}

	user := dataStr(data, "user", "")
	amount := dataStr(data, "amount", "0")

	// 幂等检查：同一笔交易不重复写入
	var count int64
	if err := s.db.WithContext(ctx).Model(&depositPO{}).
		Where("tx_hash = ?", event.TxHash).
		Count(&count).Error; err != nil {
		return fmt.Errorf("检查 deposit 幂等: %w", err)
	}
	if count > 0 {
		return nil
	}

	var sale salePO
	if err := s.db.WithContext(ctx).Where("contract_address = ?", event.ContractAddress).First(&sale).Error; err != nil {
		return fmt.Errorf("查找 sale: %w", err)
	}

	poolIndex := dataInt(data, "pool_id", 0)

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&depositPO{
			SaleID:      sale.ID,
			PoolIndex:   poolIndex,
			UserAddress: user,
			Amount:      amount,
			TxHash:      event.TxHash,
			BlockNumber: event.BlockNumber,
		}).Error; err != nil {
			return fmt.Errorf("写入 deposit: %w", err)
		}

		if err := tx.Where("sale_id = ? AND pool_index = ? AND user_address = ?",
			sale.ID, poolIndex, user).
			Assign("total_deposited", gorm.Expr("total_deposited + ?", amount)).
			FirstOrCreate(&userPoolStatePO{
				SaleID:         sale.ID,
				PoolIndex:      poolIndex,
				UserAddress:    user,
				TotalDeposited: amount,
			}).Error; err != nil {
			return fmt.Errorf("更新 user_pool_state: %w", err)
		}

		if poolIndex == 0 {
			if err := tx.Where("sale_id = ? AND user_address = ?", sale.ID, user).
				Assign("credit_used", gorm.Expr("credit_used + ?", amount)).
				FirstOrCreate(&userCreditPO{
					SaleID:      sale.ID,
					UserAddress: user,
					CreditUsed:  amount,
				}).Error; err != nil {
				return fmt.Errorf("更新 user_credit: %w", err)
			}
		}

		return nil
	})
}

// OnHarvested 处理 Harvested 事件，写入 launchpad_harvests + launchpad_vesting_schedules。
func (s *EventService) OnHarvested(ctx context.Context, event *sync.ChainEvent) error {
	var data map[string]any
	if err := json.Unmarshal([]byte(event.EventData), &data); err != nil {
		return fmt.Errorf("解析 Harvested 事件数据: %w", err)
	}

	var sale salePO
	if err := s.db.WithContext(ctx).Where("contract_address = ?", event.ContractAddress).First(&sale).Error; err != nil {
		return fmt.Errorf("查找 sale: %w", err)
	}

	user := dataStr(data, "user", "")
	poolIndex := dataInt(data, "pool_id", 0)

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		harvest := harvestPO{
			SaleID:         sale.ID,
			PoolIndex:      poolIndex,
			UserAddress:    user,
			OfferingAmount: dataStr(data, "offering_amount", "0"),
			PayAmount:      dataStr(data, "pay_amount", "0"),
			TGEAmount:      dataStr(data, "tge_amount", "0"),
			VestingAmount:  dataStr(data, "vesting_amount", "0"),
			TxHash:         event.TxHash,
			BlockNumber:    event.BlockNumber,
		}
		if err := tx.Create(&harvest).Error; err != nil {
			return fmt.Errorf("写入 harvest: %w", err)
		}

		if vestingAmt := dataStr(data, "vesting_amount", "0"); vestingAmt != "0" {
			schedule := vestingSchedulePO{
				SaleID:      sale.ID,
				PoolIndex:   poolIndex,
				Beneficiary: user,
				AmountTotal: vestingAmt,
			}
			if err := tx.Create(&schedule).Error; err != nil {
				return fmt.Errorf("写入 vesting_schedule: %w", err)
			}
		}

		return nil
	})
}

// OnReleased 处理 Released 事件，写入 launchpad_vesting_releases。
func (s *EventService) OnReleased(ctx context.Context, event *sync.ChainEvent) error {
	var data map[string]any
	if err := json.Unmarshal([]byte(event.EventData), &data); err != nil {
		return fmt.Errorf("解析 Released 事件数据: %w", err)
	}

	release := vestingReleasePO{
		Amount:      dataStr(data, "amount", "0"),
		TxHash:      event.TxHash,
		BlockNumber: event.BlockNumber,
	}
	return s.db.WithContext(ctx).Create(&release).Error
}

// OnRevoked 处理 Revoked 事件，更新 launchpad_sales。
func (s *EventService) OnRevoked(ctx context.Context, event *sync.ChainEvent) error {
	return s.db.WithContext(ctx).
		Model(&salePO{}).
		Where("contract_address = ?", event.ContractAddress).
		Update("vesting_revoked", true).Error
}

// OnFinalWithdraw 处理 FinalWithdraw 事件，仅记录。
func (s *EventService) OnFinalWithdraw(ctx context.Context, event *sync.ChainEvent) error {
	slog.Info("FinalWithdraw 事件已记录", "event_id", event.ID, "contract", event.ContractAddress)
	return nil
}

// OnWhitelistAdded 处理 WhitelistAdded 事件，写入 launchpad_whitelists。
func (s *EventService) OnWhitelistAdded(ctx context.Context, event *sync.ChainEvent) error {
	var data map[string]any
	if err := json.Unmarshal([]byte(event.EventData), &data); err != nil {
		return fmt.Errorf("解析 WhitelistAdded 事件数据: %w", err)
	}

	var sale salePO
	if err := s.db.WithContext(ctx).Where("contract_address = ?", event.ContractAddress).First(&sale).Error; err != nil {
		return fmt.Errorf("查找 sale: %w", err)
	}

	user := dataStr(data, "user", "")
	return s.db.WithContext(ctx).
		Where("sale_id = ? AND address = ? AND is_active = true", sale.ID, user).
		FirstOrCreate(&whitelistPO{
			SaleID:   sale.ID,
			Address:  user,
			IsActive: true,
		}).Error
}

// OnWhitelistRemoved 处理 WhitelistRemoved 事件，删除白名单记录。
func (s *EventService) OnWhitelistRemoved(ctx context.Context, event *sync.ChainEvent) error {
	var data map[string]any
	if err := json.Unmarshal([]byte(event.EventData), &data); err != nil {
		return fmt.Errorf("解析 WhitelistRemoved 事件数据: %w", err)
	}

	var sale salePO
	if err := s.db.WithContext(ctx).Where("contract_address = ?", event.ContractAddress).First(&sale).Error; err != nil {
		return fmt.Errorf("查找 sale: %w", err)
	}

	user := dataStr(data, "user", "")
	return s.db.WithContext(ctx).
		Where("sale_id = ? AND address = ?", sale.ID, user).
		Delete(&whitelistPO{}).Error
}

// OnUpdateCeiling 处理 UpdateCeiling 事件，更新 launchpad_tier_params。
func (s *EventService) OnUpdateCeiling(ctx context.Context, event *sync.ChainEvent) error {
	return s.updateTierParam(ctx, event, "ceiling")
}

// OnUpdateMultiplier 处理 UpdateMultiplier 事件，更新 launchpad_tier_params。
func (s *EventService) OnUpdateMultiplier(ctx context.Context, event *sync.ChainEvent) error {
	return s.updateTierParam(ctx, event, "multiplier")
}

// OnUpdateTierBaseAmount 处理 UpdateTierBaseAmount 事件，更新 launchpad_tier_params。
func (s *EventService) OnUpdateTierBaseAmount(ctx context.Context, event *sync.ChainEvent) error {
	return s.updateTierParam(ctx, event, "tier_base_amount")
}

// updateTierParam 通用参数更新。
func (s *EventService) updateTierParam(ctx context.Context, event *sync.ChainEvent, field string) error {
	var data map[string]any
	if err := json.Unmarshal([]byte(event.EventData), &data); err != nil {
		return fmt.Errorf("解析 %s 事件数据: %w", event.EventName, err)
	}

	value := dataStr(data, "value", "0")
	return s.db.WithContext(ctx).
		Model(&tierParamPO{}).
		Where("chain_id = ?", event.ChainID).
		Update(field, value).Error
}

// --- 事件处理专用 PO ---

// tierParamPO 对应 launchpad_tier_params 表，仅 event_service 使用。
// 其他 PO 类型复用 repository.go 中已有的定义。
type tierParamPO struct {
	ID             int64  `gorm:"column:id;primaryKey"`
	ChainID        int    `gorm:"column:chain_id"`
	Ceiling        string `gorm:"column:ceiling"`
	Multiplier     string `gorm:"column:multiplier"`
	TierBaseAmount string `gorm:"column:tier_base_amount"`
}

func (tierParamPO) TableName() string { return "launchpad_tier_params" }

// --- 事件数据辅助函数 ---

func dataStr(data map[string]any, key, fallback string) string {
	if v, ok := data[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return fallback
}

func dataInt(data map[string]any, key string, fallback int) int {
	if v, ok := data[key]; ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		}
	}
	return fallback
}

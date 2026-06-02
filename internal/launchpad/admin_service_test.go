package launchpad

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

func TestAdminService_CreateSale(t *testing.T) {
	repo := newMockPrepareTxRepo()
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("create_sale_data"), hash: "hash1"})
	chain := &mockChainReader{}
	prepareSvc := NewPrepareService(repo, chain, encoder, 30*time.Minute)
	svc := NewAdminService(prepareSvc, encoder, "MousePadByTierDeployer", "MousePadByTier")

	tx, err := svc.CreateSale(context.Background(), CreateSaleInput{
		CallerAddress: "0xAdmin",
		RaisingToken:  "0xRaise",
		OfferingToken: "0xOffer",
		Admin:         "0xAdminAddr",
		MouseTier:     "0xTier",
		StartBlock:    1000,
		EndBlock:      2000,
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpCreateSale, tx.OperationType)
	assert.Equal(t, "0xAdmin", tx.CallerAddress)
	assert.Equal(t, "hash1", tx.CalldataHash)
}

func TestAdminService_SetPool(t *testing.T) {
	repo := newMockPrepareTxRepo()
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("set_pool_data"), hash: "hash2"})
	chain := &mockChainReader{}
	prepareSvc := NewPrepareService(repo, chain, encoder, 30*time.Minute)
	svc := NewAdminService(prepareSvc, encoder, "MousePadByTierDeployer", "MousePadByTier")

	tx, err := svc.SetPool(context.Background(), SetPoolInput{
		CallerAddress:      "0xAdmin",
		SaleID:             1,
		PoolIndex:          0,
		OfferingAmount:     "5000000000000000000",
		RaisingAmount:      "1000000000000000000",
		LimitPerUser:       "1000000000000000000",
		IsSpecialSale:      false,
		HasTax:             false,
		VestingPercentage:  0,
		VestingCliff:       0,
		VestingSlicePeriod: 0,
		VestingDuration:    0,
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpSetPool, tx.OperationType)
	assert.NotNil(t, tx.SaleID)
	assert.Equal(t, int64(1), *tx.SaleID)
}

func TestAdminService_SetTierLimits(t *testing.T) {
	repo := newMockPrepareTxRepo()
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("tier_limits_data"), hash: "hash3"})
	chain := &mockChainReader{}
	prepareSvc := NewPrepareService(repo, chain, encoder, 30*time.Minute)
	svc := NewAdminService(prepareSvc, encoder, "MousePadByTierDeployer", "MousePadByTier")

	tx, err := svc.SetTierLimits(context.Background(), SetTierLimitsInput{
		CallerAddress: "0xAdmin",
		Tier:          1,
		Limit:         "1000",
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpSetTierLimits, tx.OperationType)
}

func TestAdminService_AddWhitelist(t *testing.T) {
	repo := newMockPrepareTxRepo()
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("add_wl_data"), hash: "hash4"})
	chain := &mockChainReader{}
	prepareSvc := NewPrepareService(repo, chain, encoder, 30*time.Minute)
	svc := NewAdminService(prepareSvc, encoder, "MousePadByTierDeployer", "MousePadByTier")

	tx, err := svc.AddWhitelist(context.Background(), WhitelistInput{
		CallerAddress: "0xAdmin",
		SaleID:        1,
		Users:         []string{"0xUser1", "0xUser2"},
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpAddWhitelist, tx.OperationType)
}

func TestAdminService_RemoveWhitelist(t *testing.T) {
	repo := newMockPrepareTxRepo()
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("rm_wl_data"), hash: "hash5"})
	chain := &mockChainReader{}
	prepareSvc := NewPrepareService(repo, chain, encoder, 30*time.Minute)
	svc := NewAdminService(prepareSvc, encoder, "MousePadByTierDeployer", "MousePadByTier")

	tx, err := svc.RemoveWhitelist(context.Background(), WhitelistInput{
		CallerAddress: "0xAdmin",
		SaleID:        1,
		Users:         []string{"0xUser1"},
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpRemoveWhitelist, tx.OperationType)
}

func TestAdminService_SetStartEndBlock(t *testing.T) {
	repo := newMockPrepareTxRepo()
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("block_data"), hash: "hash6"})
	chain := &mockChainReader{}
	prepareSvc := NewPrepareService(repo, chain, encoder, 30*time.Minute)
	svc := NewAdminService(prepareSvc, encoder, "MousePadByTierDeployer", "MousePadByTier")

	tx, err := svc.SetStartEndBlock(context.Background(), SetStartEndBlockInput{
		CallerAddress: "0xAdmin",
		SaleID:        1,
		StartBlock:    1000000,
		EndBlock:      1010000,
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpSetStartEndBlock, tx.OperationType)
}

func TestAdminService_Revoke(t *testing.T) {
	repo := newMockPrepareTxRepo()
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("revoke_data"), hash: "hash7"})
	chain := &mockChainReader{}
	prepareSvc := NewPrepareService(repo, chain, encoder, 30*time.Minute)
	svc := NewAdminService(prepareSvc, encoder, "MousePadByTierDeployer", "MousePadByTier")

	tx, err := svc.Revoke(context.Background(), RevokeInput{
		CallerAddress: "0xAdmin",
		SaleID:        1,
		PoolIndex:     0,
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpRevoke, tx.OperationType)
}

func TestAdminService_FinalWithdraw(t *testing.T) {
	repo := newMockPrepareTxRepo()
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("withdraw_data"), hash: "hash8"})
	chain := &mockChainReader{}
	prepareSvc := NewPrepareService(repo, chain, encoder, 30*time.Minute)
	svc := NewAdminService(prepareSvc, encoder, "MousePadByTierDeployer", "MousePadByTier")

	tx, err := svc.FinalWithdraw(context.Background(), FinalWithdrawInput{
		CallerAddress:  "0xAdmin",
		RaisingAmount:  "1000000000000000000",
		OfferingAmount: "2000000000000000000",
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpFinalWithdraw, tx.OperationType)
}

func TestAdminService_RecoverToken(t *testing.T) {
	repo := newMockPrepareTxRepo()
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("recover_data"), hash: "hash9"})
	chain := &mockChainReader{}
	prepareSvc := NewPrepareService(repo, chain, encoder, 30*time.Minute)
	svc := NewAdminService(prepareSvc, encoder, "MousePadByTierDeployer", "MousePadByTier")

	tx, err := svc.RecoverToken(context.Background(), RecoverTokenInput{
		CallerAddress: "0xAdmin",
		TokenAddress:  "0xToken",
		To:            "0xRecipient",
		Amount:        "1000000000000000000",
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpRecoverToken, tx.OperationType)
}

func TestAdminService_EncodeError(t *testing.T) {
	repo := newMockPrepareTxRepo()
	encoder := newMockEncoder(&mockABIEncoder{err: assert.AnError})
	chain := &mockChainReader{}
	prepareSvc := NewPrepareService(repo, chain, encoder, 30*time.Minute)
	svc := NewAdminService(prepareSvc, encoder, "MousePadByTierDeployer", "MousePadByTier")

	_, err := svc.CreateSale(context.Background(), CreateSaleInput{
		CallerAddress: "0xAdmin",
		RaisingToken:  "0xRaise",
		OfferingToken: "0xOffer",
		Admin:         "0xAdminAddr",
		MouseTier:     "0xTier",
		StartBlock:    1000,
		EndBlock:      2000,
	})
	assert.Error(t, err)
}

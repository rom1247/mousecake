package launchpad

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// newTestAdminService 创建测试用 AdminService，注入 mock sale 函数。
func newTestAdminService(
	encoder EncoderInterface,
	findSaleByID func(ctx context.Context, saleID int64) (*domain.Sale, error),
	createSale func(ctx context.Context, sale *domain.Sale) error,
) (*AdminService, *mockPrepareTxRepo) {
	repo := newMockPrepareTxRepo()
	chain := &mockChainReader{}
	prepareSvc := NewPrepareService(repo, chain, encoder, 30*time.Minute)
	svc := &AdminService{
		prepareSvc:      prepareSvc,
		encoder:         encoder,
		deployerName:    "MousePadByTierDeployer",
		padContractName: "MousePadByTier",
		findSaleByID:    findSaleByID,
		createSale:      createSale,
		deployerAddr:    "0xDeployer",
	}
	return svc, repo
}

// mockFindDeployedSale 返回一个总是返回已部署 sale 的 findSaleByID 函数。
func mockFindDeployedSale(contractAddr string) func(ctx context.Context, saleID int64) (*domain.Sale, error) {
	now := time.Now()
	return func(_ context.Context, saleID int64) (*domain.Sale, error) {
		return domain.ReconstructSale(saleID, contractAddr, domain.SaleDeployed, 1,
			"0xdep", "0xowner", "0xraise", "0xoffer", "0xtier",
			1000, 2000, 0, false, 1000, now, now, nil), nil
	}
}

// mockCreateSale 返回一个设置 sale.ID 的 createSale 函数。
func mockCreateSale() func(ctx context.Context, sale *domain.Sale) error {
	var nextID int64 = 1
	return func(_ context.Context, sale *domain.Sale) error {
		sale.ID = nextID
		nextID++
		return nil
	}
}

func TestAdminService_CreateSale(t *testing.T) {
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("create_sale_data"), hash: "hash1"})
	svc, _ := newTestAdminService(encoder, nil, mockCreateSale())

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
	assert.Equal(t, "0xDeployer", tx.TargetAddress)
	assert.NotNil(t, tx.SaleID)
	assert.Equal(t, int64(1), *tx.SaleID)
}

func TestAdminService_SetPool(t *testing.T) {
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("set_pool_data"), hash: "hash2"})
	svc, _ := newTestAdminService(encoder, mockFindDeployedSale("0xSale"), nil)

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
	assert.Equal(t, "0xSale", tx.TargetAddress)
}

func TestAdminService_SetTierLimits(t *testing.T) {
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("tier_limits_data"), hash: "hash3"})
	svc, _ := newTestAdminService(encoder, mockFindDeployedSale("0xSale"), nil)

	tx, err := svc.SetTierLimits(context.Background(), SetTierLimitsInput{
		CallerAddress: "0xAdmin",
		SaleID:        1,
		Tier:          1,
		Limit:         "1000",
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpSetTierLimits, tx.OperationType)
	assert.Equal(t, "0xSale", tx.TargetAddress)
}

func TestAdminService_AddWhitelist(t *testing.T) {
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("add_wl_data"), hash: "hash4"})
	svc, _ := newTestAdminService(encoder, mockFindDeployedSale("0xSale"), nil)

	tx, err := svc.AddWhitelist(context.Background(), WhitelistInput{
		CallerAddress: "0xAdmin",
		SaleID:        1,
		Users:         []string{"0xUser1", "0xUser2"},
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpAddWhitelist, tx.OperationType)
	assert.Equal(t, "0xSale", tx.TargetAddress)
}

func TestAdminService_RemoveWhitelist(t *testing.T) {
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("rm_wl_data"), hash: "hash5"})
	svc, _ := newTestAdminService(encoder, mockFindDeployedSale("0xSale"), nil)

	tx, err := svc.RemoveWhitelist(context.Background(), WhitelistInput{
		CallerAddress: "0xAdmin",
		SaleID:        1,
		Users:         []string{"0xUser1"},
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpRemoveWhitelist, tx.OperationType)
	assert.Equal(t, "0xSale", tx.TargetAddress)
}

func TestAdminService_SetStartEndBlock(t *testing.T) {
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("block_data"), hash: "hash6"})
	svc, _ := newTestAdminService(encoder, mockFindDeployedSale("0xSale"), nil)

	tx, err := svc.SetStartEndBlock(context.Background(), SetStartEndBlockInput{
		CallerAddress: "0xAdmin",
		SaleID:        1,
		StartBlock:    1000000,
		EndBlock:      1010000,
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpSetStartEndBlock, tx.OperationType)
	assert.Equal(t, "0xSale", tx.TargetAddress)
}

func TestAdminService_Revoke(t *testing.T) {
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("revoke_data"), hash: "hash7"})
	svc, _ := newTestAdminService(encoder, mockFindDeployedSale("0xSale"), nil)

	tx, err := svc.Revoke(context.Background(), RevokeInput{
		CallerAddress: "0xAdmin",
		SaleID:        1,
		PoolIndex:     0,
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpRevoke, tx.OperationType)
	assert.Equal(t, "0xSale", tx.TargetAddress)
}

func TestAdminService_FinalWithdraw(t *testing.T) {
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("withdraw_data"), hash: "hash8"})
	svc, _ := newTestAdminService(encoder, mockFindDeployedSale("0xSale"), nil)

	tx, err := svc.FinalWithdraw(context.Background(), FinalWithdrawInput{
		CallerAddress:  "0xAdmin",
		SaleID:         1,
		RaisingAmount:  "1000000000000000000",
		OfferingAmount: "2000000000000000000",
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpFinalWithdraw, tx.OperationType)
	assert.Equal(t, "0xSale", tx.TargetAddress)
}

func TestAdminService_RecoverToken(t *testing.T) {
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("recover_data"), hash: "hash9"})
	svc, _ := newTestAdminService(encoder, mockFindDeployedSale("0xSale"), nil)

	tx, err := svc.RecoverToken(context.Background(), RecoverTokenInput{
		CallerAddress: "0xAdmin",
		SaleID:        1,
		TokenAddress:  "0xToken",
		To:            "0xRecipient",
		Amount:        "1000000000000000000",
	})
	require.NoError(t, err)
	assert.Equal(t, domain.OpRecoverToken, tx.OperationType)
	assert.Equal(t, "0xSale", tx.TargetAddress)
}

func TestAdminService_EncodeError(t *testing.T) {
	encoder := newMockEncoder(&mockABIEncoder{err: assert.AnError})
	svc, _ := newTestAdminService(encoder, nil, mockCreateSale())

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

func TestAdminService_SetPool_Sale未部署(t *testing.T) {
	encoder := newMockEncoder(&mockABIEncoder{calldata: []byte("set_pool_data"), hash: "hash2"})
	now := time.Now()
	svc, _ := newTestAdminService(encoder, func(_ context.Context, saleID int64) (*domain.Sale, error) {
		return domain.ReconstructSale(saleID, "", domain.SaleDeploying, 1,
			"0xdep", "0xowner", "0xraise", "0xoffer", "0xtier",
			1000, 2000, 0, false, 1000, now, now, nil), nil
	}, nil)

	_, err := svc.SetPool(context.Background(), SetPoolInput{
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
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrSaleNotDeployed)
}

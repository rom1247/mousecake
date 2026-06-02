// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package launchpad

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// MousePadByTierMetaData contains all meta data concerning the MousePadByTier contract.
var MousePadByTierMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addWhitelist\",\"inputs\":[{\"name\":\"users\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimedPool\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"computeReleasableAmount\",\"inputs\":[{\"name\":\"scheduleId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pid\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"endBlock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"finalWithdraw\",\"inputs\":[{\"name\":\"raiseAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"offeringAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getUserTier\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"harvest\",\"inputs\":[{\"name\":\"pid\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_raiseToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_offeringToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_startBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_endBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_mouseTier\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mouseTier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractMouseTier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextScheduleId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"offeringToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"poolInfo\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"offeringAmountPool\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raisingAmountPool\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalAmountPool\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"limitPerUserInLP\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isSpecialSale\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"hasTax\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"vestingPercentage\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"vestingCliff\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"vestingSlicePeriodSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"vestingDuration\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"configured\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"raiseToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recoverToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"release\",\"inputs\":[{\"name\":\"scheduleId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeWhitelist\",\"inputs\":[{\"name\":\"users\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revoke\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPool\",\"inputs\":[{\"name\":\"pid\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"offeringAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raisingAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"limitPerUser\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isSpecialSale\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"hasTax\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"vestingPercentage\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"vestingCliff\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"vestingSlicePeriodSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"vestingDuration\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setStartEndBlock\",\"inputs\":[{\"name\":\"_startBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_endBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTierLimits\",\"inputs\":[{\"name\":\"tier\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"limit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"startBlock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tierLimits\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalTaxPool\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"userAmountPool\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"userCreditUsed\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"userVestingScheduleCount\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"vestingRevoked\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"vestingSchedules\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"beneficiary\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pid\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountTotal\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"released\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"vestingStartTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"viewPoolInformation\",\"inputs\":[{\"name\":\"pid\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"offeringAmountPool\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raisingAmountPool\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalAmountPool\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"limitPerUserInLP\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isSpecialSale\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"vestingPercentage\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"vestingCliff\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"vestingSlicePeriodSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"vestingDuration\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"viewPoolTaxInfo\",\"inputs\":[{\"name\":\"pid\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"hasTax\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"_totalTaxPool\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"viewUserPoolInfo\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pid\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"amountPool\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_claimedPool\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"whitelist\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Deposited\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"pid\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FinalWithdraw\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"raiseAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"offeringAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Harvested\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"pid\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"offeringAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"refundAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"taxAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PoolSet\",\"inputs\":[{\"name\":\"pid\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Released\",\"inputs\":[{\"name\":\"beneficiary\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"scheduleId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Revoked\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StartEndBlockUpdated\",\"inputs\":[{\"name\":\"startBlock\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"endBlock\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenRecovered\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WhitelistAdded\",\"inputs\":[{\"name\":\"users\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WhitelistRemoved\",\"inputs\":[{\"name\":\"users\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyClaimed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotRecoverSaleToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientOfferingBalance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidEndBlock\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPid\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSchedule\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidStartBlock\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidVestingPercentage\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotEOA\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInSaleWindow\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotParticipated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotWhitelisted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PoolNotConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SaleNotEnded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SaleStarted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TierLimitExceeded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UserLimitExceeded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
	ID:  "MousePadByTier",
	Bin: "0x608060405234801561000f575f5ffd5b50335f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610081575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016100789190610196565b60405180910390fd5b6100908161009660201b60201c565b506101af565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61018082610157565b9050919050565b61019081610176565b82525050565b5f6020820190506101a95f830184610187565b92915050565b61357f806101bc5f395ff3fe608060405234801561000f575f5ffd5b5060043610610225575f3560e01c80637b4711531161012e578063c5982051116100b6578063edac985b1161007a578063edac985b14610689578063f2fde38b146106a5578063f9488735146106c1578063f9b961d1146106f1578063feebf5291461072257610225565b8063c5982051146105e9578063d2d1016214610605578063ddc6326214610621578063e2bbb1581461063d578063e4d2620e1461065957610225565b80639f1b5248116100fd5780639f1b52481461056b578063a7229fd914610587578063a8660a78146105a3578063b6549f75146105c1578063b7813607146105cb57610225565b80637b471153146104d15780638159d84e146104ed5780638da5cb5b1461051d5780639b19251a1461053b57610225565b806337bdc99b116101b15780634da60aa5116101805780634da60aa51461042857806366ebd6dc146104465780636d3cbe2114610464578063715018a61461049757806376aaf578146104a157610225565b806337bdc99b146103855780633b340eb7146103a157806346ab91bf146103d257806348cd4cb11461040a57610225565b8063145d544b116101f8578063145d544b146102b15780631526fe27146102e1578063232452161461031b57806324a746211461033757806328fddfaf1461036757610225565b8063083c6323146102295780630db6f813146102475780630f038d411461026357806312d3047e14610281575b5f5ffd5b610231610752565b60405161023e9190612afa565b60405180910390f35b610261600480360381019061025c9190612b7a565b610758565b005b61026b610934565b6040516102789190612afa565b60405180910390f35b61029b60048036038101906102969190612cad565b61093a565b6040516102a89190612afa565b60405180910390f35b6102cb60048036038101906102c69190612cd8565b61094f565b6040516102d89190612afa565b60405180910390f35b6102fb60048036038101906102f69190612cd8565b610960565b6040516103129b9a99989796959493929190612d12565b60405180910390f35b61033560048036038101906103309190612e1c565b6109db565b005b610351600480360381019061034c9190612e67565b610ab9565b60405161035e9190612afa565b60405180910390f35b61036f610ad9565b60405161037c9190612ea5565b60405180910390f35b61039f600480360381019061039a9190612cd8565b610aeb565b005b6103bb60048036038101906103b69190612cd8565b610c9f565b6040516103c9929190612ebe565b60405180910390f35b6103ec60048036038101906103e79190612cd8565b610ce3565b60405161040199989796959493929190612ee5565b60405180910390f35b610412610d58565b60405161041f9190612afa565b60405180910390f35b610430610d5e565b60405161043d9190612fcb565b60405180910390f35b61044e610d83565b60405161045b9190613004565b60405180910390f35b61047e60048036038101906104799190612cd8565b610da8565b60405161048e949392919061302c565b60405180910390f35b61049f610df3565b005b6104bb60048036038101906104b69190612cd8565b610e06565b6040516104c89190612afa565b60405180910390f35b6104eb60048036038101906104e6919061306f565b610e1b565b005b61050760048036038101906105029190612e67565b610e3d565b6040516105149190612ea5565b60405180910390f35b610525610e67565b60405161053291906130ad565b60405180910390f35b61055560048036038101906105509190612cad565b610e8e565b6040516105629190612ea5565b60405180910390f35b6105856004803603810190610580919061306f565b610eab565b005b6105a1600480360381019061059c91906130c6565b610fc6565b005b6105ab611141565b6040516105b89190612afa565b60405180910390f35b6105c9611147565b005b6105d3611197565b6040516105e09190612fcb565b60405180910390f35b61060360048036038101906105fe919061306f565b6111bc565b005b61061f600480360381019061061a9190613116565b6112bc565b005b61063b60048036038101906106369190612cd8565b611518565b005b6106576004803603810190610652919061306f565b611ace565b005b610673600480360381019061066e9190612cad565b61219c565b6040516106809190612afa565b60405180910390f35b6106a3600480360381019061069e9190612e1c565b61223d565b005b6106bf60048036038101906106ba9190612cad565b61231c565b005b6106db60048036038101906106d69190612cd8565b6123a0565b6040516106e89190612afa565b60405180910390f35b61070b60048036038101906107069190612e67565b6123b5565b60405161071992919061319f565b60405180910390f35b61073c60048036038101906107379190612cad565b612468565b6040516107499190612afa565b60405180910390f35b60055481565b61076061247d565b600454431061079b576040517f912ee23d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60018a11156107d6576040517f87e8068300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8914806107e357505f88145b1561081a576040517f1f2a200500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6064841115610855576040517f5351790e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60065f8c81526020019081526020015f20905089815f018190555088816001018190555087816003018190555086816004015f6101000a81548160ff021916908315150217905550858160040160016101000a81548160ff0219169083151502179055508481600501819055508381600601819055508281600701819055508181600801819055506001816009015f6101000a81548160ff0219169083151502179055508a7f885555e6d772d9bb0a79bb5480d510aeb3e4cfb30010229ea19e2cfb432c2e4f60405160405180910390a25050505050505050505050565b60115481565b6010602052805f5260405f205f915090505481565b5f61095982612504565b9050919050565b6006602052805f5260405f205f91509050805f015490806001015490806002015490806003015490806004015f9054906101000a900460ff16908060040160019054906101000a900460ff1690806005015490806006015490806007015490806008015490806009015f9054906101000a900460ff1690508b565b6109e361247d565b5f5f90505b82829050811015610a7b575f600b5f858585818110610a0a57610a096131c6565b5b9050602002016020810190610a1f9190612cad565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff02191690831515021790555080806001019150506109e8565b507f1d474f57a5c483b47a8bf6006e39086f96dd040a00cb348e22f80a4ca2c6f2228282604051610aad9291906132af565b60405180910390a15050565b6007602052815f5260405f20602052805f5260405f205f91509150505481565b600e5f9054906101000a900460ff1681565b5f600f5f8381526020019081526020015f2090505f73ffffffffffffffffffffffffffffffffffffffff16815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603610b86576040517fdba16ce800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f610b9083612504565b90505f8103610ba0575050610c9c565b80826003015f828254610bb391906132fe565b92505081905550610c28825f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff168260025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1661274a9092919063ffffffff16565b82825f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f82e416ba72d10e709b5de7ac16f5f49ff1d94f22d55bf582d353d3c313a1e8dd83604051610c919190612afa565b60405180910390a350505b50565b5f5f5f60065f8581526020019081526020015f2090508060040160019054906101000a900460ff16600c5f8681526020019081526020015f20549250925050915091565b5f5f5f5f5f5f5f5f5f5f60065f8c81526020019081526020015f209050805f0154816001015482600201548360030154846004015f9054906101000a900460ff168560050154866006015487600701548860080154995099509950995099509950995099509950509193959799909294969850565b60045481565b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60035f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600f602052805f5260405f205f91509050805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060010154908060020154908060030154905084565b610dfb61247d565b610e045f61279d565b565b600c602052805f5260405f205f915090505481565b610e2361247d565b80600a5f8481526020019081526020015f20819055505050565b6008602052815f5260405f20602052805f5260405f205f915091509054906101000a900460ff1681565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b600b602052805f5260405f205f915054906101000a900460ff1681565b610eb361247d565b5f821115610f0f57610f0e610ec6610e67565b8360015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1661274a9092919063ffffffff16565b5b5f811115610f6b57610f6a610f22610e67565b8260025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1661274a9092919063ffffffff16565b5b610f73610e67565b73ffffffffffffffffffffffffffffffffffffffff167f67cf37956407117fc88903b0eb24757695487f7de7da5efdc1f907de3403ff368383604051610fba929190613331565b60405180910390a25050565b610fce61247d565b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161480611075575060025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16145b156110ac576040517f07c548bc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6110d782828573ffffffffffffffffffffffffffffffffffffffff1661274a9092919063ffffffff16565b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f879f92dded0f26b83c3e00b12e0395dc72cfc3077343d1854ed6988edd1f9096836040516111349190612afa565b60405180910390a3505050565b600d5481565b61114f61247d565b6001600e5f6101000a81548160ff0219169083151502179055507f44825a4b2df8acb19ce4e1afba9aa850c8b65cdb7942e2078f27d0b0960efee660405160405180910390a1565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6111c461247d565b60045443106111ff576040517f912ee23d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b438211611238576040517fec2caa0d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818111611271576040517f7bd4747600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81600481905550806005819055507f790971a2aa5cd27780b496658cb0fa5b661500f0963a4efd24ef00df9ab263cf82826040516112b0929190613331565b60405180910390a15050565b5f6112c561285e565b90505f815f0160089054906101000a900460ff161590505f825f015f9054906101000a900467ffffffffffffffff1690505f5f8267ffffffffffffffff1614801561130d5750825b90505f60018367ffffffffffffffff1614801561134057505f3073ffffffffffffffffffffffffffffffffffffffff163b145b90508115801561134e575080155b15611385576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001855f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555083156113d2576001855f0160086101000a81548160ff0219169083151502179055505b8a60015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508960025f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555088600481905550876005819055508560035f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060016011819055506114b18761279d565b831561150b575f855f0160086101000a81548160ff0219169083151502179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2600160405161150291906133a4565b60405180910390a15b5050505050505050505050565b6005544311611553576040517f9d98b04b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600181111561158e576040517f87e8068300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60065f8381526020019081526020015f2090505f60075f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8481526020019081526020015f205403611627576040517fd996c87f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60085f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8381526020019081526020015f205f9054906101000a900460ff16156116b7576040517f646cf55800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600160085f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8481526020019081526020015f205f6101000a81548160ff0219169083151502179055505f60075f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8481526020019081526020015f205490505f826002015490505f836001015490505f845f015490505f5f5f8486116117aa5784848861179991906133bd565b6117a3919061342b565b9250611852565b8587856117b791906133bd565b6117c1919061342b565b92505f8688876117d191906133bd565b6117db919061342b565b905080886117e9919061345b565b92508860040160019054906101000a900460ff16801561180857505f83115b1561185057611818838888612871565b915081600c5f8c81526020019081526020015f205f82825461183a91906132fe565b92505081905550818361184d919061345b565b92505b505b5f886005015490505f606482606461186a919061345b565b8661187591906133bd565b61187f919061342b565b90505f818661188e919061345b565b90505f8211156118e5576118e4338360025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1661274a9092919063ffffffff16565b5b5f85111561193a57611939338660015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1661274a9092919063ffffffff16565b5b5f811115611a6d575f600d54036119535742600d819055505b5f601154905060405180608001604052803373ffffffffffffffffffffffffffffffffffffffff1681526020018e81526020018381526020015f815250600f5f8381526020019081526020015f205f820151815f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160010155604082015181600201556060820151816003015590505060105f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f815480929190611a539061348e565b9190505550600181611a6591906132fe565b601181905550505b8b3373ffffffffffffffffffffffffffffffffffffffff167fe7d9801d1d042469575a53334cb3601a4d973e612b19b1d06b1ed01606e725ff888888604051611ab8939291906134d5565b60405180910390a3505050505050505050505050565b60045443111580611ae157506005544310155b15611b18576040517f1d0debbb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611b21336128f0565b15611b58576040517fba092d1600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001811115611b93576040517f87e8068300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8203611bcc576040517f1f2a200500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60065f8381526020019081526020015f209050806009015f9054906101000a900460ff16611c27576040517fd64e375e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60065f600181526020019081526020015f205f015460065f5f81526020019081526020015f205f0154611c5b91906132fe565b90508060025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401611cb891906130ad565b602060405180830381865afa158015611cd3573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611cf7919061351e565b1015611d2f576040517fd2874bc200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b816004015f9054906101000a900460ff1615611e5f57600b5f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16611dc5576040517f584a793800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81600301548460075f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8681526020019081526020015f2054611e2291906132fe565b1115611e5a576040517fc3dd4aea00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61207d565b5f60035f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e4d2620e336040518263ffffffff1660e01b8152600401611eba91906130ad565b602060405180830381865afa158015611ed5573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611ef9919061351e565b90505f600a5f8381526020019081526020015f20549050808660095f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054611f5a91906132fe565b1115611f92576040517ff338707100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83600301548660075f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8881526020019081526020015f2054611fef91906132fe565b1115612027576040517fc3dd4aea00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8560095f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461207391906132fe565b9250508190555050505b8360075f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8581526020019081526020015f205f8282546120d891906132fe565b9250508190555083826002015f8282546120f291906132fe565b9250508190555061214733308660015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16612901909392919063ffffffff16565b823373ffffffffffffffffffffffffffffffffffffffff167f73a19dd210f1a7f902193214c0ee91dd35ee5b4d920cba8d519eca65a7b488ca8660405161218e9190612afa565b60405180910390a350505050565b5f60035f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e4d2620e836040518263ffffffff1660e01b81526004016121f791906130ad565b602060405180830381865afa158015612212573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190612236919061351e565b9050919050565b61224561247d565b5f5f90505b828290508110156122de576001600b5f85858581811061226d5761226c6131c6565b5b90506020020160208101906122829190612cad565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff021916908315150217905550808060010191505061224a565b507ff74f148a4f930a0f67a2c33ba932a14e3e91b4e6468f21e545932fd82511153882826040516123109291906132af565b60405180910390a15050565b61232461247d565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603612394575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161238b91906130ad565b60405180910390fd5b61239d8161279d565b50565b600a602052805f5260405f205f915090505481565b5f5f60075f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8481526020019081526020015f205460085f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8581526020019081526020015f205f9054906101000a900460ff16915091509250929050565b6009602052805f5260405f205f915090505481565b612485612956565b73ffffffffffffffffffffffffffffffffffffffff166124a3610e67565b73ffffffffffffffffffffffffffffffffffffffff1614612502576124c6612956565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016124f991906130ad565b60405180910390fd5b565b5f5f600f5f8481526020019081526020015f2090505f73ffffffffffffffffffffffffffffffffffffffff16815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603612577575f915050612745565b5f816002015490505f82600301549050600e5f9054906101000a900460ff16156125b15780826125a7919061345b565b9350505050612745565b5f600d54036125c5575f9350505050612745565b5f60065f856001015481526020019081526020015f2090505f816006015490505f826008015490505f836007015490505f600d5442612604919061345b565b90508381101561261e575f98505050505050505050612745565b828110612640578587612631919061345b565b98505050505050505050612745565b5f8490505f8183612651919061345b565b90505f5f87111561267b576001858361266a919061342b565b61267491906132fe565b905061268a565b8482612687919061342b565b90505b5f5f8811156126bd5760018689896126a2919061345b565b6126ac919061342b565b6126b691906132fe565b90506126cc565b85876126c9919061342b565b90505b5f81036126e7575f9c50505050505050505050505050612745565b808211156126f3578091505b5f81838d61270191906133bd565b61270b919061342b565b90508a8111612729575f9d5050505050505050505050505050612745565b8a81612735919061345b565b9d50505050505050505050505050505b919050565b612757838383600161295d565b61279857826040517f5274afe700000000000000000000000000000000000000000000000000000000815260040161278f91906130ad565b60405180910390fd5b505050565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f5f6128686129bf565b90508091505090565b5f5f82670de0b6b3a76400008561288891906133bd565b612892919061342b565b90505f61289e826129e8565b9050606481116128ca5761271081876128b791906133bd565b6128c1919061342b565b925050506128e9565b620f424081876128da91906133bd565b6128e4919061342b565b925050505b9392505050565b5f5f823b90505f8111915050919050565b61290f848484846001612a71565b61295057836040517f5274afe700000000000000000000000000000000000000000000000000000000815260040161294791906130ad565b60405180910390fd5b50505050565b5f33905090565b5f5f63a9059cbb60e01b9050604051815f525f1960601c86166004528460245260205f60445f5f8b5af1925060015f511483166129b15783831516156129a5573d5f823e3d81fd5b5f873b113d1516831692505b806040525050949350505050565b5f7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005f1b905090565b5f671bc16d674ec80000821015612a025760329050612a6c565b673782dace9d900000821015612a1b5760199050612a6c565b676f05b59d3b200000821015612a3457600f9050612a6c565b67de0b6b3a76400000821015612a4d57600a9050612a6c565b6801bc16d674ec800000821015612a675760059050612a6c565b607d90505b919050565b5f5f6323b872dd60e01b9050604051815f525f1960601c87166004525f1960601c86166024528460445260205f60645f5f8c5af1925060015f51148316612acf578383151615612ac3573d5f823e3d81fd5b5f883b113d1516831692505b806040525f606052505095945050505050565b5f819050919050565b612af481612ae2565b82525050565b5f602082019050612b0d5f830184612aeb565b92915050565b5f5ffd5b5f5ffd5b612b2481612ae2565b8114612b2e575f5ffd5b50565b5f81359050612b3f81612b1b565b92915050565b5f8115159050919050565b612b5981612b45565b8114612b63575f5ffd5b50565b5f81359050612b7481612b50565b92915050565b5f5f5f5f5f5f5f5f5f5f6101408b8d031215612b9957612b98612b13565b5b5f612ba68d828e01612b31565b9a50506020612bb78d828e01612b31565b9950506040612bc88d828e01612b31565b9850506060612bd98d828e01612b31565b9750506080612bea8d828e01612b66565b96505060a0612bfb8d828e01612b66565b95505060c0612c0c8d828e01612b31565b94505060e0612c1d8d828e01612b31565b935050610100612c2f8d828e01612b31565b925050610120612c418d828e01612b31565b9150509295989b9194979a5092959850565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f612c7c82612c53565b9050919050565b612c8c81612c72565b8114612c96575f5ffd5b50565b5f81359050612ca781612c83565b92915050565b5f60208284031215612cc257612cc1612b13565b5b5f612ccf84828501612c99565b91505092915050565b5f60208284031215612ced57612cec612b13565b5b5f612cfa84828501612b31565b91505092915050565b612d0c81612b45565b82525050565b5f61016082019050612d265f83018e612aeb565b612d33602083018d612aeb565b612d40604083018c612aeb565b612d4d606083018b612aeb565b612d5a608083018a612d03565b612d6760a0830189612d03565b612d7460c0830188612aeb565b612d8160e0830187612aeb565b612d8f610100830186612aeb565b612d9d610120830185612aeb565b612dab610140830184612d03565b9c9b505050505050505050505050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83601f840112612ddc57612ddb612dbb565b5b8235905067ffffffffffffffff811115612df957612df8612dbf565b5b602083019150836020820283011115612e1557612e14612dc3565b5b9250929050565b5f5f60208385031215612e3257612e31612b13565b5b5f83013567ffffffffffffffff811115612e4f57612e4e612b17565b5b612e5b85828601612dc7565b92509250509250929050565b5f5f60408385031215612e7d57612e7c612b13565b5b5f612e8a85828601612c99565b9250506020612e9b85828601612b31565b9150509250929050565b5f602082019050612eb85f830184612d03565b92915050565b5f604082019050612ed15f830185612d03565b612ede6020830184612aeb565b9392505050565b5f61012082019050612ef95f83018c612aeb565b612f06602083018b612aeb565b612f13604083018a612aeb565b612f206060830189612aeb565b612f2d6080830188612d03565b612f3a60a0830187612aeb565b612f4760c0830186612aeb565b612f5460e0830185612aeb565b612f62610100830184612aeb565b9a9950505050505050505050565b5f819050919050565b5f612f93612f8e612f8984612c53565b612f70565b612c53565b9050919050565b5f612fa482612f79565b9050919050565b5f612fb582612f9a565b9050919050565b612fc581612fab565b82525050565b5f602082019050612fde5f830184612fbc565b92915050565b5f612fee82612f9a565b9050919050565b612ffe81612fe4565b82525050565b5f6020820190506130175f830184612ff5565b92915050565b61302681612c72565b82525050565b5f60808201905061303f5f83018761301d565b61304c6020830186612aeb565b6130596040830185612aeb565b6130666060830184612aeb565b95945050505050565b5f5f6040838503121561308557613084612b13565b5b5f61309285828601612b31565b92505060206130a385828601612b31565b9150509250929050565b5f6020820190506130c05f83018461301d565b92915050565b5f5f5f606084860312156130dd576130dc612b13565b5b5f6130ea86828701612c99565b93505060206130fb86828701612c99565b925050604061310c86828701612b31565b9150509250925092565b5f5f5f5f5f5f60c087890312156131305761312f612b13565b5b5f61313d89828a01612c99565b965050602061314e89828a01612c99565b955050604061315f89828a01612b31565b945050606061317089828a01612b31565b935050608061318189828a01612c99565b92505060a061319289828a01612c99565b9150509295509295509295565b5f6040820190506131b25f830185612aeb565b6131bf6020830184612d03565b9392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f82825260208201905092915050565b5f819050919050565b61321581612c72565b82525050565b5f613226838361320c565b60208301905092915050565b5f6132406020840184612c99565b905092915050565b5f602082019050919050565b5f61325f83856131f3565b935061326a82613203565b805f5b858110156132a25761327f8284613232565b613289888261321b565b975061329483613248565b92505060018101905061326d565b5085925050509392505050565b5f6020820190508181035f8301526132c8818486613254565b90509392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61330882612ae2565b915061331383612ae2565b925082820190508082111561332b5761332a6132d1565b5b92915050565b5f6040820190506133445f830185612aeb565b6133516020830184612aeb565b9392505050565b5f819050919050565b5f67ffffffffffffffff82169050919050565b5f61338e61338961338484613358565b612f70565b613361565b9050919050565b61339e81613374565b82525050565b5f6020820190506133b75f830184613395565b92915050565b5f6133c782612ae2565b91506133d283612ae2565b92508282026133e081612ae2565b915082820484148315176133f7576133f66132d1565b5b5092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f61343582612ae2565b915061344083612ae2565b9250826134505761344f6133fe565b5b828204905092915050565b5f61346582612ae2565b915061347083612ae2565b9250828203905081811115613488576134876132d1565b5b92915050565b5f61349882612ae2565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036134ca576134c96132d1565b5b600182019050919050565b5f6060820190506134e85f830186612aeb565b6134f56020830185612aeb565b6135026040830184612aeb565b949350505050565b5f8151905061351881612b1b565b92915050565b5f6020828403121561353357613532612b13565b5b5f6135408482850161350a565b9150509291505056fea26469706673582212202664a992b99e9153fc96df2c0ec9c2f7ca2fc0f7d10ea3d2e65baf7ad62226ae64736f6c634300081c0033",
}

// MousePadByTier is an auto generated Go binding around an Ethereum contract.
type MousePadByTier struct {
	abi abi.ABI
}

// NewMousePadByTier creates a new instance of MousePadByTier.
func NewMousePadByTier() *MousePadByTier {
	parsed, err := MousePadByTierMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &MousePadByTier{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *MousePadByTier) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAddWhitelist is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xedac985b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addWhitelist(address[] users) returns()
func (mousePadByTier *MousePadByTier) PackAddWhitelist(users []common.Address) []byte {
	enc, err := mousePadByTier.abi.Pack("addWhitelist", users)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddWhitelist is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xedac985b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addWhitelist(address[] users) returns()
func (mousePadByTier *MousePadByTier) TryPackAddWhitelist(users []common.Address) ([]byte, error) {
	return mousePadByTier.abi.Pack("addWhitelist", users)
}

// PackClaimedPool is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8159d84e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function claimedPool(address , uint256 ) view returns(bool)
func (mousePadByTier *MousePadByTier) PackClaimedPool(arg0 common.Address, arg1 *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("claimedPool", arg0, arg1)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClaimedPool is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8159d84e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function claimedPool(address , uint256 ) view returns(bool)
func (mousePadByTier *MousePadByTier) TryPackClaimedPool(arg0 common.Address, arg1 *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("claimedPool", arg0, arg1)
}

// UnpackClaimedPool is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8159d84e.
//
// Solidity: function claimedPool(address , uint256 ) view returns(bool)
func (mousePadByTier *MousePadByTier) UnpackClaimedPool(data []byte) (bool, error) {
	out, err := mousePadByTier.abi.Unpack("claimedPool", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackComputeReleasableAmount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x145d544b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function computeReleasableAmount(uint256 scheduleId) view returns(uint256)
func (mousePadByTier *MousePadByTier) PackComputeReleasableAmount(scheduleId *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("computeReleasableAmount", scheduleId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackComputeReleasableAmount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x145d544b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function computeReleasableAmount(uint256 scheduleId) view returns(uint256)
func (mousePadByTier *MousePadByTier) TryPackComputeReleasableAmount(scheduleId *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("computeReleasableAmount", scheduleId)
}

// UnpackComputeReleasableAmount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x145d544b.
//
// Solidity: function computeReleasableAmount(uint256 scheduleId) view returns(uint256)
func (mousePadByTier *MousePadByTier) UnpackComputeReleasableAmount(data []byte) (*big.Int, error) {
	out, err := mousePadByTier.abi.Unpack("computeReleasableAmount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe2bbb158.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function deposit(uint256 amount, uint256 pid) returns()
func (mousePadByTier *MousePadByTier) PackDeposit(amount *big.Int, pid *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("deposit", amount, pid)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe2bbb158.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function deposit(uint256 amount, uint256 pid) returns()
func (mousePadByTier *MousePadByTier) TryPackDeposit(amount *big.Int, pid *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("deposit", amount, pid)
}

// PackEndBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x083c6323.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function endBlock() view returns(uint256)
func (mousePadByTier *MousePadByTier) PackEndBlock() []byte {
	enc, err := mousePadByTier.abi.Pack("endBlock")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEndBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x083c6323.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function endBlock() view returns(uint256)
func (mousePadByTier *MousePadByTier) TryPackEndBlock() ([]byte, error) {
	return mousePadByTier.abi.Pack("endBlock")
}

// UnpackEndBlock is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x083c6323.
//
// Solidity: function endBlock() view returns(uint256)
func (mousePadByTier *MousePadByTier) UnpackEndBlock(data []byte) (*big.Int, error) {
	out, err := mousePadByTier.abi.Unpack("endBlock", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackFinalWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9f1b5248.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function finalWithdraw(uint256 raiseAmount, uint256 offeringAmount) returns()
func (mousePadByTier *MousePadByTier) PackFinalWithdraw(raiseAmount *big.Int, offeringAmount *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("finalWithdraw", raiseAmount, offeringAmount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFinalWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9f1b5248.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function finalWithdraw(uint256 raiseAmount, uint256 offeringAmount) returns()
func (mousePadByTier *MousePadByTier) TryPackFinalWithdraw(raiseAmount *big.Int, offeringAmount *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("finalWithdraw", raiseAmount, offeringAmount)
}

// PackGetUserTier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe4d2620e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getUserTier(address user) view returns(uint256)
func (mousePadByTier *MousePadByTier) PackGetUserTier(user common.Address) []byte {
	enc, err := mousePadByTier.abi.Pack("getUserTier", user)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetUserTier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe4d2620e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getUserTier(address user) view returns(uint256)
func (mousePadByTier *MousePadByTier) TryPackGetUserTier(user common.Address) ([]byte, error) {
	return mousePadByTier.abi.Pack("getUserTier", user)
}

// UnpackGetUserTier is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe4d2620e.
//
// Solidity: function getUserTier(address user) view returns(uint256)
func (mousePadByTier *MousePadByTier) UnpackGetUserTier(data []byte) (*big.Int, error) {
	out, err := mousePadByTier.abi.Unpack("getUserTier", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackHarvest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xddc63262.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function harvest(uint256 pid) returns()
func (mousePadByTier *MousePadByTier) PackHarvest(pid *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("harvest", pid)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHarvest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xddc63262.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function harvest(uint256 pid) returns()
func (mousePadByTier *MousePadByTier) TryPackHarvest(pid *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("harvest", pid)
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd2d10162.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function initialize(address _raiseToken, address _offeringToken, uint256 _startBlock, uint256 _endBlock, address _admin, address _mouseTier) returns()
func (mousePadByTier *MousePadByTier) PackInitialize(raiseToken common.Address, offeringToken common.Address, startBlock *big.Int, endBlock *big.Int, admin common.Address, mouseTier common.Address) []byte {
	enc, err := mousePadByTier.abi.Pack("initialize", raiseToken, offeringToken, startBlock, endBlock, admin, mouseTier)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd2d10162.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function initialize(address _raiseToken, address _offeringToken, uint256 _startBlock, uint256 _endBlock, address _admin, address _mouseTier) returns()
func (mousePadByTier *MousePadByTier) TryPackInitialize(raiseToken common.Address, offeringToken common.Address, startBlock *big.Int, endBlock *big.Int, admin common.Address, mouseTier common.Address) ([]byte, error) {
	return mousePadByTier.abi.Pack("initialize", raiseToken, offeringToken, startBlock, endBlock, admin, mouseTier)
}

// PackMouseTier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x66ebd6dc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mouseTier() view returns(address)
func (mousePadByTier *MousePadByTier) PackMouseTier() []byte {
	enc, err := mousePadByTier.abi.Pack("mouseTier")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMouseTier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x66ebd6dc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mouseTier() view returns(address)
func (mousePadByTier *MousePadByTier) TryPackMouseTier() ([]byte, error) {
	return mousePadByTier.abi.Pack("mouseTier")
}

// UnpackMouseTier is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x66ebd6dc.
//
// Solidity: function mouseTier() view returns(address)
func (mousePadByTier *MousePadByTier) UnpackMouseTier(data []byte) (common.Address, error) {
	out, err := mousePadByTier.abi.Unpack("mouseTier", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackNextScheduleId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0f038d41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function nextScheduleId() view returns(uint256)
func (mousePadByTier *MousePadByTier) PackNextScheduleId() []byte {
	enc, err := mousePadByTier.abi.Pack("nextScheduleId")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackNextScheduleId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0f038d41.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function nextScheduleId() view returns(uint256)
func (mousePadByTier *MousePadByTier) TryPackNextScheduleId() ([]byte, error) {
	return mousePadByTier.abi.Pack("nextScheduleId")
}

// UnpackNextScheduleId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0f038d41.
//
// Solidity: function nextScheduleId() view returns(uint256)
func (mousePadByTier *MousePadByTier) UnpackNextScheduleId(data []byte) (*big.Int, error) {
	out, err := mousePadByTier.abi.Unpack("nextScheduleId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackOfferingToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb7813607.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function offeringToken() view returns(address)
func (mousePadByTier *MousePadByTier) PackOfferingToken() []byte {
	enc, err := mousePadByTier.abi.Pack("offeringToken")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOfferingToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb7813607.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function offeringToken() view returns(address)
func (mousePadByTier *MousePadByTier) TryPackOfferingToken() ([]byte, error) {
	return mousePadByTier.abi.Pack("offeringToken")
}

// UnpackOfferingToken is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb7813607.
//
// Solidity: function offeringToken() view returns(address)
func (mousePadByTier *MousePadByTier) UnpackOfferingToken(data []byte) (common.Address, error) {
	out, err := mousePadByTier.abi.Unpack("offeringToken", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (mousePadByTier *MousePadByTier) PackOwner() []byte {
	enc, err := mousePadByTier.abi.Pack("owner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function owner() view returns(address)
func (mousePadByTier *MousePadByTier) TryPackOwner() ([]byte, error) {
	return mousePadByTier.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (mousePadByTier *MousePadByTier) UnpackOwner(data []byte) (common.Address, error) {
	out, err := mousePadByTier.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPoolInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1526fe27.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function poolInfo(uint256 ) view returns(uint256 offeringAmountPool, uint256 raisingAmountPool, uint256 totalAmountPool, uint256 limitPerUserInLP, bool isSpecialSale, bool hasTax, uint256 vestingPercentage, uint256 vestingCliff, uint256 vestingSlicePeriodSeconds, uint256 vestingDuration, bool configured)
func (mousePadByTier *MousePadByTier) PackPoolInfo(arg0 *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("poolInfo", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPoolInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1526fe27.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function poolInfo(uint256 ) view returns(uint256 offeringAmountPool, uint256 raisingAmountPool, uint256 totalAmountPool, uint256 limitPerUserInLP, bool isSpecialSale, bool hasTax, uint256 vestingPercentage, uint256 vestingCliff, uint256 vestingSlicePeriodSeconds, uint256 vestingDuration, bool configured)
func (mousePadByTier *MousePadByTier) TryPackPoolInfo(arg0 *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("poolInfo", arg0)
}

// PoolInfoOutput serves as a container for the return parameters of contract
// method PoolInfo.
type PoolInfoOutput struct {
	OfferingAmountPool        *big.Int
	RaisingAmountPool         *big.Int
	TotalAmountPool           *big.Int
	LimitPerUserInLP          *big.Int
	IsSpecialSale             bool
	HasTax                    bool
	VestingPercentage         *big.Int
	VestingCliff              *big.Int
	VestingSlicePeriodSeconds *big.Int
	VestingDuration           *big.Int
	Configured                bool
}

// UnpackPoolInfo is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1526fe27.
//
// Solidity: function poolInfo(uint256 ) view returns(uint256 offeringAmountPool, uint256 raisingAmountPool, uint256 totalAmountPool, uint256 limitPerUserInLP, bool isSpecialSale, bool hasTax, uint256 vestingPercentage, uint256 vestingCliff, uint256 vestingSlicePeriodSeconds, uint256 vestingDuration, bool configured)
func (mousePadByTier *MousePadByTier) UnpackPoolInfo(data []byte) (PoolInfoOutput, error) {
	out, err := mousePadByTier.abi.Unpack("poolInfo", data)
	outstruct := new(PoolInfoOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.OfferingAmountPool = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.RaisingAmountPool = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.TotalAmountPool = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.LimitPerUserInLP = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.IsSpecialSale = *abi.ConvertType(out[4], new(bool)).(*bool)
	outstruct.HasTax = *abi.ConvertType(out[5], new(bool)).(*bool)
	outstruct.VestingPercentage = abi.ConvertType(out[6], new(big.Int)).(*big.Int)
	outstruct.VestingCliff = abi.ConvertType(out[7], new(big.Int)).(*big.Int)
	outstruct.VestingSlicePeriodSeconds = abi.ConvertType(out[8], new(big.Int)).(*big.Int)
	outstruct.VestingDuration = abi.ConvertType(out[9], new(big.Int)).(*big.Int)
	outstruct.Configured = *abi.ConvertType(out[10], new(bool)).(*bool)
	return *outstruct, nil
}

// PackRaiseToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4da60aa5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function raiseToken() view returns(address)
func (mousePadByTier *MousePadByTier) PackRaiseToken() []byte {
	enc, err := mousePadByTier.abi.Pack("raiseToken")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRaiseToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4da60aa5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function raiseToken() view returns(address)
func (mousePadByTier *MousePadByTier) TryPackRaiseToken() ([]byte, error) {
	return mousePadByTier.abi.Pack("raiseToken")
}

// UnpackRaiseToken is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4da60aa5.
//
// Solidity: function raiseToken() view returns(address)
func (mousePadByTier *MousePadByTier) UnpackRaiseToken(data []byte) (common.Address, error) {
	out, err := mousePadByTier.abi.Unpack("raiseToken", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackRecoverToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa7229fd9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function recoverToken(address token, address to, uint256 amount) returns()
func (mousePadByTier *MousePadByTier) PackRecoverToken(token common.Address, to common.Address, amount *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("recoverToken", token, to, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRecoverToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa7229fd9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function recoverToken(address token, address to, uint256 amount) returns()
func (mousePadByTier *MousePadByTier) TryPackRecoverToken(token common.Address, to common.Address, amount *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("recoverToken", token, to, amount)
}

// PackRelease is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x37bdc99b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function release(uint256 scheduleId) returns()
func (mousePadByTier *MousePadByTier) PackRelease(scheduleId *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("release", scheduleId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRelease is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x37bdc99b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function release(uint256 scheduleId) returns()
func (mousePadByTier *MousePadByTier) TryPackRelease(scheduleId *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("release", scheduleId)
}

// PackRemoveWhitelist is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23245216.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function removeWhitelist(address[] users) returns()
func (mousePadByTier *MousePadByTier) PackRemoveWhitelist(users []common.Address) []byte {
	enc, err := mousePadByTier.abi.Pack("removeWhitelist", users)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRemoveWhitelist is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23245216.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function removeWhitelist(address[] users) returns()
func (mousePadByTier *MousePadByTier) TryPackRemoveWhitelist(users []common.Address) ([]byte, error) {
	return mousePadByTier.abi.Pack("removeWhitelist", users)
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() returns()
func (mousePadByTier *MousePadByTier) PackRenounceOwnership() []byte {
	enc, err := mousePadByTier.abi.Pack("renounceOwnership")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceOwnership() returns()
func (mousePadByTier *MousePadByTier) TryPackRenounceOwnership() ([]byte, error) {
	return mousePadByTier.abi.Pack("renounceOwnership")
}

// PackRevoke is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb6549f75.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revoke() returns()
func (mousePadByTier *MousePadByTier) PackRevoke() []byte {
	enc, err := mousePadByTier.abi.Pack("revoke")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevoke is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb6549f75.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revoke() returns()
func (mousePadByTier *MousePadByTier) TryPackRevoke() ([]byte, error) {
	return mousePadByTier.abi.Pack("revoke")
}

// PackSetPool is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0db6f813.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setPool(uint256 pid, uint256 offeringAmount, uint256 raisingAmount, uint256 limitPerUser, bool isSpecialSale, bool hasTax, uint256 vestingPercentage, uint256 vestingCliff, uint256 vestingSlicePeriodSeconds, uint256 vestingDuration) returns()
func (mousePadByTier *MousePadByTier) PackSetPool(pid *big.Int, offeringAmount *big.Int, raisingAmount *big.Int, limitPerUser *big.Int, isSpecialSale bool, hasTax bool, vestingPercentage *big.Int, vestingCliff *big.Int, vestingSlicePeriodSeconds *big.Int, vestingDuration *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("setPool", pid, offeringAmount, raisingAmount, limitPerUser, isSpecialSale, hasTax, vestingPercentage, vestingCliff, vestingSlicePeriodSeconds, vestingDuration)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetPool is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0db6f813.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setPool(uint256 pid, uint256 offeringAmount, uint256 raisingAmount, uint256 limitPerUser, bool isSpecialSale, bool hasTax, uint256 vestingPercentage, uint256 vestingCliff, uint256 vestingSlicePeriodSeconds, uint256 vestingDuration) returns()
func (mousePadByTier *MousePadByTier) TryPackSetPool(pid *big.Int, offeringAmount *big.Int, raisingAmount *big.Int, limitPerUser *big.Int, isSpecialSale bool, hasTax bool, vestingPercentage *big.Int, vestingCliff *big.Int, vestingSlicePeriodSeconds *big.Int, vestingDuration *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("setPool", pid, offeringAmount, raisingAmount, limitPerUser, isSpecialSale, hasTax, vestingPercentage, vestingCliff, vestingSlicePeriodSeconds, vestingDuration)
}

// PackSetStartEndBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc5982051.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setStartEndBlock(uint256 _startBlock, uint256 _endBlock) returns()
func (mousePadByTier *MousePadByTier) PackSetStartEndBlock(startBlock *big.Int, endBlock *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("setStartEndBlock", startBlock, endBlock)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetStartEndBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc5982051.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setStartEndBlock(uint256 _startBlock, uint256 _endBlock) returns()
func (mousePadByTier *MousePadByTier) TryPackSetStartEndBlock(startBlock *big.Int, endBlock *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("setStartEndBlock", startBlock, endBlock)
}

// PackSetTierLimits is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7b471153.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setTierLimits(uint256 tier, uint256 limit) returns()
func (mousePadByTier *MousePadByTier) PackSetTierLimits(tier *big.Int, limit *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("setTierLimits", tier, limit)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetTierLimits is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7b471153.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setTierLimits(uint256 tier, uint256 limit) returns()
func (mousePadByTier *MousePadByTier) TryPackSetTierLimits(tier *big.Int, limit *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("setTierLimits", tier, limit)
}

// PackStartBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x48cd4cb1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function startBlock() view returns(uint256)
func (mousePadByTier *MousePadByTier) PackStartBlock() []byte {
	enc, err := mousePadByTier.abi.Pack("startBlock")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStartBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x48cd4cb1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function startBlock() view returns(uint256)
func (mousePadByTier *MousePadByTier) TryPackStartBlock() ([]byte, error) {
	return mousePadByTier.abi.Pack("startBlock")
}

// UnpackStartBlock is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x48cd4cb1.
//
// Solidity: function startBlock() view returns(uint256)
func (mousePadByTier *MousePadByTier) UnpackStartBlock(data []byte) (*big.Int, error) {
	out, err := mousePadByTier.abi.Unpack("startBlock", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTierLimits is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf9488735.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function tierLimits(uint256 ) view returns(uint256)
func (mousePadByTier *MousePadByTier) PackTierLimits(arg0 *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("tierLimits", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTierLimits is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf9488735.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function tierLimits(uint256 ) view returns(uint256)
func (mousePadByTier *MousePadByTier) TryPackTierLimits(arg0 *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("tierLimits", arg0)
}

// UnpackTierLimits is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf9488735.
//
// Solidity: function tierLimits(uint256 ) view returns(uint256)
func (mousePadByTier *MousePadByTier) UnpackTierLimits(data []byte) (*big.Int, error) {
	out, err := mousePadByTier.abi.Unpack("tierLimits", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTotalTaxPool is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x76aaf578.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalTaxPool(uint256 ) view returns(uint256)
func (mousePadByTier *MousePadByTier) PackTotalTaxPool(arg0 *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("totalTaxPool", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalTaxPool is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x76aaf578.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalTaxPool(uint256 ) view returns(uint256)
func (mousePadByTier *MousePadByTier) TryPackTotalTaxPool(arg0 *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("totalTaxPool", arg0)
}

// UnpackTotalTaxPool is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x76aaf578.
//
// Solidity: function totalTaxPool(uint256 ) view returns(uint256)
func (mousePadByTier *MousePadByTier) UnpackTotalTaxPool(data []byte) (*big.Int, error) {
	out, err := mousePadByTier.abi.Unpack("totalTaxPool", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (mousePadByTier *MousePadByTier) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := mousePadByTier.abi.Pack("transferOwnership", newOwner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (mousePadByTier *MousePadByTier) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return mousePadByTier.abi.Pack("transferOwnership", newOwner)
}

// PackUserAmountPool is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x24a74621.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function userAmountPool(address , uint256 ) view returns(uint256)
func (mousePadByTier *MousePadByTier) PackUserAmountPool(arg0 common.Address, arg1 *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("userAmountPool", arg0, arg1)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUserAmountPool is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x24a74621.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function userAmountPool(address , uint256 ) view returns(uint256)
func (mousePadByTier *MousePadByTier) TryPackUserAmountPool(arg0 common.Address, arg1 *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("userAmountPool", arg0, arg1)
}

// UnpackUserAmountPool is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x24a74621.
//
// Solidity: function userAmountPool(address , uint256 ) view returns(uint256)
func (mousePadByTier *MousePadByTier) UnpackUserAmountPool(data []byte) (*big.Int, error) {
	out, err := mousePadByTier.abi.Unpack("userAmountPool", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackUserCreditUsed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfeebf529.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function userCreditUsed(address ) view returns(uint256)
func (mousePadByTier *MousePadByTier) PackUserCreditUsed(arg0 common.Address) []byte {
	enc, err := mousePadByTier.abi.Pack("userCreditUsed", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUserCreditUsed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfeebf529.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function userCreditUsed(address ) view returns(uint256)
func (mousePadByTier *MousePadByTier) TryPackUserCreditUsed(arg0 common.Address) ([]byte, error) {
	return mousePadByTier.abi.Pack("userCreditUsed", arg0)
}

// UnpackUserCreditUsed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfeebf529.
//
// Solidity: function userCreditUsed(address ) view returns(uint256)
func (mousePadByTier *MousePadByTier) UnpackUserCreditUsed(data []byte) (*big.Int, error) {
	out, err := mousePadByTier.abi.Unpack("userCreditUsed", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackUserVestingScheduleCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x12d3047e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function userVestingScheduleCount(address ) view returns(uint256)
func (mousePadByTier *MousePadByTier) PackUserVestingScheduleCount(arg0 common.Address) []byte {
	enc, err := mousePadByTier.abi.Pack("userVestingScheduleCount", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUserVestingScheduleCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x12d3047e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function userVestingScheduleCount(address ) view returns(uint256)
func (mousePadByTier *MousePadByTier) TryPackUserVestingScheduleCount(arg0 common.Address) ([]byte, error) {
	return mousePadByTier.abi.Pack("userVestingScheduleCount", arg0)
}

// UnpackUserVestingScheduleCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x12d3047e.
//
// Solidity: function userVestingScheduleCount(address ) view returns(uint256)
func (mousePadByTier *MousePadByTier) UnpackUserVestingScheduleCount(data []byte) (*big.Int, error) {
	out, err := mousePadByTier.abi.Unpack("userVestingScheduleCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackVestingRevoked is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x28fddfaf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function vestingRevoked() view returns(bool)
func (mousePadByTier *MousePadByTier) PackVestingRevoked() []byte {
	enc, err := mousePadByTier.abi.Pack("vestingRevoked")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVestingRevoked is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x28fddfaf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function vestingRevoked() view returns(bool)
func (mousePadByTier *MousePadByTier) TryPackVestingRevoked() ([]byte, error) {
	return mousePadByTier.abi.Pack("vestingRevoked")
}

// UnpackVestingRevoked is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x28fddfaf.
//
// Solidity: function vestingRevoked() view returns(bool)
func (mousePadByTier *MousePadByTier) UnpackVestingRevoked(data []byte) (bool, error) {
	out, err := mousePadByTier.abi.Unpack("vestingRevoked", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackVestingSchedules is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6d3cbe21.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function vestingSchedules(uint256 ) view returns(address beneficiary, uint256 pid, uint256 amountTotal, uint256 released)
func (mousePadByTier *MousePadByTier) PackVestingSchedules(arg0 *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("vestingSchedules", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVestingSchedules is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6d3cbe21.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function vestingSchedules(uint256 ) view returns(address beneficiary, uint256 pid, uint256 amountTotal, uint256 released)
func (mousePadByTier *MousePadByTier) TryPackVestingSchedules(arg0 *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("vestingSchedules", arg0)
}

// VestingSchedulesOutput serves as a container for the return parameters of contract
// method VestingSchedules.
type VestingSchedulesOutput struct {
	Beneficiary common.Address
	Pid         *big.Int
	AmountTotal *big.Int
	Released    *big.Int
}

// UnpackVestingSchedules is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6d3cbe21.
//
// Solidity: function vestingSchedules(uint256 ) view returns(address beneficiary, uint256 pid, uint256 amountTotal, uint256 released)
func (mousePadByTier *MousePadByTier) UnpackVestingSchedules(data []byte) (VestingSchedulesOutput, error) {
	out, err := mousePadByTier.abi.Unpack("vestingSchedules", data)
	outstruct := new(VestingSchedulesOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Beneficiary = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Pid = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.AmountTotal = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.Released = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackVestingStartTime is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa8660a78.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function vestingStartTime() view returns(uint256)
func (mousePadByTier *MousePadByTier) PackVestingStartTime() []byte {
	enc, err := mousePadByTier.abi.Pack("vestingStartTime")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVestingStartTime is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa8660a78.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function vestingStartTime() view returns(uint256)
func (mousePadByTier *MousePadByTier) TryPackVestingStartTime() ([]byte, error) {
	return mousePadByTier.abi.Pack("vestingStartTime")
}

// UnpackVestingStartTime is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa8660a78.
//
// Solidity: function vestingStartTime() view returns(uint256)
func (mousePadByTier *MousePadByTier) UnpackVestingStartTime(data []byte) (*big.Int, error) {
	out, err := mousePadByTier.abi.Unpack("vestingStartTime", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackViewPoolInformation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x46ab91bf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function viewPoolInformation(uint256 pid) view returns(uint256 offeringAmountPool, uint256 raisingAmountPool, uint256 totalAmountPool, uint256 limitPerUserInLP, bool isSpecialSale, uint256 vestingPercentage, uint256 vestingCliff, uint256 vestingSlicePeriodSeconds, uint256 vestingDuration)
func (mousePadByTier *MousePadByTier) PackViewPoolInformation(pid *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("viewPoolInformation", pid)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackViewPoolInformation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x46ab91bf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function viewPoolInformation(uint256 pid) view returns(uint256 offeringAmountPool, uint256 raisingAmountPool, uint256 totalAmountPool, uint256 limitPerUserInLP, bool isSpecialSale, uint256 vestingPercentage, uint256 vestingCliff, uint256 vestingSlicePeriodSeconds, uint256 vestingDuration)
func (mousePadByTier *MousePadByTier) TryPackViewPoolInformation(pid *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("viewPoolInformation", pid)
}

// ViewPoolInformationOutput serves as a container for the return parameters of contract
// method ViewPoolInformation.
type ViewPoolInformationOutput struct {
	OfferingAmountPool        *big.Int
	RaisingAmountPool         *big.Int
	TotalAmountPool           *big.Int
	LimitPerUserInLP          *big.Int
	IsSpecialSale             bool
	VestingPercentage         *big.Int
	VestingCliff              *big.Int
	VestingSlicePeriodSeconds *big.Int
	VestingDuration           *big.Int
}

// UnpackViewPoolInformation is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x46ab91bf.
//
// Solidity: function viewPoolInformation(uint256 pid) view returns(uint256 offeringAmountPool, uint256 raisingAmountPool, uint256 totalAmountPool, uint256 limitPerUserInLP, bool isSpecialSale, uint256 vestingPercentage, uint256 vestingCliff, uint256 vestingSlicePeriodSeconds, uint256 vestingDuration)
func (mousePadByTier *MousePadByTier) UnpackViewPoolInformation(data []byte) (ViewPoolInformationOutput, error) {
	out, err := mousePadByTier.abi.Unpack("viewPoolInformation", data)
	outstruct := new(ViewPoolInformationOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.OfferingAmountPool = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.RaisingAmountPool = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.TotalAmountPool = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.LimitPerUserInLP = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.IsSpecialSale = *abi.ConvertType(out[4], new(bool)).(*bool)
	outstruct.VestingPercentage = abi.ConvertType(out[5], new(big.Int)).(*big.Int)
	outstruct.VestingCliff = abi.ConvertType(out[6], new(big.Int)).(*big.Int)
	outstruct.VestingSlicePeriodSeconds = abi.ConvertType(out[7], new(big.Int)).(*big.Int)
	outstruct.VestingDuration = abi.ConvertType(out[8], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackViewPoolTaxInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3b340eb7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function viewPoolTaxInfo(uint256 pid) view returns(bool hasTax, uint256 _totalTaxPool)
func (mousePadByTier *MousePadByTier) PackViewPoolTaxInfo(pid *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("viewPoolTaxInfo", pid)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackViewPoolTaxInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3b340eb7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function viewPoolTaxInfo(uint256 pid) view returns(bool hasTax, uint256 _totalTaxPool)
func (mousePadByTier *MousePadByTier) TryPackViewPoolTaxInfo(pid *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("viewPoolTaxInfo", pid)
}

// ViewPoolTaxInfoOutput serves as a container for the return parameters of contract
// method ViewPoolTaxInfo.
type ViewPoolTaxInfoOutput struct {
	HasTax       bool
	TotalTaxPool *big.Int
}

// UnpackViewPoolTaxInfo is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3b340eb7.
//
// Solidity: function viewPoolTaxInfo(uint256 pid) view returns(bool hasTax, uint256 _totalTaxPool)
func (mousePadByTier *MousePadByTier) UnpackViewPoolTaxInfo(data []byte) (ViewPoolTaxInfoOutput, error) {
	out, err := mousePadByTier.abi.Unpack("viewPoolTaxInfo", data)
	outstruct := new(ViewPoolTaxInfoOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.HasTax = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.TotalTaxPool = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackViewUserPoolInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf9b961d1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function viewUserPoolInfo(address user, uint256 pid) view returns(uint256 amountPool, bool _claimedPool)
func (mousePadByTier *MousePadByTier) PackViewUserPoolInfo(user common.Address, pid *big.Int) []byte {
	enc, err := mousePadByTier.abi.Pack("viewUserPoolInfo", user, pid)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackViewUserPoolInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf9b961d1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function viewUserPoolInfo(address user, uint256 pid) view returns(uint256 amountPool, bool _claimedPool)
func (mousePadByTier *MousePadByTier) TryPackViewUserPoolInfo(user common.Address, pid *big.Int) ([]byte, error) {
	return mousePadByTier.abi.Pack("viewUserPoolInfo", user, pid)
}

// ViewUserPoolInfoOutput serves as a container for the return parameters of contract
// method ViewUserPoolInfo.
type ViewUserPoolInfoOutput struct {
	AmountPool  *big.Int
	ClaimedPool bool
}

// UnpackViewUserPoolInfo is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf9b961d1.
//
// Solidity: function viewUserPoolInfo(address user, uint256 pid) view returns(uint256 amountPool, bool _claimedPool)
func (mousePadByTier *MousePadByTier) UnpackViewUserPoolInfo(data []byte) (ViewUserPoolInfoOutput, error) {
	out, err := mousePadByTier.abi.Unpack("viewUserPoolInfo", data)
	outstruct := new(ViewUserPoolInfoOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.AmountPool = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.ClaimedPool = *abi.ConvertType(out[1], new(bool)).(*bool)
	return *outstruct, nil
}

// PackWhitelist is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9b19251a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function whitelist(address ) view returns(bool)
func (mousePadByTier *MousePadByTier) PackWhitelist(arg0 common.Address) []byte {
	enc, err := mousePadByTier.abi.Pack("whitelist", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWhitelist is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9b19251a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function whitelist(address ) view returns(bool)
func (mousePadByTier *MousePadByTier) TryPackWhitelist(arg0 common.Address) ([]byte, error) {
	return mousePadByTier.abi.Pack("whitelist", arg0)
}

// UnpackWhitelist is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9b19251a.
//
// Solidity: function whitelist(address ) view returns(bool)
func (mousePadByTier *MousePadByTier) UnpackWhitelist(data []byte) (bool, error) {
	out, err := mousePadByTier.abi.Unpack("whitelist", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// MousePadByTierDeposited represents a Deposited event raised by the MousePadByTier contract.
type MousePadByTierDeposited struct {
	User   common.Address
	Pid    *big.Int
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const MousePadByTierDepositedEventName = "Deposited"

// ContractEventName returns the user-defined event name.
func (MousePadByTierDeposited) ContractEventName() string {
	return MousePadByTierDepositedEventName
}

// UnpackDepositedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Deposited(address indexed user, uint256 indexed pid, uint256 amount)
func (mousePadByTier *MousePadByTier) UnpackDepositedEvent(log *types.Log) (*MousePadByTierDeposited, error) {
	event := "Deposited"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePadByTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePadByTierDeposited)
	if len(log.Data) > 0 {
		if err := mousePadByTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePadByTier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MousePadByTierFinalWithdraw represents a FinalWithdraw event raised by the MousePadByTier contract.
type MousePadByTierFinalWithdraw struct {
	To             common.Address
	RaiseAmount    *big.Int
	OfferingAmount *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const MousePadByTierFinalWithdrawEventName = "FinalWithdraw"

// ContractEventName returns the user-defined event name.
func (MousePadByTierFinalWithdraw) ContractEventName() string {
	return MousePadByTierFinalWithdrawEventName
}

// UnpackFinalWithdrawEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event FinalWithdraw(address indexed to, uint256 raiseAmount, uint256 offeringAmount)
func (mousePadByTier *MousePadByTier) UnpackFinalWithdrawEvent(log *types.Log) (*MousePadByTierFinalWithdraw, error) {
	event := "FinalWithdraw"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePadByTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePadByTierFinalWithdraw)
	if len(log.Data) > 0 {
		if err := mousePadByTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePadByTier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MousePadByTierHarvested represents a Harvested event raised by the MousePadByTier contract.
type MousePadByTierHarvested struct {
	User           common.Address
	Pid            *big.Int
	OfferingAmount *big.Int
	RefundAmount   *big.Int
	TaxAmount      *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const MousePadByTierHarvestedEventName = "Harvested"

// ContractEventName returns the user-defined event name.
func (MousePadByTierHarvested) ContractEventName() string {
	return MousePadByTierHarvestedEventName
}

// UnpackHarvestedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Harvested(address indexed user, uint256 indexed pid, uint256 offeringAmount, uint256 refundAmount, uint256 taxAmount)
func (mousePadByTier *MousePadByTier) UnpackHarvestedEvent(log *types.Log) (*MousePadByTierHarvested, error) {
	event := "Harvested"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePadByTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePadByTierHarvested)
	if len(log.Data) > 0 {
		if err := mousePadByTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePadByTier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MousePadByTierInitialized represents a Initialized event raised by the MousePadByTier contract.
type MousePadByTierInitialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const MousePadByTierInitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (MousePadByTierInitialized) ContractEventName() string {
	return MousePadByTierInitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (mousePadByTier *MousePadByTier) UnpackInitializedEvent(log *types.Log) (*MousePadByTierInitialized, error) {
	event := "Initialized"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePadByTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePadByTierInitialized)
	if len(log.Data) > 0 {
		if err := mousePadByTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePadByTier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MousePadByTierOwnershipTransferred represents a OwnershipTransferred event raised by the MousePadByTier contract.
type MousePadByTierOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const MousePadByTierOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (MousePadByTierOwnershipTransferred) ContractEventName() string {
	return MousePadByTierOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (mousePadByTier *MousePadByTier) UnpackOwnershipTransferredEvent(log *types.Log) (*MousePadByTierOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePadByTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePadByTierOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := mousePadByTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePadByTier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MousePadByTierPoolSet represents a PoolSet event raised by the MousePadByTier contract.
type MousePadByTierPoolSet struct {
	Pid *big.Int
	Raw *types.Log // Blockchain specific contextual infos
}

const MousePadByTierPoolSetEventName = "PoolSet"

// ContractEventName returns the user-defined event name.
func (MousePadByTierPoolSet) ContractEventName() string {
	return MousePadByTierPoolSetEventName
}

// UnpackPoolSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PoolSet(uint256 indexed pid)
func (mousePadByTier *MousePadByTier) UnpackPoolSetEvent(log *types.Log) (*MousePadByTierPoolSet, error) {
	event := "PoolSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePadByTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePadByTierPoolSet)
	if len(log.Data) > 0 {
		if err := mousePadByTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePadByTier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MousePadByTierReleased represents a Released event raised by the MousePadByTier contract.
type MousePadByTierReleased struct {
	Beneficiary common.Address
	ScheduleId  *big.Int
	Amount      *big.Int
	Raw         *types.Log // Blockchain specific contextual infos
}

const MousePadByTierReleasedEventName = "Released"

// ContractEventName returns the user-defined event name.
func (MousePadByTierReleased) ContractEventName() string {
	return MousePadByTierReleasedEventName
}

// UnpackReleasedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Released(address indexed beneficiary, uint256 indexed scheduleId, uint256 amount)
func (mousePadByTier *MousePadByTier) UnpackReleasedEvent(log *types.Log) (*MousePadByTierReleased, error) {
	event := "Released"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePadByTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePadByTierReleased)
	if len(log.Data) > 0 {
		if err := mousePadByTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePadByTier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MousePadByTierRevoked represents a Revoked event raised by the MousePadByTier contract.
type MousePadByTierRevoked struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const MousePadByTierRevokedEventName = "Revoked"

// ContractEventName returns the user-defined event name.
func (MousePadByTierRevoked) ContractEventName() string {
	return MousePadByTierRevokedEventName
}

// UnpackRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Revoked()
func (mousePadByTier *MousePadByTier) UnpackRevokedEvent(log *types.Log) (*MousePadByTierRevoked, error) {
	event := "Revoked"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePadByTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePadByTierRevoked)
	if len(log.Data) > 0 {
		if err := mousePadByTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePadByTier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MousePadByTierStartEndBlockUpdated represents a StartEndBlockUpdated event raised by the MousePadByTier contract.
type MousePadByTierStartEndBlockUpdated struct {
	StartBlock *big.Int
	EndBlock   *big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const MousePadByTierStartEndBlockUpdatedEventName = "StartEndBlockUpdated"

// ContractEventName returns the user-defined event name.
func (MousePadByTierStartEndBlockUpdated) ContractEventName() string {
	return MousePadByTierStartEndBlockUpdatedEventName
}

// UnpackStartEndBlockUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event StartEndBlockUpdated(uint256 startBlock, uint256 endBlock)
func (mousePadByTier *MousePadByTier) UnpackStartEndBlockUpdatedEvent(log *types.Log) (*MousePadByTierStartEndBlockUpdated, error) {
	event := "StartEndBlockUpdated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePadByTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePadByTierStartEndBlockUpdated)
	if len(log.Data) > 0 {
		if err := mousePadByTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePadByTier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MousePadByTierTokenRecovered represents a TokenRecovered event raised by the MousePadByTier contract.
type MousePadByTierTokenRecovered struct {
	Token  common.Address
	To     common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const MousePadByTierTokenRecoveredEventName = "TokenRecovered"

// ContractEventName returns the user-defined event name.
func (MousePadByTierTokenRecovered) ContractEventName() string {
	return MousePadByTierTokenRecoveredEventName
}

// UnpackTokenRecoveredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenRecovered(address indexed token, address indexed to, uint256 amount)
func (mousePadByTier *MousePadByTier) UnpackTokenRecoveredEvent(log *types.Log) (*MousePadByTierTokenRecovered, error) {
	event := "TokenRecovered"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePadByTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePadByTierTokenRecovered)
	if len(log.Data) > 0 {
		if err := mousePadByTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePadByTier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MousePadByTierWhitelistAdded represents a WhitelistAdded event raised by the MousePadByTier contract.
type MousePadByTierWhitelistAdded struct {
	Users []common.Address
	Raw   *types.Log // Blockchain specific contextual infos
}

const MousePadByTierWhitelistAddedEventName = "WhitelistAdded"

// ContractEventName returns the user-defined event name.
func (MousePadByTierWhitelistAdded) ContractEventName() string {
	return MousePadByTierWhitelistAddedEventName
}

// UnpackWhitelistAddedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event WhitelistAdded(address[] users)
func (mousePadByTier *MousePadByTier) UnpackWhitelistAddedEvent(log *types.Log) (*MousePadByTierWhitelistAdded, error) {
	event := "WhitelistAdded"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePadByTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePadByTierWhitelistAdded)
	if len(log.Data) > 0 {
		if err := mousePadByTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePadByTier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// MousePadByTierWhitelistRemoved represents a WhitelistRemoved event raised by the MousePadByTier contract.
type MousePadByTierWhitelistRemoved struct {
	Users []common.Address
	Raw   *types.Log // Blockchain specific contextual infos
}

const MousePadByTierWhitelistRemovedEventName = "WhitelistRemoved"

// ContractEventName returns the user-defined event name.
func (MousePadByTierWhitelistRemoved) ContractEventName() string {
	return MousePadByTierWhitelistRemovedEventName
}

// UnpackWhitelistRemovedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event WhitelistRemoved(address[] users)
func (mousePadByTier *MousePadByTier) UnpackWhitelistRemovedEvent(log *types.Log) (*MousePadByTierWhitelistRemoved, error) {
	event := "WhitelistRemoved"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePadByTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePadByTierWhitelistRemoved)
	if len(log.Data) > 0 {
		if err := mousePadByTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePadByTier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (mousePadByTier *MousePadByTier) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["AlreadyClaimed"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackAlreadyClaimedError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["CannotRecoverSaleToken"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackCannotRecoverSaleTokenError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["InsufficientOfferingBalance"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackInsufficientOfferingBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["InvalidEndBlock"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackInvalidEndBlockError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["InvalidPid"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackInvalidPidError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["InvalidSchedule"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackInvalidScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["InvalidStartBlock"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackInvalidStartBlockError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["InvalidVestingPercentage"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackInvalidVestingPercentageError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["NotEOA"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackNotEOAError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["NotInSaleWindow"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackNotInSaleWindowError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["NotParticipated"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackNotParticipatedError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["NotWhitelisted"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackNotWhitelistedError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["PoolNotConfigured"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackPoolNotConfiguredError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackSafeERC20FailedOperationError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["SaleNotEnded"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackSaleNotEndedError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["SaleStarted"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackSaleStartedError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["TierLimitExceeded"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackTierLimitExceededError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["UserLimitExceeded"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackUserLimitExceededError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePadByTier.abi.Errors["ZeroAmount"].ID.Bytes()[:4]) {
		return mousePadByTier.UnpackZeroAmountError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// MousePadByTierAlreadyClaimed represents a AlreadyClaimed error raised by the MousePadByTier contract.
type MousePadByTierAlreadyClaimed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AlreadyClaimed()
func MousePadByTierAlreadyClaimedErrorID() common.Hash {
	return common.HexToHash("0x646cf558a545d59f8a09cbf8a0eb8a9332f1d17834843b20fc8d154839dc46d7")
}

// UnpackAlreadyClaimedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AlreadyClaimed()
func (mousePadByTier *MousePadByTier) UnpackAlreadyClaimedError(raw []byte) (*MousePadByTierAlreadyClaimed, error) {
	out := new(MousePadByTierAlreadyClaimed)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "AlreadyClaimed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierCannotRecoverSaleToken represents a CannotRecoverSaleToken error raised by the MousePadByTier contract.
type MousePadByTierCannotRecoverSaleToken struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CannotRecoverSaleToken()
func MousePadByTierCannotRecoverSaleTokenErrorID() common.Hash {
	return common.HexToHash("0x07c548bceaaa1730c020d4e8a6203243cba23c51825b3d4bf9d000a7ee0a5a3d")
}

// UnpackCannotRecoverSaleTokenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CannotRecoverSaleToken()
func (mousePadByTier *MousePadByTier) UnpackCannotRecoverSaleTokenError(raw []byte) (*MousePadByTierCannotRecoverSaleToken, error) {
	out := new(MousePadByTierCannotRecoverSaleToken)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "CannotRecoverSaleToken", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierInsufficientOfferingBalance represents a InsufficientOfferingBalance error raised by the MousePadByTier contract.
type MousePadByTierInsufficientOfferingBalance struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientOfferingBalance()
func MousePadByTierInsufficientOfferingBalanceErrorID() common.Hash {
	return common.HexToHash("0xd2874bc2c33d99091fac11dfa9974a2de504dec1952385f2b13d63c9f7931643")
}

// UnpackInsufficientOfferingBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientOfferingBalance()
func (mousePadByTier *MousePadByTier) UnpackInsufficientOfferingBalanceError(raw []byte) (*MousePadByTierInsufficientOfferingBalance, error) {
	out := new(MousePadByTierInsufficientOfferingBalance)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "InsufficientOfferingBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierInvalidEndBlock represents a InvalidEndBlock error raised by the MousePadByTier contract.
type MousePadByTierInvalidEndBlock struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidEndBlock()
func MousePadByTierInvalidEndBlockErrorID() common.Hash {
	return common.HexToHash("0x7bd47476a4ad87dcdb7f93e2c43562cca33463a49755b5821c6e567311858c06")
}

// UnpackInvalidEndBlockError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidEndBlock()
func (mousePadByTier *MousePadByTier) UnpackInvalidEndBlockError(raw []byte) (*MousePadByTierInvalidEndBlock, error) {
	out := new(MousePadByTierInvalidEndBlock)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "InvalidEndBlock", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierInvalidInitialization represents a InvalidInitialization error raised by the MousePadByTier contract.
type MousePadByTierInvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func MousePadByTierInvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (mousePadByTier *MousePadByTier) UnpackInvalidInitializationError(raw []byte) (*MousePadByTierInvalidInitialization, error) {
	out := new(MousePadByTierInvalidInitialization)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierInvalidPid represents a InvalidPid error raised by the MousePadByTier contract.
type MousePadByTierInvalidPid struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPid()
func MousePadByTierInvalidPidErrorID() common.Hash {
	return common.HexToHash("0x87e80683816e96b292d47c170f9d87fc24d8a905d1e3fc3fd12a7d5da1daaf20")
}

// UnpackInvalidPidError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPid()
func (mousePadByTier *MousePadByTier) UnpackInvalidPidError(raw []byte) (*MousePadByTierInvalidPid, error) {
	out := new(MousePadByTierInvalidPid)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "InvalidPid", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierInvalidSchedule represents a InvalidSchedule error raised by the MousePadByTier contract.
type MousePadByTierInvalidSchedule struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSchedule()
func MousePadByTierInvalidScheduleErrorID() common.Hash {
	return common.HexToHash("0xdba16ce8d857d25a44e7fd7a61f7aa30811825f024212f74f25b753b90dcf591")
}

// UnpackInvalidScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSchedule()
func (mousePadByTier *MousePadByTier) UnpackInvalidScheduleError(raw []byte) (*MousePadByTierInvalidSchedule, error) {
	out := new(MousePadByTierInvalidSchedule)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "InvalidSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierInvalidStartBlock represents a InvalidStartBlock error raised by the MousePadByTier contract.
type MousePadByTierInvalidStartBlock struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidStartBlock()
func MousePadByTierInvalidStartBlockErrorID() common.Hash {
	return common.HexToHash("0xec2caa0d775080df367e244edc5b8bdb34f2ed861e41d4b2aba05f02a20fbe54")
}

// UnpackInvalidStartBlockError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidStartBlock()
func (mousePadByTier *MousePadByTier) UnpackInvalidStartBlockError(raw []byte) (*MousePadByTierInvalidStartBlock, error) {
	out := new(MousePadByTierInvalidStartBlock)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "InvalidStartBlock", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierInvalidVestingPercentage represents a InvalidVestingPercentage error raised by the MousePadByTier contract.
type MousePadByTierInvalidVestingPercentage struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidVestingPercentage()
func MousePadByTierInvalidVestingPercentageErrorID() common.Hash {
	return common.HexToHash("0x5351790e4139bbba14c2b7ed530e27bad9bdc1804ec66f57606caf9f3960132f")
}

// UnpackInvalidVestingPercentageError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidVestingPercentage()
func (mousePadByTier *MousePadByTier) UnpackInvalidVestingPercentageError(raw []byte) (*MousePadByTierInvalidVestingPercentage, error) {
	out := new(MousePadByTierInvalidVestingPercentage)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "InvalidVestingPercentage", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierNotEOA represents a NotEOA error raised by the MousePadByTier contract.
type MousePadByTierNotEOA struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotEOA()
func MousePadByTierNotEOAErrorID() common.Hash {
	return common.HexToHash("0xba092d16b2843f20004dd5dc06156734e6eb5560b8b8baad292949415a5ebffd")
}

// UnpackNotEOAError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotEOA()
func (mousePadByTier *MousePadByTier) UnpackNotEOAError(raw []byte) (*MousePadByTierNotEOA, error) {
	out := new(MousePadByTierNotEOA)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "NotEOA", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierNotInSaleWindow represents a NotInSaleWindow error raised by the MousePadByTier contract.
type MousePadByTierNotInSaleWindow struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInSaleWindow()
func MousePadByTierNotInSaleWindowErrorID() common.Hash {
	return common.HexToHash("0x1d0debbbf0f4a493fa074487f46ba315cb06c4530586c75c11cdfffeeb8ab319")
}

// UnpackNotInSaleWindowError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInSaleWindow()
func (mousePadByTier *MousePadByTier) UnpackNotInSaleWindowError(raw []byte) (*MousePadByTierNotInSaleWindow, error) {
	out := new(MousePadByTierNotInSaleWindow)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "NotInSaleWindow", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierNotInitializing represents a NotInitializing error raised by the MousePadByTier contract.
type MousePadByTierNotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func MousePadByTierNotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (mousePadByTier *MousePadByTier) UnpackNotInitializingError(raw []byte) (*MousePadByTierNotInitializing, error) {
	out := new(MousePadByTierNotInitializing)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierNotParticipated represents a NotParticipated error raised by the MousePadByTier contract.
type MousePadByTierNotParticipated struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotParticipated()
func MousePadByTierNotParticipatedErrorID() common.Hash {
	return common.HexToHash("0xd996c87f4fa4eba5aee686f58a80b538d1a7b91178f95c93026844ae7d07a8c4")
}

// UnpackNotParticipatedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotParticipated()
func (mousePadByTier *MousePadByTier) UnpackNotParticipatedError(raw []byte) (*MousePadByTierNotParticipated, error) {
	out := new(MousePadByTierNotParticipated)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "NotParticipated", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierNotWhitelisted represents a NotWhitelisted error raised by the MousePadByTier contract.
type MousePadByTierNotWhitelisted struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotWhitelisted()
func MousePadByTierNotWhitelistedErrorID() common.Hash {
	return common.HexToHash("0x584a79384c1fda40857fbf9b9412f34690de80c3d1f0a700350da3cc6e7423fe")
}

// UnpackNotWhitelistedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotWhitelisted()
func (mousePadByTier *MousePadByTier) UnpackNotWhitelistedError(raw []byte) (*MousePadByTierNotWhitelisted, error) {
	out := new(MousePadByTierNotWhitelisted)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "NotWhitelisted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the MousePadByTier contract.
type MousePadByTierOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func MousePadByTierOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (mousePadByTier *MousePadByTier) UnpackOwnableInvalidOwnerError(raw []byte) (*MousePadByTierOwnableInvalidOwner, error) {
	out := new(MousePadByTierOwnableInvalidOwner)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the MousePadByTier contract.
type MousePadByTierOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func MousePadByTierOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (mousePadByTier *MousePadByTier) UnpackOwnableUnauthorizedAccountError(raw []byte) (*MousePadByTierOwnableUnauthorizedAccount, error) {
	out := new(MousePadByTierOwnableUnauthorizedAccount)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierPoolNotConfigured represents a PoolNotConfigured error raised by the MousePadByTier contract.
type MousePadByTierPoolNotConfigured struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PoolNotConfigured()
func MousePadByTierPoolNotConfiguredErrorID() common.Hash {
	return common.HexToHash("0xd64e375efadad4ef47fb8523a2c1fd859026c68de6ac6c9e154bc74f08794702")
}

// UnpackPoolNotConfiguredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PoolNotConfigured()
func (mousePadByTier *MousePadByTier) UnpackPoolNotConfiguredError(raw []byte) (*MousePadByTierPoolNotConfigured, error) {
	out := new(MousePadByTierPoolNotConfigured)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "PoolNotConfigured", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierSafeERC20FailedOperation represents a SafeERC20FailedOperation error raised by the MousePadByTier contract.
type MousePadByTierSafeERC20FailedOperation struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeERC20FailedOperation(address token)
func MousePadByTierSafeERC20FailedOperationErrorID() common.Hash {
	return common.HexToHash("0x5274afe73c98b4749fc91ffae6b7b574e7842cb2144a159e9377a5f20b32edf9")
}

// UnpackSafeERC20FailedOperationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeERC20FailedOperation(address token)
func (mousePadByTier *MousePadByTier) UnpackSafeERC20FailedOperationError(raw []byte) (*MousePadByTierSafeERC20FailedOperation, error) {
	out := new(MousePadByTierSafeERC20FailedOperation)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "SafeERC20FailedOperation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierSaleNotEnded represents a SaleNotEnded error raised by the MousePadByTier contract.
type MousePadByTierSaleNotEnded struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SaleNotEnded()
func MousePadByTierSaleNotEndedErrorID() common.Hash {
	return common.HexToHash("0x9d98b04b43f6291fabbd02ee26b0f4b482edd819d4d60314b840e7a65dd6e624")
}

// UnpackSaleNotEndedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SaleNotEnded()
func (mousePadByTier *MousePadByTier) UnpackSaleNotEndedError(raw []byte) (*MousePadByTierSaleNotEnded, error) {
	out := new(MousePadByTierSaleNotEnded)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "SaleNotEnded", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierSaleStarted represents a SaleStarted error raised by the MousePadByTier contract.
type MousePadByTierSaleStarted struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SaleStarted()
func MousePadByTierSaleStartedErrorID() common.Hash {
	return common.HexToHash("0x912ee23dde46ec889d6748212cce445d667f7041597691dc89e8549ad8bc0acb")
}

// UnpackSaleStartedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SaleStarted()
func (mousePadByTier *MousePadByTier) UnpackSaleStartedError(raw []byte) (*MousePadByTierSaleStarted, error) {
	out := new(MousePadByTierSaleStarted)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "SaleStarted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierTierLimitExceeded represents a TierLimitExceeded error raised by the MousePadByTier contract.
type MousePadByTierTierLimitExceeded struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TierLimitExceeded()
func MousePadByTierTierLimitExceededErrorID() common.Hash {
	return common.HexToHash("0xf3387071c3c8998c78fe47b8bca07eda46ce06a04a726e4118c2de4026b95b9c")
}

// UnpackTierLimitExceededError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TierLimitExceeded()
func (mousePadByTier *MousePadByTier) UnpackTierLimitExceededError(raw []byte) (*MousePadByTierTierLimitExceeded, error) {
	out := new(MousePadByTierTierLimitExceeded)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "TierLimitExceeded", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierUserLimitExceeded represents a UserLimitExceeded error raised by the MousePadByTier contract.
type MousePadByTierUserLimitExceeded struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UserLimitExceeded()
func MousePadByTierUserLimitExceededErrorID() common.Hash {
	return common.HexToHash("0xc3dd4aea0ecd7374258f7f38ecda8af47dfd496e1b8b3451ca758d91327d3961")
}

// UnpackUserLimitExceededError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UserLimitExceeded()
func (mousePadByTier *MousePadByTier) UnpackUserLimitExceededError(raw []byte) (*MousePadByTierUserLimitExceeded, error) {
	out := new(MousePadByTierUserLimitExceeded)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "UserLimitExceeded", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePadByTierZeroAmount represents a ZeroAmount error raised by the MousePadByTier contract.
type MousePadByTierZeroAmount struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroAmount()
func MousePadByTierZeroAmountErrorID() common.Hash {
	return common.HexToHash("0x1f2a2005cb66a8e145327e8814e243a0996aec9bfe1e15a495778b1236dbd485")
}

// UnpackZeroAmountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroAmount()
func (mousePadByTier *MousePadByTier) UnpackZeroAmountError(raw []byte) (*MousePadByTierZeroAmount, error) {
	out := new(MousePadByTierZeroAmount)
	if err := mousePadByTier.abi.UnpackIntoInterface(out, "ZeroAmount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

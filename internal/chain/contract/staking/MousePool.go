// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package staking

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

// MousePoolMetaData contains all meta data concerning the MousePool contract.
var MousePoolMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"mouseToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_masterChef\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_treasury\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BOOST_WEIGHT_LIMIT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DURATION_FACTOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DURATION_FACTOR_OVERDUE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"FEE_RATE_PRECISION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_LOCK_DURATION_LIMIT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_OVERDUE_FEE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_PERFORMANCE_FEE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_WITHDRAW_FEE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_WITHDRAW_FEE_PERIOD\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_DEPOSIT_AMOUNT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_LOCK_DURATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MOUSE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UNLOCK_FREE_DURATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"admin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"availableDepositFee\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"availableOverdueFee\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"availableWithdrawFee\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"boostContract\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"boostWeight\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lockDuration\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"freeOverdueFee\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"freePerformanceFee\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"freeWithdrawFee\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPricePerFullShare\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"masterChef\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractMouseMasterChef\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"masterChefPid\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxLockDuration\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"overdueFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingReward\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"performanceFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"performanceFeeContract\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recoverToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAdmin\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBoostContract\",\"inputs\":[{\"name\":\"newBoostContract\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBoostWeight\",\"inputs\":[{\"name\":\"weight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMasterChefPid\",\"inputs\":[{\"name\":\"pid\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMaxLockDuration\",\"inputs\":[{\"name\":\"duration\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOperator\",\"inputs\":[{\"name\":\"newOperator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOverdueFee\",\"inputs\":[{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPerformanceFee\",\"inputs\":[{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPerformanceFeeContract\",\"inputs\":[{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTreasury\",\"inputs\":[{\"name\":\"newTreasury\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setVMouse\",\"inputs\":[{\"name\":\"newVMouse\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWhiteList\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeType\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"status\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWithdrawFee\",\"inputs\":[{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWithdrawFeeContract\",\"inputs\":[{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWithdrawFeePeriod\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalBoostDebt\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalLockedAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalShares\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"treasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unlock\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"userInfo\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"userBoostedShare\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lockedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lockEndTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lockStartTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"locked\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"mouseAtLastUserAction\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lastUserActionTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"vMouse\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawAll\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawByAmount\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawByShares\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawFeeContract\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawFeePeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AdminUpdated\",\"inputs\":[{\"name\":\"oldAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BoostContractUpdated\",\"inputs\":[{\"name\":\"oldBoostContract\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newBoostContract\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"lockDuration\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Lock\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"lockedAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"lockDuration\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MasterChefPidUpdated\",\"inputs\":[{\"name\":\"oldPid\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newPid\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorUpdated\",\"inputs\":[{\"name\":\"oldOperator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOperator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OverdueFeeUpdated\",\"inputs\":[{\"name\":\"oldFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PerformanceFeeUpdated\",\"inputs\":[{\"name\":\"oldFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenRecovered\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TreasuryUpdated\",\"inputs\":[{\"name\":\"oldTreasury\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newTreasury\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unlock\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"lockedAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"boostSharesCleared\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WhiteListUpdated\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"feeType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"status\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawAll\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawFeeUpdated\",\"inputs\":[{\"name\":\"oldFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CannotRecoverMouse\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DepositTooSmall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FeeExceedsMax\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientShares\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LockDurationTooLong\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LockDurationTooShort\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LockNotExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotOwnerOrAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotOwnerOrOperator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ParameterExceedsMax\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"WithdrawWhileLocked\",\"inputs\":[]}]",
	ID:  "MousePool",
	Bin: "0x60c06040526301e1338060095564e8d4a51000600a5560c8600b5560c8600c55600a600d556014600e556203f480600f5564e8d4a51000601055348015610044575f5ffd5b50604051614fce380380614fce8339818101604052810190610066919061032d565b335f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100d7575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016100ce91906103a0565b60405180910390fd5b6100e6816101d860201b60201c565b508373ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250508273ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff16815250508160075f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508060065f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505050506103b9565b60015f6101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905561020b8161020e60201b60201c565b50565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6102fc826102d3565b9050919050565b61030c816102f2565b8114610316575f5ffd5b50565b5f8151905061032781610303565b92915050565b5f5f5f5f60808587031215610345576103446102cf565b5b5f61035287828801610319565b945050602061036387828801610319565b935050604061037487828801610319565b925050606061038587828801610319565b91505092959194509250565b61039a816102f2565b82525050565b5f6020820190506103b35f830184610391565b92915050565b60805160a051614b9261043c5f395f8181610f3901528181611a0a01528181611c9c0152818161325c01528181613ee7015261419b01525f8181611d40015281816122530152818161237b015281816129fa01528181613537015281816139c801528181613a1401528181613a6401528181613b2c01526140f90152614b925ff3fe608060405234801561000f575f5ffd5b5060043610610410575f3560e01c8063853828b61161021e578063c62c610b1161012e578063e4b37ef5116100b6578063f2fde38b11610085578063f2fde38b14610bad578063f4ca4d1d14610bc9578063f786b95814610be7578063f851a44014610c03578063fd253b6414610c2157610410565b8063e4b37ef514610b37578063e73008bc14610b55578063e941fa7814610b73578063f0f4426014610b9157610410565b8063df10b4e6116100fd578063df10b4e614610aa3578063dfcedeee14610ac1578063e2bbb15814610adf578063e30c397814610afb578063e464c62314610b1957610410565b8063c62c610b14610a1b578063d4b0de2f14610a39578063d748a9a214610a57578063def7869d14610a8757610410565b8063a7229fd9116101b1578063b3ab15fb11610180578063b3ab15fb14610979578063b6ac642a14610995578063b9c62706146109b1578063bb9f408d146109e1578063bdca9165146109fd57610410565b8063a7229fd914610903578063aaada5da1461091f578063acaf88cd1461093d578063b352d2121461095b57610410565b806393c99e6a116101ed57806393c99e6a1461088d578063948a03f2146108a9578063a16cdbb1146108c7578063a5834e06146108e557610410565b8063853828b61461082b57806387788782146108355780638bb34bab146108535780638da5cb5b1461086f57610410565b80633eb7887411610324578063704b6c02116102ac5780637340e8891161027b5780637340e889146107ab57806377c7b8fc146107db57806378b4330f146107f957806379ba5097146108175780638456cb591461082157610410565b8063704b6c021461074b57806370897b2314610767578063715018a614610783578063722713f71461078d57610410565b8063570ca735116102f3578063570ca735146106b5578063575a86b2146106d35780635c975abb146106f157806361d027b31461070f5780636edce5ca1461072d57610410565b80633eb78874146106415780633f4ba83a1461065f57806351c068a1146106695780635521e9bf1461069957610410565b80631959a002116103a75780632cfc5f01116103765780632cfc5f01146105af5780632f6c493c146105cd57806339340d22146105e95780633a98ef39146106055780633d7bba0c1461062357610410565b80631959a002146105225780631ea30fef146105595780631efac1b8146105775780632ad124351461059357610410565b80630c59696b116103e35780630c59696b1461049c5780631250e7b7146104b8578063137ee36e146104e857806314ff30391461050657610410565b806301e81326146104145780630586755a1461043257806305a9f2741461044e5780630bf235531461046c575b5f5ffd5b61041c610c3f565b604051610429919061441f565b60405180910390f35b61044c60048036038101906104479190614466565b610c47565b005b610456610d58565b604051610463919061441f565b60405180910390f35b610486600480360381019061048191906144eb565b610d5e565b604051610493919061441f565b60405180910390f35b6104b660048036038101906104b19190614466565b610dd5565b005b6104d260048036038101906104cd91906144eb565b610f25565b6040516104df919061441f565b60405180910390f35b6104f0610f36565b6040516104fd919061441f565b60405180910390f35b610520600480360381019061051b9190614466565b610fd8565b005b61053c600480360381019061053791906144eb565b6110e5565b604051610550989796959493929190614530565b60405180910390f35b61056161117c565b60405161056e919061441f565b60405180910390f35b610591600480360381019061058c9190614466565b611187565b005b6105ad60048036038101906105a891906144eb565b611295565b005b6105b761139f565b6040516105c4919061441f565b60405180910390f35b6105e760048036038101906105e291906144eb565b6113a6565b005b61060360048036038101906105fe919061460c565b611669565b005b61060d6118ae565b60405161061a919061441f565b60405180910390f35b61062b6118b4565b604051610638919061466b565b60405180910390f35b6106496118d9565b604051610656919061441f565b60405180910390f35b6106676118df565b005b610683600480360381019061067e91906144eb565b6119b0565b6040516106909190614684565b60405180910390f35b6106b360048036038101906106ae9190614466565b6119cd565b005b6106bd6119e3565b6040516106ca919061466b565b60405180910390f35b6106db611a08565b6040516106e891906146f8565b60405180910390f35b6106f9611a2c565b6040516107069190614684565b60405180910390f35b610717611a42565b604051610724919061466b565b60405180910390f35b610735611a67565b604051610742919061441f565b60405180910390f35b610765600480360381019061076091906144eb565b611a6d565b005b610781600480360381019061077c9190614466565b611b38565b005b61078b611c85565b005b610795611c98565b6040516107a2919061441f565b60405180910390f35b6107c560048036038101906107c091906144eb565b611df0565b6040516107d29190614684565b60405180910390f35b6107e3611e0d565b6040516107f0919061441f565b60405180910390f35b610801611e62565b60405161080e919061441f565b60405180910390f35b61081f611e69565b005b610829611ef7565b005b610833611fc8565b005b61083d612022565b60405161084a919061441f565b60405180910390f35b61086d60048036038101906108689190614466565b612028565b005b61087761203e565b604051610884919061466b565b60405180910390f35b6108a760048036038101906108a29190614466565b612065565b005b6108b1612175565b6040516108be919061441f565b60405180910390f35b6108cf61217e565b6040516108dc919061441f565b60405180910390f35b6108ed612184565b6040516108fa919061441f565b60405180910390f35b61091d60048036038101906109189190614711565b61218a565b005b61092761236b565b604051610934919061441f565b60405180910390f35b610945612372565b604051610952919061441f565b60405180910390f35b610963612379565b6040516109709190614781565b60405180910390f35b610993600480360381019061098e91906144eb565b61239d565b005b6109af60048036038101906109aa9190614466565b612468565b005b6109cb60048036038101906109c691906144eb565b6125b5565b6040516109d8919061441f565b60405180910390f35b6109fb60048036038101906109f69190614466565b6125c6565b005b610a056126d3565b604051610a12919061441f565b60405180910390f35b610a236126d9565b604051610a30919061441f565b60405180910390f35b610a416126df565b604051610a4e919061441f565b60405180910390f35b610a716004803603810190610a6c91906144eb565b6126e5565b604051610a7e9190614684565b60405180910390f35b610aa16004803603810190610a9c91906144eb565b612702565b005b610aab61288c565b604051610ab8919061441f565b60405180910390f35b610ac9612892565b604051610ad6919061466b565b60405180910390f35b610af96004803603810190610af4919061479a565b6128b7565b005b610b03612cb5565b604051610b10919061466b565b60405180910390f35b610b21612cdd565b604051610b2e919061441f565b60405180910390f35b610b3f612ce5565b604051610b4c919061441f565b60405180910390f35b610b5d612ceb565b604051610b6a919061441f565b60405180910390f35b610b7b612cf1565b604051610b88919061441f565b60405180910390f35b610bab6004803603810190610ba691906144eb565b612cf7565b005b610bc76004803603810190610bc291906144eb565b612dc2565b005b610bd1612e6e565b604051610bde919061441f565b60405180910390f35b610c016004803603810190610bfc9190614466565b612e74565b005b610c0b612f83565b604051610c18919061466b565b60405180910390f35b610c29612fa8565b604051610c36919061441f565b60405180910390f35b6305265c0081565b610c4f61203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614158015610cd7575060065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b15610d0e576040517fdce3812500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6002549050816002819055507f3de26cad46c6b48fb1c8dc0f0f5f92316802389f907e74b04768e62e6e062f9f8183604051610d4c9291906147d8565b60405180910390a15050565b60055481565b5f5f60135f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2090505f815f015403610db2575f915050610dd0565b5f610dbf825f0154612fb1565b9050610dcb8482612ffd565b925050505b919050565b610ddd61203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614158015610e65575060065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b15610e9c576040517fdce3812500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b64e8d4a51000811115610edb576040517f5ff85e3f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6010549050816010819055507f5783a2479544f38cf6d14a081184e49d1388a0c6e2d13080ed3f827bba9f1f468183604051610f199291906147d8565b60405180910390a15050565b5f610f2f826130f6565b9050919050565b5f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166313590aa3600254306040518363ffffffff1660e01b8152600401610f949291906147ff565b602060405180830381865afa158015610faf573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610fd3919061483a565b905090565b610fe061203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614158015611068575060065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b1561109f576040517fdce3812500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6101f48111156110db576040517f5ff85e3f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80600e8190555050565b5f5f5f5f5f5f5f5f5f60135f8b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f209050805f01548160010154826002015483600301548460040154856005015f9054906101000a900460ff16866006015487600701549850985098509850985098509850985050919395975091939597565b662386f26fc1000081565b61118f61203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614158015611217575060065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b1561124e576040517fdce3812500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b62093a8081111561128b576040517fdda1d95400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80600f8190555050565b61129d61203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614158015611325575060065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b1561135c576040517fdce3812500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060125f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b62093a8081565b6113ae613219565b8073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614158015611437575060085f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b8015611476575061144661203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b156114ad576040517f98f76d0e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60135f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f209050806003015442101561152b576040517f6855a80200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806005015f9054906101000a900460ff166115465750611666565b61154e61325a565b5f611558836132e6565b90505f81111561156d5761156c838261346f565b5b5f82600101549050826001015460045f82825461158a9190614892565b925050819055505f8360010181905550826002015460055f8282546115af9190614892565b925050819055505f836005015f6101000a81548160ff0219169083151502179055505f83600201819055505f83600301819055505f6115f0845f0154612fb1565b905080846006018190555042846007018190555061160d856135a0565b8473ffffffffffffffffffffffffffffffffffffffff167ff7870c5b224cbc19873599e46ccfc7103934650509b1af0c3ce90138377c20048560020154846040516116599291906147d8565b60405180910390a2505050505b50565b61167161203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141580156116f9575060065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b15611730576040517fdce3812500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8260ff1603611793578060145f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff021916908315150217905550611859565b60018260ff16036117f7578060155f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff021916908315150217905550611858565b60028260ff1603611857578060165f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055505b5b5b8273ffffffffffffffffffffffffffffffffffffffff167fd4d43edabfa41e2128f2062be4847db83d2bb5a81242b80d6acc6c1c005baefe83836040516118a19291906148d4565b60405180910390a2505050565b60035481565b60125f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600c5481565b6118e761203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415801561196f575060065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b156119a6576040517fdce3812500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6119ae61378c565b565b6015602052805f5260405f205f915054906101000a900460ff1681565b6119d5613219565b6119e0815f5f6137ee565b50565b60085f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b7f000000000000000000000000000000000000000000000000000000000000000081565b5f600160149054906101000a900460ff16905090565b60075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600a5481565b611a75613c6c565b5f60065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160065f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f101b8081ff3b56bbf45deb824d86a3b0fd38b7e3dd42421105cf8abe9106db0b60405160405180910390a35050565b611b4061203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614158015611bc8575060065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b15611bff576040517fdce3812500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6101f4811115611c3b576040517f5ff85e3f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f600b54905081600b819055507f607b1c943753982194530bf7133a5972ea2626e028005410efa54ab20035caf88183604051611c799291906147d8565b60405180910390a15050565b611c8d613c6c565b611c965f613cf3565b565b5f5f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166393f1a40b600254306040518363ffffffff1660e01b8152600401611cf79291906147ff565b606060405180830381865afa158015611d12573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611d3691906148fb565b50509050600454817f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401611d97919061466b565b602060405180830381865afa158015611db2573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611dd6919061483a565b611de0919061494b565b611dea919061494b565b91505090565b6014602052805f5260405f205f915054906101000a900460ff1681565b5f5f60035403611e2757670de0b6b3a76400009050611e5f565b600354670de0b6b3a7640000600454611e3e611c98565b611e489190614892565b611e52919061497e565b611e5c91906149ec565b90505b90565b62093a8081565b5f611e72613d23565b90508073ffffffffffffffffffffffffffffffffffffffff16611e93612cb5565b73ffffffffffffffffffffffffffffffffffffffff1614611eeb57806040517f118cdaa7000000000000000000000000000000000000000000000000000000008152600401611ee2919061466b565b60405180910390fd5b611ef481613cf3565b50565b611eff61203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614158015611f87575060065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b15611fbe576040517fdce3812500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611fc6613d2a565b565b611fd0613219565b5f60135f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20905061201f5f825f015460016137ee565b50565b600b5481565b612030613219565b61203b5f825f6137ee565b50565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b61206d61203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141580156120f5575060065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b1561212c576040517fdce3812500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b64e8d4a5100081111561216b576040517fdda1d95400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80600a8190555050565b64e8d4a5100081565b60095481565b60105481565b61219261203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415801561221a575060065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b15612251576040517fdce3812500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036122d6576040517fa0145ecd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61230182828573ffffffffffffffffffffffffffffffffffffffff16613d8c9092919063ffffffff16565b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f879f92dded0f26b83c3e00b12e0395dc72cfc3077343d1854ed6988edd1f90968360405161235e919061441f565b60405180910390a3505050565b62278d0081565b62ed4e0081565b7f000000000000000000000000000000000000000000000000000000000000000081565b6123a5613c6c565b5f60085f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160085f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167ffbe5b6cbafb274f445d7fed869dc77a838d8243a22c460de156560e8857cad0360405160405180910390a35050565b61247061203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141580156124f8575060065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b1561252f576040517fdce3812500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6101f481111561256b576040517f5ff85e3f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f600d54905081600d819055507f733071ab8253b372ed26a6d1b04aec71c4bfcd209c93397df32bb77478cdd2c881836040516125a99291906147d8565b60405180910390a15050565b5f6125bf826132e6565b9050919050565b6125ce61203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614158015612656575060065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b1561268d576040517fdce3812500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6101f48111156126c9576040517f5ff85e3f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80600c8190555050565b6101f481565b60025481565b6101f481565b6016602052805f5260405f205f915054906101000a900460ff1681565b61270a61203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614158015612792575060065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b156127c9576040517fdce3812500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60115f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160115f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f5365992143af79f25219eb8fb9f3b26382aae3916feea1ccff3daa8ef82f5f0c60405160405180910390a35050565b600f5481565b60115f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6128bf613219565b5f821180156128d45750662386f26fc1000082105b1561290b576040517f6ba4a1c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f81111561298d5762093a80811015612950576040517f49eeb0b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60095481111561298c576040517ff761f1cd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5b61299561325a565b5f60135f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2090505f815f015411156129ea576129e933613ddf565b5b5f831115612a4057612a3f3330857f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16613e07909392919063ffffffff16565b5b5f5f5f851115612a5657612a5385613e5c565b91505b5f841115612be3575f8490505f869050846005015f9054906101000a900460ff1615612af857846001015460045f828254612a919190614892565b925050819055505f85600101819055505f42866003015411612ab3575f612ac4565b428660030154612ac39190614892565b5b90508681612ad2919061494b565b9250600954831115612ae45760095492505b878660020154612af4919061494b565b9150505b612b028183613ea8565b92506001856005015f6101000a81548160ff0219169083151502179055508085600201819055508142612b35919061494b565b856003018190555042856004018190555082856001015f828254612b59919061494b565b925050819055508260045f828254612b71919061494b565b925050819055508660055f828254612b89919061494b565b925050819055503373ffffffffffffffffffffffffffffffffffffffff167f49eaf4942f1237055eb4cfa5f31c9dfe50d5b4ade01e021f7de8be2fbbde557b8284604051612bd89291906147d8565b60405180910390a250505b81835f015f828254612bf5919061494b565b925050819055508160035f828254612c0d919061494b565b925050819055505f851115612c2657612c2585613ee5565b5b5f612c33845f0154612fb1565b9050808460060181905550428460070181905550612c5433855f0154613f72565b612c5d336135a0565b3373ffffffffffffffffffffffffffffffffffffffff167f90890809c654f11d6e72a28fa60149770a0d11ec6c92319d6ceb2bb0a4ea1a158787604051612ca59291906147d8565b60405180910390a2505050505050565b5f60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6301e1338081565b600e5481565b60045481565b600d5481565b612cff613c6c565b5f60075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160075f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f4ab5be82436d353e61ca18726e984e561f5c1cc7c6d38b29d2553c790434705a60405160405180910390a35050565b612dca613c6c565b8060015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff16612e2961203e565b73ffffffffffffffffffffffffffffffffffffffff167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a350565b61271081565b612e7c61203e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614158015612f04575060065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b15612f3b576040517fdce3812500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6305265c00811115612f79576040517fdda1d95400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060098190555050565b60065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b64e8d4a5100081565b5f5f60035403612fc3575f9050612ff8565b5f600454612fcf611c98565b612fd99190614892565b90506003548184612fea919061497e565b612ff491906149ec565b9150505b919050565b5f60155f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1615613055575f90506130f0565b5f60135f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f209050600f5481600701546130a7919061494b565b42106130b6575f9150506130f0565b5f6130c0856140e5565b6130cc57600d546130d0565b600e545b905061271081856130e1919061497e565b6130eb91906149ec565b925050505b92915050565b5f60145f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff161561314e575f9050613214565b5f60135f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2090505f815f0154036131a1575f915050613214565b5f6131ae825f0154612fb1565b9050816006015481116131c5575f92505050613214565b5f8260060154826131d69190614892565b90505f6131e2866140e5565b6131ee57600b546131f2565b600c545b90506127108183613203919061497e565b61320d91906149ec565b9450505050505b919050565b613221611a2c565b15613258576040517fd93c066500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663e2bbb1586002545f6040518363ffffffff1660e01b81526004016132b7929190614a55565b5f604051808303815f87803b1580156132ce575f5ffd5b505af11580156132e0573d5f5f3e3d5ffd5b50505050565b5f60165f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff161561333e575f905061346a565b5f60135f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f209050806005015f9054906101000a900460ff1661339c575f91505061346a565b5f62278d0082600301546133b0919061494b565b90508042116133c3575f9250505061346a565b5f81426133d09190614892565b90505f6133df845f0154612fb1565b90505f846006015482116133f3575f613404565b8460060154826134039190614892565b5b90505f62ed4e008411613417578361341c565b62ed4e005b90506402540be400606462ed4e00613434919061497e565b61343e919061497e565b816010548461344d919061497e565b613457919061497e565b61346191906149ec565b96505050505050505b919050565b5f60135f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2090505f6134b983613e5c565b90505f81036134c957505061359c565b80825f015f8282546134db9190614892565b925050819055508060035f8282546134f39190614892565b925050819055505f61350482612fb1565b905061350f816140f6565b61357b60075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16827f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16613d8c9092919063ffffffff16565b613587835f0154612fb1565b83600601819055504283600701819055505050505b5050565b5f73ffffffffffffffffffffffffffffffffffffffff1660115f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160315613789575f60135f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2090505f60115f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16838360020154846005015f9054906101000a900460ff1661368f575f6136a0565b42856003015461369f9190614892565b5b6005546009546040516024016136ba959493929190614a7c565b6040516020818303038152906040527f033fbcf0000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040516137449190614b1f565b5f604051808303815f865af19150503d805f811461377d576040519150601f19603f3d011682016040523d82523d5f602084013e613782565b606091505b5050905050505b50565b613794614233565b5f600160146101000a81548160ff0219169083151502179055507f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6137d7613d23565b6040516137e4919061466b565b60405180910390a1565b5f60135f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f209050806005015f9054906101000a900460ff16801561384e5750806003015442105b15613885576040517fc653938f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61388d61325a565b61389633613ddf565b81156138a757805f015492506138f9565b5f8411156138f8576138b884613e5c565b9250805f01548311156138f7576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5b5b805f0154831115613936576040517f3999656700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f61394084612fb1565b90505f61394d3383612ffd565b90505f818361395c9190614892565b905085845f015f8282546139709190614892565b925050819055508560035f8282546139889190614892565b92505081905550613998836140f6565b5f821115613a0d57613a0c60075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16837f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16613d8c9092919063ffffffff16565b5b613a5833827f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16613d8c9092919063ffffffff16565b5f60035403613b73575f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401613abb919061466b565b602060405180830381865afa158015613ad6573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190613afa919061483a565b90505f811115613b7157613b7060075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16827f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16613d8c9092919063ffffffff16565b5b505b5f845f01541115613b9757613b8a845f0154612fb1565b8460060181905550613ba1565b5f84600601819055505b428460070181905550613bb733855f0154613f72565b8415613c12573373ffffffffffffffffffffffffffffffffffffffff167fae7bb41990cd336548a7c61b3ffa4a8207fe5e886087c23f5739a99114be78398288604051613c059291906147d8565b60405180910390a2613c63565b3373ffffffffffffffffffffffffffffffffffffffff167ff279e6a1f5e320cca91135676d9cb6e44ca8a08c0b88342bcdb1144f6511b5688288604051613c5a9291906147d8565b60405180910390a25b50505050505050565b613c74613d23565b73ffffffffffffffffffffffffffffffffffffffff16613c9261203e565b73ffffffffffffffffffffffffffffffffffffffff1614613cf157613cb5613d23565b6040517f118cdaa7000000000000000000000000000000000000000000000000000000008152600401613ce8919061466b565b60405180910390fd5b565b60015f6101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055613d2081614273565b50565b5f33905090565b613d32613219565b60018060146101000a81548160ff0219169083151502179055507f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258613d75613d23565b604051613d82919061466b565b60405180910390a1565b613d998383836001614334565b613dda57826040517f5274afe7000000000000000000000000000000000000000000000000000000008152600401613dd1919061466b565b60405180910390fd5b505050565b5f613de9826130f6565b90505f8103613df85750613e04565b613e02828261346f565b505b50565b613e15848484846001614396565b613e5657836040517f5274afe7000000000000000000000000000000000000000000000000000000008152600401613e4d919061466b565b60405180910390fd5b50505050565b5f5f60035403613e6e57819050613ea3565b5f600454613e7a611c98565b613e849190614892565b90508060035484613e95919061497e565b613e9f91906149ec565b9150505b919050565b5f60646301e13380613eba919061497e565b600a548385613ec9919061497e565b613ed3919061497e565b613edd91906149ec565b905092915050565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663e2bbb158600254836040518363ffffffff1660e01b8152600401613f429291906147d8565b5f604051808303815f87803b158015613f59575f5ffd5b505af1158015613f6b573d5f5f3e3d5ffd5b5050505050565b5f73ffffffffffffffffffffffffffffffffffffffff1660125f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603156140e1575f60125f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168383604051602401614013929190614b35565b6040516020818303038152906040527f47e7ef24000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505060405161409d9190614b1f565b5f604051808303815f865af19150503d805f81146140d6576040519150601f19603f3d011682016040523d82523d5f602084013e6140db565b606091505b50509050505b5050565b5f5f823b90505f8111915050919050565b5f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401614150919061466b565b602060405180830381865afa15801561416b573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061418f919061483a565b90508181101561422f577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663441a3e7060025483856141e49190614892565b6040518363ffffffff1660e01b81526004016142019291906147d8565b5f604051808303815f87803b158015614218575f5ffd5b505af115801561422a573d5f5f3e3d5ffd5b505050505b5050565b61423b611a2c565b614271576040517f8dfc202b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f5f63a9059cbb60e01b9050604051815f525f1960601c86166004528460245260205f60445f5f8b5af1925060015f5114831661438857838315161561437c573d5f823e3d81fd5b5f873b113d1516831692505b806040525050949350505050565b5f5f6323b872dd60e01b9050604051815f525f1960601c87166004525f1960601c86166024528460445260205f60645f5f8c5af1925060015f511483166143f45783831516156143e8573d5f823e3d81fd5b5f883b113d1516831692505b806040525f606052505095945050505050565b5f819050919050565b61441981614407565b82525050565b5f6020820190506144325f830184614410565b92915050565b5f5ffd5b61444581614407565b811461444f575f5ffd5b50565b5f813590506144608161443c565b92915050565b5f6020828403121561447b5761447a614438565b5b5f61448884828501614452565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6144ba82614491565b9050919050565b6144ca816144b0565b81146144d4575f5ffd5b50565b5f813590506144e5816144c1565b92915050565b5f60208284031215614500576144ff614438565b5b5f61450d848285016144d7565b91505092915050565b5f8115159050919050565b61452a81614516565b82525050565b5f610100820190506145445f83018b614410565b614551602083018a614410565b61455e6040830189614410565b61456b6060830188614410565b6145786080830187614410565b61458560a0830186614521565b61459260c0830185614410565b61459f60e0830184614410565b9998505050505050505050565b5f60ff82169050919050565b6145c1816145ac565b81146145cb575f5ffd5b50565b5f813590506145dc816145b8565b92915050565b6145eb81614516565b81146145f5575f5ffd5b50565b5f81359050614606816145e2565b92915050565b5f5f5f6060848603121561462357614622614438565b5b5f614630868287016144d7565b9350506020614641868287016145ce565b9250506040614652868287016145f8565b9150509250925092565b614665816144b0565b82525050565b5f60208201905061467e5f83018461465c565b92915050565b5f6020820190506146975f830184614521565b92915050565b5f819050919050565b5f6146c06146bb6146b684614491565b61469d565b614491565b9050919050565b5f6146d1826146a6565b9050919050565b5f6146e2826146c7565b9050919050565b6146f2816146d8565b82525050565b5f60208201905061470b5f8301846146e9565b92915050565b5f5f5f6060848603121561472857614727614438565b5b5f614735868287016144d7565b9350506020614746868287016144d7565b925050604061475786828701614452565b9150509250925092565b5f61476b826146c7565b9050919050565b61477b81614761565b82525050565b5f6020820190506147945f830184614772565b92915050565b5f5f604083850312156147b0576147af614438565b5b5f6147bd85828601614452565b92505060206147ce85828601614452565b9150509250929050565b5f6040820190506147eb5f830185614410565b6147f86020830184614410565b9392505050565b5f6040820190506148125f830185614410565b61481f602083018461465c565b9392505050565b5f815190506148348161443c565b92915050565b5f6020828403121561484f5761484e614438565b5b5f61485c84828501614826565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61489c82614407565b91506148a783614407565b92508282039050818111156148bf576148be614865565b5b92915050565b6148ce816145ac565b82525050565b5f6040820190506148e75f8301856148c5565b6148f46020830184614521565b9392505050565b5f5f5f6060848603121561491257614911614438565b5b5f61491f86828701614826565b935050602061493086828701614826565b925050604061494186828701614826565b9150509250925092565b5f61495582614407565b915061496083614407565b925082820190508082111561497857614977614865565b5b92915050565b5f61498882614407565b915061499383614407565b92508282026149a181614407565b915082820484148315176149b8576149b7614865565b5b5092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f6149f682614407565b9150614a0183614407565b925082614a1157614a106149bf565b5b828204905092915050565b5f819050919050565b5f614a3f614a3a614a3584614a1c565b61469d565b614407565b9050919050565b614a4f81614a25565b82525050565b5f604082019050614a685f830185614410565b614a756020830184614a46565b9392505050565b5f60a082019050614a8f5f83018861465c565b614a9c6020830187614410565b614aa96040830186614410565b614ab66060830185614410565b614ac36080830184614410565b9695505050505050565b5f81519050919050565b5f81905092915050565b8281835e5f83830152505050565b5f614af982614acd565b614b038185614ad7565b9350614b13818560208601614ae1565b80840191505092915050565b5f614b2a8284614aef565b915081905092915050565b5f604082019050614b485f83018561465c565b614b556020830184614410565b939250505056fea2646970667358221220ed0e079f12d1cd29c904bc4fc8df0619facb07ef97f50becdde440eae6aea61d64736f6c634300081c0033",
}

// MousePool is an auto generated Go binding around an Ethereum contract.
type MousePool struct {
	abi abi.ABI
}

// NewMousePool creates a new instance of MousePool.
func NewMousePool() *MousePool {
	parsed, err := MousePoolMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &MousePool{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *MousePool) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address mouseToken, address _masterChef, address _treasury, address _admin) returns()
func (mousePool *MousePool) PackConstructor(mouseToken common.Address, _masterChef common.Address, _treasury common.Address, _admin common.Address) []byte {
	enc, err := mousePool.abi.Pack("", mouseToken, _masterChef, _treasury, _admin)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackBOOSTWEIGHTLIMIT is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfd253b64.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function BOOST_WEIGHT_LIMIT() view returns(uint256)
func (mousePool *MousePool) PackBOOSTWEIGHTLIMIT() []byte {
	enc, err := mousePool.abi.Pack("BOOST_WEIGHT_LIMIT")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBOOSTWEIGHTLIMIT is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfd253b64.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function BOOST_WEIGHT_LIMIT() view returns(uint256)
func (mousePool *MousePool) TryPackBOOSTWEIGHTLIMIT() ([]byte, error) {
	return mousePool.abi.Pack("BOOST_WEIGHT_LIMIT")
}

// UnpackBOOSTWEIGHTLIMIT is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfd253b64.
//
// Solidity: function BOOST_WEIGHT_LIMIT() view returns(uint256)
func (mousePool *MousePool) UnpackBOOSTWEIGHTLIMIT(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("BOOST_WEIGHT_LIMIT", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDURATIONFACTOR is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe464c623.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DURATION_FACTOR() view returns(uint256)
func (mousePool *MousePool) PackDURATIONFACTOR() []byte {
	enc, err := mousePool.abi.Pack("DURATION_FACTOR")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDURATIONFACTOR is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe464c623.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DURATION_FACTOR() view returns(uint256)
func (mousePool *MousePool) TryPackDURATIONFACTOR() ([]byte, error) {
	return mousePool.abi.Pack("DURATION_FACTOR")
}

// UnpackDURATIONFACTOR is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe464c623.
//
// Solidity: function DURATION_FACTOR() view returns(uint256)
func (mousePool *MousePool) UnpackDURATIONFACTOR(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("DURATION_FACTOR", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDURATIONFACTOROVERDUE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xacaf88cd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DURATION_FACTOR_OVERDUE() view returns(uint256)
func (mousePool *MousePool) PackDURATIONFACTOROVERDUE() []byte {
	enc, err := mousePool.abi.Pack("DURATION_FACTOR_OVERDUE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDURATIONFACTOROVERDUE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xacaf88cd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DURATION_FACTOR_OVERDUE() view returns(uint256)
func (mousePool *MousePool) TryPackDURATIONFACTOROVERDUE() ([]byte, error) {
	return mousePool.abi.Pack("DURATION_FACTOR_OVERDUE")
}

// UnpackDURATIONFACTOROVERDUE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xacaf88cd.
//
// Solidity: function DURATION_FACTOR_OVERDUE() view returns(uint256)
func (mousePool *MousePool) UnpackDURATIONFACTOROVERDUE(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("DURATION_FACTOR_OVERDUE", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackFEERATEPRECISION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf4ca4d1d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function FEE_RATE_PRECISION() view returns(uint256)
func (mousePool *MousePool) PackFEERATEPRECISION() []byte {
	enc, err := mousePool.abi.Pack("FEE_RATE_PRECISION")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFEERATEPRECISION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf4ca4d1d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function FEE_RATE_PRECISION() view returns(uint256)
func (mousePool *MousePool) TryPackFEERATEPRECISION() ([]byte, error) {
	return mousePool.abi.Pack("FEE_RATE_PRECISION")
}

// UnpackFEERATEPRECISION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf4ca4d1d.
//
// Solidity: function FEE_RATE_PRECISION() view returns(uint256)
func (mousePool *MousePool) UnpackFEERATEPRECISION(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("FEE_RATE_PRECISION", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMAXLOCKDURATIONLIMIT is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01e81326.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_LOCK_DURATION_LIMIT() view returns(uint256)
func (mousePool *MousePool) PackMAXLOCKDURATIONLIMIT() []byte {
	enc, err := mousePool.abi.Pack("MAX_LOCK_DURATION_LIMIT")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXLOCKDURATIONLIMIT is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01e81326.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_LOCK_DURATION_LIMIT() view returns(uint256)
func (mousePool *MousePool) TryPackMAXLOCKDURATIONLIMIT() ([]byte, error) {
	return mousePool.abi.Pack("MAX_LOCK_DURATION_LIMIT")
}

// UnpackMAXLOCKDURATIONLIMIT is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01e81326.
//
// Solidity: function MAX_LOCK_DURATION_LIMIT() view returns(uint256)
func (mousePool *MousePool) UnpackMAXLOCKDURATIONLIMIT(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("MAX_LOCK_DURATION_LIMIT", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMAXOVERDUEFEE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x948a03f2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_OVERDUE_FEE() view returns(uint256)
func (mousePool *MousePool) PackMAXOVERDUEFEE() []byte {
	enc, err := mousePool.abi.Pack("MAX_OVERDUE_FEE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXOVERDUEFEE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x948a03f2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_OVERDUE_FEE() view returns(uint256)
func (mousePool *MousePool) TryPackMAXOVERDUEFEE() ([]byte, error) {
	return mousePool.abi.Pack("MAX_OVERDUE_FEE")
}

// UnpackMAXOVERDUEFEE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x948a03f2.
//
// Solidity: function MAX_OVERDUE_FEE() view returns(uint256)
func (mousePool *MousePool) UnpackMAXOVERDUEFEE(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("MAX_OVERDUE_FEE", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMAXPERFORMANCEFEE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbdca9165.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_PERFORMANCE_FEE() view returns(uint256)
func (mousePool *MousePool) PackMAXPERFORMANCEFEE() []byte {
	enc, err := mousePool.abi.Pack("MAX_PERFORMANCE_FEE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXPERFORMANCEFEE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbdca9165.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_PERFORMANCE_FEE() view returns(uint256)
func (mousePool *MousePool) TryPackMAXPERFORMANCEFEE() ([]byte, error) {
	return mousePool.abi.Pack("MAX_PERFORMANCE_FEE")
}

// UnpackMAXPERFORMANCEFEE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbdca9165.
//
// Solidity: function MAX_PERFORMANCE_FEE() view returns(uint256)
func (mousePool *MousePool) UnpackMAXPERFORMANCEFEE(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("MAX_PERFORMANCE_FEE", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMAXWITHDRAWFEE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd4b0de2f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_WITHDRAW_FEE() view returns(uint256)
func (mousePool *MousePool) PackMAXWITHDRAWFEE() []byte {
	enc, err := mousePool.abi.Pack("MAX_WITHDRAW_FEE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXWITHDRAWFEE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd4b0de2f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_WITHDRAW_FEE() view returns(uint256)
func (mousePool *MousePool) TryPackMAXWITHDRAWFEE() ([]byte, error) {
	return mousePool.abi.Pack("MAX_WITHDRAW_FEE")
}

// UnpackMAXWITHDRAWFEE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd4b0de2f.
//
// Solidity: function MAX_WITHDRAW_FEE() view returns(uint256)
func (mousePool *MousePool) UnpackMAXWITHDRAWFEE(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("MAX_WITHDRAW_FEE", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMAXWITHDRAWFEEPERIOD is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2cfc5f01.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_WITHDRAW_FEE_PERIOD() view returns(uint256)
func (mousePool *MousePool) PackMAXWITHDRAWFEEPERIOD() []byte {
	enc, err := mousePool.abi.Pack("MAX_WITHDRAW_FEE_PERIOD")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXWITHDRAWFEEPERIOD is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2cfc5f01.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_WITHDRAW_FEE_PERIOD() view returns(uint256)
func (mousePool *MousePool) TryPackMAXWITHDRAWFEEPERIOD() ([]byte, error) {
	return mousePool.abi.Pack("MAX_WITHDRAW_FEE_PERIOD")
}

// UnpackMAXWITHDRAWFEEPERIOD is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2cfc5f01.
//
// Solidity: function MAX_WITHDRAW_FEE_PERIOD() view returns(uint256)
func (mousePool *MousePool) UnpackMAXWITHDRAWFEEPERIOD(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("MAX_WITHDRAW_FEE_PERIOD", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMINDEPOSITAMOUNT is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1ea30fef.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MIN_DEPOSIT_AMOUNT() view returns(uint256)
func (mousePool *MousePool) PackMINDEPOSITAMOUNT() []byte {
	enc, err := mousePool.abi.Pack("MIN_DEPOSIT_AMOUNT")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMINDEPOSITAMOUNT is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1ea30fef.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MIN_DEPOSIT_AMOUNT() view returns(uint256)
func (mousePool *MousePool) TryPackMINDEPOSITAMOUNT() ([]byte, error) {
	return mousePool.abi.Pack("MIN_DEPOSIT_AMOUNT")
}

// UnpackMINDEPOSITAMOUNT is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1ea30fef.
//
// Solidity: function MIN_DEPOSIT_AMOUNT() view returns(uint256)
func (mousePool *MousePool) UnpackMINDEPOSITAMOUNT(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("MIN_DEPOSIT_AMOUNT", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMINLOCKDURATION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x78b4330f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MIN_LOCK_DURATION() view returns(uint256)
func (mousePool *MousePool) PackMINLOCKDURATION() []byte {
	enc, err := mousePool.abi.Pack("MIN_LOCK_DURATION")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMINLOCKDURATION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x78b4330f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MIN_LOCK_DURATION() view returns(uint256)
func (mousePool *MousePool) TryPackMINLOCKDURATION() ([]byte, error) {
	return mousePool.abi.Pack("MIN_LOCK_DURATION")
}

// UnpackMINLOCKDURATION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x78b4330f.
//
// Solidity: function MIN_LOCK_DURATION() view returns(uint256)
func (mousePool *MousePool) UnpackMINLOCKDURATION(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("MIN_LOCK_DURATION", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMOUSE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb352d212.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MOUSE() view returns(address)
func (mousePool *MousePool) PackMOUSE() []byte {
	enc, err := mousePool.abi.Pack("MOUSE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMOUSE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb352d212.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MOUSE() view returns(address)
func (mousePool *MousePool) TryPackMOUSE() ([]byte, error) {
	return mousePool.abi.Pack("MOUSE")
}

// UnpackMOUSE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb352d212.
//
// Solidity: function MOUSE() view returns(address)
func (mousePool *MousePool) UnpackMOUSE(data []byte) (common.Address, error) {
	out, err := mousePool.abi.Unpack("MOUSE", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackUNLOCKFREEDURATION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaaada5da.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function UNLOCK_FREE_DURATION() view returns(uint256)
func (mousePool *MousePool) PackUNLOCKFREEDURATION() []byte {
	enc, err := mousePool.abi.Pack("UNLOCK_FREE_DURATION")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUNLOCKFREEDURATION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaaada5da.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function UNLOCK_FREE_DURATION() view returns(uint256)
func (mousePool *MousePool) TryPackUNLOCKFREEDURATION() ([]byte, error) {
	return mousePool.abi.Pack("UNLOCK_FREE_DURATION")
}

// UnpackUNLOCKFREEDURATION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaaada5da.
//
// Solidity: function UNLOCK_FREE_DURATION() view returns(uint256)
func (mousePool *MousePool) UnpackUNLOCKFREEDURATION(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("UNLOCK_FREE_DURATION", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackAcceptOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ba5097.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function acceptOwnership() returns()
func (mousePool *MousePool) PackAcceptOwnership() []byte {
	enc, err := mousePool.abi.Pack("acceptOwnership")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAcceptOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ba5097.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function acceptOwnership() returns()
func (mousePool *MousePool) TryPackAcceptOwnership() ([]byte, error) {
	return mousePool.abi.Pack("acceptOwnership")
}

// PackAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf851a440.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function admin() view returns(address)
func (mousePool *MousePool) PackAdmin() []byte {
	enc, err := mousePool.abi.Pack("admin")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf851a440.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function admin() view returns(address)
func (mousePool *MousePool) TryPackAdmin() ([]byte, error) {
	return mousePool.abi.Pack("admin")
}

// UnpackAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (mousePool *MousePool) UnpackAdmin(data []byte) (common.Address, error) {
	out, err := mousePool.abi.Unpack("admin", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackAvailableDepositFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1250e7b7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function availableDepositFee(address user) view returns(uint256)
func (mousePool *MousePool) PackAvailableDepositFee(user common.Address) []byte {
	enc, err := mousePool.abi.Pack("availableDepositFee", user)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAvailableDepositFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1250e7b7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function availableDepositFee(address user) view returns(uint256)
func (mousePool *MousePool) TryPackAvailableDepositFee(user common.Address) ([]byte, error) {
	return mousePool.abi.Pack("availableDepositFee", user)
}

// UnpackAvailableDepositFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1250e7b7.
//
// Solidity: function availableDepositFee(address user) view returns(uint256)
func (mousePool *MousePool) UnpackAvailableDepositFee(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("availableDepositFee", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackAvailableOverdueFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb9c62706.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function availableOverdueFee(address user) view returns(uint256)
func (mousePool *MousePool) PackAvailableOverdueFee(user common.Address) []byte {
	enc, err := mousePool.abi.Pack("availableOverdueFee", user)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAvailableOverdueFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb9c62706.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function availableOverdueFee(address user) view returns(uint256)
func (mousePool *MousePool) TryPackAvailableOverdueFee(user common.Address) ([]byte, error) {
	return mousePool.abi.Pack("availableOverdueFee", user)
}

// UnpackAvailableOverdueFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb9c62706.
//
// Solidity: function availableOverdueFee(address user) view returns(uint256)
func (mousePool *MousePool) UnpackAvailableOverdueFee(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("availableOverdueFee", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackAvailableWithdrawFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0bf23553.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function availableWithdrawFee(address user) view returns(uint256)
func (mousePool *MousePool) PackAvailableWithdrawFee(user common.Address) []byte {
	enc, err := mousePool.abi.Pack("availableWithdrawFee", user)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAvailableWithdrawFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0bf23553.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function availableWithdrawFee(address user) view returns(uint256)
func (mousePool *MousePool) TryPackAvailableWithdrawFee(user common.Address) ([]byte, error) {
	return mousePool.abi.Pack("availableWithdrawFee", user)
}

// UnpackAvailableWithdrawFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0bf23553.
//
// Solidity: function availableWithdrawFee(address user) view returns(uint256)
func (mousePool *MousePool) UnpackAvailableWithdrawFee(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("availableWithdrawFee", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x722713f7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOf() view returns(uint256)
func (mousePool *MousePool) PackBalanceOf() []byte {
	enc, err := mousePool.abi.Pack("balanceOf")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x722713f7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function balanceOf() view returns(uint256)
func (mousePool *MousePool) TryPackBalanceOf() ([]byte, error) {
	return mousePool.abi.Pack("balanceOf")
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x722713f7.
//
// Solidity: function balanceOf() view returns(uint256)
func (mousePool *MousePool) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBoostContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdfcedeee.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function boostContract() view returns(address)
func (mousePool *MousePool) PackBoostContract() []byte {
	enc, err := mousePool.abi.Pack("boostContract")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBoostContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdfcedeee.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function boostContract() view returns(address)
func (mousePool *MousePool) TryPackBoostContract() ([]byte, error) {
	return mousePool.abi.Pack("boostContract")
}

// UnpackBoostContract is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdfcedeee.
//
// Solidity: function boostContract() view returns(address)
func (mousePool *MousePool) UnpackBoostContract(data []byte) (common.Address, error) {
	out, err := mousePool.abi.Unpack("boostContract", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackBoostWeight is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6edce5ca.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function boostWeight() view returns(uint256)
func (mousePool *MousePool) PackBoostWeight() []byte {
	enc, err := mousePool.abi.Pack("boostWeight")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBoostWeight is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6edce5ca.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function boostWeight() view returns(uint256)
func (mousePool *MousePool) TryPackBoostWeight() ([]byte, error) {
	return mousePool.abi.Pack("boostWeight")
}

// UnpackBoostWeight is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6edce5ca.
//
// Solidity: function boostWeight() view returns(uint256)
func (mousePool *MousePool) UnpackBoostWeight(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("boostWeight", data)
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
// Solidity: function deposit(uint256 amount, uint256 lockDuration) returns()
func (mousePool *MousePool) PackDeposit(amount *big.Int, lockDuration *big.Int) []byte {
	enc, err := mousePool.abi.Pack("deposit", amount, lockDuration)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe2bbb158.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function deposit(uint256 amount, uint256 lockDuration) returns()
func (mousePool *MousePool) TryPackDeposit(amount *big.Int, lockDuration *big.Int) ([]byte, error) {
	return mousePool.abi.Pack("deposit", amount, lockDuration)
}

// PackFreeOverdueFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd748a9a2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function freeOverdueFee(address ) view returns(bool)
func (mousePool *MousePool) PackFreeOverdueFee(arg0 common.Address) []byte {
	enc, err := mousePool.abi.Pack("freeOverdueFee", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFreeOverdueFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd748a9a2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function freeOverdueFee(address ) view returns(bool)
func (mousePool *MousePool) TryPackFreeOverdueFee(arg0 common.Address) ([]byte, error) {
	return mousePool.abi.Pack("freeOverdueFee", arg0)
}

// UnpackFreeOverdueFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd748a9a2.
//
// Solidity: function freeOverdueFee(address ) view returns(bool)
func (mousePool *MousePool) UnpackFreeOverdueFee(data []byte) (bool, error) {
	out, err := mousePool.abi.Unpack("freeOverdueFee", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackFreePerformanceFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7340e889.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function freePerformanceFee(address ) view returns(bool)
func (mousePool *MousePool) PackFreePerformanceFee(arg0 common.Address) []byte {
	enc, err := mousePool.abi.Pack("freePerformanceFee", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFreePerformanceFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7340e889.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function freePerformanceFee(address ) view returns(bool)
func (mousePool *MousePool) TryPackFreePerformanceFee(arg0 common.Address) ([]byte, error) {
	return mousePool.abi.Pack("freePerformanceFee", arg0)
}

// UnpackFreePerformanceFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7340e889.
//
// Solidity: function freePerformanceFee(address ) view returns(bool)
func (mousePool *MousePool) UnpackFreePerformanceFee(data []byte) (bool, error) {
	out, err := mousePool.abi.Unpack("freePerformanceFee", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackFreeWithdrawFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x51c068a1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function freeWithdrawFee(address ) view returns(bool)
func (mousePool *MousePool) PackFreeWithdrawFee(arg0 common.Address) []byte {
	enc, err := mousePool.abi.Pack("freeWithdrawFee", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFreeWithdrawFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x51c068a1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function freeWithdrawFee(address ) view returns(bool)
func (mousePool *MousePool) TryPackFreeWithdrawFee(arg0 common.Address) ([]byte, error) {
	return mousePool.abi.Pack("freeWithdrawFee", arg0)
}

// UnpackFreeWithdrawFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x51c068a1.
//
// Solidity: function freeWithdrawFee(address ) view returns(bool)
func (mousePool *MousePool) UnpackFreeWithdrawFee(data []byte) (bool, error) {
	out, err := mousePool.abi.Unpack("freeWithdrawFee", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackGetPricePerFullShare is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x77c7b8fc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPricePerFullShare() view returns(uint256)
func (mousePool *MousePool) PackGetPricePerFullShare() []byte {
	enc, err := mousePool.abi.Pack("getPricePerFullShare")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPricePerFullShare is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x77c7b8fc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPricePerFullShare() view returns(uint256)
func (mousePool *MousePool) TryPackGetPricePerFullShare() ([]byte, error) {
	return mousePool.abi.Pack("getPricePerFullShare")
}

// UnpackGetPricePerFullShare is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x77c7b8fc.
//
// Solidity: function getPricePerFullShare() view returns(uint256)
func (mousePool *MousePool) UnpackGetPricePerFullShare(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("getPricePerFullShare", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMasterChef is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x575a86b2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function masterChef() view returns(address)
func (mousePool *MousePool) PackMasterChef() []byte {
	enc, err := mousePool.abi.Pack("masterChef")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMasterChef is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x575a86b2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function masterChef() view returns(address)
func (mousePool *MousePool) TryPackMasterChef() ([]byte, error) {
	return mousePool.abi.Pack("masterChef")
}

// UnpackMasterChef is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x575a86b2.
//
// Solidity: function masterChef() view returns(address)
func (mousePool *MousePool) UnpackMasterChef(data []byte) (common.Address, error) {
	out, err := mousePool.abi.Unpack("masterChef", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackMasterChefPid is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc62c610b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function masterChefPid() view returns(uint256)
func (mousePool *MousePool) PackMasterChefPid() []byte {
	enc, err := mousePool.abi.Pack("masterChefPid")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMasterChefPid is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc62c610b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function masterChefPid() view returns(uint256)
func (mousePool *MousePool) TryPackMasterChefPid() ([]byte, error) {
	return mousePool.abi.Pack("masterChefPid")
}

// UnpackMasterChefPid is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc62c610b.
//
// Solidity: function masterChefPid() view returns(uint256)
func (mousePool *MousePool) UnpackMasterChefPid(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("masterChefPid", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMaxLockDuration is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa16cdbb1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxLockDuration() view returns(uint256)
func (mousePool *MousePool) PackMaxLockDuration() []byte {
	enc, err := mousePool.abi.Pack("maxLockDuration")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxLockDuration is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa16cdbb1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxLockDuration() view returns(uint256)
func (mousePool *MousePool) TryPackMaxLockDuration() ([]byte, error) {
	return mousePool.abi.Pack("maxLockDuration")
}

// UnpackMaxLockDuration is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa16cdbb1.
//
// Solidity: function maxLockDuration() view returns(uint256)
func (mousePool *MousePool) UnpackMaxLockDuration(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("maxLockDuration", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackOperator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x570ca735.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function operator() view returns(address)
func (mousePool *MousePool) PackOperator() []byte {
	enc, err := mousePool.abi.Pack("operator")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOperator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x570ca735.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function operator() view returns(address)
func (mousePool *MousePool) TryPackOperator() ([]byte, error) {
	return mousePool.abi.Pack("operator")
}

// UnpackOperator is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (mousePool *MousePool) UnpackOperator(data []byte) (common.Address, error) {
	out, err := mousePool.abi.Unpack("operator", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackOverdueFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa5834e06.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function overdueFee() view returns(uint256)
func (mousePool *MousePool) PackOverdueFee() []byte {
	enc, err := mousePool.abi.Pack("overdueFee")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOverdueFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa5834e06.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function overdueFee() view returns(uint256)
func (mousePool *MousePool) TryPackOverdueFee() ([]byte, error) {
	return mousePool.abi.Pack("overdueFee")
}

// UnpackOverdueFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa5834e06.
//
// Solidity: function overdueFee() view returns(uint256)
func (mousePool *MousePool) UnpackOverdueFee(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("overdueFee", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (mousePool *MousePool) PackOwner() []byte {
	enc, err := mousePool.abi.Pack("owner")
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
func (mousePool *MousePool) TryPackOwner() ([]byte, error) {
	return mousePool.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (mousePool *MousePool) UnpackOwner(data []byte) (common.Address, error) {
	out, err := mousePool.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8456cb59.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pause() returns()
func (mousePool *MousePool) PackPause() []byte {
	enc, err := mousePool.abi.Pack("pause")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8456cb59.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pause() returns()
func (mousePool *MousePool) TryPackPause() ([]byte, error) {
	return mousePool.abi.Pack("pause")
}

// PackPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c975abb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function paused() view returns(bool)
func (mousePool *MousePool) PackPaused() []byte {
	enc, err := mousePool.abi.Pack("paused")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c975abb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function paused() view returns(bool)
func (mousePool *MousePool) TryPackPaused() ([]byte, error) {
	return mousePool.abi.Pack("paused")
}

// UnpackPaused is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (mousePool *MousePool) UnpackPaused(data []byte) (bool, error) {
	out, err := mousePool.abi.Unpack("paused", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackPendingOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe30c3978.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingOwner() view returns(address)
func (mousePool *MousePool) PackPendingOwner() []byte {
	enc, err := mousePool.abi.Pack("pendingOwner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe30c3978.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingOwner() view returns(address)
func (mousePool *MousePool) TryPackPendingOwner() ([]byte, error) {
	return mousePool.abi.Pack("pendingOwner")
}

// UnpackPendingOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (mousePool *MousePool) UnpackPendingOwner(data []byte) (common.Address, error) {
	out, err := mousePool.abi.Unpack("pendingOwner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPendingReward is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x137ee36e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingReward() view returns(uint256)
func (mousePool *MousePool) PackPendingReward() []byte {
	enc, err := mousePool.abi.Pack("pendingReward")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingReward is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x137ee36e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingReward() view returns(uint256)
func (mousePool *MousePool) TryPackPendingReward() ([]byte, error) {
	return mousePool.abi.Pack("pendingReward")
}

// UnpackPendingReward is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x137ee36e.
//
// Solidity: function pendingReward() view returns(uint256)
func (mousePool *MousePool) UnpackPendingReward(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("pendingReward", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPerformanceFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x87788782.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function performanceFee() view returns(uint256)
func (mousePool *MousePool) PackPerformanceFee() []byte {
	enc, err := mousePool.abi.Pack("performanceFee")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPerformanceFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x87788782.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function performanceFee() view returns(uint256)
func (mousePool *MousePool) TryPackPerformanceFee() ([]byte, error) {
	return mousePool.abi.Pack("performanceFee")
}

// UnpackPerformanceFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x87788782.
//
// Solidity: function performanceFee() view returns(uint256)
func (mousePool *MousePool) UnpackPerformanceFee(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("performanceFee", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPerformanceFeeContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3eb78874.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function performanceFeeContract() view returns(uint256)
func (mousePool *MousePool) PackPerformanceFeeContract() []byte {
	enc, err := mousePool.abi.Pack("performanceFeeContract")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPerformanceFeeContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3eb78874.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function performanceFeeContract() view returns(uint256)
func (mousePool *MousePool) TryPackPerformanceFeeContract() ([]byte, error) {
	return mousePool.abi.Pack("performanceFeeContract")
}

// UnpackPerformanceFeeContract is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3eb78874.
//
// Solidity: function performanceFeeContract() view returns(uint256)
func (mousePool *MousePool) UnpackPerformanceFeeContract(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("performanceFeeContract", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRecoverToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa7229fd9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function recoverToken(address token, address to, uint256 amount) returns()
func (mousePool *MousePool) PackRecoverToken(token common.Address, to common.Address, amount *big.Int) []byte {
	enc, err := mousePool.abi.Pack("recoverToken", token, to, amount)
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
func (mousePool *MousePool) TryPackRecoverToken(token common.Address, to common.Address, amount *big.Int) ([]byte, error) {
	return mousePool.abi.Pack("recoverToken", token, to, amount)
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() returns()
func (mousePool *MousePool) PackRenounceOwnership() []byte {
	enc, err := mousePool.abi.Pack("renounceOwnership")
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
func (mousePool *MousePool) TryPackRenounceOwnership() ([]byte, error) {
	return mousePool.abi.Pack("renounceOwnership")
}

// PackSetAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x704b6c02.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setAdmin(address newAdmin) returns()
func (mousePool *MousePool) PackSetAdmin(newAdmin common.Address) []byte {
	enc, err := mousePool.abi.Pack("setAdmin", newAdmin)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x704b6c02.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setAdmin(address newAdmin) returns()
func (mousePool *MousePool) TryPackSetAdmin(newAdmin common.Address) ([]byte, error) {
	return mousePool.abi.Pack("setAdmin", newAdmin)
}

// PackSetBoostContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdef7869d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setBoostContract(address newBoostContract) returns()
func (mousePool *MousePool) PackSetBoostContract(newBoostContract common.Address) []byte {
	enc, err := mousePool.abi.Pack("setBoostContract", newBoostContract)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetBoostContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdef7869d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setBoostContract(address newBoostContract) returns()
func (mousePool *MousePool) TryPackSetBoostContract(newBoostContract common.Address) ([]byte, error) {
	return mousePool.abi.Pack("setBoostContract", newBoostContract)
}

// PackSetBoostWeight is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x93c99e6a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setBoostWeight(uint256 weight) returns()
func (mousePool *MousePool) PackSetBoostWeight(weight *big.Int) []byte {
	enc, err := mousePool.abi.Pack("setBoostWeight", weight)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetBoostWeight is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x93c99e6a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setBoostWeight(uint256 weight) returns()
func (mousePool *MousePool) TryPackSetBoostWeight(weight *big.Int) ([]byte, error) {
	return mousePool.abi.Pack("setBoostWeight", weight)
}

// PackSetMasterChefPid is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0586755a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setMasterChefPid(uint256 pid) returns()
func (mousePool *MousePool) PackSetMasterChefPid(pid *big.Int) []byte {
	enc, err := mousePool.abi.Pack("setMasterChefPid", pid)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetMasterChefPid is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0586755a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setMasterChefPid(uint256 pid) returns()
func (mousePool *MousePool) TryPackSetMasterChefPid(pid *big.Int) ([]byte, error) {
	return mousePool.abi.Pack("setMasterChefPid", pid)
}

// PackSetMaxLockDuration is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf786b958.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setMaxLockDuration(uint256 duration) returns()
func (mousePool *MousePool) PackSetMaxLockDuration(duration *big.Int) []byte {
	enc, err := mousePool.abi.Pack("setMaxLockDuration", duration)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetMaxLockDuration is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf786b958.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setMaxLockDuration(uint256 duration) returns()
func (mousePool *MousePool) TryPackSetMaxLockDuration(duration *big.Int) ([]byte, error) {
	return mousePool.abi.Pack("setMaxLockDuration", duration)
}

// PackSetOperator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb3ab15fb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setOperator(address newOperator) returns()
func (mousePool *MousePool) PackSetOperator(newOperator common.Address) []byte {
	enc, err := mousePool.abi.Pack("setOperator", newOperator)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetOperator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb3ab15fb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setOperator(address newOperator) returns()
func (mousePool *MousePool) TryPackSetOperator(newOperator common.Address) ([]byte, error) {
	return mousePool.abi.Pack("setOperator", newOperator)
}

// PackSetOverdueFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0c59696b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setOverdueFee(uint256 fee) returns()
func (mousePool *MousePool) PackSetOverdueFee(fee *big.Int) []byte {
	enc, err := mousePool.abi.Pack("setOverdueFee", fee)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetOverdueFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0c59696b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setOverdueFee(uint256 fee) returns()
func (mousePool *MousePool) TryPackSetOverdueFee(fee *big.Int) ([]byte, error) {
	return mousePool.abi.Pack("setOverdueFee", fee)
}

// PackSetPerformanceFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70897b23.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setPerformanceFee(uint256 fee) returns()
func (mousePool *MousePool) PackSetPerformanceFee(fee *big.Int) []byte {
	enc, err := mousePool.abi.Pack("setPerformanceFee", fee)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetPerformanceFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70897b23.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setPerformanceFee(uint256 fee) returns()
func (mousePool *MousePool) TryPackSetPerformanceFee(fee *big.Int) ([]byte, error) {
	return mousePool.abi.Pack("setPerformanceFee", fee)
}

// PackSetPerformanceFeeContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbb9f408d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setPerformanceFeeContract(uint256 fee) returns()
func (mousePool *MousePool) PackSetPerformanceFeeContract(fee *big.Int) []byte {
	enc, err := mousePool.abi.Pack("setPerformanceFeeContract", fee)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetPerformanceFeeContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbb9f408d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setPerformanceFeeContract(uint256 fee) returns()
func (mousePool *MousePool) TryPackSetPerformanceFeeContract(fee *big.Int) ([]byte, error) {
	return mousePool.abi.Pack("setPerformanceFeeContract", fee)
}

// PackSetTreasury is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf0f44260.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setTreasury(address newTreasury) returns()
func (mousePool *MousePool) PackSetTreasury(newTreasury common.Address) []byte {
	enc, err := mousePool.abi.Pack("setTreasury", newTreasury)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetTreasury is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf0f44260.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setTreasury(address newTreasury) returns()
func (mousePool *MousePool) TryPackSetTreasury(newTreasury common.Address) ([]byte, error) {
	return mousePool.abi.Pack("setTreasury", newTreasury)
}

// PackSetVMouse is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2ad12435.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setVMouse(address newVMouse) returns()
func (mousePool *MousePool) PackSetVMouse(newVMouse common.Address) []byte {
	enc, err := mousePool.abi.Pack("setVMouse", newVMouse)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetVMouse is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2ad12435.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setVMouse(address newVMouse) returns()
func (mousePool *MousePool) TryPackSetVMouse(newVMouse common.Address) ([]byte, error) {
	return mousePool.abi.Pack("setVMouse", newVMouse)
}

// PackSetWhiteList is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x39340d22.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setWhiteList(address user, uint8 feeType, bool status) returns()
func (mousePool *MousePool) PackSetWhiteList(user common.Address, feeType uint8, status bool) []byte {
	enc, err := mousePool.abi.Pack("setWhiteList", user, feeType, status)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetWhiteList is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x39340d22.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setWhiteList(address user, uint8 feeType, bool status) returns()
func (mousePool *MousePool) TryPackSetWhiteList(user common.Address, feeType uint8, status bool) ([]byte, error) {
	return mousePool.abi.Pack("setWhiteList", user, feeType, status)
}

// PackSetWithdrawFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb6ac642a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setWithdrawFee(uint256 fee) returns()
func (mousePool *MousePool) PackSetWithdrawFee(fee *big.Int) []byte {
	enc, err := mousePool.abi.Pack("setWithdrawFee", fee)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetWithdrawFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb6ac642a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setWithdrawFee(uint256 fee) returns()
func (mousePool *MousePool) TryPackSetWithdrawFee(fee *big.Int) ([]byte, error) {
	return mousePool.abi.Pack("setWithdrawFee", fee)
}

// PackSetWithdrawFeeContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x14ff3039.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setWithdrawFeeContract(uint256 fee) returns()
func (mousePool *MousePool) PackSetWithdrawFeeContract(fee *big.Int) []byte {
	enc, err := mousePool.abi.Pack("setWithdrawFeeContract", fee)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetWithdrawFeeContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x14ff3039.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setWithdrawFeeContract(uint256 fee) returns()
func (mousePool *MousePool) TryPackSetWithdrawFeeContract(fee *big.Int) ([]byte, error) {
	return mousePool.abi.Pack("setWithdrawFeeContract", fee)
}

// PackSetWithdrawFeePeriod is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1efac1b8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setWithdrawFeePeriod(uint256 period) returns()
func (mousePool *MousePool) PackSetWithdrawFeePeriod(period *big.Int) []byte {
	enc, err := mousePool.abi.Pack("setWithdrawFeePeriod", period)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetWithdrawFeePeriod is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1efac1b8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setWithdrawFeePeriod(uint256 period) returns()
func (mousePool *MousePool) TryPackSetWithdrawFeePeriod(period *big.Int) ([]byte, error) {
	return mousePool.abi.Pack("setWithdrawFeePeriod", period)
}

// PackTotalBoostDebt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe73008bc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalBoostDebt() view returns(uint256)
func (mousePool *MousePool) PackTotalBoostDebt() []byte {
	enc, err := mousePool.abi.Pack("totalBoostDebt")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalBoostDebt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe73008bc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalBoostDebt() view returns(uint256)
func (mousePool *MousePool) TryPackTotalBoostDebt() ([]byte, error) {
	return mousePool.abi.Pack("totalBoostDebt")
}

// UnpackTotalBoostDebt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe73008bc.
//
// Solidity: function totalBoostDebt() view returns(uint256)
func (mousePool *MousePool) UnpackTotalBoostDebt(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("totalBoostDebt", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTotalLockedAmount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x05a9f274.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalLockedAmount() view returns(uint256)
func (mousePool *MousePool) PackTotalLockedAmount() []byte {
	enc, err := mousePool.abi.Pack("totalLockedAmount")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalLockedAmount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x05a9f274.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalLockedAmount() view returns(uint256)
func (mousePool *MousePool) TryPackTotalLockedAmount() ([]byte, error) {
	return mousePool.abi.Pack("totalLockedAmount")
}

// UnpackTotalLockedAmount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x05a9f274.
//
// Solidity: function totalLockedAmount() view returns(uint256)
func (mousePool *MousePool) UnpackTotalLockedAmount(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("totalLockedAmount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTotalShares is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3a98ef39.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalShares() view returns(uint256)
func (mousePool *MousePool) PackTotalShares() []byte {
	enc, err := mousePool.abi.Pack("totalShares")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalShares is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3a98ef39.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalShares() view returns(uint256)
func (mousePool *MousePool) TryPackTotalShares() ([]byte, error) {
	return mousePool.abi.Pack("totalShares")
}

// UnpackTotalShares is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3a98ef39.
//
// Solidity: function totalShares() view returns(uint256)
func (mousePool *MousePool) UnpackTotalShares(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("totalShares", data)
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
func (mousePool *MousePool) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := mousePool.abi.Pack("transferOwnership", newOwner)
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
func (mousePool *MousePool) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return mousePool.abi.Pack("transferOwnership", newOwner)
}

// PackTreasury is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x61d027b3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function treasury() view returns(address)
func (mousePool *MousePool) PackTreasury() []byte {
	enc, err := mousePool.abi.Pack("treasury")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTreasury is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x61d027b3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function treasury() view returns(address)
func (mousePool *MousePool) TryPackTreasury() ([]byte, error) {
	return mousePool.abi.Pack("treasury")
}

// UnpackTreasury is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (mousePool *MousePool) UnpackTreasury(data []byte) (common.Address, error) {
	out, err := mousePool.abi.Unpack("treasury", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackUnlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f6c493c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function unlock(address user) returns()
func (mousePool *MousePool) PackUnlock(user common.Address) []byte {
	enc, err := mousePool.abi.Pack("unlock", user)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUnlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f6c493c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function unlock(address user) returns()
func (mousePool *MousePool) TryPackUnlock(user common.Address) ([]byte, error) {
	return mousePool.abi.Pack("unlock", user)
}

// PackUnpause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f4ba83a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function unpause() returns()
func (mousePool *MousePool) PackUnpause() []byte {
	enc, err := mousePool.abi.Pack("unpause")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUnpause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f4ba83a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function unpause() returns()
func (mousePool *MousePool) TryPackUnpause() ([]byte, error) {
	return mousePool.abi.Pack("unpause")
}

// PackUserInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1959a002.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function userInfo(address user) view returns(uint256 shares, uint256 userBoostedShare, uint256 lockedAmount, uint256 lockEndTime, uint256 lockStartTime, bool locked, uint256 mouseAtLastUserAction, uint256 lastUserActionTime)
func (mousePool *MousePool) PackUserInfo(user common.Address) []byte {
	enc, err := mousePool.abi.Pack("userInfo", user)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUserInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1959a002.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function userInfo(address user) view returns(uint256 shares, uint256 userBoostedShare, uint256 lockedAmount, uint256 lockEndTime, uint256 lockStartTime, bool locked, uint256 mouseAtLastUserAction, uint256 lastUserActionTime)
func (mousePool *MousePool) TryPackUserInfo(user common.Address) ([]byte, error) {
	return mousePool.abi.Pack("userInfo", user)
}

// UserInfoOutput serves as a container for the return parameters of contract
// method UserInfo.
type UserInfoOutput struct {
	Shares                *big.Int
	UserBoostedShare      *big.Int
	LockedAmount          *big.Int
	LockEndTime           *big.Int
	LockStartTime         *big.Int
	Locked                bool
	MouseAtLastUserAction *big.Int
	LastUserActionTime    *big.Int
}

// UnpackUserInfo is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1959a002.
//
// Solidity: function userInfo(address user) view returns(uint256 shares, uint256 userBoostedShare, uint256 lockedAmount, uint256 lockEndTime, uint256 lockStartTime, bool locked, uint256 mouseAtLastUserAction, uint256 lastUserActionTime)
func (mousePool *MousePool) UnpackUserInfo(data []byte) (UserInfoOutput, error) {
	out, err := mousePool.abi.Unpack("userInfo", data)
	outstruct := new(UserInfoOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Shares = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.UserBoostedShare = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.LockedAmount = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.LockEndTime = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.LockStartTime = abi.ConvertType(out[4], new(big.Int)).(*big.Int)
	outstruct.Locked = *abi.ConvertType(out[5], new(bool)).(*bool)
	outstruct.MouseAtLastUserAction = abi.ConvertType(out[6], new(big.Int)).(*big.Int)
	outstruct.LastUserActionTime = abi.ConvertType(out[7], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackVMouse is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3d7bba0c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function vMouse() view returns(address)
func (mousePool *MousePool) PackVMouse() []byte {
	enc, err := mousePool.abi.Pack("vMouse")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVMouse is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3d7bba0c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function vMouse() view returns(address)
func (mousePool *MousePool) TryPackVMouse() ([]byte, error) {
	return mousePool.abi.Pack("vMouse")
}

// UnpackVMouse is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3d7bba0c.
//
// Solidity: function vMouse() view returns(address)
func (mousePool *MousePool) UnpackVMouse(data []byte) (common.Address, error) {
	out, err := mousePool.abi.Unpack("vMouse", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackWithdrawAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x853828b6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawAll() returns()
func (mousePool *MousePool) PackWithdrawAll() []byte {
	enc, err := mousePool.abi.Pack("withdrawAll")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawAll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x853828b6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawAll() returns()
func (mousePool *MousePool) TryPackWithdrawAll() ([]byte, error) {
	return mousePool.abi.Pack("withdrawAll")
}

// PackWithdrawByAmount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5521e9bf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawByAmount(uint256 amount) returns()
func (mousePool *MousePool) PackWithdrawByAmount(amount *big.Int) []byte {
	enc, err := mousePool.abi.Pack("withdrawByAmount", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawByAmount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5521e9bf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawByAmount(uint256 amount) returns()
func (mousePool *MousePool) TryPackWithdrawByAmount(amount *big.Int) ([]byte, error) {
	return mousePool.abi.Pack("withdrawByAmount", amount)
}

// PackWithdrawByShares is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8bb34bab.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawByShares(uint256 shares) returns()
func (mousePool *MousePool) PackWithdrawByShares(shares *big.Int) []byte {
	enc, err := mousePool.abi.Pack("withdrawByShares", shares)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawByShares is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8bb34bab.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawByShares(uint256 shares) returns()
func (mousePool *MousePool) TryPackWithdrawByShares(shares *big.Int) ([]byte, error) {
	return mousePool.abi.Pack("withdrawByShares", shares)
}

// PackWithdrawFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe941fa78.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawFee() view returns(uint256)
func (mousePool *MousePool) PackWithdrawFee() []byte {
	enc, err := mousePool.abi.Pack("withdrawFee")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe941fa78.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawFee() view returns(uint256)
func (mousePool *MousePool) TryPackWithdrawFee() ([]byte, error) {
	return mousePool.abi.Pack("withdrawFee")
}

// UnpackWithdrawFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe941fa78.
//
// Solidity: function withdrawFee() view returns(uint256)
func (mousePool *MousePool) UnpackWithdrawFee(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("withdrawFee", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackWithdrawFeeContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe4b37ef5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawFeeContract() view returns(uint256)
func (mousePool *MousePool) PackWithdrawFeeContract() []byte {
	enc, err := mousePool.abi.Pack("withdrawFeeContract")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawFeeContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe4b37ef5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawFeeContract() view returns(uint256)
func (mousePool *MousePool) TryPackWithdrawFeeContract() ([]byte, error) {
	return mousePool.abi.Pack("withdrawFeeContract")
}

// UnpackWithdrawFeeContract is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe4b37ef5.
//
// Solidity: function withdrawFeeContract() view returns(uint256)
func (mousePool *MousePool) UnpackWithdrawFeeContract(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("withdrawFeeContract", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackWithdrawFeePeriod is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdf10b4e6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawFeePeriod() view returns(uint256)
func (mousePool *MousePool) PackWithdrawFeePeriod() []byte {
	enc, err := mousePool.abi.Pack("withdrawFeePeriod")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawFeePeriod is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdf10b4e6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawFeePeriod() view returns(uint256)
func (mousePool *MousePool) TryPackWithdrawFeePeriod() ([]byte, error) {
	return mousePool.abi.Pack("withdrawFeePeriod")
}

// UnpackWithdrawFeePeriod is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdf10b4e6.
//
// Solidity: function withdrawFeePeriod() view returns(uint256)
func (mousePool *MousePool) UnpackWithdrawFeePeriod(data []byte) (*big.Int, error) {
	out, err := mousePool.abi.Unpack("withdrawFeePeriod", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// MousePoolAdminUpdated represents a AdminUpdated event raised by the MousePool contract.
type MousePoolAdminUpdated struct {
	OldAdmin common.Address
	NewAdmin common.Address
	Raw      *types.Log // Blockchain specific contextual infos
}

const MousePoolAdminUpdatedEventName = "AdminUpdated"

// ContractEventName returns the user-defined event name.
func (MousePoolAdminUpdated) ContractEventName() string {
	return MousePoolAdminUpdatedEventName
}

// UnpackAdminUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AdminUpdated(address indexed oldAdmin, address indexed newAdmin)
func (mousePool *MousePool) UnpackAdminUpdatedEvent(log *types.Log) (*MousePoolAdminUpdated, error) {
	event := "AdminUpdated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolAdminUpdated)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolBoostContractUpdated represents a BoostContractUpdated event raised by the MousePool contract.
type MousePoolBoostContractUpdated struct {
	OldBoostContract common.Address
	NewBoostContract common.Address
	Raw              *types.Log // Blockchain specific contextual infos
}

const MousePoolBoostContractUpdatedEventName = "BoostContractUpdated"

// ContractEventName returns the user-defined event name.
func (MousePoolBoostContractUpdated) ContractEventName() string {
	return MousePoolBoostContractUpdatedEventName
}

// UnpackBoostContractUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BoostContractUpdated(address indexed oldBoostContract, address indexed newBoostContract)
func (mousePool *MousePool) UnpackBoostContractUpdatedEvent(log *types.Log) (*MousePoolBoostContractUpdated, error) {
	event := "BoostContractUpdated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolBoostContractUpdated)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolDeposit represents a Deposit event raised by the MousePool contract.
type MousePoolDeposit struct {
	User         common.Address
	Amount       *big.Int
	LockDuration *big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const MousePoolDepositEventName = "Deposit"

// ContractEventName returns the user-defined event name.
func (MousePoolDeposit) ContractEventName() string {
	return MousePoolDepositEventName
}

// UnpackDepositEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Deposit(address indexed user, uint256 amount, uint256 lockDuration)
func (mousePool *MousePool) UnpackDepositEvent(log *types.Log) (*MousePoolDeposit, error) {
	event := "Deposit"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolDeposit)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolLock represents a Lock event raised by the MousePool contract.
type MousePoolLock struct {
	User         common.Address
	LockedAmount *big.Int
	LockDuration *big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const MousePoolLockEventName = "Lock"

// ContractEventName returns the user-defined event name.
func (MousePoolLock) ContractEventName() string {
	return MousePoolLockEventName
}

// UnpackLockEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Lock(address indexed user, uint256 lockedAmount, uint256 lockDuration)
func (mousePool *MousePool) UnpackLockEvent(log *types.Log) (*MousePoolLock, error) {
	event := "Lock"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolLock)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolMasterChefPidUpdated represents a MasterChefPidUpdated event raised by the MousePool contract.
type MousePoolMasterChefPidUpdated struct {
	OldPid *big.Int
	NewPid *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const MousePoolMasterChefPidUpdatedEventName = "MasterChefPidUpdated"

// ContractEventName returns the user-defined event name.
func (MousePoolMasterChefPidUpdated) ContractEventName() string {
	return MousePoolMasterChefPidUpdatedEventName
}

// UnpackMasterChefPidUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MasterChefPidUpdated(uint256 oldPid, uint256 newPid)
func (mousePool *MousePool) UnpackMasterChefPidUpdatedEvent(log *types.Log) (*MousePoolMasterChefPidUpdated, error) {
	event := "MasterChefPidUpdated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolMasterChefPidUpdated)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolOperatorUpdated represents a OperatorUpdated event raised by the MousePool contract.
type MousePoolOperatorUpdated struct {
	OldOperator common.Address
	NewOperator common.Address
	Raw         *types.Log // Blockchain specific contextual infos
}

const MousePoolOperatorUpdatedEventName = "OperatorUpdated"

// ContractEventName returns the user-defined event name.
func (MousePoolOperatorUpdated) ContractEventName() string {
	return MousePoolOperatorUpdatedEventName
}

// UnpackOperatorUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OperatorUpdated(address indexed oldOperator, address indexed newOperator)
func (mousePool *MousePool) UnpackOperatorUpdatedEvent(log *types.Log) (*MousePoolOperatorUpdated, error) {
	event := "OperatorUpdated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolOperatorUpdated)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolOverdueFeeUpdated represents a OverdueFeeUpdated event raised by the MousePool contract.
type MousePoolOverdueFeeUpdated struct {
	OldFee *big.Int
	NewFee *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const MousePoolOverdueFeeUpdatedEventName = "OverdueFeeUpdated"

// ContractEventName returns the user-defined event name.
func (MousePoolOverdueFeeUpdated) ContractEventName() string {
	return MousePoolOverdueFeeUpdatedEventName
}

// UnpackOverdueFeeUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OverdueFeeUpdated(uint256 oldFee, uint256 newFee)
func (mousePool *MousePool) UnpackOverdueFeeUpdatedEvent(log *types.Log) (*MousePoolOverdueFeeUpdated, error) {
	event := "OverdueFeeUpdated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolOverdueFeeUpdated)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the MousePool contract.
type MousePoolOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const MousePoolOwnershipTransferStartedEventName = "OwnershipTransferStarted"

// ContractEventName returns the user-defined event name.
func (MousePoolOwnershipTransferStarted) ContractEventName() string {
	return MousePoolOwnershipTransferStartedEventName
}

// UnpackOwnershipTransferStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (mousePool *MousePool) UnpackOwnershipTransferStartedEvent(log *types.Log) (*MousePoolOwnershipTransferStarted, error) {
	event := "OwnershipTransferStarted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolOwnershipTransferStarted)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolOwnershipTransferred represents a OwnershipTransferred event raised by the MousePool contract.
type MousePoolOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const MousePoolOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (MousePoolOwnershipTransferred) ContractEventName() string {
	return MousePoolOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (mousePool *MousePool) UnpackOwnershipTransferredEvent(log *types.Log) (*MousePoolOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolPaused represents a Paused event raised by the MousePool contract.
type MousePoolPaused struct {
	Account common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const MousePoolPausedEventName = "Paused"

// ContractEventName returns the user-defined event name.
func (MousePoolPaused) ContractEventName() string {
	return MousePoolPausedEventName
}

// UnpackPausedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Paused(address account)
func (mousePool *MousePool) UnpackPausedEvent(log *types.Log) (*MousePoolPaused, error) {
	event := "Paused"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolPaused)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolPerformanceFeeUpdated represents a PerformanceFeeUpdated event raised by the MousePool contract.
type MousePoolPerformanceFeeUpdated struct {
	OldFee *big.Int
	NewFee *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const MousePoolPerformanceFeeUpdatedEventName = "PerformanceFeeUpdated"

// ContractEventName returns the user-defined event name.
func (MousePoolPerformanceFeeUpdated) ContractEventName() string {
	return MousePoolPerformanceFeeUpdatedEventName
}

// UnpackPerformanceFeeUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PerformanceFeeUpdated(uint256 oldFee, uint256 newFee)
func (mousePool *MousePool) UnpackPerformanceFeeUpdatedEvent(log *types.Log) (*MousePoolPerformanceFeeUpdated, error) {
	event := "PerformanceFeeUpdated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolPerformanceFeeUpdated)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolTokenRecovered represents a TokenRecovered event raised by the MousePool contract.
type MousePoolTokenRecovered struct {
	Token  common.Address
	To     common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const MousePoolTokenRecoveredEventName = "TokenRecovered"

// ContractEventName returns the user-defined event name.
func (MousePoolTokenRecovered) ContractEventName() string {
	return MousePoolTokenRecoveredEventName
}

// UnpackTokenRecoveredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenRecovered(address indexed token, address indexed to, uint256 amount)
func (mousePool *MousePool) UnpackTokenRecoveredEvent(log *types.Log) (*MousePoolTokenRecovered, error) {
	event := "TokenRecovered"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolTokenRecovered)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolTreasuryUpdated represents a TreasuryUpdated event raised by the MousePool contract.
type MousePoolTreasuryUpdated struct {
	OldTreasury common.Address
	NewTreasury common.Address
	Raw         *types.Log // Blockchain specific contextual infos
}

const MousePoolTreasuryUpdatedEventName = "TreasuryUpdated"

// ContractEventName returns the user-defined event name.
func (MousePoolTreasuryUpdated) ContractEventName() string {
	return MousePoolTreasuryUpdatedEventName
}

// UnpackTreasuryUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TreasuryUpdated(address indexed oldTreasury, address indexed newTreasury)
func (mousePool *MousePool) UnpackTreasuryUpdatedEvent(log *types.Log) (*MousePoolTreasuryUpdated, error) {
	event := "TreasuryUpdated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolTreasuryUpdated)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolUnlock represents a Unlock event raised by the MousePool contract.
type MousePoolUnlock struct {
	User               common.Address
	LockedAmount       *big.Int
	BoostSharesCleared *big.Int
	Raw                *types.Log // Blockchain specific contextual infos
}

const MousePoolUnlockEventName = "Unlock"

// ContractEventName returns the user-defined event name.
func (MousePoolUnlock) ContractEventName() string {
	return MousePoolUnlockEventName
}

// UnpackUnlockEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Unlock(address indexed user, uint256 lockedAmount, uint256 boostSharesCleared)
func (mousePool *MousePool) UnpackUnlockEvent(log *types.Log) (*MousePoolUnlock, error) {
	event := "Unlock"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolUnlock)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolUnpaused represents a Unpaused event raised by the MousePool contract.
type MousePoolUnpaused struct {
	Account common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const MousePoolUnpausedEventName = "Unpaused"

// ContractEventName returns the user-defined event name.
func (MousePoolUnpaused) ContractEventName() string {
	return MousePoolUnpausedEventName
}

// UnpackUnpausedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Unpaused(address account)
func (mousePool *MousePool) UnpackUnpausedEvent(log *types.Log) (*MousePoolUnpaused, error) {
	event := "Unpaused"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolUnpaused)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolWhiteListUpdated represents a WhiteListUpdated event raised by the MousePool contract.
type MousePoolWhiteListUpdated struct {
	User    common.Address
	FeeType uint8
	Status  bool
	Raw     *types.Log // Blockchain specific contextual infos
}

const MousePoolWhiteListUpdatedEventName = "WhiteListUpdated"

// ContractEventName returns the user-defined event name.
func (MousePoolWhiteListUpdated) ContractEventName() string {
	return MousePoolWhiteListUpdatedEventName
}

// UnpackWhiteListUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event WhiteListUpdated(address indexed user, uint8 feeType, bool status)
func (mousePool *MousePool) UnpackWhiteListUpdatedEvent(log *types.Log) (*MousePoolWhiteListUpdated, error) {
	event := "WhiteListUpdated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolWhiteListUpdated)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolWithdraw represents a Withdraw event raised by the MousePool contract.
type MousePoolWithdraw struct {
	User   common.Address
	Amount *big.Int
	Shares *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const MousePoolWithdrawEventName = "Withdraw"

// ContractEventName returns the user-defined event name.
func (MousePoolWithdraw) ContractEventName() string {
	return MousePoolWithdrawEventName
}

// UnpackWithdrawEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdraw(address indexed user, uint256 amount, uint256 shares)
func (mousePool *MousePool) UnpackWithdrawEvent(log *types.Log) (*MousePoolWithdraw, error) {
	event := "Withdraw"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolWithdraw)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolWithdrawAll represents a WithdrawAll event raised by the MousePool contract.
type MousePoolWithdrawAll struct {
	User   common.Address
	Amount *big.Int
	Shares *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const MousePoolWithdrawAllEventName = "WithdrawAll"

// ContractEventName returns the user-defined event name.
func (MousePoolWithdrawAll) ContractEventName() string {
	return MousePoolWithdrawAllEventName
}

// UnpackWithdrawAllEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event WithdrawAll(address indexed user, uint256 amount, uint256 shares)
func (mousePool *MousePool) UnpackWithdrawAllEvent(log *types.Log) (*MousePoolWithdrawAll, error) {
	event := "WithdrawAll"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolWithdrawAll)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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

// MousePoolWithdrawFeeUpdated represents a WithdrawFeeUpdated event raised by the MousePool contract.
type MousePoolWithdrawFeeUpdated struct {
	OldFee *big.Int
	NewFee *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const MousePoolWithdrawFeeUpdatedEventName = "WithdrawFeeUpdated"

// ContractEventName returns the user-defined event name.
func (MousePoolWithdrawFeeUpdated) ContractEventName() string {
	return MousePoolWithdrawFeeUpdatedEventName
}

// UnpackWithdrawFeeUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event WithdrawFeeUpdated(uint256 oldFee, uint256 newFee)
func (mousePool *MousePool) UnpackWithdrawFeeUpdatedEvent(log *types.Log) (*MousePoolWithdrawFeeUpdated, error) {
	event := "WithdrawFeeUpdated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mousePool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MousePoolWithdrawFeeUpdated)
	if len(log.Data) > 0 {
		if err := mousePool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mousePool.abi.Events[event].Inputs {
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
func (mousePool *MousePool) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], mousePool.abi.Errors["CannotRecoverMouse"].ID.Bytes()[:4]) {
		return mousePool.UnpackCannotRecoverMouseError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["DepositTooSmall"].ID.Bytes()[:4]) {
		return mousePool.UnpackDepositTooSmallError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["EnforcedPause"].ID.Bytes()[:4]) {
		return mousePool.UnpackEnforcedPauseError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["ExpectedPause"].ID.Bytes()[:4]) {
		return mousePool.UnpackExpectedPauseError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["FeeExceedsMax"].ID.Bytes()[:4]) {
		return mousePool.UnpackFeeExceedsMaxError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["InsufficientBalance"].ID.Bytes()[:4]) {
		return mousePool.UnpackInsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["InsufficientShares"].ID.Bytes()[:4]) {
		return mousePool.UnpackInsufficientSharesError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["LockDurationTooLong"].ID.Bytes()[:4]) {
		return mousePool.UnpackLockDurationTooLongError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["LockDurationTooShort"].ID.Bytes()[:4]) {
		return mousePool.UnpackLockDurationTooShortError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["LockNotExpired"].ID.Bytes()[:4]) {
		return mousePool.UnpackLockNotExpiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["NotAdmin"].ID.Bytes()[:4]) {
		return mousePool.UnpackNotAdminError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["NotOwnerOrAdmin"].ID.Bytes()[:4]) {
		return mousePool.UnpackNotOwnerOrAdminError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["NotOwnerOrOperator"].ID.Bytes()[:4]) {
		return mousePool.UnpackNotOwnerOrOperatorError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return mousePool.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return mousePool.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["ParameterExceedsMax"].ID.Bytes()[:4]) {
		return mousePool.UnpackParameterExceedsMaxError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]) {
		return mousePool.UnpackSafeERC20FailedOperationError(raw[4:])
	}
	if bytes.Equal(raw[:4], mousePool.abi.Errors["WithdrawWhileLocked"].ID.Bytes()[:4]) {
		return mousePool.UnpackWithdrawWhileLockedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// MousePoolCannotRecoverMouse represents a CannotRecoverMouse error raised by the MousePool contract.
type MousePoolCannotRecoverMouse struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CannotRecoverMouse()
func MousePoolCannotRecoverMouseErrorID() common.Hash {
	return common.HexToHash("0xa0145ecd2567bf7bb1b2261870c02cc65948f58bcd296ae6dda8697db0fe41b7")
}

// UnpackCannotRecoverMouseError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CannotRecoverMouse()
func (mousePool *MousePool) UnpackCannotRecoverMouseError(raw []byte) (*MousePoolCannotRecoverMouse, error) {
	out := new(MousePoolCannotRecoverMouse)
	if err := mousePool.abi.UnpackIntoInterface(out, "CannotRecoverMouse", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolDepositTooSmall represents a DepositTooSmall error raised by the MousePool contract.
type MousePoolDepositTooSmall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DepositTooSmall()
func MousePoolDepositTooSmallErrorID() common.Hash {
	return common.HexToHash("0x6ba4a1c70440d46aadec62c30af648c585540392d13a2ec5280fdd24098c0e89")
}

// UnpackDepositTooSmallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DepositTooSmall()
func (mousePool *MousePool) UnpackDepositTooSmallError(raw []byte) (*MousePoolDepositTooSmall, error) {
	out := new(MousePoolDepositTooSmall)
	if err := mousePool.abi.UnpackIntoInterface(out, "DepositTooSmall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolEnforcedPause represents a EnforcedPause error raised by the MousePool contract.
type MousePoolEnforcedPause struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EnforcedPause()
func MousePoolEnforcedPauseErrorID() common.Hash {
	return common.HexToHash("0xd93c0665d6c96d04a8f174024fc4ddd66c250604aff22bbec808de86dd3637e3")
}

// UnpackEnforcedPauseError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EnforcedPause()
func (mousePool *MousePool) UnpackEnforcedPauseError(raw []byte) (*MousePoolEnforcedPause, error) {
	out := new(MousePoolEnforcedPause)
	if err := mousePool.abi.UnpackIntoInterface(out, "EnforcedPause", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolExpectedPause represents a ExpectedPause error raised by the MousePool contract.
type MousePoolExpectedPause struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ExpectedPause()
func MousePoolExpectedPauseErrorID() common.Hash {
	return common.HexToHash("0x8dfc202bcfe9a735b559bee70674422512bc5c30f687046ae8778315fb81da44")
}

// UnpackExpectedPauseError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ExpectedPause()
func (mousePool *MousePool) UnpackExpectedPauseError(raw []byte) (*MousePoolExpectedPause, error) {
	out := new(MousePoolExpectedPause)
	if err := mousePool.abi.UnpackIntoInterface(out, "ExpectedPause", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolFeeExceedsMax represents a FeeExceedsMax error raised by the MousePool contract.
type MousePoolFeeExceedsMax struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FeeExceedsMax()
func MousePoolFeeExceedsMaxErrorID() common.Hash {
	return common.HexToHash("0x5ff85e3f2cf9745712c5f72f0a9fdbb4d9a38d4c5ee5c35168e27cec99ac6ed0")
}

// UnpackFeeExceedsMaxError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FeeExceedsMax()
func (mousePool *MousePool) UnpackFeeExceedsMaxError(raw []byte) (*MousePoolFeeExceedsMax, error) {
	out := new(MousePoolFeeExceedsMax)
	if err := mousePool.abi.UnpackIntoInterface(out, "FeeExceedsMax", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolInsufficientBalance represents a InsufficientBalance error raised by the MousePool contract.
type MousePoolInsufficientBalance struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientBalance()
func MousePoolInsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xf4d678b8ce6b5157126b1484a53523762a93571537a7d5ae97d8014a44715c94")
}

// UnpackInsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientBalance()
func (mousePool *MousePool) UnpackInsufficientBalanceError(raw []byte) (*MousePoolInsufficientBalance, error) {
	out := new(MousePoolInsufficientBalance)
	if err := mousePool.abi.UnpackIntoInterface(out, "InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolInsufficientShares represents a InsufficientShares error raised by the MousePool contract.
type MousePoolInsufficientShares struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientShares()
func MousePoolInsufficientSharesErrorID() common.Hash {
	return common.HexToHash("0x399965675cfec4301cbe5ec24fb407575c5a7e4f40d219532068c8e5b35040f9")
}

// UnpackInsufficientSharesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientShares()
func (mousePool *MousePool) UnpackInsufficientSharesError(raw []byte) (*MousePoolInsufficientShares, error) {
	out := new(MousePoolInsufficientShares)
	if err := mousePool.abi.UnpackIntoInterface(out, "InsufficientShares", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolLockDurationTooLong represents a LockDurationTooLong error raised by the MousePool contract.
type MousePoolLockDurationTooLong struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LockDurationTooLong()
func MousePoolLockDurationTooLongErrorID() common.Hash {
	return common.HexToHash("0xf761f1cd9a3753f43eb6d5d92fd1d5743cdf6baf62b6a831c6e3367886c86bf6")
}

// UnpackLockDurationTooLongError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LockDurationTooLong()
func (mousePool *MousePool) UnpackLockDurationTooLongError(raw []byte) (*MousePoolLockDurationTooLong, error) {
	out := new(MousePoolLockDurationTooLong)
	if err := mousePool.abi.UnpackIntoInterface(out, "LockDurationTooLong", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolLockDurationTooShort represents a LockDurationTooShort error raised by the MousePool contract.
type MousePoolLockDurationTooShort struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LockDurationTooShort()
func MousePoolLockDurationTooShortErrorID() common.Hash {
	return common.HexToHash("0x49eeb0b3c8f64bf75daa046f145d76a285ffa603de1619741c4e879b5c88b8c4")
}

// UnpackLockDurationTooShortError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LockDurationTooShort()
func (mousePool *MousePool) UnpackLockDurationTooShortError(raw []byte) (*MousePoolLockDurationTooShort, error) {
	out := new(MousePoolLockDurationTooShort)
	if err := mousePool.abi.UnpackIntoInterface(out, "LockDurationTooShort", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolLockNotExpired represents a LockNotExpired error raised by the MousePool contract.
type MousePoolLockNotExpired struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LockNotExpired()
func MousePoolLockNotExpiredErrorID() common.Hash {
	return common.HexToHash("0x6855a8023be3603052b1aa9cd2b668dd018abaa9f4c14ca7e48fd5cf530155c2")
}

// UnpackLockNotExpiredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LockNotExpired()
func (mousePool *MousePool) UnpackLockNotExpiredError(raw []byte) (*MousePoolLockNotExpired, error) {
	out := new(MousePoolLockNotExpired)
	if err := mousePool.abi.UnpackIntoInterface(out, "LockNotExpired", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolNotAdmin represents a NotAdmin error raised by the MousePool contract.
type MousePoolNotAdmin struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotAdmin()
func MousePoolNotAdminErrorID() common.Hash {
	return common.HexToHash("0x7bfa4b9fb0cd3687c1f539f384b3f3f258f2c9aa9186353d0815413b508ed97d")
}

// UnpackNotAdminError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotAdmin()
func (mousePool *MousePool) UnpackNotAdminError(raw []byte) (*MousePoolNotAdmin, error) {
	out := new(MousePoolNotAdmin)
	if err := mousePool.abi.UnpackIntoInterface(out, "NotAdmin", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolNotOwnerOrAdmin represents a NotOwnerOrAdmin error raised by the MousePool contract.
type MousePoolNotOwnerOrAdmin struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotOwnerOrAdmin()
func MousePoolNotOwnerOrAdminErrorID() common.Hash {
	return common.HexToHash("0xdce381251f41c62fb17cbdacf34723d9217f4babc58b7f69521cfb5ae7975ccc")
}

// UnpackNotOwnerOrAdminError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotOwnerOrAdmin()
func (mousePool *MousePool) UnpackNotOwnerOrAdminError(raw []byte) (*MousePoolNotOwnerOrAdmin, error) {
	out := new(MousePoolNotOwnerOrAdmin)
	if err := mousePool.abi.UnpackIntoInterface(out, "NotOwnerOrAdmin", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolNotOwnerOrOperator represents a NotOwnerOrOperator error raised by the MousePool contract.
type MousePoolNotOwnerOrOperator struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotOwnerOrOperator()
func MousePoolNotOwnerOrOperatorErrorID() common.Hash {
	return common.HexToHash("0x98f76d0e59b6c832118deb8fe48bbe79412668756e8d10653b38bfbab51d7a89")
}

// UnpackNotOwnerOrOperatorError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotOwnerOrOperator()
func (mousePool *MousePool) UnpackNotOwnerOrOperatorError(raw []byte) (*MousePoolNotOwnerOrOperator, error) {
	out := new(MousePoolNotOwnerOrOperator)
	if err := mousePool.abi.UnpackIntoInterface(out, "NotOwnerOrOperator", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the MousePool contract.
type MousePoolOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func MousePoolOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (mousePool *MousePool) UnpackOwnableInvalidOwnerError(raw []byte) (*MousePoolOwnableInvalidOwner, error) {
	out := new(MousePoolOwnableInvalidOwner)
	if err := mousePool.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the MousePool contract.
type MousePoolOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func MousePoolOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (mousePool *MousePool) UnpackOwnableUnauthorizedAccountError(raw []byte) (*MousePoolOwnableUnauthorizedAccount, error) {
	out := new(MousePoolOwnableUnauthorizedAccount)
	if err := mousePool.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolParameterExceedsMax represents a ParameterExceedsMax error raised by the MousePool contract.
type MousePoolParameterExceedsMax struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ParameterExceedsMax()
func MousePoolParameterExceedsMaxErrorID() common.Hash {
	return common.HexToHash("0xdda1d95418286d8e457fe5fa24fdcd40c4eed9955c9919a1383e8952fc90a0a9")
}

// UnpackParameterExceedsMaxError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ParameterExceedsMax()
func (mousePool *MousePool) UnpackParameterExceedsMaxError(raw []byte) (*MousePoolParameterExceedsMax, error) {
	out := new(MousePoolParameterExceedsMax)
	if err := mousePool.abi.UnpackIntoInterface(out, "ParameterExceedsMax", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolSafeERC20FailedOperation represents a SafeERC20FailedOperation error raised by the MousePool contract.
type MousePoolSafeERC20FailedOperation struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeERC20FailedOperation(address token)
func MousePoolSafeERC20FailedOperationErrorID() common.Hash {
	return common.HexToHash("0x5274afe73c98b4749fc91ffae6b7b574e7842cb2144a159e9377a5f20b32edf9")
}

// UnpackSafeERC20FailedOperationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeERC20FailedOperation(address token)
func (mousePool *MousePool) UnpackSafeERC20FailedOperationError(raw []byte) (*MousePoolSafeERC20FailedOperation, error) {
	out := new(MousePoolSafeERC20FailedOperation)
	if err := mousePool.abi.UnpackIntoInterface(out, "SafeERC20FailedOperation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MousePoolWithdrawWhileLocked represents a WithdrawWhileLocked error raised by the MousePool contract.
type MousePoolWithdrawWhileLocked struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WithdrawWhileLocked()
func MousePoolWithdrawWhileLockedErrorID() common.Hash {
	return common.HexToHash("0xc653938f9fa6be384e1cea3b98b88741b405381f418ad2579f8acbfd3b0ae09a")
}

// UnpackWithdrawWhileLockedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WithdrawWhileLocked()
func (mousePool *MousePool) UnpackWithdrawWhileLockedError(raw []byte) (*MousePoolWithdrawWhileLocked, error) {
	out := new(MousePoolWithdrawWhileLocked)
	if err := mousePool.abi.UnpackIntoInterface(out, "WithdrawWhileLocked", raw); err != nil {
		return nil, err
	}
	return out, nil
}

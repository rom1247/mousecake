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

// MouseTierMetaData contains all meta data concerning the MouseTier contract.
var MouseTierMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_mousePool\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_ceiling\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_multiplier\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_tierBaseAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MIN_CEILING\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"X_FACTOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ceiling\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserCredit\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserTier\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mousePool\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMousePool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"multiplier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tierBaseAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateCeiling\",\"inputs\":[{\"name\":\"newCeiling\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMultiplier\",\"inputs\":[{\"name\":\"newMultiplier\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateTierBaseAmount\",\"inputs\":[{\"name\":\"newTierBaseAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpdateCeiling\",\"inputs\":[{\"name\":\"oldCeiling\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newCeiling\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpdateMultiplier\",\"inputs\":[{\"name\":\"oldMultiplier\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newMultiplier\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpdateTierBaseAmount\",\"inputs\":[{\"name\":\"oldTierBaseAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newTierBaseAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CeilingTooLow\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SameValue\",\"inputs\":[]}]",
	ID:  "MouseTier",
	Bin: "0x60a060405234801561000f575f5ffd5b50604051610fc2380380610fc283398181016040528101906100319190610256565b335f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100a2575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161009991906102c9565b60405180910390fd5b6100b18161010460201b60201c565b508373ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff1681525050826001819055508160028190555080600381905550505050506102e2565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6101f2826101c9565b9050919050565b610202816101e8565b811461020c575f5ffd5b50565b5f8151905061021d816101f9565b92915050565b5f819050919050565b61023581610223565b811461023f575f5ffd5b50565b5f815190506102508161022c565b92915050565b5f5f5f5f6080858703121561026e5761026d6101c5565b5b5f61027b8782880161020f565b945050602061028c87828801610242565b935050604061029d87828801610242565b92505060606102ae87828801610242565b91505092959194509250565b6102c3816101e8565b82525050565b5f6020820190506102dc5f8301846102ba565b92915050565b608051610cc16103015f395f818161034d01526104f40152610cc15ff3fe608060405234801561000f575f5ffd5b50600436106100e8575f3560e01c8063b1463a6a1161008a578063de05fd1e11610064578063de05fd1e14610210578063e3e5edf21461022e578063e4d2620e1461024c578063f2fde38b1461027c576100e8565b8063b1463a6a146101a8578063d15f9a3f146101c4578063dad6f0dc146101e0576100e8565b8063715018a6116100c6578063715018a614610144578063753ed1bd1461014e5780638b4d97a31461016c5780638da5cb5b1461018a576100e8565b80631b3ed722146100ec5780633a513af11461010a5780635ffe614614610128575b5f5ffd5b6100f4610298565b60405161010191906108d9565b60405180910390f35b61011261029e565b60405161011f91906108d9565b60405180910390f35b610142600480360381019061013d9190610920565b6102a5565b005b61014c610332565b005b610156610345565b60405161016391906108d9565b60405180910390f35b61017461034b565b60405161018191906109c5565b60405180910390f35b61019261036f565b60405161019f91906109fe565b60405180910390f35b6101c260048036038101906101bd9190610920565b610396565b005b6101de60048036038101906101d99190610920565b610460565b005b6101fa60048036038101906101f59190610a41565b6104ed565b60405161020791906108d9565b60405180910390f35b610218610625565b60405161022591906108d9565b60405180910390f35b61023661062b565b60405161024391906108d9565b60405180910390f35b61026660048036038101906102619190610a41565b610637565b60405161027391906108d9565b60405180910390f35b61029660048036038101906102919190610a41565b6106ee565b005b60025481565b622e248081565b6102ad610772565b60025481036102e8576040517fc23f6ccb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6002549050816002819055507f81a6f26d76b9353671a761143de53391280fa6ee9159d0710645e672670a57ef8183604051610326929190610a6c565b60405180910390a15050565b61033a610772565b6103435f6107f9565b565b60015481565b7f000000000000000000000000000000000000000000000000000000000000000081565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b61039e610772565b622e24808110156103db576040517f1c916a0d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001548103610416576040517fc23f6ccb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6001549050816001819055507f0236fce09d7de285f15fcf552a4f5e9356728e5193efb5168278417a84432c348183604051610454929190610a6c565b60405180910390a15050565b610468610772565b60035481036104a3576040517fc23f6ccb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6003549050816003819055507f14c041708be03e2ae64204404f585e159fb5b6861e88e275ccedb4e6096a606781836040516104e1929190610a6c565b60405180910390a15050565b5f5f5f5f5f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16631959a002876040518263ffffffff1660e01b815260040161054b91906109fe565b61010060405180830381865afa158015610567573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061058b9190610adc565b5050955095509550955050508015806105a357508242115b156105b4575f945050505050610620565b5f82846105c19190610bba565b90505f670de0b6b3a7640000600254876105db9190610bed565b6105e59190610c5b565b905060015482106105fe57809650505050505050610620565b600154828261060d9190610bed565b6106179190610c5b565b96505050505050505b919050565b60035481565b670de0b6b3a764000081565b5f5f610642836104ed565b905060035460016106539190610bed565b811015610663575f9150506106e9565b60035460036106729190610bed565b8110156106835760019150506106e9565b600354600a6106929190610bed565b8110156106a35760029150506106e9565b60035460146106b29190610bed565b8110156106c35760039150506106e9565b60035460236106d29190610bed565b8110156106e35760049150506106e9565b60059150505b919050565b6106f6610772565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610766575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161075d91906109fe565b60405180910390fd5b61076f816107f9565b50565b61077a6108ba565b73ffffffffffffffffffffffffffffffffffffffff1661079861036f565b73ffffffffffffffffffffffffffffffffffffffff16146107f7576107bb6108ba565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016107ee91906109fe565b60405180910390fd5b565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f33905090565b5f819050919050565b6108d3816108c1565b82525050565b5f6020820190506108ec5f8301846108ca565b92915050565b5f5ffd5b6108ff816108c1565b8114610909575f5ffd5b50565b5f8135905061091a816108f6565b92915050565b5f60208284031215610935576109346108f2565b5b5f6109428482850161090c565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f819050919050565b5f61098d6109886109838461094b565b61096a565b61094b565b9050919050565b5f61099e82610973565b9050919050565b5f6109af82610994565b9050919050565b6109bf816109a5565b82525050565b5f6020820190506109d85f8301846109b6565b92915050565b5f6109e88261094b565b9050919050565b6109f8816109de565b82525050565b5f602082019050610a115f8301846109ef565b92915050565b610a20816109de565b8114610a2a575f5ffd5b50565b5f81359050610a3b81610a17565b92915050565b5f60208284031215610a5657610a556108f2565b5b5f610a6384828501610a2d565b91505092915050565b5f604082019050610a7f5f8301856108ca565b610a8c60208301846108ca565b9392505050565b5f81519050610aa1816108f6565b92915050565b5f8115159050919050565b610abb81610aa7565b8114610ac5575f5ffd5b50565b5f81519050610ad681610ab2565b92915050565b5f5f5f5f5f5f5f5f610100898b031215610af957610af86108f2565b5b5f610b068b828c01610a93565b9850506020610b178b828c01610a93565b9750506040610b288b828c01610a93565b9650506060610b398b828c01610a93565b9550506080610b4a8b828c01610a93565b94505060a0610b5b8b828c01610ac8565b93505060c0610b6c8b828c01610a93565b92505060e0610b7d8b828c01610a93565b9150509295985092959890939650565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f610bc4826108c1565b9150610bcf836108c1565b9250828203905081811115610be757610be6610b8d565b5b92915050565b5f610bf7826108c1565b9150610c02836108c1565b9250828202610c10816108c1565b91508282048414831517610c2757610c26610b8d565b5b5092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f610c65826108c1565b9150610c70836108c1565b925082610c8057610c7f610c2e565b5b82820490509291505056fea2646970667358221220ed19c546f203edba15bfebaac5611d15d27bb1120602b7590703f62b8f7b003d64736f6c634300081c0033",
}

// MouseTier is an auto generated Go binding around an Ethereum contract.
type MouseTier struct {
	abi abi.ABI
}

// NewMouseTier creates a new instance of MouseTier.
func NewMouseTier() *MouseTier {
	parsed, err := MouseTierMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &MouseTier{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *MouseTier) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address _mousePool, uint256 _ceiling, uint256 _multiplier, uint256 _tierBaseAmount) returns()
func (mouseTier *MouseTier) PackConstructor(_mousePool common.Address, _ceiling *big.Int, _multiplier *big.Int, _tierBaseAmount *big.Int) []byte {
	enc, err := mouseTier.abi.Pack("", _mousePool, _ceiling, _multiplier, _tierBaseAmount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackMINCEILING is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3a513af1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MIN_CEILING() view returns(uint256)
func (mouseTier *MouseTier) PackMINCEILING() []byte {
	enc, err := mouseTier.abi.Pack("MIN_CEILING")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMINCEILING is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3a513af1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MIN_CEILING() view returns(uint256)
func (mouseTier *MouseTier) TryPackMINCEILING() ([]byte, error) {
	return mouseTier.abi.Pack("MIN_CEILING")
}

// UnpackMINCEILING is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3a513af1.
//
// Solidity: function MIN_CEILING() view returns(uint256)
func (mouseTier *MouseTier) UnpackMINCEILING(data []byte) (*big.Int, error) {
	out, err := mouseTier.abi.Unpack("MIN_CEILING", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackXFACTOR is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe3e5edf2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function X_FACTOR() view returns(uint256)
func (mouseTier *MouseTier) PackXFACTOR() []byte {
	enc, err := mouseTier.abi.Pack("X_FACTOR")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackXFACTOR is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe3e5edf2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function X_FACTOR() view returns(uint256)
func (mouseTier *MouseTier) TryPackXFACTOR() ([]byte, error) {
	return mouseTier.abi.Pack("X_FACTOR")
}

// UnpackXFACTOR is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe3e5edf2.
//
// Solidity: function X_FACTOR() view returns(uint256)
func (mouseTier *MouseTier) UnpackXFACTOR(data []byte) (*big.Int, error) {
	out, err := mouseTier.abi.Unpack("X_FACTOR", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCeiling is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x753ed1bd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ceiling() view returns(uint256)
func (mouseTier *MouseTier) PackCeiling() []byte {
	enc, err := mouseTier.abi.Pack("ceiling")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCeiling is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x753ed1bd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ceiling() view returns(uint256)
func (mouseTier *MouseTier) TryPackCeiling() ([]byte, error) {
	return mouseTier.abi.Pack("ceiling")
}

// UnpackCeiling is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x753ed1bd.
//
// Solidity: function ceiling() view returns(uint256)
func (mouseTier *MouseTier) UnpackCeiling(data []byte) (*big.Int, error) {
	out, err := mouseTier.abi.Unpack("ceiling", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetUserCredit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdad6f0dc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getUserCredit(address user) view returns(uint256)
func (mouseTier *MouseTier) PackGetUserCredit(user common.Address) []byte {
	enc, err := mouseTier.abi.Pack("getUserCredit", user)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetUserCredit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdad6f0dc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getUserCredit(address user) view returns(uint256)
func (mouseTier *MouseTier) TryPackGetUserCredit(user common.Address) ([]byte, error) {
	return mouseTier.abi.Pack("getUserCredit", user)
}

// UnpackGetUserCredit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdad6f0dc.
//
// Solidity: function getUserCredit(address user) view returns(uint256)
func (mouseTier *MouseTier) UnpackGetUserCredit(data []byte) (*big.Int, error) {
	out, err := mouseTier.abi.Unpack("getUserCredit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetUserTier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe4d2620e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getUserTier(address user) view returns(uint256)
func (mouseTier *MouseTier) PackGetUserTier(user common.Address) []byte {
	enc, err := mouseTier.abi.Pack("getUserTier", user)
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
func (mouseTier *MouseTier) TryPackGetUserTier(user common.Address) ([]byte, error) {
	return mouseTier.abi.Pack("getUserTier", user)
}

// UnpackGetUserTier is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe4d2620e.
//
// Solidity: function getUserTier(address user) view returns(uint256)
func (mouseTier *MouseTier) UnpackGetUserTier(data []byte) (*big.Int, error) {
	out, err := mouseTier.abi.Unpack("getUserTier", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMousePool is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8b4d97a3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mousePool() view returns(address)
func (mouseTier *MouseTier) PackMousePool() []byte {
	enc, err := mouseTier.abi.Pack("mousePool")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMousePool is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8b4d97a3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mousePool() view returns(address)
func (mouseTier *MouseTier) TryPackMousePool() ([]byte, error) {
	return mouseTier.abi.Pack("mousePool")
}

// UnpackMousePool is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8b4d97a3.
//
// Solidity: function mousePool() view returns(address)
func (mouseTier *MouseTier) UnpackMousePool(data []byte) (common.Address, error) {
	out, err := mouseTier.abi.Unpack("mousePool", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackMultiplier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1b3ed722.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function multiplier() view returns(uint256)
func (mouseTier *MouseTier) PackMultiplier() []byte {
	enc, err := mouseTier.abi.Pack("multiplier")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMultiplier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1b3ed722.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function multiplier() view returns(uint256)
func (mouseTier *MouseTier) TryPackMultiplier() ([]byte, error) {
	return mouseTier.abi.Pack("multiplier")
}

// UnpackMultiplier is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1b3ed722.
//
// Solidity: function multiplier() view returns(uint256)
func (mouseTier *MouseTier) UnpackMultiplier(data []byte) (*big.Int, error) {
	out, err := mouseTier.abi.Unpack("multiplier", data)
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
func (mouseTier *MouseTier) PackOwner() []byte {
	enc, err := mouseTier.abi.Pack("owner")
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
func (mouseTier *MouseTier) TryPackOwner() ([]byte, error) {
	return mouseTier.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (mouseTier *MouseTier) UnpackOwner(data []byte) (common.Address, error) {
	out, err := mouseTier.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() returns()
func (mouseTier *MouseTier) PackRenounceOwnership() []byte {
	enc, err := mouseTier.abi.Pack("renounceOwnership")
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
func (mouseTier *MouseTier) TryPackRenounceOwnership() ([]byte, error) {
	return mouseTier.abi.Pack("renounceOwnership")
}

// PackTierBaseAmount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xde05fd1e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function tierBaseAmount() view returns(uint256)
func (mouseTier *MouseTier) PackTierBaseAmount() []byte {
	enc, err := mouseTier.abi.Pack("tierBaseAmount")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTierBaseAmount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xde05fd1e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function tierBaseAmount() view returns(uint256)
func (mouseTier *MouseTier) TryPackTierBaseAmount() ([]byte, error) {
	return mouseTier.abi.Pack("tierBaseAmount")
}

// UnpackTierBaseAmount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xde05fd1e.
//
// Solidity: function tierBaseAmount() view returns(uint256)
func (mouseTier *MouseTier) UnpackTierBaseAmount(data []byte) (*big.Int, error) {
	out, err := mouseTier.abi.Unpack("tierBaseAmount", data)
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
func (mouseTier *MouseTier) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := mouseTier.abi.Pack("transferOwnership", newOwner)
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
func (mouseTier *MouseTier) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return mouseTier.abi.Pack("transferOwnership", newOwner)
}

// PackUpdateCeiling is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb1463a6a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateCeiling(uint256 newCeiling) returns()
func (mouseTier *MouseTier) PackUpdateCeiling(newCeiling *big.Int) []byte {
	enc, err := mouseTier.abi.Pack("updateCeiling", newCeiling)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateCeiling is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb1463a6a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateCeiling(uint256 newCeiling) returns()
func (mouseTier *MouseTier) TryPackUpdateCeiling(newCeiling *big.Int) ([]byte, error) {
	return mouseTier.abi.Pack("updateCeiling", newCeiling)
}

// PackUpdateMultiplier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5ffe6146.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateMultiplier(uint256 newMultiplier) returns()
func (mouseTier *MouseTier) PackUpdateMultiplier(newMultiplier *big.Int) []byte {
	enc, err := mouseTier.abi.Pack("updateMultiplier", newMultiplier)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateMultiplier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5ffe6146.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateMultiplier(uint256 newMultiplier) returns()
func (mouseTier *MouseTier) TryPackUpdateMultiplier(newMultiplier *big.Int) ([]byte, error) {
	return mouseTier.abi.Pack("updateMultiplier", newMultiplier)
}

// PackUpdateTierBaseAmount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd15f9a3f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateTierBaseAmount(uint256 newTierBaseAmount) returns()
func (mouseTier *MouseTier) PackUpdateTierBaseAmount(newTierBaseAmount *big.Int) []byte {
	enc, err := mouseTier.abi.Pack("updateTierBaseAmount", newTierBaseAmount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateTierBaseAmount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd15f9a3f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateTierBaseAmount(uint256 newTierBaseAmount) returns()
func (mouseTier *MouseTier) TryPackUpdateTierBaseAmount(newTierBaseAmount *big.Int) ([]byte, error) {
	return mouseTier.abi.Pack("updateTierBaseAmount", newTierBaseAmount)
}

// MouseTierOwnershipTransferred represents a OwnershipTransferred event raised by the MouseTier contract.
type MouseTierOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const MouseTierOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (MouseTierOwnershipTransferred) ContractEventName() string {
	return MouseTierOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (mouseTier *MouseTier) UnpackOwnershipTransferredEvent(log *types.Log) (*MouseTierOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mouseTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MouseTierOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := mouseTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mouseTier.abi.Events[event].Inputs {
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

// MouseTierUpdateCeiling represents a UpdateCeiling event raised by the MouseTier contract.
type MouseTierUpdateCeiling struct {
	OldCeiling *big.Int
	NewCeiling *big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const MouseTierUpdateCeilingEventName = "UpdateCeiling"

// ContractEventName returns the user-defined event name.
func (MouseTierUpdateCeiling) ContractEventName() string {
	return MouseTierUpdateCeilingEventName
}

// UnpackUpdateCeilingEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event UpdateCeiling(uint256 oldCeiling, uint256 newCeiling)
func (mouseTier *MouseTier) UnpackUpdateCeilingEvent(log *types.Log) (*MouseTierUpdateCeiling, error) {
	event := "UpdateCeiling"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mouseTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MouseTierUpdateCeiling)
	if len(log.Data) > 0 {
		if err := mouseTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mouseTier.abi.Events[event].Inputs {
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

// MouseTierUpdateMultiplier represents a UpdateMultiplier event raised by the MouseTier contract.
type MouseTierUpdateMultiplier struct {
	OldMultiplier *big.Int
	NewMultiplier *big.Int
	Raw           *types.Log // Blockchain specific contextual infos
}

const MouseTierUpdateMultiplierEventName = "UpdateMultiplier"

// ContractEventName returns the user-defined event name.
func (MouseTierUpdateMultiplier) ContractEventName() string {
	return MouseTierUpdateMultiplierEventName
}

// UnpackUpdateMultiplierEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event UpdateMultiplier(uint256 oldMultiplier, uint256 newMultiplier)
func (mouseTier *MouseTier) UnpackUpdateMultiplierEvent(log *types.Log) (*MouseTierUpdateMultiplier, error) {
	event := "UpdateMultiplier"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mouseTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MouseTierUpdateMultiplier)
	if len(log.Data) > 0 {
		if err := mouseTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mouseTier.abi.Events[event].Inputs {
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

// MouseTierUpdateTierBaseAmount represents a UpdateTierBaseAmount event raised by the MouseTier contract.
type MouseTierUpdateTierBaseAmount struct {
	OldTierBaseAmount *big.Int
	NewTierBaseAmount *big.Int
	Raw               *types.Log // Blockchain specific contextual infos
}

const MouseTierUpdateTierBaseAmountEventName = "UpdateTierBaseAmount"

// ContractEventName returns the user-defined event name.
func (MouseTierUpdateTierBaseAmount) ContractEventName() string {
	return MouseTierUpdateTierBaseAmountEventName
}

// UnpackUpdateTierBaseAmountEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event UpdateTierBaseAmount(uint256 oldTierBaseAmount, uint256 newTierBaseAmount)
func (mouseTier *MouseTier) UnpackUpdateTierBaseAmountEvent(log *types.Log) (*MouseTierUpdateTierBaseAmount, error) {
	event := "UpdateTierBaseAmount"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mouseTier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MouseTierUpdateTierBaseAmount)
	if len(log.Data) > 0 {
		if err := mouseTier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mouseTier.abi.Events[event].Inputs {
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
func (mouseTier *MouseTier) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], mouseTier.abi.Errors["CeilingTooLow"].ID.Bytes()[:4]) {
		return mouseTier.UnpackCeilingTooLowError(raw[4:])
	}
	if bytes.Equal(raw[:4], mouseTier.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return mouseTier.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], mouseTier.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return mouseTier.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], mouseTier.abi.Errors["SameValue"].ID.Bytes()[:4]) {
		return mouseTier.UnpackSameValueError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// MouseTierCeilingTooLow represents a CeilingTooLow error raised by the MouseTier contract.
type MouseTierCeilingTooLow struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CeilingTooLow()
func MouseTierCeilingTooLowErrorID() common.Hash {
	return common.HexToHash("0x1c916a0dfb2db2565e2192adf6222467950984080095bf7b9c16ac66bfd7f16d")
}

// UnpackCeilingTooLowError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CeilingTooLow()
func (mouseTier *MouseTier) UnpackCeilingTooLowError(raw []byte) (*MouseTierCeilingTooLow, error) {
	out := new(MouseTierCeilingTooLow)
	if err := mouseTier.abi.UnpackIntoInterface(out, "CeilingTooLow", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MouseTierOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the MouseTier contract.
type MouseTierOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func MouseTierOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (mouseTier *MouseTier) UnpackOwnableInvalidOwnerError(raw []byte) (*MouseTierOwnableInvalidOwner, error) {
	out := new(MouseTierOwnableInvalidOwner)
	if err := mouseTier.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MouseTierOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the MouseTier contract.
type MouseTierOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func MouseTierOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (mouseTier *MouseTier) UnpackOwnableUnauthorizedAccountError(raw []byte) (*MouseTierOwnableUnauthorizedAccount, error) {
	out := new(MouseTierOwnableUnauthorizedAccount)
	if err := mouseTier.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MouseTierSameValue represents a SameValue error raised by the MouseTier contract.
type MouseTierSameValue struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SameValue()
func MouseTierSameValueErrorID() common.Hash {
	return common.HexToHash("0xc23f6ccb763204b707e413082835ff553e80f1d949cec72488d8558030f5eb8a")
}

// UnpackSameValueError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SameValue()
func (mouseTier *MouseTier) UnpackSameValueError(raw []byte) (*MouseTierSameValue, error) {
	out := new(MouseTierSameValue)
	if err := mouseTier.abi.UnpackIntoInterface(out, "SameValue", raw); err != nil {
		return nil, err
	}
	return out, nil
}

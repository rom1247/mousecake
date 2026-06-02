// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package token

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

// MouseTokenMetaData contains all meta data concerning the MouseToken contract.
var MouseTokenMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MOUSE_TOTAL_SUPPLY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"masterChef\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mintTo\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMasterChef\",\"inputs\":[{\"name\":\"newMasterChef\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MasterChefUpdated\",\"inputs\":[{\"name\":\"oldMasterChef\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newMasterChef\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MintTo\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ExceedsTotalSupply\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotMasterChef\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	ID:  "MouseToken",
	Bin: "0x608060405234801561000f575f5ffd5b50336040518060400160405280600581526020017f4d6f7573650000000000000000000000000000000000000000000000000000008152506040518060400160405280600581526020017f4d4f555345000000000000000000000000000000000000000000000000000000815250816003908161008c919061045a565b50806004908161009c919061045a565b5050505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361010f575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016101069190610568565b60405180910390fd5b61011e8161012460201b60201c565b50610581565b60065f6101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556101578161015a60201b60201c565b50565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061029857607f821691505b6020821081036102ab576102aa610254565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f6008830261030d7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826102d2565b61031786836102d2565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f61035b6103566103518461032f565b610338565b61032f565b9050919050565b5f819050919050565b61037483610341565b61038861038082610362565b8484546102de565b825550505050565b5f5f905090565b61039f610390565b6103aa81848461036b565b505050565b5b818110156103cd576103c25f82610397565b6001810190506103b0565b5050565b601f821115610412576103e3816102b1565b6103ec846102c3565b810160208510156103fb578190505b61040f610407856102c3565b8301826103af565b50505b505050565b5f82821c905092915050565b5f6104325f1984600802610417565b1980831691505092915050565b5f61044a8383610423565b9150826002028217905092915050565b6104638261021d565b67ffffffffffffffff81111561047c5761047b610227565b5b6104868254610281565b6104918282856103d1565b5f60209050601f8311600181146104c2575f84156104b0578287015190505b6104ba858261043f565b865550610521565b601f1984166104d0866102b1565b5f5b828110156104f7578489015182556001820191506020850194506020810190506104d2565b868310156105145784890151610510601f891682610423565b8355505b6001600288020188555050505b505050505050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61055282610529565b9050919050565b61056281610548565b82525050565b5f60208201905061057b5f830184610559565b92915050565b6115138061058e5f395ff3fe608060405234801561000f575f5ffd5b5060043610610114575f3560e01c8063715018a6116100a0578063a2d9f4dc1161006f578063a2d9f4dc146102aa578063a9059cbb146102c6578063dd62ed3e146102f6578063e30c397814610326578063f2fde38b1461034457610114565b8063715018a61461025a57806379ba5097146102645780638da5cb5b1461026e57806395d89b411461028c57610114565b8063313ce567116100e7578063313ce567146101b4578063449a52f8146101d2578063575a86b2146101ee578063636c86eb1461020c57806370a082311461022a57610114565b806306fdde0314610118578063095ea7b31461013657806318160ddd1461016657806323b872dd14610184575b5f5ffd5b610120610360565b60405161012d919061118c565b60405180910390f35b610150600480360381019061014b919061123d565b6103f0565b60405161015d9190611295565b60405180910390f35b61016e610412565b60405161017b91906112bd565b60405180910390f35b61019e600480360381019061019991906112d6565b61041b565b6040516101ab9190611295565b60405180910390f35b6101bc610449565b6040516101c99190611341565b60405180910390f35b6101ec60048036038101906101e7919061123d565b610451565b005b6101f661058a565b6040516102039190611369565b60405180910390f35b6102146105af565b60405161022191906112bd565b60405180910390f35b610244600480360381019061023f9190611382565b6105be565b60405161025191906112bd565b60405180910390f35b610262610603565b005b61026c610616565b005b6102766106a4565b6040516102839190611369565b60405180910390f35b6102946106cc565b6040516102a1919061118c565b60405180910390f35b6102c460048036038101906102bf9190611382565b61075c565b005b6102e060048036038101906102db919061123d565b610827565b6040516102ed9190611295565b60405180910390f35b610310600480360381019061030b91906113ad565b610849565b60405161031d91906112bd565b60405180910390f35b61032e6108cb565b60405161033b9190611369565b60405180910390f35b61035e60048036038101906103599190611382565b6108f3565b005b60606003805461036f90611418565b80601f016020809104026020016040519081016040528092919081815260200182805461039b90611418565b80156103e65780601f106103bd576101008083540402835291602001916103e6565b820191905f5260205f20905b8154815290600101906020018083116103c957829003601f168201915b5050505050905090565b5f5f6103fa61099f565b90506104078185856109a6565b600191505092915050565b5f600254905090565b5f5f61042561099f565b90506104328582856109b8565b61043d858585610a4b565b60019150509392505050565b5f6012905090565b60075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146104d7576040517f0868337f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6a5f20327de60ebf53000000816104ec610412565b6104f69190611475565b111561052e576040517f177e3fc300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105388282610b3b565b8173ffffffffffffffffffffffffffffffffffffffff167f2ccee9bef0756ce350dcfbf029ffa186a686bc7f7f391ff02e6b686f8650ab258260405161057e91906112bd565b60405180910390a25050565b60075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6a5f20327de60ebf5300000081565b5f5f5f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b61060b610bba565b6106145f610c41565b565b5f61061f61099f565b90508073ffffffffffffffffffffffffffffffffffffffff166106406108cb565b73ffffffffffffffffffffffffffffffffffffffff161461069857806040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161068f9190611369565b60405180910390fd5b6106a181610c41565b50565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6060600480546106db90611418565b80601f016020809104026020016040519081016040528092919081815260200182805461070790611418565b80156107525780601f1061072957610100808354040283529160200191610752565b820191905f5260205f20905b81548152906001019060200180831161073557829003601f168201915b5050505050905090565b610764610bba565b5f60075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160075f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f6c895b98ec822731298010df62b63c1495016883018b708debe95c645e26348a60405160405180910390a35050565b5f5f61083161099f565b905061083e818585610a4b565b600191505092915050565b5f60015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905092915050565b5f60065f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6108fb610bba565b8060065f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff1661095a6106a4565b73ffffffffffffffffffffffffffffffffffffffff167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a350565b5f33905090565b6109b38383836001610c71565b505050565b5f6109c38484610849565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811015610a455781811015610a36578281836040517ffb8f41b2000000000000000000000000000000000000000000000000000000008152600401610a2d939291906114a8565b60405180910390fd5b610a4484848484035f610c71565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610abb575f6040517f96c6fd1e000000000000000000000000000000000000000000000000000000008152600401610ab29190611369565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610b2b575f6040517fec442f05000000000000000000000000000000000000000000000000000000008152600401610b229190611369565b60405180910390fd5b610b36838383610e40565b505050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610bab575f6040517fec442f05000000000000000000000000000000000000000000000000000000008152600401610ba29190611369565b60405180910390fd5b610bb65f8383610e40565b5050565b610bc261099f565b73ffffffffffffffffffffffffffffffffffffffff16610be06106a4565b73ffffffffffffffffffffffffffffffffffffffff1614610c3f57610c0361099f565b6040517f118cdaa7000000000000000000000000000000000000000000000000000000008152600401610c369190611369565b60405180910390fd5b565b60065f6101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055610c6e81611059565b50565b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603610ce1575f6040517fe602df05000000000000000000000000000000000000000000000000000000008152600401610cd89190611369565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610d51575f6040517f94280d62000000000000000000000000000000000000000000000000000000008152600401610d489190611369565b60405180910390fd5b8160015f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508015610e3a578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92584604051610e3191906112bd565b60405180910390a35b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610e90578060025f828254610e849190611475565b92505081905550610f5e565b5f5f5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905081811015610f19578381836040517fe450d38c000000000000000000000000000000000000000000000000000000008152600401610f10939291906114a8565b60405180910390fd5b8181035f5f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550505b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610fa5578060025f8282540392505081905550610fef565b805f5f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161104c91906112bd565b60405180910390a3505050565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f61115e8261111c565b6111688185611126565b9350611178818560208601611136565b61118181611144565b840191505092915050565b5f6020820190508181035f8301526111a48184611154565b905092915050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6111d9826111b0565b9050919050565b6111e9816111cf565b81146111f3575f5ffd5b50565b5f81359050611204816111e0565b92915050565b5f819050919050565b61121c8161120a565b8114611226575f5ffd5b50565b5f8135905061123781611213565b92915050565b5f5f60408385031215611253576112526111ac565b5b5f611260858286016111f6565b925050602061127185828601611229565b9150509250929050565b5f8115159050919050565b61128f8161127b565b82525050565b5f6020820190506112a85f830184611286565b92915050565b6112b78161120a565b82525050565b5f6020820190506112d05f8301846112ae565b92915050565b5f5f5f606084860312156112ed576112ec6111ac565b5b5f6112fa868287016111f6565b935050602061130b868287016111f6565b925050604061131c86828701611229565b9150509250925092565b5f60ff82169050919050565b61133b81611326565b82525050565b5f6020820190506113545f830184611332565b92915050565b611363816111cf565b82525050565b5f60208201905061137c5f83018461135a565b92915050565b5f60208284031215611397576113966111ac565b5b5f6113a4848285016111f6565b91505092915050565b5f5f604083850312156113c3576113c26111ac565b5b5f6113d0858286016111f6565b92505060206113e1858286016111f6565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061142f57607f821691505b602082108103611442576114416113eb565b5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61147f8261120a565b915061148a8361120a565b92508282019050808211156114a2576114a1611448565b5b92915050565b5f6060820190506114bb5f83018661135a565b6114c860208301856112ae565b6114d560408301846112ae565b94935050505056fea264697066735822122032b1c897bc4754a98c8d3f89cccf7eb39ac9f02107389d77f2d5774cd46ede4664736f6c634300081c0033",
}

// MouseToken is an auto generated Go binding around an Ethereum contract.
type MouseToken struct {
	abi abi.ABI
}

// NewMouseToken creates a new instance of MouseToken.
func NewMouseToken() *MouseToken {
	parsed, err := MouseTokenMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &MouseToken{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *MouseToken) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackMOUSETOTALSUPPLY is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x636c86eb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MOUSE_TOTAL_SUPPLY() view returns(uint256)
func (mouseToken *MouseToken) PackMOUSETOTALSUPPLY() []byte {
	enc, err := mouseToken.abi.Pack("MOUSE_TOTAL_SUPPLY")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMOUSETOTALSUPPLY is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x636c86eb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MOUSE_TOTAL_SUPPLY() view returns(uint256)
func (mouseToken *MouseToken) TryPackMOUSETOTALSUPPLY() ([]byte, error) {
	return mouseToken.abi.Pack("MOUSE_TOTAL_SUPPLY")
}

// UnpackMOUSETOTALSUPPLY is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x636c86eb.
//
// Solidity: function MOUSE_TOTAL_SUPPLY() view returns(uint256)
func (mouseToken *MouseToken) UnpackMOUSETOTALSUPPLY(data []byte) (*big.Int, error) {
	out, err := mouseToken.abi.Unpack("MOUSE_TOTAL_SUPPLY", data)
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
func (mouseToken *MouseToken) PackAcceptOwnership() []byte {
	enc, err := mouseToken.abi.Pack("acceptOwnership")
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
func (mouseToken *MouseToken) TryPackAcceptOwnership() ([]byte, error) {
	return mouseToken.abi.Pack("acceptOwnership")
}

// PackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (mouseToken *MouseToken) PackAllowance(owner common.Address, spender common.Address) []byte {
	enc, err := mouseToken.abi.Pack("allowance", owner, spender)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (mouseToken *MouseToken) TryPackAllowance(owner common.Address, spender common.Address) ([]byte, error) {
	return mouseToken.abi.Pack("allowance", owner, spender)
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (mouseToken *MouseToken) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := mouseToken.abi.Unpack("allowance", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (mouseToken *MouseToken) PackApprove(spender common.Address, value *big.Int) []byte {
	enc, err := mouseToken.abi.Pack("approve", spender, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (mouseToken *MouseToken) TryPackApprove(spender common.Address, value *big.Int) ([]byte, error) {
	return mouseToken.abi.Pack("approve", spender, value)
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (mouseToken *MouseToken) UnpackApprove(data []byte) (bool, error) {
	out, err := mouseToken.abi.Unpack("approve", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (mouseToken *MouseToken) PackBalanceOf(account common.Address) []byte {
	enc, err := mouseToken.abi.Pack("balanceOf", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (mouseToken *MouseToken) TryPackBalanceOf(account common.Address) ([]byte, error) {
	return mouseToken.abi.Pack("balanceOf", account)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (mouseToken *MouseToken) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := mouseToken.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function decimals() view returns(uint8)
func (mouseToken *MouseToken) PackDecimals() []byte {
	enc, err := mouseToken.abi.Pack("decimals")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function decimals() view returns(uint8)
func (mouseToken *MouseToken) TryPackDecimals() ([]byte, error) {
	return mouseToken.abi.Pack("decimals")
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (mouseToken *MouseToken) UnpackDecimals(data []byte) (uint8, error) {
	out, err := mouseToken.abi.Unpack("decimals", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackMasterChef is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x575a86b2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function masterChef() view returns(address)
func (mouseToken *MouseToken) PackMasterChef() []byte {
	enc, err := mouseToken.abi.Pack("masterChef")
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
func (mouseToken *MouseToken) TryPackMasterChef() ([]byte, error) {
	return mouseToken.abi.Pack("masterChef")
}

// UnpackMasterChef is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x575a86b2.
//
// Solidity: function masterChef() view returns(address)
func (mouseToken *MouseToken) UnpackMasterChef(data []byte) (common.Address, error) {
	out, err := mouseToken.abi.Unpack("masterChef", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackMintTo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x449a52f8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mintTo(address to, uint256 amount) returns()
func (mouseToken *MouseToken) PackMintTo(to common.Address, amount *big.Int) []byte {
	enc, err := mouseToken.abi.Pack("mintTo", to, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMintTo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x449a52f8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mintTo(address to, uint256 amount) returns()
func (mouseToken *MouseToken) TryPackMintTo(to common.Address, amount *big.Int) ([]byte, error) {
	return mouseToken.abi.Pack("mintTo", to, amount)
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function name() view returns(string)
func (mouseToken *MouseToken) PackName() []byte {
	enc, err := mouseToken.abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function name() view returns(string)
func (mouseToken *MouseToken) TryPackName() ([]byte, error) {
	return mouseToken.abi.Pack("name")
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (mouseToken *MouseToken) UnpackName(data []byte) (string, error) {
	out, err := mouseToken.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (mouseToken *MouseToken) PackOwner() []byte {
	enc, err := mouseToken.abi.Pack("owner")
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
func (mouseToken *MouseToken) TryPackOwner() ([]byte, error) {
	return mouseToken.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (mouseToken *MouseToken) UnpackOwner(data []byte) (common.Address, error) {
	out, err := mouseToken.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPendingOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe30c3978.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingOwner() view returns(address)
func (mouseToken *MouseToken) PackPendingOwner() []byte {
	enc, err := mouseToken.abi.Pack("pendingOwner")
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
func (mouseToken *MouseToken) TryPackPendingOwner() ([]byte, error) {
	return mouseToken.abi.Pack("pendingOwner")
}

// UnpackPendingOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (mouseToken *MouseToken) UnpackPendingOwner(data []byte) (common.Address, error) {
	out, err := mouseToken.abi.Unpack("pendingOwner", data)
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
func (mouseToken *MouseToken) PackRenounceOwnership() []byte {
	enc, err := mouseToken.abi.Pack("renounceOwnership")
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
func (mouseToken *MouseToken) TryPackRenounceOwnership() ([]byte, error) {
	return mouseToken.abi.Pack("renounceOwnership")
}

// PackSetMasterChef is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa2d9f4dc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setMasterChef(address newMasterChef) returns()
func (mouseToken *MouseToken) PackSetMasterChef(newMasterChef common.Address) []byte {
	enc, err := mouseToken.abi.Pack("setMasterChef", newMasterChef)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetMasterChef is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa2d9f4dc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setMasterChef(address newMasterChef) returns()
func (mouseToken *MouseToken) TryPackSetMasterChef(newMasterChef common.Address) ([]byte, error) {
	return mouseToken.abi.Pack("setMasterChef", newMasterChef)
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function symbol() view returns(string)
func (mouseToken *MouseToken) PackSymbol() []byte {
	enc, err := mouseToken.abi.Pack("symbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function symbol() view returns(string)
func (mouseToken *MouseToken) TryPackSymbol() ([]byte, error) {
	return mouseToken.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (mouseToken *MouseToken) UnpackSymbol(data []byte) (string, error) {
	out, err := mouseToken.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalSupply() view returns(uint256)
func (mouseToken *MouseToken) PackTotalSupply() []byte {
	enc, err := mouseToken.abi.Pack("totalSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalSupply() view returns(uint256)
func (mouseToken *MouseToken) TryPackTotalSupply() ([]byte, error) {
	return mouseToken.abi.Pack("totalSupply")
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (mouseToken *MouseToken) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := mouseToken.abi.Unpack("totalSupply", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (mouseToken *MouseToken) PackTransfer(to common.Address, value *big.Int) []byte {
	enc, err := mouseToken.abi.Pack("transfer", to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (mouseToken *MouseToken) TryPackTransfer(to common.Address, value *big.Int) ([]byte, error) {
	return mouseToken.abi.Pack("transfer", to, value)
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (mouseToken *MouseToken) UnpackTransfer(data []byte) (bool, error) {
	out, err := mouseToken.abi.Unpack("transfer", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (mouseToken *MouseToken) PackTransferFrom(from common.Address, to common.Address, value *big.Int) []byte {
	enc, err := mouseToken.abi.Pack("transferFrom", from, to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (mouseToken *MouseToken) TryPackTransferFrom(from common.Address, to common.Address, value *big.Int) ([]byte, error) {
	return mouseToken.abi.Pack("transferFrom", from, to, value)
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (mouseToken *MouseToken) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := mouseToken.abi.Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (mouseToken *MouseToken) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := mouseToken.abi.Pack("transferOwnership", newOwner)
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
func (mouseToken *MouseToken) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return mouseToken.abi.Pack("transferOwnership", newOwner)
}

// MouseTokenApproval represents a Approval event raised by the MouseToken contract.
type MouseTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const MouseTokenApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (MouseTokenApproval) ContractEventName() string {
	return MouseTokenApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (mouseToken *MouseToken) UnpackApprovalEvent(log *types.Log) (*MouseTokenApproval, error) {
	event := "Approval"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mouseToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MouseTokenApproval)
	if len(log.Data) > 0 {
		if err := mouseToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mouseToken.abi.Events[event].Inputs {
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

// MouseTokenMasterChefUpdated represents a MasterChefUpdated event raised by the MouseToken contract.
type MouseTokenMasterChefUpdated struct {
	OldMasterChef common.Address
	NewMasterChef common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const MouseTokenMasterChefUpdatedEventName = "MasterChefUpdated"

// ContractEventName returns the user-defined event name.
func (MouseTokenMasterChefUpdated) ContractEventName() string {
	return MouseTokenMasterChefUpdatedEventName
}

// UnpackMasterChefUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MasterChefUpdated(address indexed oldMasterChef, address indexed newMasterChef)
func (mouseToken *MouseToken) UnpackMasterChefUpdatedEvent(log *types.Log) (*MouseTokenMasterChefUpdated, error) {
	event := "MasterChefUpdated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mouseToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MouseTokenMasterChefUpdated)
	if len(log.Data) > 0 {
		if err := mouseToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mouseToken.abi.Events[event].Inputs {
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

// MouseTokenMintTo represents a MintTo event raised by the MouseToken contract.
type MouseTokenMintTo struct {
	To     common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const MouseTokenMintToEventName = "MintTo"

// ContractEventName returns the user-defined event name.
func (MouseTokenMintTo) ContractEventName() string {
	return MouseTokenMintToEventName
}

// UnpackMintToEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MintTo(address indexed to, uint256 amount)
func (mouseToken *MouseToken) UnpackMintToEvent(log *types.Log) (*MouseTokenMintTo, error) {
	event := "MintTo"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mouseToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MouseTokenMintTo)
	if len(log.Data) > 0 {
		if err := mouseToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mouseToken.abi.Events[event].Inputs {
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

// MouseTokenOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the MouseToken contract.
type MouseTokenOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const MouseTokenOwnershipTransferStartedEventName = "OwnershipTransferStarted"

// ContractEventName returns the user-defined event name.
func (MouseTokenOwnershipTransferStarted) ContractEventName() string {
	return MouseTokenOwnershipTransferStartedEventName
}

// UnpackOwnershipTransferStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (mouseToken *MouseToken) UnpackOwnershipTransferStartedEvent(log *types.Log) (*MouseTokenOwnershipTransferStarted, error) {
	event := "OwnershipTransferStarted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mouseToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MouseTokenOwnershipTransferStarted)
	if len(log.Data) > 0 {
		if err := mouseToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mouseToken.abi.Events[event].Inputs {
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

// MouseTokenOwnershipTransferred represents a OwnershipTransferred event raised by the MouseToken contract.
type MouseTokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const MouseTokenOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (MouseTokenOwnershipTransferred) ContractEventName() string {
	return MouseTokenOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (mouseToken *MouseToken) UnpackOwnershipTransferredEvent(log *types.Log) (*MouseTokenOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mouseToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MouseTokenOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := mouseToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mouseToken.abi.Events[event].Inputs {
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

// MouseTokenTransfer represents a Transfer event raised by the MouseToken contract.
type MouseTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const MouseTokenTransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (MouseTokenTransfer) ContractEventName() string {
	return MouseTokenTransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (mouseToken *MouseToken) UnpackTransferEvent(log *types.Log) (*MouseTokenTransfer, error) {
	event := "Transfer"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != mouseToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MouseTokenTransfer)
	if len(log.Data) > 0 {
		if err := mouseToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mouseToken.abi.Events[event].Inputs {
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
func (mouseToken *MouseToken) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], mouseToken.abi.Errors["ERC20InsufficientAllowance"].ID.Bytes()[:4]) {
		return mouseToken.UnpackERC20InsufficientAllowanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], mouseToken.abi.Errors["ERC20InsufficientBalance"].ID.Bytes()[:4]) {
		return mouseToken.UnpackERC20InsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], mouseToken.abi.Errors["ERC20InvalidApprover"].ID.Bytes()[:4]) {
		return mouseToken.UnpackERC20InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], mouseToken.abi.Errors["ERC20InvalidReceiver"].ID.Bytes()[:4]) {
		return mouseToken.UnpackERC20InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], mouseToken.abi.Errors["ERC20InvalidSender"].ID.Bytes()[:4]) {
		return mouseToken.UnpackERC20InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], mouseToken.abi.Errors["ERC20InvalidSpender"].ID.Bytes()[:4]) {
		return mouseToken.UnpackERC20InvalidSpenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], mouseToken.abi.Errors["ExceedsTotalSupply"].ID.Bytes()[:4]) {
		return mouseToken.UnpackExceedsTotalSupplyError(raw[4:])
	}
	if bytes.Equal(raw[:4], mouseToken.abi.Errors["NotMasterChef"].ID.Bytes()[:4]) {
		return mouseToken.UnpackNotMasterChefError(raw[4:])
	}
	if bytes.Equal(raw[:4], mouseToken.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return mouseToken.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], mouseToken.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return mouseToken.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// MouseTokenERC20InsufficientAllowance represents a ERC20InsufficientAllowance error raised by the MouseToken contract.
type MouseTokenERC20InsufficientAllowance struct {
	Spender   common.Address
	Allowance *big.Int
	Needed    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func MouseTokenERC20InsufficientAllowanceErrorID() common.Hash {
	return common.HexToHash("0xfb8f41b23e99d2101d86da76cdfa87dd51c82ed07d3cb62cbc473e469dbc75c3")
}

// UnpackERC20InsufficientAllowanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func (mouseToken *MouseToken) UnpackERC20InsufficientAllowanceError(raw []byte) (*MouseTokenERC20InsufficientAllowance, error) {
	out := new(MouseTokenERC20InsufficientAllowance)
	if err := mouseToken.abi.UnpackIntoInterface(out, "ERC20InsufficientAllowance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MouseTokenERC20InsufficientBalance represents a ERC20InsufficientBalance error raised by the MouseToken contract.
type MouseTokenERC20InsufficientBalance struct {
	Sender  common.Address
	Balance *big.Int
	Needed  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func MouseTokenERC20InsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xe450d38cd8d9f7d95077d567d60ed49c7254716e6ad08fc9872816c97e0ffec6")
}

// UnpackERC20InsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func (mouseToken *MouseToken) UnpackERC20InsufficientBalanceError(raw []byte) (*MouseTokenERC20InsufficientBalance, error) {
	out := new(MouseTokenERC20InsufficientBalance)
	if err := mouseToken.abi.UnpackIntoInterface(out, "ERC20InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MouseTokenERC20InvalidApprover represents a ERC20InvalidApprover error raised by the MouseToken contract.
type MouseTokenERC20InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidApprover(address approver)
func MouseTokenERC20InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0xe602df05cc75712490294c6c104ab7c17f4030363910a7a2626411c6d3118847")
}

// UnpackERC20InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidApprover(address approver)
func (mouseToken *MouseToken) UnpackERC20InvalidApproverError(raw []byte) (*MouseTokenERC20InvalidApprover, error) {
	out := new(MouseTokenERC20InvalidApprover)
	if err := mouseToken.abi.UnpackIntoInterface(out, "ERC20InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MouseTokenERC20InvalidReceiver represents a ERC20InvalidReceiver error raised by the MouseToken contract.
type MouseTokenERC20InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func MouseTokenERC20InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0xec442f055133b72f3b2f9f0bb351c406b178527de2040a7d1feb4e058771f613")
}

// UnpackERC20InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func (mouseToken *MouseToken) UnpackERC20InvalidReceiverError(raw []byte) (*MouseTokenERC20InvalidReceiver, error) {
	out := new(MouseTokenERC20InvalidReceiver)
	if err := mouseToken.abi.UnpackIntoInterface(out, "ERC20InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MouseTokenERC20InvalidSender represents a ERC20InvalidSender error raised by the MouseToken contract.
type MouseTokenERC20InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSender(address sender)
func MouseTokenERC20InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x96c6fd1edd0cd6ef7ff0ecc0facdf53148dc0048b57fe58af65755250a7a96bd")
}

// UnpackERC20InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSender(address sender)
func (mouseToken *MouseToken) UnpackERC20InvalidSenderError(raw []byte) (*MouseTokenERC20InvalidSender, error) {
	out := new(MouseTokenERC20InvalidSender)
	if err := mouseToken.abi.UnpackIntoInterface(out, "ERC20InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MouseTokenERC20InvalidSpender represents a ERC20InvalidSpender error raised by the MouseToken contract.
type MouseTokenERC20InvalidSpender struct {
	Spender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSpender(address spender)
func MouseTokenERC20InvalidSpenderErrorID() common.Hash {
	return common.HexToHash("0x94280d62c347d8d9f4d59a76ea321452406db88df38e0c9da304f58b57b373a2")
}

// UnpackERC20InvalidSpenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSpender(address spender)
func (mouseToken *MouseToken) UnpackERC20InvalidSpenderError(raw []byte) (*MouseTokenERC20InvalidSpender, error) {
	out := new(MouseTokenERC20InvalidSpender)
	if err := mouseToken.abi.UnpackIntoInterface(out, "ERC20InvalidSpender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MouseTokenExceedsTotalSupply represents a ExceedsTotalSupply error raised by the MouseToken contract.
type MouseTokenExceedsTotalSupply struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ExceedsTotalSupply()
func MouseTokenExceedsTotalSupplyErrorID() common.Hash {
	return common.HexToHash("0x177e3fc3206c5e813dba3ed57dce92e4456b978aa76321d4f35a1188a809d8a2")
}

// UnpackExceedsTotalSupplyError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ExceedsTotalSupply()
func (mouseToken *MouseToken) UnpackExceedsTotalSupplyError(raw []byte) (*MouseTokenExceedsTotalSupply, error) {
	out := new(MouseTokenExceedsTotalSupply)
	if err := mouseToken.abi.UnpackIntoInterface(out, "ExceedsTotalSupply", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MouseTokenNotMasterChef represents a NotMasterChef error raised by the MouseToken contract.
type MouseTokenNotMasterChef struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotMasterChef()
func MouseTokenNotMasterChefErrorID() common.Hash {
	return common.HexToHash("0x0868337f8a8a9c9d866aea06afae409ee26fb409701dd2dcfe7f6cdb810ed444")
}

// UnpackNotMasterChefError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotMasterChef()
func (mouseToken *MouseToken) UnpackNotMasterChefError(raw []byte) (*MouseTokenNotMasterChef, error) {
	out := new(MouseTokenNotMasterChef)
	if err := mouseToken.abi.UnpackIntoInterface(out, "NotMasterChef", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MouseTokenOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the MouseToken contract.
type MouseTokenOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func MouseTokenOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (mouseToken *MouseToken) UnpackOwnableInvalidOwnerError(raw []byte) (*MouseTokenOwnableInvalidOwner, error) {
	out := new(MouseTokenOwnableInvalidOwner)
	if err := mouseToken.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MouseTokenOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the MouseToken contract.
type MouseTokenOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func MouseTokenOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (mouseToken *MouseToken) UnpackOwnableUnauthorizedAccountError(raw []byte) (*MouseTokenOwnableUnauthorizedAccount, error) {
	out := new(MouseTokenOwnableUnauthorizedAccount)
	if err := mouseToken.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

package chain

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"

	"github.com/mousecake-go/mousecake-go/internal/chain/contract/launchpad"
	"github.com/mousecake-go/mousecake-go/internal/chain/contract/masterchef"
	"github.com/mousecake-go/mousecake-go/internal/chain/contract/staking"
	"github.com/mousecake-go/mousecake-go/internal/chain/contract/token"
)

// RegisterAllContracts 将所有已知合约的 ABI 注册到注册中心。
// addresses 为可选的合约地址映射，key 为合约名，value 为十六进制地址。
func RegisterAllContracts(r *ABIRegistry, addresses map[string]string) {
	register := func(name string, metaData *bind.MetaData, addrHex string) {
		var addr common.Address
		if addrHex != "" {
			addr = common.HexToAddress(addrHex)
		}
		_ = r.Register(name, metaData, addr)
	}

	register("MousePadByTier", &launchpad.MousePadByTierMetaData, addresses["MousePadByTier"])
	register("MousePadByTierDeployer", &launchpad.MousePadByTierDeployerMetaData, addresses["MousePadByTierDeployer"])
	register("MouseTier", &launchpad.MouseTierMetaData, addresses["MouseTier"])
	register("MouseToken", &token.MouseTokenMetaData, addresses["MouseToken"])
	register("MousePool", &staking.MousePoolMetaData, addresses["MousePool"])
	register("MouseMasterChef", &masterchef.MouseMasterChefMetaData, addresses["MouseMasterChef"])
}

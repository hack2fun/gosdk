package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func ConvertCosmosAddressToEthAddress(addrString string) (string, error) {
	if addrString == "" {
		return "", fmt.Errorf("address is empty")
	}
	conf := sdk.GetConfig()

	var addr []byte
	switch {
	case strings.HasPrefix(addrString, conf.GetBech32ValidatorAddrPrefix()):
		addr, _ = sdk.ValAddressFromBech32(addrString)
	case strings.HasPrefix(addrString, conf.GetBech32AccountAddrPrefix()):
		addr, _ = sdk.AccAddressFromBech32(addrString)
	default:
		return "", fmt.Errorf("expected a valid bech32 address (acc prefix %s), got '%s'", conf.GetBech32AccountAddrPrefix(), addrString)
	}
	return common.BytesToAddress(addr).String(), nil
}

func ConvertEthAddressToCosmosAddress(addrString string) (string, error) {
	if addrString == "" {
		return "", fmt.Errorf("address is empty")
	}

	var addr []byte
	switch {
	case common.IsHexAddress(addrString):
		addr = common.HexToAddress(addrString).Bytes()
	default:
		return "", fmt.Errorf("expected a valid hex address, got '%s'", addrString)
	}
	return sdk.AccAddress(addr).String(), nil
}

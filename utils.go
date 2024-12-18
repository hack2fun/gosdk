package gosdk

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func ConvertAddress(addrString string) (ethAddr, cosmosAddr string, err error) {
	if addrString == "" {
		return "", "", fmt.Errorf("addr can't be empty")
	}

	var addr []byte
	switch {
	case common.IsHexAddress(addrString):
		addr = common.HexToAddress(addrString).Bytes()
	case strings.HasPrefix(addrString, Bech32PrefixValAddr):
		addr, _ = sdk.ValAddressFromBech32(addrString)
	case strings.HasPrefix(addrString, Bech32PrefixAccAddr):
		addr, _ = sdk.AccAddressFromBech32(addrString)
	default:
		return "", "", fmt.Errorf("expected a valid hex or bech32 address (acc prefix %s), got '%s'",
			Bech32PrefixValAddr, addrString)
	}

	return common.BytesToAddress(addr).Hex(), sdk.AccAddress(addr).String(), nil
}

func ConvertToCosmosAddress(addrString string) (string, error) {
	if addrString == "" {
		return "", fmt.Errorf("addr can't be empty")
	}

	var addr []byte
	switch {
	case common.IsHexAddress(addrString):
		addr = common.HexToAddress(addrString).Bytes()
	case strings.HasPrefix(addrString, Bech32PrefixValAddr):
		addr, _ = sdk.ValAddressFromBech32(addrString)
	case strings.HasPrefix(addrString, Bech32PrefixAccAddr):
		addr, _ = sdk.AccAddressFromBech32(addrString)
	default:
		return "", fmt.Errorf("expected a valid hex or bech32 address (acc prefix %s), got '%s'",
			Bech32PrefixValAddr, addrString)
	}

	return sdk.AccAddress(addr).String(), nil
}

func ConvertToETHAddress(addrString string) (string, error) {
	if addrString == "" {
		return "", fmt.Errorf("addr can't be empty")
	}

	var addr []byte
	switch {
	case common.IsHexAddress(addrString):
		addr = common.HexToAddress(addrString).Bytes()
	case strings.HasPrefix(addrString, Bech32PrefixValAddr):
		addr, _ = sdk.ValAddressFromBech32(addrString)
	case strings.HasPrefix(addrString, Bech32PrefixAccAddr):
		addr, _ = sdk.AccAddressFromBech32(addrString)
	default:
		return "", fmt.Errorf("expected a valid hex or bech32 address (acc prefix %s), got '%s'",
			Bech32PrefixValAddr, addrString)
	}

	return common.BytesToAddress(addr).Hex(), nil
}

func ConvertCosmosAddressToEthAddress(addrString string) (string, error) {
	if addrString == "" {
		return "", fmt.Errorf("address is empty")
	}

	var addr []byte
	switch {
	case strings.HasPrefix(addrString, Bech32PrefixValAddr):
		addr, _ = sdk.ValAddressFromBech32(addrString)
	case strings.HasPrefix(addrString, Bech32PrefixAccAddr):
		addr, _ = sdk.AccAddressFromBech32(addrString)
	default:
		return "", fmt.Errorf("expected a valid bech32 address (acc prefix %s), got '%s'", Bech32PrefixAccAddr, addrString)
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

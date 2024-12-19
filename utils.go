package gosdk

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// ConvertAddress converts an address string to both Ethereum and Cysic addresses.
//
// @param addrString the address string to convert
// @return the Ethereum address, Cysic address, or an error if conversion fails
func ConvertAddress(addrString string) (ethAddr, cysicAddr string, err error) {
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

// ConvertToCysicAddress converts an address string to a Cysic address.
//
// @param addrString the address string to convert
// @return the Cysic address, or an error if conversion fails
func ConvertToCysicAddress(addrString string) (string, error) {
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

// ConvertToETHAddress converts an address string to an Ethereum address.
//
// @param addrString the address string to convert
// @return the Ethereum address, or an error if conversion fails
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

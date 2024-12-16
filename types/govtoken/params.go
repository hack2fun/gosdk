package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

func validateExchangeRate(i interface{}) error {
	rate, ok := i.(uint64)
	if !ok || rate == 0 {
		return ErrInvalidRate
	}
	return nil
}

package gosdk

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func init() {
	conf := sdk.GetConfig()
	SetBech32Prefixes(conf)
	SetBip44CoinType(conf)
}

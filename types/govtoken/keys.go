package types

const (
	ModuleName   = "govtoken"
	StoreKey     = "gvtkn"
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	ModuleAccountName = "govtoken"
)

var (
	GovTokenKey     = []byte{0x01}
	ParamsKey       = []byte{0x02}
	ExchangeableKey = []byte{0x03}

	ExchangeRateKey = []byte("govtoken_exchange_rate")
)

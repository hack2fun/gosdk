package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrUnauthorized      = errorsmod.Register(ModuleName, 1, "unauthorized")
	ErrInvalidDenom      = errorsmod.Register(ModuleName, 2, "invalid denomination")
	ErrInsufficientFunds = errorsmod.Register(ModuleName, 3, "insufficient funds")
	ErrInvalidOwner      = errorsmod.Register(ModuleName, 4, "ivalid owner")
	ErrInvalidRate       = errorsmod.Register(ModuleName, 5, "ivalid exchange rate")
	ErrNonExchangeable   = errorsmod.Register(ModuleName, 6, "tokens are not exchangeable")
)

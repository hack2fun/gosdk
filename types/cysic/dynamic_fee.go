// Copyright 2024 Cysic Labs
// This file is part of Cysic Labs' Cysicmint library.
//
// The Cysicmint library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Cysicmint library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Cysicmint library. If not, see https://github.com/cysic-labs/cysic-network/blob/main/LICENSE
package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// HasDynamicFeeExtensionOption returns true if the tx implements the `ExtensionOptionDynamicFeeTx` extension option.
func HasDynamicFeeExtensionOption(any *codectypes.Any) bool {
	_, ok := any.GetCachedValue().(*ExtensionOptionDynamicFeeTx)
	return ok
}

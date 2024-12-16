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
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ sdk.Msg = &MsgDelegate{}
)

// GetSigners returns the expected signers for a MsgDelegate message.
func (m MsgDelegate) GetSigners() []sdk.AccAddress {
	sender := common.HexToAddress(m.Worker)
	signer := sdk.AccAddress(sender.Bytes())
	return []sdk.AccAddress{signer}
}

// ValidateBasic does a sanity check of the provided data
func (m *MsgDelegate) ValidateBasic() error {
	if !common.IsHexAddress(m.Worker) {
		return fmt.Errorf("worker cannot be empty")
	}
	if len(m.Validator) == 0 {
		return fmt.Errorf("validator cannot be empty")
	}

	if len(m.Token) == 0 {
		return fmt.Errorf("token cannot be empty")
	}

	if len(m.Amount) == 0 {
		return fmt.Errorf("amount cannot be empty")
	}

	return nil
}

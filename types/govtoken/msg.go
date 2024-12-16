package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgMint{}
var _ sdk.Msg = &MsgBurn{}
var _ sdk.Msg = &MsgChangeOwner{}
var _ sdk.Msg = &MsgExchangeToGovToken{}
var _ sdk.Msg = &MsgExchangeToPlatformToken{}
var _ sdk.Msg = &MsgSetExchangeRate{}
var _ sdk.Msg = &MsgStakeAsValidator{}
var _ sdk.Msg = &MsgDelegateToValidator{}

// MsgMint implements sdk.Msg
func (msg *MsgMint) Route() string { return RouterKey }

func (msg *MsgMint) Type() string { return "Mint" }

func (msg *MsgMint) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

func (msg *MsgMint) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMint) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}
	if msg.Amount.IsNegative() || msg.Amount.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount cannot be zero or negative")
	}
	return nil
}

// MsgBurn implements sdk.Msg
func (msg *MsgBurn) Route() string { return RouterKey }

func (msg *MsgBurn) Type() string { return "Burn" }

func (msg *MsgBurn) GetSigners() []sdk.AccAddress {
	burner, err := sdk.AccAddressFromBech32(msg.Burner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{burner}
}

func (msg *MsgBurn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBurn) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Burner); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid burner address (%s)", err)
	}
	if msg.Amount.IsNegative() || msg.Amount.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount cannot be zero or negative")
	}
	return nil
}

// MsgChangeOwner implements sdk.Msg
func (msg *MsgChangeOwner) Route() string { return RouterKey }

func (msg *MsgChangeOwner) Type() string { return "ChangeOwner" }

func (msg *MsgChangeOwner) GetSigners() []sdk.AccAddress {
	oldOwner, err := sdk.AccAddressFromBech32(msg.OldOwner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{oldOwner}
}

func (msg *MsgChangeOwner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgChangeOwner) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.OldOwner); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid old owner address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.NewOwner); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid new owner address (%s)", err)
	}
	return nil
}

// MsgExchangeToGovToken implements sdk.Msg
func (msg *MsgExchangeToGovToken) Route() string { return RouterKey }

func (msg *MsgExchangeToGovToken) Type() string { return "ExchangeToGovToken" }

func (msg *MsgExchangeToGovToken) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgExchangeToGovToken) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgExchangeToGovToken) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if msg.Amount.IsNegative() || msg.Amount.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount cannot be zero or negative")
	}
	return nil
}

// MsgExchangeToPlatformToken implements sdk.Msg
func (msg *MsgExchangeToPlatformToken) Route() string { return RouterKey }

func (msg *MsgExchangeToPlatformToken) Type() string { return "ExchangeToPlatformToken" }

func (msg *MsgExchangeToPlatformToken) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgExchangeToPlatformToken) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgExchangeToPlatformToken) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if msg.Amount.IsNegative() || msg.Amount.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount cannot be zero or negative")
	}
	return nil
}

// MsgSetExchangeRate implements sdk.Msg
func (msg *MsgSetExchangeRate) Route() string { return RouterKey }

func (msg *MsgSetExchangeRate) Type() string { return "SetExchangeRate" }

func (msg *MsgSetExchangeRate) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

func (msg *MsgSetExchangeRate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetExchangeRate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	if msg.Rate == 0 {
		return errorsmod.Wrap(ErrInvalidRate, "exchange rate must be positive")
	}
	return nil
}

// MsgStakeAsValidator implements sdk.Msg
func (msg *MsgStakeAsValidator) Route() string { return RouterKey }

func (msg *MsgStakeAsValidator) Type() string { return "StakeAsValidator" }

func (msg *MsgStakeAsValidator) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgStakeAsValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgStakeAsValidator) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if msg.Amount.IsNegative() || msg.Amount.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount cannot be zero or negative")
	}
	if msg.ValidatorName == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "validator name cannot be empty")
	}
	if msg.ValidatorDescription == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "validator description cannot be empty")
	}
	if msg.CommissionRate == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "commission rate cannot be empty")
	}
	if msg.MaxCommissionRate == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "max commission rate cannot be empty")
	}
	if msg.MaxChangeCommissionRate == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "max change commission rate cannot be empty")
	}
	if msg.MinSelfDelegation.IsNegative() || msg.MinSelfDelegation.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "minimum self delegation cannot be zero or negative")
	}
	if msg.ValidatorPubkey == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "validator pubkey cannot be empty")
	}
	return nil
}

// MsgDelegateToValidator implements sdk.Msg
func (msg *MsgDelegateToValidator) Route() string { return RouterKey }

func (msg *MsgDelegateToValidator) Type() string { return "DelegateToValidator" }

func (msg *MsgDelegateToValidator) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgDelegateToValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDelegateToValidator) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if msg.Amount.IsNegative() || msg.Amount.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount cannot be zero or negative")
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}
	return nil
}

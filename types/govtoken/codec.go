package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	amino = codec.NewLegacyAmino()
	// ModuleCdc references the global cysic  module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding.
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	// AminoCdc is a amino codec created to support amino JSON compatible msgs.
	AminoCdc = codec.NewAminoCodec(amino)
)

const (
	// Amino names
	paramsName                  = "cysicmint/govtoken/Params"
	setExchangeRateName         = "cysicmint/govtoken/MsgSetExchangeRate"
	exchangeToPlatformTokenName = "cysicmint/govtoken/MsgExchangeToPlatformToken"
	exchangeToGovTokenName      = "cysicmint/govtoken/MsgExchangeToGovToken"
	changeOwnerName             = "cysicmint/govtoken/MsgChangeOwner"
	burnName                    = "cysicmint/govtoken/MsgBurn"
	mintName                    = "cysicmint/govtoken/MsgMint"
)

// RegisterInterfaces registers the x/cysic interfaces types with the interface registry
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMint{},
		&MsgBurn{},
		&MsgChangeOwner{},
		&MsgExchangeToGovToken{},
		&MsgExchangeToPlatformToken{},
		&MsgSetExchangeRate{},
	)

	registry.RegisterImplementations(
		(*paramtypes.ParamSet)(nil),
		&Params{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterLegacyAminoCodec required
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgMint{}, mintName, nil)
	cdc.RegisterConcrete(&MsgBurn{}, burnName, nil)
	cdc.RegisterConcrete(&MsgChangeOwner{}, changeOwnerName, nil)
	cdc.RegisterConcrete(&MsgExchangeToGovToken{}, exchangeToGovTokenName, nil)
	cdc.RegisterConcrete(&MsgExchangeToPlatformToken{}, exchangeToPlatformTokenName, nil)
	cdc.RegisterConcrete(&MsgSetExchangeRate{}, setExchangeRateName, nil)
	cdc.RegisterConcrete(&Params{}, paramsName, nil)
}

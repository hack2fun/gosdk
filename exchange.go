package gosdk

import (
	"fmt"
	"log"

	govTokenTypes "github.com/cysic-tech/gosdk/types/govtoken"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ExchangeToCGT exchanges tokens to governance tokens.
//
// @param signer the Signer instance used to sign the transaction
// @param exchangeDetail the exchange details
// @return the transaction hash as a string, or an error if the exchange fails
func (s *Server) ExchangeToCGT(signer Signer, exchangeDetail *govTokenTypes.MsgExchangeToGovToken) (string, error) {
	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return "", err
	}

	if exchangeDetail == nil {
		return "", fmt.Errorf("exchangeDetail is nil")
	}

	if err := exchangeDetail.ValidateBasic(); err != nil {
		return "", err
	}

	txHash, err := s.buildAndBroadcastCosmosTx(signer, []sdk.Msg{exchangeDetail})
	if err != nil {
		return "", fmt.Errorf("create token failed, err: %v", err)
	}

	return txHash, nil
}

// ExchangeToCYS exchanges tokens to platform tokens.
//
// @param signer the Signer instance used to sign the transaction
// @param exchangeDetail the exchange details
// @return the transaction hash as a string, or an error if the exchange fails
func (s *Server) ExchangeToCYS(signer Signer, exchangeDetail *govTokenTypes.MsgExchangeToPlatformToken) (string, error) {
	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return "", err
	}

	if exchangeDetail == nil {
		return "", fmt.Errorf("exchangeDetail is nil")
	}

	if err := exchangeDetail.ValidateBasic(); err != nil {
		return "", err
	}

	txHash, err := s.buildAndBroadcastCosmosTx(signer, []sdk.Msg{exchangeDetail})
	if err != nil {
		return "", fmt.Errorf("create token failed, err: %v", err)
	}

	return txHash, nil
}

package gosdk

import (
	"context"
	"log"

	delegatetypes "github.com/cysic-tech/gosdk/types/delegate"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// QueryDelegatorDelegations
func (s *Server) QueryDelegatorDelegations(delegatorAddress string) ([]sdk.Coin, error) {
	result := make([]sdk.Coin, 0)
	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return result, err
	}

	client := stakingtypes.NewQueryClient(s.Conn)

	targetAddr, err := ConvertToCosmosAddress(delegatorAddress)
	if err != nil {
		log.Printf("error when convert addr: %v to cosmosAddr, err: %v", delegatorAddress, err.Error())
		return result, err
	}

	req := &stakingtypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: targetAddr,
	}

	resp, err := client.DelegatorDelegations(context.Background(), req)
	if err != nil {
		log.Printf("could not query delegateReward: %v", err)
		return result, err
	}

	resultMap := make(map[string]sdk.Coin)
	for _, info := range resp.DelegationResponses {
		balance := info.Balance
		if old, exist := resultMap[balance.Denom]; !exist {
			resultMap[balance.Denom] = sdk.Coin{
				Denom:  balance.Denom,
				Amount: balance.Amount,
			}
		} else {
			resultMap[balance.Denom] = sdk.Coin{
				Denom:  balance.Denom,
				Amount: old.Amount.Add(balance.Amount),
			}
		}
	}

	for _, info := range resultMap {
		result = append(result, info)
	}

	return result, nil
}

// QueryDelegateReward
func (s *Server) QueryDelegateReward(delegatorAddress string) ([]sdk.Coin, error) {
	result := make([]sdk.Coin, 0)
	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return result, err
	}

	client := distributiontypes.NewQueryClient(s.Conn)

	targetAddr, err := ConvertToCosmosAddress(delegatorAddress)
	if err != nil {
		log.Printf("error when convert addr: %v to cosmosAddr, err: %v", delegatorAddress, err.Error())
		return result, err
	}

	req := &distributiontypes.QueryDelegationTotalRewardsRequest{
		DelegatorAddress: targetAddr,
	}

	resp, err := client.DelegationTotalRewards(context.Background(), req)
	if err != nil {
		log.Printf("could not query delegateReward: %v", err)
		return result, err
	}

	resultMap := make(map[string]sdk.Coin)
	for _, balance := range resp.Total {
		if old, exist := resultMap[balance.Denom]; !exist {
			resultMap[balance.Denom] = sdk.Coin{
				Denom:  balance.Denom,
				Amount: balance.Amount.RoundInt(),
			}
		} else {
			resultMap[balance.Denom] = sdk.Coin{
				Denom:  balance.Denom,
				Amount: old.Amount.Add(balance.Amount.RoundInt()),
			}
		}
	}

	for _, info := range resultMap {
		result = append(result, info)
	}

	return result, nil
}

// WithdrawDelegatorReward
func (s *Server) WithdrawDelegatorReward(signer Signer, validatorAddress string) (string, error) {
	delegatorAddr := signer.CosmosAddr.String()

	msg := &distributiontypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddress,
	}

	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return "", err
	}

	if err := msg.ValidateBasic(); err != nil {
		return "", err
	}

	return s.buildAndBroadcastCosmosTx(signer, []sdk.Msg{msg})
}

// DelegateVeToken
func (s *Server) DelegateVeToken(signer Signer, validatorAddress string, coin string, amount math.Int) (string, error) {
	msg := &delegatetypes.MsgDelegate{
		Worker:    signer.EthAddr.String(),
		Validator: validatorAddress,
		Token:     coin,
		Amount:    amount.String(),
	}

	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return "", err
	}

	if err := msg.ValidateBasic(); err != nil {
		return "", err
	}

	return s.buildAndBroadcastCosmosTx(signer, []sdk.Msg{msg})
}

// DelegateCGT
func (s *Server) DelegateCGT(signer Signer, validatorAddress string, amount math.Int) (string, error) {
	msg := &stakingtypes.MsgDelegate{
		DelegatorAddress: signer.CosmosAddr.String(),
		ValidatorAddress: validatorAddress,
		Amount: sdk.Coin{
			Denom:  GovToken,
			Amount: amount,
		},
	}

	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return "", err
	}

	if err := msg.ValidateBasic(); err != nil {
		return "", err
	}

	return s.buildAndBroadcastCosmosTx(signer, []sdk.Msg{msg})
}

// UnDelegateCGT
func (s *Server) UnDelegateCGT(signer Signer, validatorAddress string, amount math.Int) (string, error) {
	msg := &stakingtypes.MsgUndelegate{
		DelegatorAddress: signer.CosmosAddr.String(),
		ValidatorAddress: validatorAddress,
		Amount: sdk.Coin{
			Denom:  GovToken,
			Amount: amount,
		},
	}

	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return "", err
	}

	if err := msg.ValidateBasic(); err != nil {
		return "", err
	}

	return s.buildAndBroadcastCosmosTx(signer, []sdk.Msg{msg})
}

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

// QueryDelegatorDelegations queries the delegations of a delegator.
//
// @param delegatorAddress the address of the delegator
// @return a list of coins representing the delegations, or an error if the query fails
func (s *Server) QueryDelegatorDelegations(delegatorAddress string) ([]sdk.Coin, error) {
	result := make([]sdk.Coin, 0)
	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return result, err
	}

	client := stakingtypes.NewQueryClient(s.Conn)

	targetAddr, err := ConvertToCysicAddress(delegatorAddress)
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

// QueryDelegateReward queries the rewards for a delegator.
//
// @param delegatorAddress the address of the delegator
// @return a list of coins representing the rewards, or an error if the query fails
func (s *Server) QueryDelegateReward(delegatorAddress string) ([]sdk.Coin, error) {
	result := make([]sdk.Coin, 0)
	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return result, err
	}

	client := distributiontypes.NewQueryClient(s.Conn)

	targetAddr, err := ConvertToCysicAddress(delegatorAddress)
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

// WithdrawDelegatorReward withdraws rewards for a delegator.
//
// @param signer the Signer instance used to sign the transaction
// @param validatorAddress the address of the validator
// @return the transaction hash as a string, or an error if the withdrawal fails
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

// DelegateVeToken delegates veTokens to a validator.
//
// @param signer the Signer instance used to sign the transaction
// @param validatorAddress the address of the validator
// @param coin the token to delegate
// @param amount the amount to delegate
// @return the transaction hash as a string, or an error if the delegation fails
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

// DelegateCGT delegates CGT tokens to a validator.
//
// @param signer the Signer instance used to sign the transaction
// @param validatorAddress the address of the validator
// @param amount the amount to delegate
// @return the transaction hash as a string, or an error if the delegation fails
func (s *Server) DelegateCGT(signer Signer, validatorAddress string, amount math.Int) (string, error) {
	msg := &stakingtypes.MsgDelegate{
		DelegatorAddress: signer.CosmosAddr.String(),
		ValidatorAddress: validatorAddress,
		Amount: sdk.Coin{
			Denom:  CGTToken,
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

// UnDelegateCGT undelegates CGT tokens from a validator.
//
// @param signer the Signer instance used to sign the transaction
// @param validatorAddress the address of the validator
// @param amount the amount to undelegate
// @return the transaction hash as a string, or an error if the undelegation fails
func (s *Server) UnDelegateCGT(signer Signer, validatorAddress string, amount math.Int) (string, error) {
	msg := &stakingtypes.MsgUndelegate{
		DelegatorAddress: signer.CosmosAddr.String(),
		ValidatorAddress: validatorAddress,
		Amount: sdk.Coin{
			Denom:  CGTToken,
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

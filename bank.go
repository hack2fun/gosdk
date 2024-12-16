package gosdk

import (
	"context"
	"fmt"
	"log"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/shopspring/decimal"
)

func (s *Server) GetBalanceList(address string) (sdk.Coins, error) {
	result := sdk.Coins{}
	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return result, err
	}

	client := banktypes.NewQueryClient(s.Conn)

	targetAddr, err := ConvertToCosmosAddress(address)
	if err != nil {
		log.Printf("error when convert addr: %v to cosmosAddr, err: %v", address, err.Error())
		return result, err
	}
	req := &banktypes.QueryAllBalancesRequest{Address: targetAddr}

	// 发送查询请求
	resp, err := client.AllBalances(context.Background(), req)
	if err != nil {
		log.Printf("could not query balances: %v", err)
		return result, err
	}

	return resp.Balances, nil
}

func (s *Server) GetBalance(address string, coin string) (string, error) {
	result := "0"
	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return "", err
	}

	client := banktypes.NewQueryClient(s.Conn)

	targetAddr, err := ConvertToCosmosAddress(address)
	if err != nil {
		log.Printf("error when convert addr: %v to cosmosAddr, err: %v", address, err.Error())
		return "", err
	}
	req := &banktypes.QueryAllBalancesRequest{Address: targetAddr}

	resp, err := client.AllBalances(context.Background(), req)
	if err != nil {
		log.Printf("could not query balances: %v", err)
		return "", err
	}

	exist, amount := resp.Balances.Find(coin)
	if exist {
		result = decimal.NewFromBigInt(amount.Amount.BigInt(), -18).String()
	}

	return result, nil
}

func (s *Server) Send(signer Signer, toAddrStr string, coin string, amount sdkmath.Int) (string, error) {
	targetAddr, err := ConvertToCosmosAddress(toAddrStr)
	if err != nil {
		log.Printf("error when convert addr: %v to cosmosAddr, err: %v", toAddrStr, err.Error())
		return "", err
	}

	toAddr, err := sdk.AccAddressFromBech32(targetAddr)
	if err != nil {
		log.Printf("error when convert toAddr: %v, afterConvertAddr: %v, err: %v\n", toAddrStr, targetAddr, err.Error())
		return "", err
	}

	sendMsg := banktypes.NewMsgSend(signer.CosmosAddr, toAddr, sdk.Coins{sdk.NewCoin(coin, amount)})
	return s.buildAndBroadcastCosmosTx(signer, []sdk.Msg{sendMsg})
}

func (s *Server) MultiSend(signer Signer, toAddrList []string, coin string, amount sdkmath.Int) (string, error) {
	coins := sdk.NewCoins(sdk.NewCoin(coin, amount.Mul(sdkmath.NewInt(int64(len(toAddrList))))))
	in := []banktypes.Input{banktypes.NewInput(signer.CosmosAddr, coins)}
	var out []banktypes.Output
	for _, toAddr := range toAddrList {
		toAddrCosmos, err := ConvertToCosmosAddress(toAddr)
		if err != nil {
			log.Printf("error when convert to addr: %v, err: %v", toAddr, err.Error())
			return "", err
		}
		to, err := sdk.AccAddressFromBech32(toAddrCosmos)
		if err != nil {
			log.Printf("error when conert addr to accAddr, addr: %v, err: %v", toAddrCosmos, err.Error())
			return "", err
		}

		sendCoins := sdk.NewCoins(sdk.NewCoin(coin, amount))
		out = append(out, banktypes.NewOutput(to, sendCoins))
	}

	msg := banktypes.NewMsgMultiSend(in, out)
	return s.buildAndBroadcastCosmosTx(signer, []sdk.Msg{msg})
}

func (s *Server) MultiSendWithDiffAmount(signer Signer, toAddrList []string, coinList []string, amountList []sdkmath.Int) (string, error) {
	if len(toAddrList) != len(coinList) || len(coinList) != len(amountList) {
		return "", fmt.Errorf("params length not equal, len(toAddr): %v, len(coinList): %v, len(amountList): %v",
			len(toAddrList), len(coinList), len(amountList))
	}

	amountMap := map[string]sdkmath.Int{}
	for i := 0; i < len(toAddrList); i++ {
		coin := coinList[i]
		amount := amountList[i]
		if _, exist := amountMap[coin]; !exist {
			amountMap[coin] = sdkmath.NewInt(0)
		}

		amountMap[coin] = amountMap[coin].Add(amount)
	}

	coins := sdk.NewCoins()
	for coin, amount := range amountMap {
		newCoin := sdk.NewCoin(coin, amount)
		coins = coins.Add(newCoin)
	}
	in := []banktypes.Input{banktypes.NewInput(signer.CosmosAddr, coins)}
	var out []banktypes.Output
	for i, toAddr := range toAddrList {
		toAddrCosmos, err := ConvertToCosmosAddress(toAddr)
		if err != nil {
			log.Printf("error when convert to addr: %v, err: %v", toAddr, err.Error())
			return "", err
		}
		to, err := sdk.AccAddressFromBech32(toAddrCosmos)
		if err != nil {
			log.Printf("error when conert addr to accAddr, addr: %v, err: %v", toAddrCosmos, err.Error())
			return "", err
		}

		sendCoins := sdk.NewCoins(sdk.NewCoin(coinList[i], amountList[i]))
		out = append(out, banktypes.NewOutput(to, sendCoins))
	}

	msg := banktypes.NewMsgMultiSend(in, out)
	return s.buildAndBroadcastCosmosTx(signer, []sdk.Msg{msg})
}

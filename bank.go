package gosdk

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"log"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// GetBalance retrieves the balance of a specific coin for an address.
//
// @param address the address to query the balance for
// @param coin the coin denomination to retrieve
// @return the balance as a string, or an error if the retrieval fails
func (s *Server) GetBalance(address string, coin string) (string, error) {
	result := "0"
	coins, err := s.GetBalanceList(address)
	if err != nil {
		log.Printf("error when GetBalanceList by addr: %v, err: %v", address, err.Error())
		return result, err
	}

	exist, amount := coins.Find(coin)
	if exist {
		result = decimal.NewFromBigInt(amount.Amount.BigInt(), -18).String()
	}

	return result, nil
}

// GetBalanceList retrieves the full list of balances for an address.
//
// @param address the address to query the balances for
// @return a list of coins representing the balances, or an error if the retrieval fails
func (s *Server) GetBalanceList(address string) (sdk.Coins, error) {
	result := sdk.Coins{}
	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return result, err
	}

	client := banktypes.NewQueryClient(s.Conn)

	targetAddr, err := ConvertToCysicAddress(address)
	if err != nil {
		log.Printf("error when convert addr: %v to cosmosAddr, err: %v", address, err.Error())
		return result, err
	}
	req := &banktypes.QueryAllBalancesRequest{Address: targetAddr}

	resp, err := client.AllBalances(context.Background(), req)
	if err != nil {
		log.Printf("could not query balances: %v", err)
		return result, err
	}

	return resp.Balances, nil
}

// Send facilitates the sending of coins from one address to another.
//
// @param signer the Signer instance used to sign the transaction
// @param toAddrStr the address to send coins to
// @param coin the coin denomination to send
// @param amount the amount of coins to send
// @return the transaction hash as a string, or an error if the send operation fails
func (s *Server) Send(signer Signer, toAddrStr string, coin string, amount sdkmath.Int) (string, error) {
	targetAddr, err := ConvertToCysicAddress(toAddrStr)
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

// MultiSend facilitates the sending of coins to multiple addresses in a single transaction.
//
// @param signer the Signer instance used to sign the transaction
// @param toAddrList the list of addresses to send coins to
// @param coin the coin denomination to send
// @param amount the amount of coins to send to each address
// @return the transaction hash as a string, or an error if the send operation fails
func (s *Server) MultiSend(signer Signer, toAddrList []string, coin string, amount sdkmath.Int) (string, error) {
	coins := sdk.NewCoins(sdk.NewCoin(coin, amount.Mul(sdkmath.NewInt(int64(len(toAddrList))))))
	in := []banktypes.Input{banktypes.NewInput(signer.CosmosAddr, coins)}
	var out []banktypes.Output
	for _, toAddr := range toAddrList {
		toAddrCosmos, err := ConvertToCysicAddress(toAddr)
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

// MultiSendWithDiffAmount facilitates the sending of different amounts of coins to multiple addresses in a single transaction.
//
// @param signer the Signer instance used to sign the transaction
// @param toAddrList the list of addresses to send coins to
// @param coinList the list of coin denominations to send
// @param amountList the list of amounts to send for each coin denomination
// @return the transaction hash as a string, or an error if the send operation fails
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
		toAddrCosmos, err := ConvertToCysicAddress(toAddr)
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

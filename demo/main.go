package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	cysicSDK "github.com/cysic-tech/gosdk"
	govTokenTypes "github.com/cysic-tech/gosdk/types/govtoken"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	txTypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	defaultServer *cysicSDK.Server
	signer        *cysicSDK.Signer

	defaultEndpoint = ""
)

func init() {
	chainEndpoint := ""
	chainID := "cysicmint_9001-1" // testnet
	gasCoin := "CYS"
	gasPrice := int64(10)
	defaultEndpoint = chainEndpoint

	mnemonic := ""
	coinType := uint32(1)
	account := uint32(1)
	index := uint32(1)
	var err error
	defaultServer, err = cysicSDK.NewServerWithGRPCAndGasLimit(chainEndpoint, chainID, gasCoin, gasPrice, 30_000_000)
	if err != nil {
		log.Printf("error while creating server: %v", err.Error())
		return
	}

	hdPath := hd.CreateHDPath(coinType, account, index).String()
	signer, err = cysicSDK.NewSignerWithMnemonic(mnemonic, "", hdPath, "eth_secp256k1")
	if err != nil {
		log.Printf("error while creating sender: %v", err.Error())
		return
	}
}

func main() {
	validatorAddrList := getValidatorAddrList()
	fmt.Printf("signer: %v , cosmos: %v\n", signer.EthAddr.String(), signer.CosmosAddr.String())

	if false {
		getDelegatorDelegation(signer.EthAddr.String())
	}

	if false {
		delegateCGT(*signer, validatorAddrList[0], sdkmath.NewIntFromBigInt(decimal.New(100, 18).BigInt()))
	}

	if false {
		getReward(signer.EthAddr.String())
	}

	if false {
		getReward(signer.EthAddr.String())
		claimReward(*signer, validatorAddrList[0])
		getReward(signer.EthAddr.String())
	}

	if false {
		undelegateCGT(*signer, validatorAddrList[0], sdkmath.NewIntFromBigInt(decimal.New(100, 18).BigInt()))
	}

	if false {
		delegateVeToken(*signer, validatorAddrList[0], "veSCR", sdkmath.NewIntFromBigInt(decimal.New(1, 18).BigInt()))
	}

	if false {
		convertToCGT(*signer, sdkmath.NewIntFromBigInt(decimal.New(1, 18).BigInt()))
	}

	if false {
		convertToCYS(*signer, sdkmath.NewIntFromBigInt(decimal.New(1, 18).BigInt()))
	}
}

func printBalance(addr string) {
	fmt.Print("addr: ", addr, "balance: ")
	balanceList, _ := defaultServer.GetBalanceList(addr)
	for i, coin := range balanceList {
		fmt.Printf("%v%v", decimal.NewFromBigInt(coin.Amount.BigInt(), -18), coin.Denom)
		if i != len(balanceList)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println()
}

func printValidatorList() {
	fmt.Println("validator: ")
	validatorList, _, _ := defaultServer.GetValidatorList(0, 100)
	for _, info := range validatorList {
		fmt.Printf("%v(%v): %v\n", info.GetMoniker(), info.OperatorAddress, info.Tokens.String())
	}
}

func getValidatorAddrList() []string {
	result := make([]string, 0)
	validatorList, _, _ := defaultServer.GetValidatorList(0, 100)
	for _, info := range validatorList {
		result = append(result, info.OperatorAddress)
	}

	return result
}

func delegateCGT(signer cysicSDK.Signer, validatorAddr string, amount sdkmath.Int) {
	fmt.Print("before delegate CGT: ")
	printBalance(signer.EthAddr.String())
	printValidatorList()

	txHash, err := defaultServer.DelegateCGT(signer, validatorAddr, amount)
	if err != nil {
		fmt.Println("error when delegateCGT: ", err.Error())
		return
	}

	fmt.Printf("delegateCGT finish, txHash: %v\n", txHash)
	waitTxFinish(txHash)
	getTx(txHash)

	fmt.Print("after delegate CGT: ")
	printBalance(signer.EthAddr.String())
	printValidatorList()
}

func delegateVeToken(signer cysicSDK.Signer, validatorAddr string, coin string, amount sdkmath.Int) {
	fmt.Print("before delegate : ", coin)
	printBalance(signer.EthAddr.String())
	printValidatorList()

	txHash, err := defaultServer.DelegateVeToken(signer, validatorAddr, coin, amount)
	if err != nil {
		fmt.Println("error when delegateCGT: ", err.Error())
		return
	}

	fmt.Printf("delegateVeToken finish, txHash: %v\n", txHash)
	waitTxFinish(txHash)
	getTx(txHash)

	fmt.Print("after delegate : ", coin)
	printBalance(signer.EthAddr.String())
	printValidatorList()
}

func undelegateCGT(signer cysicSDK.Signer, validatorAddr string, amount sdkmath.Int) {
	fmt.Print("before undelegate to CGT: ")
	getDelegatorDelegation(signer.EthAddr.String())
	printBalance(signer.EthAddr.String())
	printValidatorList()

	txHash, err := defaultServer.UnDelegateCGT(signer, validatorAddr, amount)
	if err != nil {
		fmt.Println("error when undelegateCGT: ", err.Error())
		return
	}
	fmt.Printf("undelegateCGT finish, txHash: %v\n", txHash)
	waitTxFinish(txHash)
	getTx(txHash)

	fmt.Print("after undelegate to CGT: ")
	getDelegatorDelegation(signer.EthAddr.String())
	printBalance(signer.EthAddr.String())
	printValidatorList()
}

func getDelegatorDelegation(addr string) {
	fmt.Print("current delegation: ")
	fmt.Println(defaultServer.QueryDelegatorDelegations(addr))
}

func getReward(addr string) {
	fmt.Print("current reward: ")
	fmt.Println(defaultServer.QueryDelegateReward(addr))
}

func claimReward(signer cysicSDK.Signer, validatorAddr string) {
	fmt.Print("before claimReward: ")
	printBalance(signer.EthAddr.String())

	txHash, err := defaultServer.WithdrawDelegatorReward(signer, validatorAddr)
	if err != nil {
		fmt.Println("error when WithdrawDelegatorReward: ", err.Error())
		return
	}
	waitTxFinish(txHash)
	getTx(txHash)

	fmt.Print("after claimReward: ")
	printBalance(signer.EthAddr.String())
}

func convertToCGT(signer cysicSDK.Signer, amount sdkmath.Int) {
	fmt.Print("before swap to CGT: ")
	printBalance(signer.EthAddr.String())
	txHash, err := defaultServer.ExchangeToGovToken(signer, &govTokenTypes.MsgExchangeToGovToken{
		Sender: signer.CosmosAddr.String(),
		Amount: amount,
	})
	if err != nil {
		fmt.Println("error when exchange to CGT: ", err.Error())
		return
	}

	fmt.Printf("convertToCGT finish, txHash: %v\n", txHash)
	waitTxFinish(txHash)
	fmt.Print("after swap to CGT: ")
	printBalance(signer.EthAddr.String())
}

func convertToCYS(signer cysicSDK.Signer, amount sdkmath.Int) {
	fmt.Print("before swap to CYS: ")
	printBalance(signer.EthAddr.String())
	txHash, err := defaultServer.ExchangeToPlatformToken(signer, &govTokenTypes.MsgExchangeToPlatformToken{
		Sender: signer.CosmosAddr.String(),
		Amount: amount,
	})
	if err != nil {
		fmt.Println("error when exchange to CYS: ", err.Error())
		return
	}

	fmt.Printf("convertToCYS finish, txHash: %v\n", txHash)
	waitTxFinish(txHash)
	fmt.Print("after swap to CYS: ")
	printBalance(signer.EthAddr.String())
}

func waitTxFinish(txHash string) {
	conn, err := grpc.Dial(defaultEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error when new grpc client: %s", err.Error())
		return
	}

	txClient := txTypes.NewServiceClient(conn)
	defer conn.Close()

	for i := 0; i < 15; i++ {
		_, err = txClient.GetTx(context.Background(), &txTypes.GetTxRequest{Hash: txHash})
		if err != nil {
			if strings.Index(err.Error(), "tx not found") >= 0 {
				log.Printf("start wait tx: %v packed", txHash)
				<-time.NewTimer(time.Second).C
				continue
			}

			log.Printf("error when get tx: %v, err: %v", txHash, err.Error())
			return
		}
	}
}

func getTx(txHash string) {
	conn, err := grpc.Dial(defaultServer.EndPoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(fmt.Sprintf("error when new grpc client: %s\n", err.Error()))
		return
	}
	defer conn.Close()
	txClient := txTypes.NewServiceClient(conn)

	txResp, err := txClient.GetTx(context.Background(), &txTypes.GetTxRequest{Hash: txHash})
	if err != nil {
		fmt.Println(fmt.Sprintf("error when get tx: %v, err: %v", txHash, err.Error()))
		return
	}

	fmt.Println("tx: ", txResp.TxResponse.TxHash, " in", txResp.TxResponse.Height,
		", fee payer:", txResp.Tx.FeePayer().String(), ", seq:", txResp.Tx.AuthInfo.SignerInfos[0].Sequence, ", is:", txResp.Tx.Body.Messages[0].TypeUrl)

	// not packed
	if txResp.TxResponse.Height == 0 {
		fmt.Println(txHash, " not packed")
		return
	}

	// tx failed
	if txResp.TxResponse.Code != 0 {
		fmt.Println(txHash, " failed, ", txResp.TxResponse.RawLog)
		return
	}

	fmt.Println(txHash, "logs ", txResp.TxResponse.RawLog)
	return
}

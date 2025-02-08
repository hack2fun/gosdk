package gosdk

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	sdkClient "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authSigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cysicTypes "github.com/hack2fun/gosdk/types/cysic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GetAccount retrieves account information from the chain for a given signer.
//
// @param signer the Signer instance to retrieve the account for
// @return the account information as an EthAccount struct, or an error if retrieval fails
func (s *Server) GetAccount(signer Signer) (*cysicTypes.EthAccount, error) {
	return s.GetAccountByAddr(signer.CosmosAddr.String())
}

// GetAccountByAddr retrieves account information from the chain for a given address.
//
// @param addr the address to retrieve the account information for
// @return the account information as an EthAccount struct, or an error if retrieval fails
func (s *Server) GetAccountByAddr(addr string) (*cysicTypes.EthAccount, error) {
	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return nil, err
	}

	cosmosAddr, err := ConvertToCysicAddress(addr)
	if err != nil {
		log.Printf("error when convert addr: %v to cosmos addr, err: %v", addr, err.Error())
		return nil, err
	}

	client := authTypes.NewQueryClient(s.Conn)
	req := authTypes.QueryAccountRequest{Address: cosmosAddr}
	res, err := client.Account(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	temp := &cysicTypes.EthAccount{}
	err = temp.XXX_Unmarshal(res.Account.GetValue())
	if err != nil {
		log.Printf("error when new grpc client: %s\n", err.Error())
		return nil, err
	}

	return temp, nil
}

// BroadcastTx broadcasts a signed transaction to the network.
//
// @param txBytes the signed transaction bytes
// @return the transaction response, or an error if broadcasting fails
func (s *Server) BroadcastTx(txBytes []byte) (*sdk.TxResponse, error) {
	conn, err := grpc.Dial(s.EndPoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error when new grpc client: %s\n", err.Error())
		return nil, err
	}
	defer conn.Close()

	client := sdkTx.NewServiceClient(conn)

	res, err := client.BroadcastTx(context.Background(), &sdkTx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    sdkTx.BroadcastMode_BROADCAST_MODE_SYNC,
	})
	if errRes := sdkClient.CheckTendermintError(err, txBytes); errRes != nil {
		return errRes, nil
	}
	if err != nil {
		return nil, err
	}

	return res.TxResponse, err
}

func (s *Server) getAccountNumberAndSequenceOnChain(address sdk.AccAddress) (exist bool, accNumber uint64, sequence uint64, err error) {
	temp, err := s.GetAccountByAddr(address.String())
	if err != nil {
		log.Printf("error when GetAccountByAddr: %v, err: %v", address.String(), err.Error())
		return false, 0, 0, err
	}

	accNumber = temp.AccountNumber
	sequence = temp.Sequence
	return true, accNumber, sequence, nil
}

// buildAndBroadcastCosmosTx builds a Cosmos transaction with the provided signer and messages, then broadcasts it to the network.
//
// @param signer the Signer instance used to sign the transaction
// @param msgList list of messages to include in the transaction
// @return the transaction hash as a string, or an error if the transaction fails
func (s *Server) buildAndBroadcastCosmosTx(signer Signer, msgList []sdk.Msg) (string, error) {
	for _, msg := range msgList {
		if err := msg.ValidateBasic(); err != nil {
			log.Printf("error when validate basic for msg: %v, err: %v\n", msg, err.Error())
			return "", err
		}
	}

	signerPriv := signer.privateKey
	accAddr := signer.CosmosAddr
	signerPubKey := signerPriv.PubKey()

	exist, accNumber, sequence, err := s.getAccountNumberAndSequenceOnChain(accAddr)
	if err != nil {
		log.Printf("error when get accInfo on chain, addr: %v, err: %v\n", accAddr.String(), err.Error())
		return "", err
	}
	if !exist {
		log.Printf("error account: %v not exist", accAddr.String())
		return "", fmt.Errorf("account %v not exist", accAddr.String())
	}

	if signer.Nonce != 0 && signer.Nonce > sequence {
		sequence = signer.Nonce
	}
	txBuilder, bytesToSign, err := s.GetBytesToSign(signer, accNumber, sequence, msgList)
	if err != nil {
		log.Printf("error when get wait sign tx, err: %v\n", err.Error())
		return "", err
	}

	sigBytes, err := signerPriv.Sign(bytesToSign)
	if err != nil {
		log.Printf("error when sign msg, err: %v\n", err.Error())
		return "", err
	}

	// Construct the SignatureV2 struct
	sigData := signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
		Signature: sigBytes,
	}
	sig := signing.SignatureV2{
		PubKey:   signerPubKey,
		Data:     &sigData,
		Sequence: sequence,
	}

	err = txBuilder.SetSignatures(sig)
	if err != nil {
		log.Printf("error when set signed bytes to tx, err: %v\n", err.Error())
		return "", err
	}

	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		log.Printf("error when get signed tx bytes, err: %v\n", err.Error())
		return "", err
	}

	resp, err := s.BroadcastTx(txBytes)
	if err != nil {
		log.Printf("error when broadcast tx, err: %v\n", err.Error())
		return "", err
	}
	if resp.Code != 0 {
		log.Printf("resp code not zero, log: %v\n", resp.RawLog)
		return "", fmt.Errorf(resp.RawLog)
	}

	return resp.TxHash, nil
}

// waitTxPacked waits for a transaction to be packed into a block.
//
// @param txHash the hash of the transaction to wait for
func (s *Server) waitTxPacked(txHash string) {
	conn, err := grpc.Dial(s.EndPoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error when new grpc client: %s\n", err.Error())
		return
	}
	defer conn.Close()

	txClient := sdkTx.NewServiceClient(conn)
	defer conn.Close()
	defer time.Sleep(500 * time.Millisecond)

	for i := 0; i < 10; i++ {
		resp, err := txClient.GetTx(context.Background(), &sdkTx.GetTxRequest{Hash: txHash})
		if err != nil {
			if strings.Index(err.Error(), "tx not found") >= 0 {
				log.Printf("wait tx %v packed", txHash)
				<-time.NewTimer(time.Second).C
				log.Printf("wait finish, try again")
				continue
			}

			log.Printf("error when get tx: %v, err: %v", txHash, err.Error())
			return
		}

		if resp != nil && resp.TxResponse != nil && resp.TxResponse.Height != 0 {
			break
		}
	}
}

// GetBytesToSign generates the bytes to sign for a transaction.
//
// @param signer the Signer instance used to sign the transaction
// @param accNumber the account number of the signer
// @param sequence the sequence number of the signer
// @param msgList list of messages to include in the transaction
// @return the transaction builder, bytes to sign, or an error if generation fails
func (s *Server) GetBytesToSign(signer Signer, accNumber, sequence uint64, msgList []sdk.Msg) (sdkClient.TxBuilder, []byte, error) {
	for _, msg := range msgList {
		if err := msg.ValidateBasic(); err != nil {
			log.Printf("error when validate basic for msg: %v, err: %v\n", msg, err.Error())
			return nil, nil, err
		}
	}

	txBuilder := txConfig.NewTxBuilder()
	err := txBuilder.SetMsgs(msgList...)
	if err != nil {
		log.Printf("error when set msg, err: %s", err.Error())
		return nil, nil, err
	}

	fees := make(sdk.Coins, 1)
	fees[0] = sdk.NewCoin(
		s.GasCoin,
		sdk.NewDec(s.GasPrice*int64(s.GasLimit)).Ceil().RoundInt(),
	)

	txBuilder.SetFeeAmount(fees)
	txBuilder.SetGasLimit(s.GasLimit)
	txBuilder.SetFeePayer(signer.CosmosAddr)

	signerData := authSigning.SignerData{
		ChainID:       s.ChainID,
		AccountNumber: accNumber,
		Sequence:      sequence,
		PubKey:        signer.publicKey,
		Address:       signer.CosmosAddr.String(),
	}

	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   signer.publicKey,
		Data:     &sigData,
		Sequence: sequence,
	}
	if err := txBuilder.SetSignatures(sig); err != nil {
		log.Printf("error when set signatures, err: %s\n", err.Error())
		return nil, nil, err
	}

	bytesToSign, err := txConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
	if err != nil {
		log.Printf("error when get wait sign tx, err: %v\n", err.Error())
		return nil, nil, err
	}

	return txBuilder, bytesToSign, nil
}

// broadcastMsg broadcasts a single message as a transaction.
//
// @param signer the Signer instance used to sign the transaction
// @param msg the message to broadcast
// @return the transaction hash as a string, or an error if broadcasting fails
func (s *Server) broadcastMsg(signer Signer, msg sdk.Msg) (string, error) {
	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return "", err
	}

	if msg == nil {
		return "", fmt.Errorf("msg is nil")
	}

	if err := msg.ValidateBasic(); err != nil {
		return "", err
	}

	return s.buildAndBroadcastCosmosTx(signer, []sdk.Msg{msg})
}

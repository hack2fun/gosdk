package gosdk

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	
	cysicTypes "github.com/cysic-tech/gosdk/types/cysic"

	sdkClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authSigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	interfaceRegistry = codecTypes.NewInterfaceRegistry()
	cdc               = codec.NewProtoCodec(interfaceRegistry)
	txConfig          = tx.NewTxConfig(cdc, tx.DefaultSignModes)

	gasLimit = uint64(15_000_000)
	signMode = signing.SignMode_SIGN_MODE_DIRECT
)

type Server struct {
	EndPoint string
	Conn     *grpc.ClientConn
	ChainID  string
	GasCoin  string
	GasPrice int64
	GasLimit uint64
}

func NewServerWithGRPC(endPoint string, chainID string, gasCoin string, gasPrice int64) (*Server, error) {
	conn, err := grpc.Dial(endPoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error when new grpc client: %s\n", err.Error())
		return nil, err
	}

	return &Server{
		EndPoint: endPoint,
		ChainID:  chainID,
		GasPrice: gasPrice,
		GasCoin:  gasCoin,
		GasLimit: gasLimit,
		Conn:     conn,
	}, nil
}

func NewServerWithGRPCAndGasLimit(endPoint string, chainID string, gasCoin string, gasPrice int64, _gasLimit uint64) (*Server, error) {
	conn, err := grpc.Dial(endPoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error when new grpc client: %s\n", err.Error())
		return nil, err
	}

	return &Server{
		EndPoint: endPoint,
		ChainID:  chainID,
		GasPrice: gasPrice,
		GasCoin:  gasCoin,
		GasLimit: _gasLimit,
		Conn:     conn,
	}, nil
}

func (s *Server) KeepGrpcConn() error {
	if s.Conn != nil && s.Conn.GetState() != connectivity.Shutdown {
		return nil
	}

	conn, err := grpc.Dial(s.EndPoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error when new grpc client: %s\n", err.Error())
		return err
	}

	s.Conn = conn
	return nil
}

func (s *Server) Close() error {
	return s.Conn.Close()
}

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

func (s *Server) GetAccount(signer Signer) (*cysicTypes.EthAccount, error) {
	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return nil, err
	}

	// 创建 gRPC 客户端
	client := authTypes.NewQueryClient(s.Conn)
	req := authTypes.QueryAccountRequest{Address: signer.CosmosAddr.String()}
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

func (s *Server) GetAccountByAddr(addr string) (*cysicTypes.EthAccount, error) {
	if err := s.KeepGrpcConn(); err != nil {
		log.Printf("error when keep grpc conn, endpoint: %v, err: %v", s.EndPoint, err.Error())
		return nil, err
	}

	cosmosAddr, err := ConvertToCosmosAddress(addr)
	if err != nil {
		log.Printf("error when convert addr: %v to cosmos addr, err: %v", addr, err.Error())
		return nil, err
	}

	// 创建 gRPC 客户端
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

func (s *Server) getAccountNumberAndSequenceOnChain(address sdk.AccAddress) (exist bool, accNumber uint64, sequence uint64, err error) {
	conn, err := grpc.Dial(s.EndPoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error when new grpc client: %s\n", err.Error())
		return
	}
	defer conn.Close()

	authClient := authTypes.NewQueryClient(conn)
	accountReq := authTypes.QueryAccountRequest{Address: address.String()}
	accountRes, err := authClient.Account(context.Background(), &accountReq)
	if err != nil {
		exist = false
		return
	}

	temp := &cysicTypes.EthAccount{}
	err = temp.XXX_Unmarshal(accountRes.Account.GetValue())
	if err != nil {
		log.Printf("error when new grpc client: %s\n", err.Error())
		return
	}

	accNumber = temp.AccountNumber
	sequence = temp.Sequence
	return true, accNumber, sequence, nil
}

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

	return res.TxResponse, err
}

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

package gosdk

import (
	"log"

	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
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

// NewServerWithGRPC creates a new Server instance with a gRPC connection.
//
// @param endPoint the endpoint of the gRPC server
// @param chainID the chain ID of the blockchain
// @param gasCoin the coin to use for gas fees
// @param gasPrice the price of gas
// @return a new Server instance, or an error if the connection fails
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

// NewServerWithGRPCAndGasLimit creates a new Server instance with a gRPC connection and custom gas limit.
//
// @param endPoint the endpoint of the gRPC server
// @param chainID the chain ID of the blockchain
// @param gasCoin the coin to use for gas fees
// @param gasPrice the price of gas
// @param _gasLimit the custom gas limit
// @return a new Server instance, or an error if the connection fails
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

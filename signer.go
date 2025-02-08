package gosdk

import (
	"log"

	"github.com/hack2fun/gosdk/crypto/ethsecp256k1"
	etherminthd "github.com/hack2fun/gosdk/crypto/hd"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Signer struct {
	CosmosAddr sdk.AccAddress
	EthAddr    common.Address
	privateKey types.PrivKey
	publicKey  types.PubKey
	Nonce      uint64
}

// NewSignerWithPrivateKey creates a new Signer instance from a private key.
//
// @param bz the private key bytes
// @return a new Signer instance
func NewSignerWithPrivateKey(bz []byte) *Signer {
	bzArr := make([]byte, ethsecp256k1.PrivKeySize)
	copy(bzArr, bz)
	privKey := &ethsecp256k1.PrivKey{
		Key: bzArr,
	}
	cosmosAddr := sdk.AccAddress(privKey.PubKey().Address())
	ethAddr := common.BytesToAddress(cosmosAddr)

	return &Signer{
		CosmosAddr: cosmosAddr,
		EthAddr:    ethAddr,
		privateKey: privKey,
		publicKey:  privKey.PubKey(),
	}
}

// NewSignerWithMnemonic creates a new Signer instance from a mnemonic and HD path.
//
// @param mnemonic the mnemonic phrase
// @param passphrase the passphrase for the mnemonic
// @param hdPath the HD path to derive the key
// @param algo the signing algorithm to use
// @return a new Signer instance, or an error if derivation fails
func NewSignerWithMnemonic(mnemonic string, passphrase string, hdPath string, algo string) (*Signer, error) {
	kb := keyring.NewInMemory(simapp.MakeTestEncodingConfig().Codec, etherminthd.EthSecp256k1Option())
	keyringAlgoList, _ := kb.SupportedAlgorithms()
	signAlgo, err := keyring.NewSigningAlgoFromString(algo, keyringAlgoList)
	if err != nil {
		log.Printf("error when new signing algo from string, err: %v", err.Error())
		return nil, err
	}

	derivedPriv, err := signAlgo.Derive()(mnemonic, passphrase, hdPath)
	if err != nil {
		log.Printf("error when get derive private key, err: %v\n", err.Error())
		return nil, err
	}
	privKey := signAlgo.Generate()(derivedPriv)

	cosmosAddr := sdk.AccAddress(privKey.PubKey().Address())
	ethAddr := common.BytesToAddress(cosmosAddr)

	// log.Printf("addr: %v, privKey: %v", ethAddr, hex.EncodeToString(privKey.Bytes()))
	return &Signer{
		CosmosAddr: cosmosAddr,
		EthAddr:    ethAddr,
		privateKey: privKey,
		publicKey:  privKey.PubKey(),
	}, nil
}

// VerifyEthPersonalSignature verifies an Ethereum personal signature.
//
// @param address the address that signed the data
// @param data the data that was signed
// @param sig the signature to verify
// @return true if the signature is valid, false otherwise
func VerifyEthPersonalSignature(address string, data []byte, sig []byte) bool {
	sigHash, _ := accounts.TextAndHash(data)

	if sig[64] > 1 {
		sig[64] -= 27
	}

	sigPublicKey, err := crypto.SigToPub(sigHash, sig)
	if err != nil {
		log.Println("invalid signature crypto.SigToPub err: ", err)

		return false
	}

	recoverAddress := crypto.PubkeyToAddress(*sigPublicKey)
	return recoverAddress == common.HexToAddress(address)
}

// EthPersonalSign signs data using the Ethereum personal sign algorithm.
//
// @param s the Signer instance
// @param data the data to sign
// @return the signature, or an error if signing fails
func (s *Signer) EthPersonalSign(data []byte) ([]byte, error) {
	waitSignHash, _ := accounts.TextAndHash(data)

	priv := s.privateKey.(*ethsecp256k1.PrivKey)
	privKey, err := crypto.ToECDSA(priv.Key)
	if err != nil {
		log.Printf("error when convert private key, err: %v", err.Error())
		return nil, err
	}

	return crypto.Sign(waitSignHash, privKey)
}

// VerifyEthPersonalSignature verifies an Ethereum personal signature for the signer.
//
// @param s the Signer instance
// @param data the data that was signed
// @param sig the signature to verify
// @return true if the signature is valid, false otherwise
func (s *Signer) VerifyEthPersonalSignature(data []byte, sig []byte) bool {
	return VerifyEthPersonalSignature(s.EthAddr.String(), data, sig)
}

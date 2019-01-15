// Package dappauth provides the ability to check if an Ethereume address is an authorized signer for a signature generated via eth_sign rpc method.
// Supports both external wallets and contract wallets.
package dappauth

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/dapperlabs/dappauth/ERCs"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

var (
	_ERC1271MagicValue = [4]byte{32, 193, 59, 11} // 0x20c13b0b
)

// Authenticator is the instance that holds the ethclient.Client .
type Authenticator struct {
	cc  bind.ContractCaller
	ctx context.Context // Network context to support cancellation and timeouts (nil = no timeout)
}

// NewAuthenticator creates a new Authenticator .
func NewAuthenticator(ctx context.Context, cc bind.ContractCaller) *Authenticator {
	return &Authenticator{
		ctx: ctx,
		cc:  cc,
	}
}

// IsAuthorizedSigner implements the logic to check if an address is an authorized signer for a signature and challenge.
func (a *Authenticator) IsAuthorizedSigner(challenge, signature, addrHex string) (bool, error) {

	addr := common.HexToAddress(addrHex)
	origSigBytes := common.FromHex(signature)

	adjSigBytes := make([]byte, len(origSigBytes))
	copy(adjSigBytes, origSigBytes)
	adjSigBytes[64] -= 27 // Transform V from 27/28 to 0/1 according to the yellow paper

	// retrieve public key from signature
	var challengeHash []byte
	challengeHash = personalMessageHash(challenge)
	recoveredKey, err := ethCrypto.SigToPub(challengeHash, adjSigBytes)
	if err != nil {
		return false, err
	}

	recoveredAddress := ethCrypto.PubkeyToAddress(*recoveredKey)

	// try direct-keyed wallet
	if bytes.Compare(addr.Bytes(), recoveredAddress.Bytes()) == 0 {
		return true, nil
	}

	// try smart-contract wallet
	_ERC1271Caller, err := ERCs.NewERC1271Caller(addr, a.cc)
	if err != nil {
		return false, err
	}

	_ERC1271CallerSession := ERCs.ERC1271CallerSession{
		Contract: _ERC1271Caller,
		CallOpts: bind.CallOpts{
			Pending: false,
			Context: a.ctx,
		},
	}

	magicValue, err := _ERC1271CallerSession.IsValidSignature(challengeHash, origSigBytes)
	if err != nil {
		return false, err
	}

	return magicValue == _ERC1271MagicValue, nil
}

func publicKeyToHash(key *ecdsa.PublicKey) []byte {
	return ethCrypto.Keccak256(ethCrypto.FromECDSAPub(key)[1:])
}

func personalMessageHash(message string) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	return ethCrypto.Keccak256([]byte(msg))
}

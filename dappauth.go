// Package dappauth provides the ability to check if an Ethereume address is an authorized signer for a signature generated via eth_sign rpc method.
// Supports both external wallets and contract wallets.
package dappauth

import (
	"bytes"
	"context"
	"fmt"

	"github.com/dapperlabs/dappauth/ERCs"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

var (
	_ERC1271MagicValue = [4]byte{22, 38, 186, 126} // 0x1626ba7e
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
	var personalChallengeHash []byte
	personalChallengeHash = personalMessageHash(challenge)

	// error is expected when multi sig ("invalid signature length")
	recoveredKey, err := ethCrypto.SigToPub(personalChallengeHash, adjSigBytes)

	// procced with EOA check if no error
	if err == nil {
		recoveredAddress := ethCrypto.PubkeyToAddress(*recoveredKey)

		// try direct-keyed wallet
		if bytes.Compare(addr.Bytes(), recoveredAddress.Bytes()) == 0 {
			return true, nil
		}
		return false, nil
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

	// we send just a regular hash, which then the smart contract hashes ontop to an erc191 hash
	var challengeHash [32]byte
	copy(challengeHash[:], ethCrypto.Keccak256([]byte(challenge)))
	magicValue, err := _ERC1271CallerSession.IsValidSignature(challengeHash, origSigBytes)
	if err != nil {
		return false, err
	}

	return magicValue == _ERC1271MagicValue, nil
}

func personalMessageHash(message string) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	return ethCrypto.Keccak256([]byte(msg))
}

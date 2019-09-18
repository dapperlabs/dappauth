// Package dappauth provides the ability to check if an Ethereume address is an authorized signer for a signature generated via eth_sign rpc method.
// Supports both external wallets and contract wallets.
package dappauth

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strings"

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
	recoveredKey, errEOA := ethCrypto.SigToPub(personalChallengeHash, adjSigBytes)

	// procced with EOA check if no error
	if errEOA == nil {
		recoveredAddress := ethCrypto.PubkeyToAddress(*recoveredKey)

		// try direct-keyed wallet
		if bytes.Compare(addr.Bytes(), recoveredAddress.Bytes()) == 0 {
			return true, nil
		}
	}

	// try smart-contract wallet
	_ERC1271Caller, errCA := ERCs.NewERC1271Caller(addr, a.cc)
	if errCA != nil {
		return false, mergeErrors(errEOA, errCA)
	}

	_ERC1271CallerSession := ERCs.ERC1271CallerSession{
		Contract: _ERC1271Caller,
		CallOpts: bind.CallOpts{
			Pending: false,
			Context: a.ctx,
		},
	}

	// we send just a regular hash, which then the smart contract hashes ontop to an erc191 hash
	magicValue, errCA := _ERC1271CallerSession.IsValidSignature(scMessageHash(challenge), origSigBytes)
	if errCA != nil {
		return false, mergeErrors(errEOA, errCA)
	}

	return magicValue == _ERC1271MagicValue, nil
}

func personalMessageHash(challenge string) []byte {
	b := decodeChallenge(challenge)
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(b), b)
	return ethCrypto.Keccak256([]byte(msg))
}

// This is a hash just over the challenge. The smart contract takes this result and hashes on top to an erc191 hash.
func scMessageHash(challenge string) [32]byte {
	decodedChallenge := decodeChallenge(challenge)
	var challengeHash [32]byte
	copy(challengeHash[:], ethCrypto.Keccak256(decodedChallenge))
	return challengeHash
}

// See https://github.com/MetaMask/eth-sig-util/issues/60
func decodeChallenge(challenge string) []byte {
	b, err := hex.DecodeString(strings.TrimPrefix(challenge, "0x"))
	// if hex decode was successful, then treat is as a hex string
	if err == nil {
		return b
	}
	return []byte(challenge)
}

func mergeErrors(errEOA error, errCA error) error {
	var msgEOA string
	if errEOA == nil {
		msgEOA = "returned false"
	} else {
		msgEOA = fmt.Sprintf("errored with: '%v'", errEOA)
	}

	return fmt.Errorf("Authorisation check failed and errored in 2 alternative flows. 'External Owned Account' check %s. 'Contract Account' check errored with: '%v'", msgEOA, errCA)
}

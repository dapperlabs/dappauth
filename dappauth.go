package dappauth

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/dapperlabs/dappauth/ERCs"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

var (
	_ERC725CoreInterfaceID = [4]byte{210, 2, 21, 141}  // 0xd202158d
	_ERC725InterfaceID     = [4]byte{220, 61, 42, 123} // 0xdc3d2a7b
	_ERC725ActionPurpose   = big.NewInt(2)
)

// Authenticator ..
type Authenticator struct {
	cc  bind.ContractCaller
	ctx context.Context // Network context to support cancellation and timeouts (nil = no timeout)
}

// NewAuthenticator ..
func NewAuthenticator(ctx context.Context, cc bind.ContractCaller) *Authenticator {
	return &Authenticator{
		ctx: ctx,
		cc:  cc,
	}
}

// IsSignerActionableOnAddress ..
func (a *Authenticator) IsSignerActionableOnAddress(challenge, signature, addrHex string) (bool, error) {
	addr := common.HexToAddress(addrHex)
	sigBytes := common.FromHex(signature)
	sigBytes[64] -= 27 // Transform V from 27/28 to 0/1 according to the yellow paper

	// retrieve public key from signature
	var challengeHash []byte
	challengeHash = personalMessageHash(challenge)
	recoveredKey, err := ethCrypto.SigToPub(challengeHash, sigBytes)
	if err != nil {
		return false, err
	}

	recoveredAddress := ethCrypto.PubkeyToAddress(*recoveredKey)

	// try direct-keyed wallet
	if bytes.Compare(addr.Bytes(), recoveredAddress.Bytes()) == 0 {
		return true, nil
	}

	// try smart-contract wallet
	isSupportedContract, err := a.isSupportedContract(addr)
	if err != nil {
		return false, err
	}
	if !isSupportedContract {
		return false, nil
	}

	keyHasActionPurpose, err := a.keyHasActionPurpose(addr, recoveredKey)
	if err != nil {
		return false, err
	}

	return keyHasActionPurpose, err
}

func (a *Authenticator) keyHasActionPurpose(contractAddr common.Address, key *ecdsa.PublicKey) (bool, error) {

	_ERC725CoreCaller, err := ERCs.NewERC725CoreCaller(contractAddr, a.cc)
	if err != nil {
		return false, err
	}

	_ERC725CoreSession := ERCs.ERC725CoreCallerSession{
		Contract: _ERC725CoreCaller,
		CallOpts: bind.CallOpts{
			Pending: false,
			Context: a.ctx,
		},
	}

	var keyBytes [32]byte
	copy(keyBytes[:], publicKeyToHash(key))

	hasActionPurpose, err := _ERC725CoreSession.KeyHasPurpose(keyBytes, _ERC725ActionPurpose)
	if err != nil {
		return false, err
	}

	return hasActionPurpose, nil

}

func (a *Authenticator) isSupportedContract(contractAddr common.Address) (bool, error) {

	_ERC165Caller, err := ERCs.NewERC165Caller(contractAddr, a.cc)
	if err != nil {
		return false, err
	}

	_ERC165CallerSession := ERCs.ERC165CallerSession{
		Contract: _ERC165Caller,
		CallOpts: bind.CallOpts{
			Pending: false,
			Context: a.ctx,
		},
	}

	isSupportsERC725CoreInterface, err := _ERC165CallerSession.SupportsInterface(_ERC725CoreInterfaceID)
	if err != nil {
		return false, err
	}

	if isSupportsERC725CoreInterface {
		return true, nil
	}

	isSupportsERC725Interface, err := _ERC165CallerSession.SupportsInterface(_ERC725InterfaceID)
	if err != nil {
		return false, err
	}

	if isSupportsERC725Interface {
		return true, nil
	}

	return false, nil

}

func publicKeyToHash(key *ecdsa.PublicKey) []byte {
	return ethCrypto.Keccak256(ethCrypto.FromECDSAPub(key)[1:])
}

func personalMessageHash(message string) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	return ethCrypto.Keccak256([]byte(msg))
}

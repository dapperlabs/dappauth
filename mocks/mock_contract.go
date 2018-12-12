// Package mocks holds mocks used for testing
package mocks

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

var (
	_ERC725CoreInterfaceID = [4]byte{210, 2, 21, 141}  // 0xd202158d
	_ERC725InterfaceID     = [4]byte{220, 61, 42, 123} // 0xdc3d2a7b
)

type Contract struct {
	IsSupportsERC725CoreInterface bool
	IsSupportsERC725Interface     bool
	ActionableKey                 *ecdsa.PublicKey
	ErrorOnIsSupportedContract    bool
	ErrorOnKeyHasPurpose          bool
}

func (m *Contract) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return nil, fmt.Errorf("CodeAt not supported")
}

func (m *Contract) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	methodCall := hex.EncodeToString(call.Data[:4])
	methodParams := call.Data[4:]
	switch methodCall {
	case "01ffc9a7":
		return m._01ffc9a7(methodParams)
	case "d202158d":
		return m._d202158d(methodParams)
	default:
		return nil, fmt.Errorf("Unexpected method %v", methodCall)
	}
}

// "isSupportedContract" method call
func (m *Contract) _01ffc9a7(methodParams []byte) ([]byte, error) {

	if m.ErrorOnIsSupportedContract {
		return nil, fmt.Errorf("isSupportedContract call returned an error")
	}

	var interfaceID [4]byte
	copy(interfaceID[:], methodParams[:4])

	if m.IsSupportsERC725CoreInterface && interfaceID == _ERC725CoreInterfaceID {
		return _true()
	}

	if m.IsSupportsERC725Interface && interfaceID == _ERC725InterfaceID {
		return _true()
	}

	return _false()
}

// "KeyHasPurpose" method call
func (m *Contract) _d202158d(methodParams []byte) ([]byte, error) {

	if m.ErrorOnKeyHasPurpose {
		return nil, fmt.Errorf("ErrorOnKeyHasPurpose call returned an error")
	}

	keyBytes := publicKeyToHash(m.ActionableKey)
	if bytes.Compare(keyBytes, methodParams[:32]) == 0 {
		return _true()
	}

	return _false()
}

func _true() ([]byte, error) {
	return hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
}

func _false() ([]byte, error) {
	return hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
}

func publicKeyToHash(key *ecdsa.PublicKey) []byte {
	return ethCrypto.Keccak256(ethCrypto.FromECDSAPub(key)[1:])
}

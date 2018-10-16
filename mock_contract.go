package dappauth

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type mockContract struct {
	isSupportsERC725CoreInterface bool
	isSupportsERC725Interface     bool
	actionableKey                 *ecdsa.PublicKey
	ErrorOnIsSupportedContract    bool
	ErrorOnKeyHasPurpose          bool
}

func (m *mockContract) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return nil, fmt.Errorf("CodeAt not supported")
}

func (m *mockContract) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
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
func (m *mockContract) _01ffc9a7(methodParams []byte) ([]byte, error) {

	if m.ErrorOnIsSupportedContract {
		return nil, fmt.Errorf("isSupportedContract call returned an error")
	}

	var interfaceID [4]byte
	copy(interfaceID[:], methodParams[:4])

	if m.isSupportsERC725CoreInterface && interfaceID == _ERC725CoreInterfaceID {
		return _true()
	}

	if m.isSupportsERC725Interface && interfaceID == _ERC725InterfaceID {
		return _true()
	}

	return _false()
}

// "KeyHasPurpose" method call
func (m *mockContract) _d202158d(methodParams []byte) ([]byte, error) {

	if m.ErrorOnKeyHasPurpose {
		return nil, fmt.Errorf("ErrorOnKeyHasPurpose call returned an error")
	}

	keyBytes := publicKeyToHash(m.actionableKey)
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

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
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

type mockContract struct {
	authorizedKey         *ecdsa.PublicKey
	ErrorIsValidSignature bool
}

func (m *mockContract) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return nil, fmt.Errorf("CodeAt not supported")
}

func (m *mockContract) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	methodCall := hex.EncodeToString(call.Data[:4])
	methodParams := call.Data[4:]
	switch methodCall {
	case "20c13b0b":
		return m._20c13b0b(methodParams)
	default:
		return nil, fmt.Errorf("Unexpected method %v", methodCall)
	}
}

// "IsValidSignature" method call
func (m *mockContract) _20c13b0b(methodParams []byte) ([]byte, error) {

	data := methodParams[96:128]
	sig := methodParams[160:225]

	sig[64] -= 27

	if m.ErrorIsValidSignature {
		return nil, fmt.Errorf("IsValidSignature call returned an error")
	}

	if m.authorizedKey == nil {
		return _false()
	}

	recoveredKey, err := ethCrypto.SigToPub(data, sig)
	if err != nil {
		return nil, err
	}

	recoveredAddress := ethCrypto.PubkeyToAddress(*recoveredKey)
	authorizedKeyAddr := ethCrypto.PubkeyToAddress(*m.authorizedKey)

	if bytes.Compare(authorizedKeyAddr.Bytes(), recoveredAddress.Bytes()) == 0 {
		return _true()
	}

	return _false()
}

func _true() ([]byte, error) {
	// magic value is 0x20c13b0b
	return hex.DecodeString("20c13b0b00000000000000000000000000000000000000000000000000000000")
}

func _false() ([]byte, error) {
	return hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
}

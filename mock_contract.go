package dappauth

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

type mockContract struct {
	address               common.Address
	authorizedKey         *ecdsa.PublicKey
	errorIsValidSignature bool
}

func (m *mockContract) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return nil, fmt.Errorf("CodeAt not supported")
}

func (m *mockContract) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	methodCall := hex.EncodeToString(call.Data[:4])
	methodParams := call.Data[4:]
	switch methodCall {
	case "1626ba7e":
		return m._1626ba7e(methodParams)
	default:
		return nil, fmt.Errorf("Unexpected method %v", methodCall)
	}
}

// "IsValidSignature" method call
func (m *mockContract) _1626ba7e(methodParams []byte) ([]byte, error) {
	// TODO: refactor out of method
	const definition = `[
	{ "name" : "mixedBytes", "constant" : true, "outputs": [{ "name": "a", "type": "bytes32" }, { "name": "b", "type": "bytes" } ] }]`

	abi, err := ethAbi.JSON(strings.NewReader(definition))
	if err != nil {
		return nil, err
	}

	data := [32]byte{}
	sig := []byte{}

	mixedBytes := []interface{}{&data, &sig}
	err = abi.Unpack(&mixedBytes, "mixedBytes", methodParams)
	if err != nil {
		return nil, err
	}

	sig[64] -= 27 // Transform V from 27/28 to 0/1 according to the yellow paper

	if m.errorIsValidSignature {
		return nil, errors.New("Dummy error")
	}

	dataErc191Hash := erc191MessageHash(data[:], m.address)
	recoveredKey, err := ethCrypto.SigToPub(dataErc191Hash, sig)
	if err != nil {
		return nil, err
	}

	if m.authorizedKey == nil {
		return _false()
	}

	recoveredAddress := ethCrypto.PubkeyToAddress(*recoveredKey)
	authorizedKeyAddr := ethCrypto.PubkeyToAddress(*m.authorizedKey)

	if bytes.Compare(authorizedKeyAddr.Bytes(), recoveredAddress.Bytes()) == 0 {
		return _true()
	}

	return _false()
}

func _true() ([]byte, error) {
	// magic value is 0x1626ba7e
	return hex.DecodeString("1626ba7e00000000000000000000000000000000000000000000000000000000")
}

func _false() ([]byte, error) {
	return hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
}

func erc191MessageHash(msg []byte, address common.Address) []byte {

	b := append([]byte{}, 25, 0)
	b = append(b, address.Bytes()...)
	b = append(b, msg...)

	return ethCrypto.Keccak256(b)
}

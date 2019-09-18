package dappauth

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

type dappauthTest struct {
	title                         string // the test's title
	isEOA                         bool
	challenge                     string
	challengeSign                 string // challengeSign should always equal to challenge unless testing for specific key recovery logic
	signingKeys                   []*ecdsa.PrivateKey
	authAddr                      common.Address
	mockContract                  *mockContract
	expectedAuthorizedSignerError bool
	expectedAuthorizedSigner      bool
}

func TestDappAuth(t *testing.T) {

	keyA, err := ethCrypto.GenerateKey()
	checkError(err, t)
	keyB, err := ethCrypto.GenerateKey()
	checkError(err, t)
	keyC, err := ethCrypto.GenerateKey()
	checkError(err, t)

	dappauthTests := []dappauthTest{
		{
			title:         "External wallets should be authorized signers over their address",
			isEOA:         true,
			challenge:     "foo",
			challengeSign: "foo",
			signingKeys:   []*ecdsa.PrivateKey{keyA},
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				address:               [20]byte{},
				authorizedKey:         nil,
				errorIsValidSignature: false,
			},
			expectedAuthorizedSignerError: false,
			expectedAuthorizedSigner:      true,
		},
		{
			title:         "External wallets should NOT be authorized signers when signing the wrong challenge",
			isEOA:         true,
			challenge:     "foo",
			challengeSign: "bar",
			signingKeys:   []*ecdsa.PrivateKey{keyA},
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				address:               [20]byte{},
				authorizedKey:         nil,
				errorIsValidSignature: false,
			},
			expectedAuthorizedSignerError: false,
			expectedAuthorizedSigner:      false,
		},
		{
			title:         "External wallets should NOT be authorized signers over OTHER addresses",
			isEOA:         true,
			challenge:     "foo",
			challengeSign: "foo",
			signingKeys:   []*ecdsa.PrivateKey{keyA},
			authAddr:      ethCrypto.PubkeyToAddress(keyB.PublicKey),
			mockContract: &mockContract{
				address:               [20]byte{},
				authorizedKey:         nil,
				errorIsValidSignature: false,
			},
			expectedAuthorizedSignerError: false,
			expectedAuthorizedSigner:      false,
		},
		{
			title:         "Smart-contract wallets with a 1-of-1 correct internal key should be authorized signers over their address",
			isEOA:         false,
			challenge:     "foo",
			challengeSign: "foo",
			signingKeys:   []*ecdsa.PrivateKey{keyB},
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				address:               ethCrypto.PubkeyToAddress(keyA.PublicKey),
				authorizedKey:         &keyB.PublicKey,
				errorIsValidSignature: false,
			},
			expectedAuthorizedSignerError: false,
			expectedAuthorizedSigner:      true,
		},
		{
			title:         "Smart-contract wallets with a 1-of-2 correct internal key should be authorized signers over their address",
			isEOA:         false,
			challenge:     "foo",
			challengeSign: "foo",
			signingKeys:   []*ecdsa.PrivateKey{keyB, keyC},
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				address:               ethCrypto.PubkeyToAddress(keyA.PublicKey),
				authorizedKey:         &keyB.PublicKey,
				errorIsValidSignature: false,
			},
			expectedAuthorizedSignerError: false,
			expectedAuthorizedSigner:      true,
		},
		{
			title:         "Smart-contract wallets with a 1-of-1 incorrect internal key should NOT be authorized signers over their address",
			isEOA:         false,
			challenge:     "foo",
			challengeSign: "foo",
			signingKeys:   []*ecdsa.PrivateKey{keyB},
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				address:               ethCrypto.PubkeyToAddress(keyA.PublicKey),
				authorizedKey:         &keyC.PublicKey,
				errorIsValidSignature: false,
			},
			expectedAuthorizedSignerError: false,
			expectedAuthorizedSigner:      false,
		},
		{
			title:         "IsAuthorizedSigner should error when smart-contract call errors",
			isEOA:         false,
			challenge:     "foo",
			challengeSign: "foo",
			signingKeys:   []*ecdsa.PrivateKey{keyB},
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				address:               ethCrypto.PubkeyToAddress(keyA.PublicKey),
				authorizedKey:         &keyB.PublicKey,
				errorIsValidSignature: true,
			},
			expectedAuthorizedSignerError: true,
			expectedAuthorizedSigner:      false,
		},
	}

	// iterate over main test cases
	for _, test := range dappauthTests {
		t.Run(test.title, func(t *testing.T) {
			authenticator := NewAuthenticator(nil, test.mockContract)

			var sig string
			for _, signingKey := range test.signingKeys {
				sig += generateSignature(test.isEOA, test.challengeSign, signingKey, test.authAddr, t)
			}
			isAuthorizedSigner, err := authenticator.IsAuthorizedSigner(test.challenge, sig, test.authAddr.Hex())

			expectBool(err != nil, test.expectedAuthorizedSignerError, t)
			expectBool(isAuthorizedSigner, test.expectedAuthorizedSigner, t)
		})
	}

	// this test adds to test coverage
	t.Run("Invalid signature should fail", func(t *testing.T) {
		authenticator := NewAuthenticator(nil, &mockContract{
			address:               [20]byte{},
			authorizedKey:         nil,
			errorIsValidSignature: false,
		})

		invalidSigBytes := [65]byte{}
		invalidSig := hex.EncodeToString(invalidSigBytes[:])
		isAuthorizedSigner, err := authenticator.IsAuthorizedSigner("foo", invalidSig, ethCrypto.PubkeyToAddress(keyA.PublicKey).Hex())

		expectBool(err != nil, true, t)
		expectBool(isAuthorizedSigner, false, t)
	})

}

// It should decode challenge as utf8 by default when computing EOA personal messages hash
func TestPersonalMessageDecodeUTF8(t *testing.T) {
	eoaHash := hex.EncodeToString(personalMessageHash("foo"))
	expectString(fmt.Sprintf("0x%s", eoaHash), "0x76b2e96714d3b5e6eb1d1c509265430b907b44f72b2a22b06fcd4d96372b8565", t)

	scHashBytes := scMessageHash("foo")
	scHash := hex.EncodeToString(scHashBytes[:])
	expectString(fmt.Sprintf("0x%s", scHash), "0x41b1a0649752af1b28b3dc29a1556eee781e4a4c3a1f7f53f90fa834de098c4d", t)
}

// It should decode challenge as hex if hex is detected when computing EOA personal messages hash
// See https://github.com/MetaMask/eth-sig-util/issues/60
func TestPersonalMessageDecodeHex(t *testing.T) {

	// result if 0xffff is decoded as hex:  13a6aa3102b2d639f36804a2d7c31469618fd7a7907c658a7b2aa91a06e31e47
	// result if 0xffff is decoded as utf8: 247aefb5d2e5b17fca61f786c779f7388485460c13e51308f88b2ff84ffa6851
	eoaHash := hex.EncodeToString(personalMessageHash("0xffff"))
	expectString(fmt.Sprintf("0x%s", eoaHash), "0x13a6aa3102b2d639f36804a2d7c31469618fd7a7907c658a7b2aa91a06e31e47", t)

	// result if 0xffff is decoded as hex:  06d41322d79dfed27126569cb9a80eb0967335bf2f3316359d2a93c779fcd38a
	// result if 0xffff is decoded as utf8: f0443ea82539c5136844b0a175f544b7ee7bc0fc5ce940bad19f08eaf618af71
	scHashBytes := scMessageHash("0xffff")
	scHash := hex.EncodeToString(scHashBytes[:])
	expectString(fmt.Sprintf("0x%s", scHash), "0x06d41322d79dfed27126569cb9a80eb0967335bf2f3316359d2a93c779fcd38a", t)
}

func generateSignature(isEOA bool, msg string, key *ecdsa.PrivateKey, address common.Address, t *testing.T) string {
	if isEOA {
		return signEOAPersonalMessage(msg, key, t)
	}
	return signERC1654PersonalMessage(msg, key, address, t)
}

// emulates what EOA wallets like MetaMask perform
func signEOAPersonalMessage(msg string, key *ecdsa.PrivateKey, t *testing.T) string {
	ethMsgHash := personalMessageHash(msg)
	sig, err := ethCrypto.Sign(ethMsgHash, key)
	checkError(err, t)

	sig[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return hex.EncodeToString(sig)
}

func signERC1654PersonalMessage(msg string, key *ecdsa.PrivateKey, address common.Address, t *testing.T) string {

	// we hash once before erc191MessageHash as it will be transmitted to Ethereum nodes and potentially logged
	msgHash := erc191MessageHash(ethCrypto.Keccak256([]byte(msg)), address)

	sig, err := ethCrypto.Sign(msgHash, key)
	checkError(err, t)

	sig[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return hex.EncodeToString(sig)
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

func expectBool(actual, expected bool, t *testing.T) {
	if actual != expected {
		t.Errorf("expected %v to be %v", actual, expected)
	}
}

func expectString(actual, expected string, t *testing.T) {
	if actual != expected {
		t.Errorf("expected %s to be %s", actual, expected)
	}
}

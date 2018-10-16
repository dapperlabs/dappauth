package dappauth

import (
	"crypto/ecdsa"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

type dappauthTest struct {
	title                                    string // the test's title
	challenge                                string
	challengeSign                            string // challengeSign should always equal to challenge unless testing for specific key recovery logic
	signingKey                               *ecdsa.PrivateKey
	authAddr                                 common.Address
	mockContract                             *mockContract
	expectedIsSignerActionableOnAddressError bool
	expectedIsSignerActionableOnAddress      bool
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
			title:         "Direct-keyed wallets should have actionable control over their address",
			challenge:     "foo",
			challengeSign: "foo",
			signingKey:    keyA,
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				isSupportsERC725CoreInterface: false,
				isSupportsERC725Interface:     false,
				actionableKey:                 nil,
				ErrorOnIsSupportedContract:    false,
				ErrorOnKeyHasPurpose:          false,
			},
			expectedIsSignerActionableOnAddressError: false,
			expectedIsSignerActionableOnAddress:      true,
		},
		{
			title:         "Direct-keyed wallets should NOT have actionable control when signing the wrong challenge",
			challenge:     "foo",
			challengeSign: "bar",
			signingKey:    keyA,
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				isSupportsERC725CoreInterface: false,
				isSupportsERC725Interface:     false,
				actionableKey:                 nil,
				ErrorOnIsSupportedContract:    false,
				ErrorOnKeyHasPurpose:          false,
			},
			expectedIsSignerActionableOnAddressError: false,
			expectedIsSignerActionableOnAddress:      false,
		},
		{
			title:         "Direct-keyed wallets should NOT have actionable control over OTHER addresses",
			challenge:     "foo",
			challengeSign: "foo",
			signingKey:    keyA,
			authAddr:      ethCrypto.PubkeyToAddress(keyB.PublicKey),
			mockContract: &mockContract{
				isSupportsERC725CoreInterface: false,
				isSupportsERC725Interface:     false,
				actionableKey:                 nil,
				ErrorOnIsSupportedContract:    false,
				ErrorOnKeyHasPurpose:          false,
			},
			expectedIsSignerActionableOnAddressError: false,
			expectedIsSignerActionableOnAddress:      false,
		},
		{
			title:         "Smart-contract wallets with support for ERC725Core interface and action key should have actionable control over their address",
			challenge:     "foo",
			challengeSign: "foo",
			signingKey:    keyB,
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				isSupportsERC725CoreInterface: true,
				isSupportsERC725Interface:     false,
				actionableKey:                 &keyB.PublicKey,
				ErrorOnIsSupportedContract:    false,
				ErrorOnKeyHasPurpose:          false,
			},
			expectedIsSignerActionableOnAddressError: false,
			expectedIsSignerActionableOnAddress:      true,
		},
		{
			title:         "Smart-contract wallets with support for ERC725 interface and action key should have actionable control over their address",
			challenge:     "foo",
			challengeSign: "foo",
			signingKey:    keyB,
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				isSupportsERC725CoreInterface: false,
				isSupportsERC725Interface:     true,
				actionableKey:                 &keyB.PublicKey,
				ErrorOnIsSupportedContract:    false,
				ErrorOnKeyHasPurpose:          false,
			},
			expectedIsSignerActionableOnAddressError: false,
			expectedIsSignerActionableOnAddress:      true,
		},
		{
			title:         "Smart-contract wallets WITHOUT support for ERC725/ERC725Core interface and action key should NOT have actionable control over their address",
			challenge:     "foo",
			challengeSign: "foo",
			signingKey:    keyB,
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				isSupportsERC725CoreInterface: false,
				isSupportsERC725Interface:     false,
				actionableKey:                 &keyB.PublicKey,
				ErrorOnIsSupportedContract:    false,
				ErrorOnKeyHasPurpose:          false,
			},
			expectedIsSignerActionableOnAddressError: false,
			expectedIsSignerActionableOnAddress:      false,
		},
		{
			title:         "Smart-contract wallets with support for ERC725 interface and incorrect action key should have NO actionable control over their address",
			challenge:     "foo",
			challengeSign: "foo",
			signingKey:    keyB,
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				isSupportsERC725CoreInterface: false,
				isSupportsERC725Interface:     true,
				actionableKey:                 &keyC.PublicKey,
				ErrorOnIsSupportedContract:    false,
				ErrorOnKeyHasPurpose:          false,
			},
			expectedIsSignerActionableOnAddressError: false,
			expectedIsSignerActionableOnAddress:      false,
		},
		{
			title:         "IsSignerActionableOnAddress should error when smart-contract call errors",
			challenge:     "foo",
			challengeSign: "foo",
			signingKey:    keyB,
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				isSupportsERC725CoreInterface: true,
				isSupportsERC725Interface:     false,
				actionableKey:                 &keyB.PublicKey,
				ErrorOnIsSupportedContract:    true,
				ErrorOnKeyHasPurpose:          false,
			},
			expectedIsSignerActionableOnAddressError: true,
			expectedIsSignerActionableOnAddress:      false,
		},
		{
			title:         "IsSignerActionableOnAddress should error when smart-contract call errors",
			challenge:     "foo",
			challengeSign: "foo",
			signingKey:    keyB,
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				isSupportsERC725CoreInterface: true,
				isSupportsERC725Interface:     false,
				actionableKey:                 &keyB.PublicKey,
				ErrorOnIsSupportedContract:    false,
				ErrorOnKeyHasPurpose:          true,
			},
			expectedIsSignerActionableOnAddressError: true,
			expectedIsSignerActionableOnAddress:      false,
		},
	}

	// iterate over main test cases
	for _, test := range dappauthTests {
		t.Run(test.title, func(t *testing.T) {

			authenticator := NewAuthenticator(nil, test.mockContract)

			sig := signPersonalMessage(test.challengeSign, test.signingKey, t)
			isSignerActionableOnAddress, err := authenticator.IsSignerActionableOnAddress(test.challenge, sig, test.authAddr.Hex())
			expectBool(err != nil, test.expectedIsSignerActionableOnAddressError, t)
			if err != nil {
				expectBool(isSignerActionableOnAddress, test.expectedIsSignerActionableOnAddress, t)
			}
		})
	}

	// extra test for extra code-coverage
	t.Run("Expect messy signature to cause an error", func(t *testing.T) {
		test := dappauthTests[0]

		authenticator := NewAuthenticator(nil, test.mockContract)

		sig := signPersonalMessage(test.challengeSign, test.signingKey, t)
		sig = sig + "ffff" // mess the signature
		_, err := authenticator.IsSignerActionableOnAddress(test.challenge, sig, test.authAddr.Hex())
		expectBool(err != nil, true, t)
	})

}

func signPersonalMessage(msg string, key *ecdsa.PrivateKey, t *testing.T) string {
	ethMsgHashBytes := personalMessageHash(msg)
	sig, err := ethCrypto.Sign(ethMsgHashBytes, key)
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

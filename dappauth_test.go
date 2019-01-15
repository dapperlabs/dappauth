package dappauth

import (
	"crypto/ecdsa"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

type dappauthTest struct {
	title                         string // the test's title
	challenge                     string
	challengeSign                 string // challengeSign should always equal to challenge unless testing for specific key recovery logic
	signingKey                    *ecdsa.PrivateKey
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
			challenge:     "foo",
			challengeSign: "foo",
			signingKey:    keyA,
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				authorizedKey:         nil,
				ErrorIsValidSignature: false,
			},
			expectedAuthorizedSignerError: false,
			expectedAuthorizedSigner:      true,
		},
		{
			title:         "External wallets should NOT be authorized signers when signing the wrong challenge",
			challenge:     "foo",
			challengeSign: "bar",
			signingKey:    keyA,
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				authorizedKey:         nil,
				ErrorIsValidSignature: false,
			},
			expectedAuthorizedSignerError: false,
			expectedAuthorizedSigner:      false,
		},
		{
			title:         "External wallets should NOT be authorized signers over OTHER addresses",
			challenge:     "foo",
			challengeSign: "foo",
			signingKey:    keyA,
			authAddr:      ethCrypto.PubkeyToAddress(keyB.PublicKey),
			mockContract: &mockContract{
				authorizedKey:         nil,
				ErrorIsValidSignature: false,
			},
			expectedAuthorizedSignerError: false,
			expectedAuthorizedSigner:      false,
		},
		{
			title:         "Smart-contract wallets with a 1-of-1 correct internal key should be authorized signers over their address",
			challenge:     "foo",
			challengeSign: "foo",
			signingKey:    keyB,
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				authorizedKey:         &keyB.PublicKey,
				ErrorIsValidSignature: false,
			},
			expectedAuthorizedSignerError: false,
			expectedAuthorizedSigner:      true,
		},
		{
			title:         "Smart-contract wallets with a 1-of-1 incorrect internal key should be authorized signers over their address",
			challenge:     "foo",
			challengeSign: "foo",
			signingKey:    keyB,
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				authorizedKey:         &keyC.PublicKey,
				ErrorIsValidSignature: false,
			},
			expectedAuthorizedSignerError: false,
			expectedAuthorizedSigner:      false,
		},
		{
			title:         "IsAuthorizedSigner should error when smart-contract call errors",
			challenge:     "foo",
			challengeSign: "foo",
			signingKey:    keyB,
			authAddr:      ethCrypto.PubkeyToAddress(keyA.PublicKey),
			mockContract: &mockContract{
				authorizedKey:         &keyB.PublicKey,
				ErrorIsValidSignature: true,
			},
			expectedAuthorizedSignerError: true,
			expectedAuthorizedSigner:      false,
		},
	}

	// iterate over main test cases
	for _, test := range dappauthTests {
		t.Run(test.title, func(t *testing.T) {

			authenticator := NewAuthenticator(nil, test.mockContract)

			sig := signPersonalMessage(test.challengeSign, test.signingKey, t)
			isAuthorizedSigner, err := authenticator.IsAuthorizedSigner(test.challenge, sig, test.authAddr.Hex())
			expectBool(err != nil, test.expectedAuthorizedSignerError, t)
			expectBool(isAuthorizedSigner, test.expectedAuthorizedSigner, t)
		})
	}

	// extra test for extra code-coverage
	t.Run("Expect messy signature to cause an error", func(t *testing.T) {
		test := dappauthTests[0]

		authenticator := NewAuthenticator(nil, test.mockContract)

		sig := signPersonalMessage(test.challengeSign, test.signingKey, t)
		sig = sig + "ffff" // mess the signature
		_, err := authenticator.IsAuthorizedSigner(test.challenge, sig, test.authAddr.Hex())
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

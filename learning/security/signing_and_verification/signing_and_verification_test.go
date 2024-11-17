package signingandverification_test

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSigningAndVerification(t *testing.T) {
	t.Parallel()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	digest := sha256.Sum256([]byte("some data"))

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest[:])
	require.NoError(t, err)

	err = rsa.VerifyPKCS1v15(&privateKey.PublicKey, crypto.SHA256, digest[:], signature)
	require.NoError(t, err)
}

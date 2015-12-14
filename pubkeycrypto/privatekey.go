package pubkeycrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
)

// PrivateKey is a wrapper over rsa.PrivateKey.
type PrivateKey struct {
	inner *rsa.PrivateKey
}

// GeneratePrivateKey generates a PrivateKey.
func GeneratePrivateKey() (*PrivateKey, error) {
	// NOTE: rand.Reader uses /dev/urandom, which is not recommended for generating long-term cryptographic keys
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	privateKey := &PrivateKey{
		inner: rsaPrivateKey,
	}
	return privateKey, nil
}

// GetPublicKey returns the public key based on this private key
func (pk *PrivateKey) GetPublicKey() *PublicKey {
	return publicKeyFromRsa(&pk.inner.PublicKey)
}

// DecryptOAEP decrypts the ciphertext using RSA-OAEP.
func (pk *PrivateKey) DecryptOAEP(ciphertext []byte) ([]byte, error) {
	// TODO: what hashing algo to use?
	// TODO: why do we need rand.Reader in decrypt?
	return rsa.DecryptOAEP(sha512.New(), rand.Reader, pk.inner, ciphertext, nil)
}

// Sign signs the given message using this private key and returns the signature.
func (pk *PrivateKey) Sign(message []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, pk.inner, sha512HashFunc, sha512Sum(message))
}

// ToString produces a string representation of this private key.
func (pk *PrivateKey) ToString() string {
	bytes := x509.MarshalPKCS1PrivateKey(pk.inner)
	return base64.StdEncoding.EncodeToString(bytes)
}

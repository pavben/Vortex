package pubkeycrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"errors"
)

// Errors
var (
	ErrUnexpectedPublicKeyFormat = errors.New("Unexpected public key format")
)

// PublicKey is a wrapper over rsa.PublicKey.
type PublicKey struct {
	inner *rsa.PublicKey
}

func publicKeyFromRsa(rsaPublicKey *rsa.PublicKey) *PublicKey {
	return &PublicKey{
		inner: rsaPublicKey,
	}
}

// PublicKeyFromBytes returns a PublicKey given a slice of bytes in X509 PKIX public key format.
func PublicKeyFromBytes(b []byte) (*PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	rsaPublicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, ErrUnexpectedPublicKeyFormat
	}
	return publicKeyFromRsa(rsaPublicKey), nil
}

// EncryptOAEP encrypts the message using RSA-OAEP.
func (pk *PublicKey) EncryptOAEP(message []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha512.New(), rand.Reader, pk.inner, message, nil)
}

// VerifySignature verifies that the signature for this message was produced with the private key for which this is the public key.
func (pk *PublicKey) VerifySignature(message, signature []byte) error {
	return rsa.VerifyPKCS1v15(pk.inner, sha512HashFunc, sha512Sum(message), signature)
}

// ToBytes returns a slice of bytes representing this public key in X509 PKIX public key format.
func (pk *PublicKey) ToBytes() []byte {
	bytes, _ := x509.MarshalPKIXPublicKey(pk.inner)
	// TODO: Find out how this can fail
	return bytes
}

// ToString produces a string representation of this public key.
func (pk *PublicKey) ToString() string {
	return base64.StdEncoding.EncodeToString(pk.ToBytes())
}

// Sha1Hash returns the SHA1 hash of ToBytes() of this public key.
func (pk *PublicKey) Sha1Hash() string {
	return base64.StdEncoding.EncodeToString(sha1Sum(pk.ToBytes()))
}

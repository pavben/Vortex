package vortexconn

import (
	"crypto/aes"
	"fmt"
	"net"

	"github.com/pavben/Vortex/aesstream"
	"github.com/pavben/Vortex/pubkeycrypto"
)

// Connect establishes an encrypted connection and returns it.
func Connect(addr string, keyPair *pubkeycrypto.KeyPair) (*Connection, error) {
	tcpConn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("dial error: %v", err)
	}
	// Send the server our public key
	err = writeByteChunkPlain(tcpConn, keyPair.PublicKey.ToBytes())
	if err != nil {
		return nil, fmt.Errorf("error sending our public key: %v", err)
	}
	// Read the server's public key
	serverPublicKeyBytes, err := readByteChunkPlain(tcpConn)
	if err != nil {
		return nil, fmt.Errorf("error reading server public key: %v", err)
	}
	serverPublicKey, err := pubkeycrypto.PublicKeyFromBytes(serverPublicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing server public key: %v", err)
	}
	// Read the RSA-encrypted AES key
	encryptedAesKey, err := readByteChunkPlain(tcpConn)
	if err != nil {
		return nil, fmt.Errorf("error reading the AES key from server: %v", err)
	}
	aesKey, err := keyPair.PrivateKey.DecryptOAEP(encryptedAesKey)
	if err != nil {
		return nil, fmt.Errorf("error decrypting the AES key from server: %v", err)
	}
	// Require AES256
	if len(aesKey) != 32 {
		return nil, fmt.Errorf("server sent us an unexpected-length AES key: %d", len(aesKey))
	}
	// Read the server-generated IV
	iv, err := readByteChunkPlain(tcpConn)
	if err != nil {
		return nil, fmt.Errorf("error reading the IV from server: %v", err)
	}
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("IV length %d must equal to the AES block size %d", len(iv), aes.BlockSize)
	}
	// Create the AES stream
	aesStream, err := aesstream.NewAesStream(tcpConn, aesKey, iv)
	if err != nil {
		return nil, err
	}
	return &Connection{
		tcpConn:        tcpConn,
		aesStream:      aesStream,
		theirPublicKey: serverPublicKey,
	}, nil
}

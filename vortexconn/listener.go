package vortexconn

import (
	"crypto/aes"
	"fmt"
	"net"

	"github.com/pavben/Vortex/aesstream"
	"github.com/pavben/Vortex/pubkeycrypto"
)

// Listener is a listener for our custom AES-encrypted connections.
type Listener struct {
	tcpListener                net.Listener
	keyPair                    *pubkeycrypto.KeyPair
	establishedConnectionsChan chan *Connection
	shutdownChan               chan struct{}
}

// Listen creates and returns the Listener.
func Listen(laddr string, keyPair *pubkeycrypto.KeyPair) (*Listener, error) {
	tcpListener, err := net.Listen("tcp", laddr)
	if err != nil {
		return nil, err
	}
	listener := &Listener{
		tcpListener:                tcpListener,
		keyPair:                    keyPair,
		establishedConnectionsChan: make(chan *Connection),
		shutdownChan:               make(chan struct{}),
	}
	go func() {
		for {
			tcpConn, err2 := tcpListener.Accept()
			if err2 != nil {
				return
			}
			go func() {
				vortexConn, err3 := initConnectionAsListener(tcpConn, listener)
				if err3 != nil {
					fmt.Println("initConnectionAsListener failed:", err3)
					return
				}
				listener.establishedConnectionsChan <- vortexConn
			}()
		}
	}()
	return listener, nil
}

// Accept returns the next Connection that has been established on our listener.
func (l *Listener) Accept() *Connection {
	select {
	case newConn := <-l.establishedConnectionsChan:
		return newConn
	case <-l.shutdownChan:
		return nil
	}
}

func initConnectionAsListener(tcpConn net.Conn, listener *Listener) (*Connection, error) {
	// Send our public key
	err := writeByteChunkPlain(tcpConn, listener.keyPair.PublicKey.ToBytes())
	if err != nil {
		return nil, err
	}
	// Receive and parse the client's public key
	clientPublicKeyBytes, err := readByteChunkPlain(tcpConn)
	if err != nil {
		return nil, fmt.Errorf("error reading client public key: %v", err)
	}
	clientPublicKey, err := pubkeycrypto.PublicKeyFromBytes(clientPublicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing client public key: %v", err)
	}
	// Generate 32-byte AES key
	aesKey, err := generateRandomBytes(32)
	if err != nil {
		return nil, err
	}
	// Encrypt it with their public key using RSA
	aesKeyForClient, err := clientPublicKey.EncryptOAEP(aesKey)
	if err != nil {
		return nil, err
	}
	// Send the encrypted AES key
	err = writeByteChunkPlain(tcpConn, aesKeyForClient)
	if err != nil {
		return nil, err
	}
	// Generate 16 bytes IV and send it
	iv, err := generateRandomBytes(aes.BlockSize)
	if err != nil {
		return nil, err
	}
	err = writeByteChunkPlain(tcpConn, iv)
	if err != nil {
		return nil, err
	}
	// Create the AES stream
	aesStream, err := aesstream.NewAesStream(tcpConn, aesKey, iv)
	if err != nil {
		return nil, err
	}
	return &Connection{
		tcpConn:        tcpConn,
		aesStream:      aesStream,
		theirPublicKey: clientPublicKey,
	}, nil
}

// Close closes this listener. It will not implicitly close existing connections.
func (l *Listener) Close() {
	close(l.shutdownChan)
	l.tcpListener.Close()
}

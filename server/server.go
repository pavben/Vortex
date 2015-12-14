package main

import (
	"fmt"

	"github.com/pavben/Vortex/pubkeycrypto"
	"github.com/pavben/Vortex/vortexconn"
)

func main() {
	keyPair, err := pubkeycrypto.GenerateKeyPair()
	if err != nil {
		fmt.Println("Error generating keypair:", err)
	}
	listener, err := vortexconn.Listen(":27805", keyPair)
	if err != nil {
		fmt.Println("Listener error:", err)
	}
	clientConn := listener.Accept()
	fmt.Println(clientConn)
	listener.Close()
}

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
	conn, err := vortexconn.Connect("127.0.0.1:27805", keyPair)
	fmt.Println(conn, err)
}

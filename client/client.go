package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/pavben/Vortex/natpmp"
	"github.com/pavben/Vortex/pubkeycrypto"
	"github.com/pavben/Vortex/try"
	"github.com/pavben/Vortex/vortexconn"
)

func main() {
	keyPair, err := pubkeycrypto.GenerateKeyPair()
	if err != nil {
		fmt.Println("Error generating keypair:", err)
		return
	}
	var listenerPort uint16
	listenerI, err := try.Do(func() (interface{}, error) {
		listenerPort = randomPort()
		return vortexconn.Listen(":"+strconv.Itoa(int(listenerPort)), keyPair)
	}, 5)
	if err != nil {
		fmt.Println("Listener error:", err)
		return
	}
	listener := listenerI.(*vortexconn.Listener)
	fmt.Println("Listening on port", listenerPort)
	portMap, err := natpmp.AddPortMappingForAnyExternalPort(listenerPort, nil)
	if err != nil {
		fmt.Println("Port mapping error:", err)
		return
	}
	defer portMap.Close()
	fmt.Println("Port map result:", portMap, portMap.State)
	clientConn := listener.Accept()
	fmt.Println(clientConn)
	listener.Close()
}

func randomPort() uint16 {
	rand.Seed(time.Now().UnixNano())
	var base = 1024
	var max = 65535
	return uint16(base + rand.Intn(max-base))
}

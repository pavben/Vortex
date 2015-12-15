package natpmp

// This package is a wrapper over github.com/jackpal/go-nat-pmp

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	jpnatpmp "github.com/jackpal/go-nat-pmp"
	"github.com/pavben/Vortex/assert"
	"github.com/pavben/Vortex/try"
)

// TODO: Make this a lot higher (1 hour)
const targetPortMapLifetime = time.Duration(15) * time.Second

func targetPortMapLifetimeSeconds() uint32 {
	return uint32(targetPortMapLifetime.Seconds())
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type NatPmpPortMap struct {
	State                 *PortMapState
	requestedInternalPort uint16
	onStateChanging       *func(*PortMapState)
	releaseChan           chan struct{}
}

type PortMapState struct {
	ExternalIp   string
	ExternalPort uint16
	InternalPort uint16
	lifetime     uint32
}

// TODO: IsNatPmpAvailable possible?

// onStateChanging callback is invoked synchronously just before PortMapState is about to change, so your existing NatPmpPortMap.State will represent the old state, and the arg to onStateChanging will be the new.
func AddPortMappingForAnyExternalPort(internalPort uint16, onStateChanging *func(*PortMapState)) (*NatPmpPortMap, error) {
	portMapResultI, err := try.Do(func() (interface{}, error) {
		return addNatPmpPortMapping(randomPort(), internalPort)
	}, 3)
	if err != nil {
		return nil, fmt.Errorf("error adding port mapping: %v", err)
	}
	portMapState := portMapResultI.(*PortMapState)
	natPmpPortMap := &NatPmpPortMap{
		State: portMapState,
		requestedInternalPort: internalPort,
		onStateChanging:       onStateChanging,
		releaseChan:           make(chan struct{}),
	}
	go natPmpPortMap.startMaintainer()
	return natPmpPortMap, nil
}

func (pm *NatPmpPortMap) startMaintainer() {
	for {
		var renewDuration time.Duration
		if pm.State == nil {
			renewDuration = targetPortMapLifetime
		} else {
			renewDuration = time.Duration(pm.State.lifetime) * time.Second
		}
		// Renew 10 seconds before the stated expiry
		assert.Assert(targetPortMapLifetimeSeconds() >= 15, "Renew duration is too short: %v", renewDuration)
		renewDuration -= time.Duration(10) * time.Second
		fmt.Println("Maintainer waiting for", renewDuration)
		select {
		case <-time.Tick(renewDuration):
			// Start with trying to remap the port we already have (if present). If unable, try random ports for a total of 3 attempts.
			var externalPortToTryIfNonZero uint16
			if pm.State != nil {
				externalPortToTryIfNonZero = pm.State.ExternalPort
			}
			portMapResultI, err := try.Do(func() (interface{}, error) {
				var portToTry uint16
				if externalPortToTryIfNonZero != 0 {
					portToTry = externalPortToTryIfNonZero
					externalPortToTryIfNonZero = 0
				} else {
					portToTry = randomPort()
				}
				return addNatPmpPortMapping(portToTry, pm.requestedInternalPort)
			}, 3)
			var newPortMapState *PortMapState
			if err == nil {
				newPortMapState = portMapResultI.(*PortMapState)
				fmt.Println("Renewed port mapping. New state:", newPortMapState)
			} else {
				fmt.Println("error renewing port mapping: %v", err)
				newPortMapState = nil
			}
			if !portMapStateSame(pm.State, newPortMapState) {
				fmt.Println("Port map state changed!")
				// Synchronously invoke the callback before applying the new state
				if pm.onStateChanging != nil {
					(*pm.onStateChanging)(newPortMapState)
				}
			}
			pm.State = newPortMapState
		case <-pm.releaseChan:
			// Release the mapping and return
			if pm.State != nil {
				removeNatPmpPortMapping(pm.State.InternalPort)
			}
			pm.State = nil
			return
		}
	}
}

// Close releases the port mapping. The onStateChanging callback will not be called for this operation.
func (pm *NatPmpPortMap) Close() {
	pm.releaseChan <- struct{}{}
}

func addNatPmpPortMapping(externalPort, internalPort uint16) (*PortMapState, error) {
	fmt.Println("Adding mapping for internal port", internalPort)
	natPmpClient, err := jpnatpmp.NewClientForDefaultGateway()
	if err != nil {
		return nil, fmt.Errorf("error creating the client for the default gateway: %v", err)
	}
	externalAddressResult, err := natPmpClient.GetExternalAddress()
	if err != nil {
		return nil, fmt.Errorf("error getting the external address: %v", err)
	}
	externalIp := ipv4FromBytes(externalAddressResult.ExternalIPAddress)
	addPortMappingResult, err := natPmpClient.AddPortMapping("tcp", int(internalPort), int(externalPort), int(targetPortMapLifetimeSeconds()))
	if err != nil {
		return nil, fmt.Errorf("error adding port mapping: %v", err)
	}
	assert.Assert(addPortMappingResult.InternalPort == internalPort, "addPortMappingResult contains a different internal port from what we asked for: %d vs %d", addPortMappingResult.InternalPort, internalPort)
	assert.Assert(addPortMappingResult.MappedExternalPort == externalPort, "addPortMappingResult contains a different external port from what we asked for: %d vs %d", addPortMappingResult.MappedExternalPort, externalPort)
	assert.Assert(addPortMappingResult.PortMappingLifetimeInSeconds == targetPortMapLifetimeSeconds() || addPortMappingResult.PortMappingLifetimeInSeconds >= 15, "addPortMappingResult contains a lifetime that is too low: %d seconds", addPortMappingResult.PortMappingLifetimeInSeconds)
	portMapState := &PortMapState{
		ExternalIp:   externalIp.String(),
		ExternalPort: externalPort,
		InternalPort: internalPort,
		lifetime:     addPortMappingResult.PortMappingLifetimeInSeconds,
	}
	return portMapState, nil
}

func removeNatPmpPortMapping(internalPort uint16) error {
	fmt.Println("Removing mapping for internal port", internalPort)
	natPmpClient, err := jpnatpmp.NewClientForDefaultGateway()
	if err != nil {
		return fmt.Errorf("error creating the client for the default gateway: %v", err)
	}
	_, err = natPmpClient.AddPortMapping("tcp", int(internalPort), 0, 0)
	return err
}

func portMapStateSame(a, b *PortMapState) bool {
	if a != nil && b != nil {
		if a.ExternalIp != b.ExternalIp {
			return false
		}
		if a.ExternalPort != b.ExternalPort {
			return false
		}
		// Should never change, but adding it here for completeness
		if a.InternalPort != b.InternalPort {
			return false
		}
	} else if a != b { // One or both are null
		return false
	}
	return true
}

func randomPort() uint16 {
	var base = 1024
	var max = 65535
	return uint16(base + rng.Intn(max-base))
}

func ipv4FromBytes(b [4]byte) net.IP {
	return net.IPv4(b[0], b[1], b[2], b[3])
}

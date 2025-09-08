package service

import (
	"sync"

	// No longer importing interfaces.ProtocolFactory directly here
)

var (
	protocolRegistry     = make(map[string]ProtocolFactory) // Use service.ProtocolFactory
	protocolRegistryLock sync.RWMutex
)

// RegisterProtocol registers a new protocol factory with the service module.
// This function is safe for concurrent use.
func RegisterProtocol(name string, factory ProtocolFactory) { // Use service.ProtocolFactory
	protocolRegistryLock.Lock()
	defer protocolRegistryLock.Unlock()
	protocolRegistry[name] = factory
}

// getProtocolFactory retrieves a registered protocol factory by name.
// This function is safe for concurrent use.
func getProtocolFactory(name string) (ProtocolFactory, bool) { // Use service.ProtocolFactory
	protocolRegistryLock.RLock()
	defer protocolRegistryLock.RUnlock()
	factory, ok := protocolRegistry[name]
	return factory, ok
}

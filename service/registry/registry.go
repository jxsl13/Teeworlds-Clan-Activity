package registry

import (
	"sync"

	"github.com/jxsl13/Teeworlds-Clan-Activity/dto"
)

// NewPlayerListConsumerRegistry creates a new registry for channels that are used as primary
// data transfer entities.
func NewPlayerListConsumerRegistry() PlayerListConsumerRegistry {
	return PlayerListConsumerRegistry{
		m: make(map[string]ResponseChannel),
	}
}

// ResponseChannel contains dto.Respons objects that can be processed for further
// notification services
type ResponseChannel chan dto.Response

// PlayerListConsumerRegistry is the central object that contains all of the channels
// that are used by different services.
// it is a besic fan-out pattern where the api_cosumer distributes data
// to all of the other services that do their jobs respectively.
type PlayerListConsumerRegistry struct {
	mu sync.RWMutex
	m  map[string]ResponseChannel
}

// Register adds a new channel to the map
func (cr *PlayerListConsumerRegistry) Register(key string, size int) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.m[key] = make(ResponseChannel, size)
}

// Unregister closes and removes the channel from the map again.
func (cr *PlayerListConsumerRegistry) Unregister(key string) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	close(cr.m[key])
	delete(cr.m, key)
}

// Get returns the response channel that can be read from or pushed to
func (cr *PlayerListConsumerRegistry) Get(key string) ResponseChannel {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	return cr.m[key]
}

// GetAll returns a list of all channels from the map
func (cr *PlayerListConsumerRegistry) GetAll() []ResponseChannel {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	result := make([]ResponseChannel, 0, len(cr.m))
	for _, rChannel := range cr.m {
		result = append(result, rChannel)
	}
	return result
}

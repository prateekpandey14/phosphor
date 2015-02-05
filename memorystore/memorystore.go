package memorystore

import (
	"sync"
	"time"

	"github.com/mattheath/phosphor/domain"
)

type MemoryStore struct {
	sync.RWMutex
	store map[string]domain.Trace
}

// New initialises and returns a new MemoryStore
func New() *MemoryStore {
	s := &MemoryStore{
		store: make(map[string]domain.Trace),
	}

	// run stats worker
	go s.statsLoop()

	return s
}

// GetTrace retrieves a full Trace, composed of Frames from the store by ID
func (s *MemoryStore) GetTrace(id string) (domain.Trace, error) {
	s.RLock()
	defer s.RUnlock()

	return s.store[id], nil
}

// StoreTraceFrame into the store, if the trace doesn't not already exist
// this will be created for the global trace ID
func (s *MemoryStore) StoreTraceFrame(f domain.Frame) error {
	s.Lock()
	defer s.Unlock()

	// Load our current trace
	t := s.store[f.TraceId]

	// Add the new frame to this
	t = append(t, f)

	// Store it back
	s.store[f.TraceId] = t

	return nil
}

// statsLoop loops and outputs stats every 5 seconds
func (s *MemoryStore) statsLoop() {

	tick := time.NewTicker(5 * time.Second)

	// @todo listen for shutdown, stop ticker and exit cleanly
	for {
		<-tick.C // block until tick

		s.printStats()
	}
}

// printStats about the status of the memorystore to stdout
func (s *MemoryStore) printStats() {

	// Get some data while under the mutex
	s.RLock()
	count := len(s.store)
	s.RUnlock()

	// Separate processing and logging outside of mutex
	log.Infof("[Phosphor] Traces stored: %v", count)
}

package node

import (
	"context"
	"sync"
)

type state interface {
	get(ctx context.Context, key []byte) (error, []byte)
	put(ctx context.Context, key, value []byte) error
	delete(ctx context.Context, key []byte) error
}

type memoryState struct {
	mu    sync.RWMutex
	state map[string][]byte
}

func newMemoryState() *memoryState {
	return &memoryState{
		mu:    sync.RWMutex{},
		state: make(map[string][]byte),
	}
}

func (s *memoryState) get(ctx context.Context, key []byte) (error, []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	select {
	case <-ctx.Done():
		return ctx.Err(), nil
	default:
	}

	return nil, s.state[string(key)]
}

func (s *memoryState) put(ctx context.Context, key, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.state[string(key)] = value

	return nil
}

func (s *memoryState) delete(ctx context.Context, key []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	delete(s.state, string(key))

	return nil
}

var _ state = (*memoryState)(nil)

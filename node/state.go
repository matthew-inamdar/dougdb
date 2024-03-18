package node

import (
	"context"
	"sync"
)

type State interface {
	Get(ctx context.Context, key []byte) ([]byte, error)
	Put(ctx context.Context, key, value []byte) error
	Delete(ctx context.Context, key []byte) error
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

func (s *memoryState) Get(ctx context.Context, key []byte) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return s.state[string(key)], nil
}

func (s *memoryState) Put(ctx context.Context, key, value []byte) error {
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

func (s *memoryState) Delete(ctx context.Context, key []byte) error {
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

var _ State = (*memoryState)(nil)

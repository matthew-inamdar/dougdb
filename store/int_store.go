package store

import "encoding/binary"

type IntStore struct {
	s *Store
}

func NewIntStore(name, nodeID string) (*IntStore, error) {
	s, err := NewStore(name, nodeID)
	if err != nil {
		return nil, err
	}

	return &IntStore{s: s}, nil
}

func (s *IntStore) Get() uint64 {
	b := s.s.Get()
	return binary.LittleEndian.Uint64(b)
}

func (s *IntStore) Set(value uint64) error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, value)
	return s.s.Set(b)
}

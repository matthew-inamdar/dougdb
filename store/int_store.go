package store

import "encoding/binary"

type IntStore struct {
	s     *Store
	value int
}

func NewIntStore(name, nodeID string) (*IntStore, error) {
	s, err := NewStore(name, nodeID)
	if err != nil {
		return nil, err
	}

	i := &IntStore{s: s}
	b := s.Get()
	i.value = int(binary.LittleEndian.Uint64(b))

	return i, nil
}

func (s *IntStore) Get() int {
	return s.value
}

func (s *IntStore) Set(value int) error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(value))
	if err := s.s.Set(b); err != nil {
		return err
	}

	s.value = value
	return nil
}

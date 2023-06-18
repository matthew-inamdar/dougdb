package store

import (
	"fmt"
	"github.com/matthew-inamdar/dougdb/filesystem"
	"io"
	"os"
	"sync"
)

type Store struct {
	value []byte

	mut sync.RWMutex
	f   *os.File
}

func NewStore(name, nodeID string) (*Store, error) {
	f, err := filesystem.CreateFile(fmt.Sprintf("%s.store", name), nodeID)
	if err != nil {
		return nil, err
	}

	v, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return &Store{value: v, f: f}, nil
}

func (s *Store) Get() []byte {
	return s.value
}

func (s *Store) Set(value []byte) error {
	if err := s.f.Truncate(0); err != nil {
		return err
	}
	if _, err := s.f.Write(value); err != nil {
		return err
	}

	s.value = value
	return nil
}

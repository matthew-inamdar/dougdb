package store

import (
	"fmt"
	"github.com/matthew-inamdar/dougdb/filesystem"
	"os"
	"sync"
)

type Store struct {
	mut sync.RWMutex
	f   *os.File
}

func NewStore(name, nodeID string) (*Store, error) {
	f, err := filesystem.CreateFile(fmt.Sprintf("%s.store", name), nodeID)
	if err != nil {
		return nil, err
	}

	return &Store{f: f}, nil
}

func (s *Store) Get() []byte {}

func (s *Store) Set(value []byte) []byte {}

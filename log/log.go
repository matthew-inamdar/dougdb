package log

import (
	"github.com/matthew-inamdar/dougdb/filesystem"
	"os"
	"sync"
)

type Log struct {
	entries []Entry

	mut sync.RWMutex
	wal *os.File
}

func NewLog(nodeID string) (*Log, error) {
	f, err := filesystem.CreateFile("wal.log", nodeID)
	if err != nil {
		return nil, err
	}

	return &Log{wal: f}, nil
}

package log

import (
	"fmt"
	"os"
	"sync"
)

const walPath = "/var/lib/dougdb/%s/wal.log"

type Log struct {
	entries []Entry

	mut sync.RWMutex
	wal *os.File
}

func NewLog(nodeID string) (*Log, error) {
	f, err := os.Create(fmt.Sprintf(walPath, nodeID))
	if err != nil {
		return nil, err
	}

	return &Log{
		wal: f,
	}, nil
}

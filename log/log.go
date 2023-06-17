package log

import (
	"fmt"
	"os"
	"runtime"
	"sync"
)

type Log struct {
	entries []Entry

	mut sync.RWMutex
	wal *os.File
}

func NewLog(nodeID string) (*Log, error) {
	f, err := os.Create(getWALPath(nodeID))
	if err != nil {
		return nil, err
	}

	return &Log{
		wal: f,
	}, nil
}

func getWALPath(nodeID string) string {
	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("/var/lib/dougdb/%s/wal.log", nodeID)
	case "darwin":
		return fmt.Sprintf("/Library/Application Support/dougdb/%s/wal.log", nodeID)
	default:
		panic("unsupported operating system, only Linux and macOS supported")
	}
}

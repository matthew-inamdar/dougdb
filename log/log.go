package log

import (
	"bufio"
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

	l := &Log{wal: f}
	l.hydrateEntries()

	return l, nil
}

func (l *Log) hydrateEntries() {
	scanner := bufio.NewScanner(l.wal)
	for scanner.Scan() {
		l.entries = append(l.entries, Entry{Bytes: scanner.Bytes()})
	}
}

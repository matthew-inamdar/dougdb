package log

import (
	"bufio"
	"github.com/matthew-inamdar/dougdb/filesystem"
	"os"
	"sync"
)

type Log struct {
	entries         []Entry
	entryStartIndex int

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
	var startIndexSet bool
	for scanner.Scan() {
		// TODO: Get Index, Term and parse bytes.
		b := scanner.Bytes()
		l.entries = append(l.entries, Entry{Command: b})
		if !startIndexSet {
			startIndexSet = true
			l.entryStartIndex = 0 // TODO: Set from parsed bytes.
		}
	}
}

func (l *Log) At(index int) (Entry, bool) {
	if index >= l.entryStartIndex+len(l.entries) {
		return Entry{}, false
	}

	return l.entries[index-l.entryStartIndex], true
}

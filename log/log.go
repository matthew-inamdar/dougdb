package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type Log struct {
	entries []Entry

	mut sync.RWMutex
	wal *os.File
}

func NewLog(nodeID string) (*Log, error) {
	walPath, err := getWALPath(nodeID)
	if err != nil {
		return nil, err
	}
	if err := ensureDirectoryExists(walPath); err != nil {
		return nil, err
	}
	f, err := os.Create(walPath)
	if err != nil {
		return nil, err
	}

	return &Log{
		wal: f,
	}, nil
}

func getWALPath(nodeID string) (string, error) {
	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("/var/lib/dougdb/%s/wal.log", nodeID), nil
	case "darwin":
		// TODO: Decide how to deal with permissions.
		//return fmt.Sprintf("/Library/Application Support/dougdb/%s/wal.log", nodeID)
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, fmt.Sprintf(".dougdb/%s/wal.log", nodeID)), nil
	default:
		return "", errors.New("unsupported operating system, only Linux and macOS supported")
	}
}

func ensureDirectoryExists(fp string) error {
	dir := filepath.Dir(fp)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

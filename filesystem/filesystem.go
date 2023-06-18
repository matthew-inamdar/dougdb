package filesystem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func CreateFile(filename, nodeID string) (*os.File, error) {
	p, err := getPath(filename, nodeID)
	if err != nil {
		return nil, err
	}
	if err := ensureDirectoryExists(p); err != nil {
		return nil, err
	}
	f, err := os.Create(p)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func getPath(filename, nodeID string) (string, error) {
	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("/var/lib/dougdb/%s/%s", nodeID, filename), nil
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, fmt.Sprintf(".dougdb/%s/%s", nodeID, filename)), nil
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

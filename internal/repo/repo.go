package repo

import (
	"errors"
	"os"
	"path/filepath"
)

var ErrNotFound = errors.New("not a got repository")

func Find() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return FindFrom(dir)
}

func FindFrom(start string) (string, error) {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, ".got")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrNotFound
		}
		dir = parent
	}
}

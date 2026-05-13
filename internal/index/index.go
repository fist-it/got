package index

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Index struct {
	Version uint32
	Entries []*Entry
}

func New() *Index {
	return &Index{Version: 2}
}

func unixTime(sec, nano uint32) time.Time {
	return time.Unix(int64(sec), int64(nano))
}

// gotDir is supposed to be a /.../.got path
func Load(gotDir string) (*Index, error) {
	path := filepath.Join(gotDir, "index")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return New(), nil
		}
		return nil, err
	}
	return Decode(data)
}

func (i *Index) Save(gotDir string) error {
	i.sort()
	data, err := i.Encode()
	if err != nil {
		return err
	}
	path := filepath.Join(gotDir, "index")
	return os.WriteFile(path, data, 0644)
}

func (i *Index) sort() {
	sort.Slice(i.Entries, func(a, b int) bool {
		return i.Entries[a].Path < i.Entries[b].Path
	})
}

func (i *Index) Add(e *Entry) {
	for k, existing := range i.Entries {
		if existing.Path == e.Path {
			i.Entries[k] = e
			return
		}
	}
	i.Entries = append(i.Entries, e)
}

func (i *Index) Remove(path string) bool {
	for k, e := range i.Entries {
		if e.Path == path {
			i.Entries = append(i.Entries[:k], i.Entries[k+1:]...)
			return true
		}
	}
	return false
}

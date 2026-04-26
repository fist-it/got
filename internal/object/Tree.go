package object

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

type TreeEntry struct {
	Mode string
	Name string
	Hash [32]byte
}

type Tree struct {
	Entries []TreeEntry
}

func (t *TreeEntry) Serialize() []byte {
	result := []byte(t.Mode + " ")
	result = append(result, []byte(t.Name)...)
	result = append(result, 0)
	result = append(result, t.Hash[:]...)

	return result
}

func (t *TreeEntry) Type() string {
	return "tree_entry"
}

func DeserializeTreeEntry(data []byte) (*TreeEntry, error) {
	header := data[0:6]

	i := 0
	// Header validation
	for i < 6 {
		d_int, err := strconv.Atoi(string(header[i]))
		if err != nil {
			return nil, errors.New("Invalid TreeEntry mode")
		}

		if d_int < 0 || d_int > 7 {
			return nil, errors.New("Invalid TreeEntry mode")
		}
	}

	nullIdx := bytes.IndexByte(data, 0)
	mode := string(header)
	name := string(data[7 : nullIdx+1])
	hash := [32]byte(data[nullIdx+1 : nullIdx+1+32])

	return &TreeEntry{
		Name: name,
		Mode: mode,
		Hash: hash,
	}, nil
}

func (t *Tree) Serialize() []byte {
	var entries []byte
	for _, entry := range t.Entries {
		entries = append(entries, []byte(entry.Mode)...)
		entries = append(entries, ' ')
		entries = append(entries, []byte(entry.Name)...)
		entries = append(entries, 0)
		entries = append(entries, entry.Hash[:]...)
	}
	header := fmt.Sprintf("tree %d\x00", len(entries))
	result := append([]byte(header), entries...)

	return result
}

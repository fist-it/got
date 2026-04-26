package object

import (
	"bytes"
	"testing"
)

func TestTreeEntryType(t *testing.T) {
	te := &TreeEntry{
		Name: "testing.go",
		Mode: "100777",
		Hash: [32]byte{},
	}
	if te.Type() != "tree_entry" {
		t.Errorf("expected Mode \"tree_entry\", got %q", te.Type())
	}
}

func TestTreeEntrySerialize(t *testing.T) {
	hash := [32]byte{0xff}
	te := &TreeEntry{
		Name: "testing.go",
		Mode: "100644",
		Hash: hash,
	}
	got := te.Serialize()

	nullIdx := bytes.IndexByte(got, 0)
	if string(got[:nullIdx]) != "100644 testing.go" {
		t.Errorf("expected \"100644 testing.go\", got %q", string(got[:nullIdx]))
	}
	if !bytes.Equal(got[nullIdx+1:], hash[:]) {
		t.Errorf("hash mismatch")
	}
}

func TestTreeSerializeSingleEntry(t *testing.T) {
	hash := [32]byte{0xaa, 0xbb, 0xcc}
	tree := &Tree{
		Entries: []TreeEntry{
			{Mode: "100644", Name: "hello.go", Hash: hash},
		},
	}

	got := tree.Serialize()

	// entry: "100644 hello.go\x00" + 32 hash bytes = 16 + 32 = 48
	wantHeader := []byte("tree 48\x00")
	if !bytes.HasPrefix(got, wantHeader) {
		t.Errorf("expected header %q, got prefix %q", wantHeader, got[:len(wantHeader)])
	}
}

func TestTreeSerializeMultipleEntries(t *testing.T) {
	hash1 := [32]byte{0x01}
	hash2 := [32]byte{0x02}
	tree := &Tree{
		Entries: []TreeEntry{
			{Mode: "100644", Name: "a.go", Hash: hash1},
			{Mode: "040000", Name: "dir", Hash: hash2},
		},
	}

	got := tree.Serialize()

	// entry1: "100644 a.go\x00" + 32 = 12 + 32 = 44
	// entry2: "040000 dir\x00" + 32 = 11 + 32 = 43
	wantHeader := []byte("tree 87\x00")
	if !bytes.HasPrefix(got, wantHeader) {
		t.Errorf("expected header %q, got prefix %q", wantHeader, got[:len(wantHeader)])
	}
}

func TestTreeSerializeEmpty(t *testing.T) {
	tree := &Tree{Entries: []TreeEntry{}}
	got := tree.Serialize()
	want := []byte("tree 0\x00")
	if !bytes.Equal(got, want) {
		t.Errorf("expected %q, got %q", want, got)
	}
}

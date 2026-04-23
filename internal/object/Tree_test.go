package object

import (
	"fmt"
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
	te := &TreeEntry{
		Name: "testing.go",
		Mode: "100777",
		Hash: [32]byte{0},
	}
	result := te.Serialize()
	fmt.Print(len(result))
	got := string(result) //      | -> 32 bytes
	want := "100777 testing.go\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

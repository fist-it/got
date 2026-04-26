package object

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"
	"time"
)

func makeTestAuthor(name, email string) Author {
	loc := time.FixedZone("+0200", 2*3600)
	return Author{
		Name:      name,
		Email:     email,
		Timestamp: time.Unix(1714000000, 0).In(loc),
	}
}

func TestCommitType(t *testing.T) {
	c := &Commit{}
	if c.Type() != "commit" {
		t.Errorf("expected \"commit\", got %q", c.Type())
	}
}

func TestCommitSerialize(t *testing.T) {
	treeHash := [32]byte{0xaa, 0xbb, 0xcc}
	author := makeTestAuthor("John Doe", "john@example.com")
	c := &Commit{
		Tree:      treeHash,
		Parents:   nil,
		Author:    author,
		Committer: author,
		Message:   "initial commit",
	}

	got := string(c.Serialize())

	// check header
	if !strings.HasPrefix(got, "commit ") {
		t.Errorf("missing commit header")
	}

	// find body after null byte
	nullIdx := strings.IndexByte(got, 0)
	body := got[nullIdx+1:]

	if !strings.HasPrefix(body, "tree ") {
		t.Errorf("body should start with tree line, got %q", body[:20])
	}
	if !strings.Contains(body, "author John Doe <john@example.com>") {
		t.Errorf("missing author line")
	}
	if !strings.Contains(body, "committer John Doe <john@example.com>") {
		t.Errorf("missing committer line")
	}
	if !strings.HasSuffix(body, "\n\ninitial commit") {
		t.Errorf("expected message at end, got %q", body[len(body)-30:])
	}
}

func TestCommitSerializeWithParents(t *testing.T) {
	treeHash := [32]byte{0xaa}
	parent1 := [32]byte{0xbb}
	parent2 := [32]byte{0xcc}
	author := makeTestAuthor("John Doe", "john@example.com")

	c := &Commit{
		Tree:      treeHash,
		Parents:   [][32]byte{parent1, parent2},
		Author:    author,
		Committer: author,
		Message:   "merge commit",
	}

	got := string(c.Serialize())
	nullIdx := strings.IndexByte(got, 0)
	body := got[nullIdx+1:]

	parentCount := strings.Count(body, "parent ")
	if parentCount != 2 {
		t.Errorf("expected 2 parent lines, got %d", parentCount)
	}
}

func TestCommitSerializeTreeHash(t *testing.T) {
	treeHash := [32]byte{0xaa, 0xbb, 0xcc}
	author := makeTestAuthor("John Doe", "john@example.com")
	c := &Commit{
		Tree:      treeHash,
		Author:    author,
		Committer: author,
		Message:   "test",
	}

	got := string(c.Serialize())
	nullIdx := strings.IndexByte(got, 0)
	body := got[nullIdx+1:]

	wantHex := hex.EncodeToString(treeHash[:])
	treeLine := strings.SplitN(body, "\n", 2)[0]
	if treeLine != "tree "+wantHex {
		t.Errorf("expected tree line %q, got %q", "tree "+wantHex, treeLine)
	}
}

func TestCommitRoundTrip(t *testing.T) {
	treeHash := [32]byte{0xaa, 0xbb, 0xcc, 0xdd}
	parentHash := [32]byte{0x11, 0x22, 0x33}
	author := makeTestAuthor("John Doe", "john@example.com")

	original := &Commit{
		Tree:      treeHash,
		Parents:   [][32]byte{parentHash},
		Author:    author,
		Committer: author,
		Message:   "test commit message",
	}

	serialized := original.Serialize()
	t.Log(string(serialized))
	restored, err := DeserializeCommit(serialized)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(restored.Tree[:], original.Tree[:]) {
		t.Errorf("tree hash mismatch")
	}
	if len(restored.Parents) != 1 {
		t.Fatalf("expected 1 parent, got %d", len(restored.Parents))
	}
	if !bytes.Equal(restored.Parents[0][:], parentHash[:]) {
		t.Errorf("parent hash mismatch")
	}
	if restored.Author.Name != original.Author.Name {
		t.Errorf("author name: expected %q, got %q", original.Author.Name, restored.Author.Name)
	}
	if restored.Author.Email != original.Author.Email {
		t.Errorf("author email: expected %q, got %q", original.Author.Email, restored.Author.Email)
	}
	if restored.Message != original.Message {
		t.Errorf("message: expected %q, got %q", original.Message, restored.Message)
	}
}

func TestCommitRoundTripNoParents(t *testing.T) {
	treeHash := [32]byte{0xff}
	author := makeTestAuthor("Jane Smith", "jane@example.com")

	original := &Commit{
		Tree:      treeHash,
		Parents:   nil,
		Author:    author,
		Committer: author,
		Message:   "initial commit",
	}

	serialized := original.Serialize()
	restored, err := DeserializeCommit(serialized)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(restored.Parents) != 0 {
		t.Errorf("expected 0 parents, got %d", len(restored.Parents))
	}
	if restored.Message != "initial commit" {
		t.Errorf("message: expected \"initial commit\", got %q", restored.Message)
	}
}

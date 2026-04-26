package object

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"
	"time"
)

func makeTestTagger() Author {
	loc := time.FixedZone("+0200", 2*3600)
	return Author{
		Name:      "John Doe",
		Email:     "john@example.com",
		Timestamp: time.Unix(1714000000, 0).In(loc),
	}
}

func TestTagType(t *testing.T) {
	tag := &Tag{}
	if tag.Type() != "tag" {
		t.Errorf("expected \"tag\", got %q", tag.Type())
	}
}

func TestTagSerialize(t *testing.T) {
	objHash := [32]byte{0xaa, 0xbb, 0xcc}
	tag := &Tag{
		Object:  objHash,
		ObjType: "commit",
		Name:    "v1.0",
		Tagger:  makeTestTagger(),
		Message: "release v1.0",
	}

	got := string(tag.Serialize())

	// check outer header
	if !strings.HasPrefix(got, "tag ") {
		t.Errorf("missing tag header")
	}

	_, after, _ := strings.Cut(got, "\x00")
	body := after

	// check object line
	wantHex := hex.EncodeToString(objHash[:])
	if !strings.HasPrefix(body, "object "+wantHex+"\n") {
		t.Errorf("expected object line with %s, got %q", wantHex, body[:70])
	}

	if !strings.Contains(body, "type commit\n") {
		t.Errorf("missing type line")
	}
	if !strings.Contains(body, "tag v1.0\n") {
		t.Errorf("missing tag name line")
	}
	if !strings.Contains(body, "tagger John Doe <john@example.com>") {
		t.Errorf("missing tagger line")
	}
	if !strings.HasSuffix(body, "\n\nrelease v1.0") {
		t.Errorf("expected message at end, got %q", body[len(body)-30:])
	}
}

func TestTagSerializeImplementsObject(t *testing.T) {
	var _ Object = &Tag{}
}

func TestTagDeserialize(t *testing.T) {
	objHash := [32]byte{0xaa, 0xbb, 0xcc}
	tagger := makeTestTagger()
	original := &Tag{
		Object:  objHash,
		ObjType: "commit",
		Name:    "v1.0",
		Tagger:  tagger,
		Message: "release v1.0",
	}

	serialized := original.Serialize()
	restored, err := Deserialize(serialized)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(restored.Object[:], original.Object[:]) {
		t.Errorf("object hash mismatch")
	}
	if restored.ObjType != "commit" {
		t.Errorf("expected type \"commit\", got %q", restored.ObjType)
	}
	if restored.Name != "v1.0" {
		t.Errorf("expected name \"v1.0\", got %q", restored.Name)
	}
	if restored.Tagger.Name != tagger.Name {
		t.Errorf("tagger name: expected %q, got %q", tagger.Name, restored.Tagger.Name)
	}
	if restored.Tagger.Email != tagger.Email {
		t.Errorf("tagger email: expected %q, got %q", tagger.Email, restored.Tagger.Email)
	}
	if restored.Message != "release v1.0" {
		t.Errorf("expected message \"release v1.0\", got %q", restored.Message)
	}
}

func TestTagDeserializeInvalidHeader(t *testing.T) {
	input := []byte("blob 5\x00hello")
	_, err := Deserialize(input)
	if err == nil {
		t.Error("expected error for invalid header, got nil")
	}
}

func TestTagRoundTripBlobType(t *testing.T) {
	objHash := [32]byte{0xff, 0xee, 0xdd}
	original := &Tag{
		Object:  objHash,
		ObjType: "blob",
		Name:    "data-v2",
		Tagger:  makeTestTagger(),
		Message: "tagging a blob",
	}

	serialized := original.Serialize()
	restored, err := Deserialize(serialized)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if restored.ObjType != "blob" {
		t.Errorf("expected type \"blob\", got %q", restored.ObjType)
	}
	if restored.Name != "data-v2" {
		t.Errorf("expected name \"data-v2\", got %q", restored.Name)
	}
}

func TestTagRoundTripMultilineMessage(t *testing.T) {
	original := &Tag{
		Object:  [32]byte{0x01},
		ObjType: "commit",
		Name:    "v2.0",
		Tagger:  makeTestTagger(),
		Message: "line one\nline two\nline three",
	}

	serialized := original.Serialize()
	restored, err := Deserialize(serialized)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if restored.Message != original.Message {
		t.Errorf("expected message %q, got %q", original.Message, restored.Message)
	}
}

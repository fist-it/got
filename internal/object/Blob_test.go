package object

import "testing"

func TestBlobType(t *testing.T) {
	b := &Blob{Content: []byte("hello")}
	if b.Type() != "blob" {
		t.Errorf("expected type \"blob\", got %q", b.Type())
	}
}

func TestBlobSerialize(t *testing.T) {
	b := &Blob{Content: []byte("hello")}
	got := string(b.Serialize())
	want := "blob 5\x00hello"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestBlobDeserialize(t *testing.T) {
	input := []byte("blob 5\x00hello")
	b, err := DeserializeBlob(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(b.Content) != "hello" {
		t.Errorf("expected content \"hello\", got %q", string(b.Content))
	}
}

func TestBlobDeserializeInvalidHeader(t *testing.T) {
	input := []byte("tree 5\x00hello")
	_, err := DeserializeBlob(input)
	if err == nil {
		t.Error("expected error for invalid header, got nil")
	}
}

func TestBlobDeserializeRoundTrip(t *testing.T) {
	original := &Blob{Content: []byte("round trip test")}
	serialized := original.Serialize()
	restored, err := DeserializeBlob(serialized)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(restored.Content) != string(original.Content) {
		t.Errorf("expected %q, got %q", string(original.Content), string(restored.Content))
	}
}

func TestBlobSerializeEmpty(t *testing.T) {
	b := &Blob{Content: []byte{}}
	got := string(b.Serialize())
	want := "blob 0\x00"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestBlobImplementsObject(t *testing.T) {
	var _ Object = &Blob{}
}

package object

import (
	"testing"
	"time"
)

func TestAuthorString(t *testing.T) {
	loc := time.FixedZone("+0200", 2*3600)
	a := &Author{
		Name:      "John Doe",
		Email:     "john@example.com",
		Timestamp: time.Unix(1714000000, 0).In(loc),
	}

	got := a.String()
	want := "John Doe <john@example.com> 1714000000 +0200"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestAuthorStringNegativeTimezone(t *testing.T) {
	loc := time.FixedZone("-0530", -(5*3600 + 30*60))
	a := &Author{
		Name:      "Jane Smith",
		Email:     "jane@example.com",
		Timestamp: time.Unix(1714000000, 0).In(loc),
	}

	got := a.String()
	want := "Jane Smith <jane@example.com> 1714000000 -0530"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestParseAuthor(t *testing.T) {
	input := "John Doe <john@example.com> 1714000000 +0200"
	a, err := ParseAuthor(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Name != "John Doe" {
		t.Errorf("expected name \"John Doe\", got %q", a.Name)
	}
	if a.Email != "john@example.com" {
		t.Errorf("expected email \"john@example.com\", got %q", a.Email)
	}
	if a.Timestamp.Unix() != 1714000000 {
		t.Errorf("expected timestamp 1714000000, got %d", a.Timestamp.Unix())
	}
	_, offset := a.Timestamp.Zone()
	if offset != 7200 {
		t.Errorf("expected timezone offset 7200, got %d", offset)
	}
}

func TestParseAuthorCommitter(t *testing.T) {
	input := "Jane Smith <jane@example.com> 1714000000 -0500"
	a, err := ParseAuthor(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Name != "Jane Smith" {
		t.Errorf("expected name \"Jane Smith\", got %q", a.Name)
	}
	if a.Email != "jane@example.com" {
		t.Errorf("expected email \"jane@example.com\", got %q", a.Email)
	}
	_, offset := a.Timestamp.Zone()
	if offset != -18000 {
		t.Errorf("expected timezone offset -18000, got %d", offset)
	}
}

func TestParseAuthorRoundTrip(t *testing.T) {
	loc := time.FixedZone("+0200", 2*3600)
	original := &Author{
		Name:      "John Doe",
		Email:     "john@example.com",
		Timestamp: time.Unix(1714000000, 0).In(loc),
	}

	str := original.String()
	parsed, err := ParseAuthor(str)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed.Name != original.Name {
		t.Errorf("name: expected %q, got %q", original.Name, parsed.Name)
	}
	if parsed.Email != original.Email {
		t.Errorf("email: expected %q, got %q", original.Email, parsed.Email)
	}
	if parsed.Timestamp.Unix() != original.Timestamp.Unix() {
		t.Errorf("timestamp: expected %d, got %d", original.Timestamp.Unix(), parsed.Timestamp.Unix())
	}
}

func TestParseAuthorInvalid(t *testing.T) {
	_, err := ParseAuthor("garbage input")
	if err == nil {
		t.Error("expected error for invalid input, got nil")
	}
}

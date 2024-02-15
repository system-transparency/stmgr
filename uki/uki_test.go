package uki

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetStubFile(t *testing.T) {
	const contents = "dummy stub"

	fileName := filepath.Join(t.TempDir(), "stub.efi")
	if err := os.WriteFile(fileName, []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}

	if got := string(getStub(fileName)); got != contents {
		t.Errorf("Unexpected stub contents, got %q, wanted %q", got, contents)
	}
}

func TestGetStubEmbdded(t *testing.T) {
	fileName := filepath.Join(t.TempDir(), "no-such-file")

	if got := len(getStub(fileName)); got < 10000 {
		t.Errorf("Failed to fall back to embedded stub, got only %d bytes, expected at least 10000", got)
	}
}

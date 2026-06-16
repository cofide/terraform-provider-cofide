package credentials

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromFile_validFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)

	credDir := filepath.Join(dir, ".cofide")
	if err := os.MkdirAll(credDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(credDir, "credentials"), []byte(`{"access_token":"my-token"}`), 0600); err != nil {
		t.Fatal(err)
	}

	token, err := LoadFromFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "my-token" {
		t.Fatalf("expected %q, got %q", "my-token", token)
	}
}

func TestLoadFromFile_fileNotExist(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)

	token, err := LoadFromFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "" {
		t.Fatalf("expected empty token, got %q", token)
	}
}

func TestLoadFromFile_invalidJSON(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)

	credDir := filepath.Join(dir, ".cofide")
	if err := os.MkdirAll(credDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(credDir, "credentials"), []byte(`not-json`), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := LoadFromFile()
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

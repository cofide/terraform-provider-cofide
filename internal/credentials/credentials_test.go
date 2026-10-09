package credentials

import (
	"os"
	"path/filepath"
	"testing"
)

const testTarget = "connect.example.com:443"

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestLoadFromFile(t *testing.T) {
	targetFile := filepath.Join("credentials.d", "connect.example.com%3A443.json")
	tests := []struct {
		name      string
		files     map[string]string
		wantToken string
		wantPath  string
		wantErr   bool
	}{
		{
			name:      "target file",
			files:     map[string]string{targetFile: `{"access_token":"target-token"}`},
			wantToken: "target-token",
			wantPath:  targetFile,
		},
		{
			name:  "legacy file is not used",
			files: map[string]string{"credentials": `{"access_token":"legacy-token"}`},
		},
		{
			name: "other target's file is not used",
			files: map[string]string{
				filepath.Join("credentials.d", "other.example.com%3A443.json"): `{"access_token":"other-token"}`,
			},
		},
		{
			name: "no files",
		},
		{
			name:    "invalid JSON",
			files:   map[string]string{targetFile: `not-json`},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv(DirectoryEnvVar, dir)
			for name, content := range tt.files {
				writeFile(t, filepath.Join(dir, name), content)
			}

			token, path, err := LoadFromFile(testTarget)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if token != tt.wantToken {
				t.Errorf("expected token %q, got %q", tt.wantToken, token)
			}
			wantPath := ""
			if tt.wantPath != "" {
				wantPath = filepath.Join(dir, tt.wantPath)
			}
			if path != wantPath {
				t.Errorf("expected path %q, got %q", wantPath, path)
			}
		})
	}
}

func TestLoadFromFile_homeDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv(DirectoryEnvVar, "")
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	writeFile(t, filepath.Join(home, ".cofide", "credentials.d", "connect.example.com%3A443.json"), `{"access_token":"my-token"}`)

	token, _, err := LoadFromFile(testTarget)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "my-token" {
		t.Fatalf("expected %q, got %q", "my-token", token)
	}
}

func TestLoadFromFile_noTarget(t *testing.T) {
	if _, _, err := LoadFromFile(""); err == nil {
		t.Fatal("expected error for empty target, got nil")
	}
}

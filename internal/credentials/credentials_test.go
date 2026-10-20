package credentials

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
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
	t.Logf("Wrote file %s", path)
}

// setUserConfigDir points os.UserConfigDir at a `dir`. It is used for
// testing credential loading without affecting the user's real config
// dir.
func setUserConfigDir(t *testing.T, dir string) {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", dir) // Unix, including macOS
	t.Setenv("AppData", dir)         // Windows
}

func TestLoadFromFile(t *testing.T) {
	credsDir := filepath.Join(appDirName, credentialsDirName)
	targetFile := filepath.Join(credsDir, "connect.example.com%3A443.json")
	tests := []struct {
		name      string
		files     map[string]string
		wantToken string
		wantPath  string
		wantErr   bool
	}{
		{
			name: "target file is used",
			files: map[string]string{
				targetFile: `{"access_token":"target-token"}`,
				filepath.Join(credsDir, "other.example.com%3A443.json"): `{"access_token":"other-token"}`,
			},
			wantToken: "target-token",
			wantPath:  targetFile,
		},
		{
			name:    "no files",
			wantErr: true,
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
			setUserConfigDir(t, dir)
			for name, content := range tt.files {
				writeFile(t, filepath.Join(dir, name), content)
			}

			token, path, err := LoadFromFile(testTarget)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.wantToken, token)
			if tt.wantPath != "" {
				require.Equal(t, filepath.Join(dir, tt.wantPath), path)
			}
		})
	}
}

// clearConfigEnv unsets the variables that influence the config dir, so
// that only HOME is consulted. Not applicable on Windows.
func clearConfigEnv(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("HOME is not used on Windows")
	}
	t.Setenv("XDG_CONFIG_HOME", "")
}

func TestLoadFromFile_homeDirectory(t *testing.T) {
	clearConfigEnv(t)
	home := t.TempDir()
	t.Setenv("HOME", home)
	writeFile(t, filepath.Join(home, ".config", appDirName, credentialsDirName, "connect.example.com%3A443.json"), `{"access_token":"my-token"}`)

	token, _, err := LoadFromFile(testTarget)
	require.NoError(t, err)
	require.Equal(t, "my-token", token)
}

func TestLoadFromFile_relativeXDGConfigHome(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("XDG_CONFIG_HOME", "relative/path")

	_, _, err := LoadFromFile(testTarget)
	require.Error(t, err)
}

func TestLoadFromFile_noHome(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("HOME", "")

	_, _, err := LoadFromFile(testTarget)
	require.Error(t, err)
}

func TestLoadFromFile_noTarget(t *testing.T) {
	if _, _, err := LoadFromFile(""); err == nil {
		t.Fatal("expected error for empty target, got nil")
	}
}

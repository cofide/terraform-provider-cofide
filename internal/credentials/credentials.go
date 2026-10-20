package credentials

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// appDirName is the Cofide directory within the user config directory.
	appDirName = "cofide"
	// credentialsDirName holds one credentials file per Connect TLS gRPC target.
	credentialsDirName = "credentials.d"
)

type credentialsFile struct {
	AccessToken string `json:"access_token"`
}

// LoadFromFile reads the API token cached by `cofidectl connect
// login` for the Connect TLS gRPC target. It returns the path the
// token was read from plus any error encountered while locating or
// reading the file.
func LoadFromFile(target string) (token, path string, err error) {
	if target == "" {
		return "", "", errors.New("a Connect TLS gRPC target is required to locate cached credentials")
	}

	path, err = credentialsFilePath(target)
	if err != nil {
		return "", "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("reading %s: %w", path, err)
	}

	var cf credentialsFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return "", "", fmt.Errorf("parsing %s: %w", path, err)
	}
	return cf.AccessToken, path, nil
}

// credentialsFilePath returns the path to the credentials file for
// `target`. It matches the path used by cofidectl.
func credentialsFilePath(target string) (string, error) {
	dir, err := credentialsDir()
	if err != nil {
		return "", fmt.Errorf("could not get credentials directory: %w", err)
	}
	return filepath.Join(dir, credentialsFileName(target)), nil
}

// credentialsDir returns the local Cofide config dir used by this
// Terraform provider. Like cofidectl, it uses os.UserConfigDir,
// except that macOS follows the XDG convention like other Unix
// systems, rather than using ~/Library/Application Support/.
func credentialsDir() (string, error) {
	if runtime.GOOS != "darwin" {
		base, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(base, appDirName, credentialsDirName), nil
	}
	// Mirror Linux config dir pattern on MacOS instead of using os.UserConfigDir
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New("neither $XDG_CONFIG_HOME nor $HOME are defined")
		}
		return filepath.Join(home, ".config", appDirName, credentialsDirName), nil
	}
	if !filepath.IsAbs(base) {
		return "", errors.New("path in $XDG_CONFIG_HOME is relative")
	}
	return filepath.Join(base, appDirName, credentialsDirName), nil
}

// credentialsFileName returns the name of the credentials file for
// `target`. Characters other than ASCII letters, digits, '.', '_' and
// '-' are percent-encoded.
func credentialsFileName(target string) string {
	var b strings.Builder
	for i := range len(target) {
		c := target[i]
		switch {
		case 'a' <= c && c <= 'z', 'A' <= c && c <= 'Z', '0' <= c && c <= '9', c == '.', c == '_', c == '-':
			b.WriteByte(c)
		default:
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String() + ".json"
}

package credentials

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// These match the locations used by `cofidectl connect login`, which writes the credentials file.
const (
	credentialsDirectory = ".cofide"
	// DirectoryEnvVar overrides the directory holding cached credentials, which is ~/.cofide by
	// default.
	DirectoryEnvVar = "COFIDE_CREDENTIALS_DIR"
	// credentialsSubdirectory holds one credentials file per Connect TLS gRPC target.
	credentialsSubdirectory = "credentials.d"
	// legacyCredentialsFile is the single credentials file, shared by every Connect target, used by
	// earlier versions of cofidectl. It is read if there is no credentials file for the target.
	legacyCredentialsFile = "credentials"
)

type credentialsFile struct {
	AccessToken string `json:"access_token"`
}

// LoadFromFile reads the API token cached by `cofidectl connect login` for the Connect TLS gRPC
// target. It reads <dir>/credentials.d/<target>.json, falling back to the legacy <dir>/credentials
// file if that does not exist, where <dir> is $COFIDE_CREDENTIALS_DIR, or ~/.cofide if unset.
// It returns the path the token was read from, and ("", "", nil) if neither file exists.
func LoadFromFile(target string) (token, path string, err error) {
	if target == "" {
		return "", "", errors.New("a Connect TLS gRPC target is required to locate cached credentials")
	}
	dir := os.Getenv(DirectoryEnvVar)
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", "", fmt.Errorf("could not get user home directory: %w", err)
		}
		dir = filepath.Join(home, credentialsDirectory)
	}

	for _, path := range []string{
		filepath.Join(dir, credentialsSubdirectory, credentialsFileName(target)),
		filepath.Join(dir, legacyCredentialsFile),
	} {
		data, err := os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return "", "", err
		}
		var cf credentialsFile
		if err := json.Unmarshal(data, &cf); err != nil {
			return "", "", fmt.Errorf("parsing %s: %w", path, err)
		}
		return cf.AccessToken, path, nil
	}
	return "", "", nil
}

// credentialsFileName returns the name of the credentials file for target, as cofidectl names it.
// Characters other than ASCII letters, digits, '.', '_' and '-' are percent-encoded.
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

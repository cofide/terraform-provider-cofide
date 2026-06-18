package credentials

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type credentialsFile struct {
	AccessToken string `json:"access_token"`
}

// LoadFromFile reads the API token from ~/.cofide/credentials.
// Returns ("", nil) if the file does not exist.
func LoadFromFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, ".cofide", "credentials")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	var cf credentialsFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return "", fmt.Errorf("parsing %s: %w", path, err)
	}
	return cf.AccessToken, nil
}

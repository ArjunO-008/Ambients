//go:build darwin

package config

import (
	"os"
	"path/filepath"
)

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Application Support", "ambientspace")
}

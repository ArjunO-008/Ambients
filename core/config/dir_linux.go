//go:build linux

package config

import (
	"os"
	"path/filepath"
)

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "ambientspace")
}

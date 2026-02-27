//go:build windows

package config

import (
	"os"
	"path/filepath"
)

func configDir() string {
	return filepath.Join(os.Getenv("APPDATA"), "ambientspace")
}

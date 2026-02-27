//go:build linux

package skin

import (
	"os"
	"path/filepath"
)

func customSkinsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "ambientspace", "skins")
}

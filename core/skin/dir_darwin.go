//go:build darwin

package skin

import (
	"os"
	"path/filepath"
)

func customSkinsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Application Support", "ambientspace", "skins")
}

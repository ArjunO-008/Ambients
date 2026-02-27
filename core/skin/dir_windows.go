//go:build windows

package skin

import (
	"os"
	"path/filepath"
)

func customSkinsDir() string {
	appData := os.Getenv("APPDATA")
	return filepath.Join(appData, "ambientspace", "skins")
}

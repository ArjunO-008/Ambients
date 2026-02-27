//go:build linux

package media

import (
	"os/exec"
)

func (m *MediaService) PlayPause() { xdoKey("XF86AudioPlay") }
func (m *MediaService) Next()      { xdoKey("XF86AudioNext") }
func (m *MediaService) Previous()  { xdoKey("XF86AudioPrev") }
func (m *MediaService) Stop()      { xdoKey("XF86AudioStop") }

func xdoKey(key string) {
	exec.Command("xdotool", "key", key).Run()
}

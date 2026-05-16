//go:build linux

package media

import (
	"os/exec"
)

func (m *MediaService) PlayPause() { sendMediaKey("XF86AudioPlay", "0xcd") }
func (m *MediaService) Next()      { sendMediaKey("XF86AudioNext", "0xb0") }
func (m *MediaService) Previous()  { sendMediaKey("XF86AudioPrev", "0xb1") }
func (m *MediaService) Stop()      { sendMediaKey("XF86AudioStop", "0xb2") }

func sendMediaKey(xdotoolKey string, ydotoolCode string) {

	if isAvailable("xdotool") {
		exec.Command("xdotool", "key", xdotoolKey).Run()
		return
	}

	if isAvailable("ydotool") {
		exec.Command("ydotool", "key", ydotoolCode).Run()
		return
	}

}

func isAvailable(bin string) bool {
	_, err := exec.LookPath(bin)
	return err == nil
}

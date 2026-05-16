//go:build darwin

package media

import "os/exec"

func (m *MediaService) PlayPause() { appleScript(`key code 100`) }
func (m *MediaService) Next()      { appleScript(`key code 101`) }
func (m *MediaService) Previous()  { appleScript(`key code 98`) }
func (m *MediaService) Stop()      { appleScript(`key code 100`) }

func appleScript(keyEvent string) {
	script := `tell application "System Events" to ` + keyEvent
	exec.Command("osascript", "-e", script).Run()
}

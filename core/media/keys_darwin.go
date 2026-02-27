//go:build darwin

package media

import (
	"os/exec"
)

func (m *MediaService) PlayPause() {
	runAppleScript(`
		tell application "System Events"
			key code 100 using {}
		end tell
	`)
}

func (m *MediaService) Next() {
	runAppleScript(`
		tell application "System Events"
			key code 101 using {}
		end tell
	`)
}

func (m *MediaService) Previous() {
	runAppleScript(`
		tell application "System Events"
			key code 98 using {}
		end tell
	`)
}

func (m *MediaService) Stop() {
	runAppleScript(`
		tell application "System Events"
			key code 100 using {}
		end tell
	`)
}

func runAppleScript(script string) {
	exec.Command("osascript", "-e", script).Run()
}

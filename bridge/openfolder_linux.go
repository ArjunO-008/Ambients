//go:build linux

package bridge

import "os/exec"

func openFolder(path string) {
	exec.Command("xdg-open", path).Start()
}

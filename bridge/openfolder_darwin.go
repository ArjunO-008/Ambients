//go:build darwin

package bridge

import "os/exec"

func openFolder(path string) {
	exec.Command("open", path).Start()
}

//go:build windows

package bridge

import "os/exec"

func openFolder(path string) {
	exec.Command("explorer", path).Start()
}

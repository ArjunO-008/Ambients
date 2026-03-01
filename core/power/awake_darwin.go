//go:build darwin

package power

import "os/exec"

// On macOS we use caffeinate — a built-in CLI tool that prevents sleep.
// We store the process so we can kill it to restore sleep behavior.
type PowerService struct {
	cmd *exec.Cmd
}

func NewPowerService() *PowerService {
	return &PowerService{}
}

// Prevent starts caffeinate in the background to block display sleep.
// -d prevents display sleep, -i prevents system idle sleep.
func (p *PowerService) Prevent() {
	if p.cmd != nil {
		return // already running
	}
	p.cmd = exec.Command("caffeinate", "-d", "-i")
	p.cmd.Start()
}

// Restore kills caffeinate, re-enabling normal sleep behavior.
func (p *PowerService) Restore() {
	if p.cmd != nil {
		p.cmd.Process.Kill()
		p.cmd = nil
	}
}

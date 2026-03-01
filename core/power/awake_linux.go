//go:build linux

package power

import "os/exec"

// On Linux we use xdg-screensaver or systemd-inhibit to prevent sleep.
// systemd-inhibit is more reliable across desktop environments.
type PowerService struct {
	cmd *exec.Cmd
}

func NewPowerService() *PowerService {
	return &PowerService{}
}

// Prevent uses systemd-inhibit to block sleep and screen lock.
func (p *PowerService) Prevent() {
	if p.cmd != nil {
		return
	}
	// systemd-inhibit blocks sleep/idle until the process exits
	p.cmd = exec.Command(
		"systemd-inhibit",
		"--what=sleep:idle:handle-lid-switch",
		"--who=AmbientSpace",
		"--why=Overlay active",
		"--mode=block",
		"sleep", "infinity", // keep blocking indefinitely
	)
	p.cmd.Start()
}

// Restore kills the inhibitor, re-enabling normal sleep behavior.
func (p *PowerService) Restore() {
	if p.cmd != nil {
		p.cmd.Process.Kill()
		p.cmd = nil
	}
}

//go:build linux

package power

import "os/exec"

type PowerService struct {
	cmd *exec.Cmd
}

func NewPowerService() *PowerService {
	return &PowerService{}
}

func (p *PowerService) Prevent() {
	if p.cmd != nil {
		return
	}
	p.cmd = exec.Command(
		"systemd-inhibit",
		"--what=sleep:idle:handle-lid-switch",
		"--who=AmbientSpace",
		"--why=Overlay active",
		"--mode=block",
		"sleep", "infinity",
	)
	p.cmd.Start()
}

func (p *PowerService) Restore() {
	if p.cmd != nil {
		p.cmd.Process.Kill()
		p.cmd = nil
	}
}

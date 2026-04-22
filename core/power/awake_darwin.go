//go:build darwin

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
	p.cmd = exec.Command("caffeinate", "-d", "-i")
	p.cmd.Start()
}

func (p *PowerService) Restore() {
	if p.cmd != nil {
		p.cmd.Process.Kill()
		p.cmd = nil
	}
}

//go:build linux

package power

import (
	"os/exec"
	"time"
)

type PowerService struct {
	cmd       *exec.Cmd
	stopReset chan struct{}
	backend   string
}

func NewPowerService() *PowerService {
	return &PowerService{}
}

func (p *PowerService) Prevent() {
	if p.backend != "" {
		return
	}

	if _, err := exec.LookPath("systemd-inhibit"); err == nil {
		cmd := exec.Command(
			"systemd-inhibit",
			"--what=sleep:idle:handle-lid-switch",
			"--who=Ambients",
			"--why=Ambient overlay is active",
			"--mode=block",
			"sleep", "infinity",
		)
		if cmd.Start() == nil {
			p.cmd = cmd
			p.backend = "systemd"
			return
		}
	}

	if _, err := exec.LookPath("xdg-screensaver"); err == nil {
		p.stopReset = make(chan struct{})
		p.backend = "xdg"
		go func() {
			exec.Command("xdg-screensaver", "reset").Run()
			ticker := time.NewTicker(50 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-p.stopReset:
					return
				case <-ticker.C:
					exec.Command("xdg-screensaver", "reset").Run()
				}
			}
		}()
		return
	}

	p.backend = "none"
}

func (p *PowerService) Restore() {
	switch p.backend {
	case "systemd":
		if p.cmd != nil && p.cmd.Process != nil {
			p.cmd.Process.Kill()
			p.cmd.Wait()
			p.cmd = nil
		}
	case "xdg":
		if p.stopReset != nil {
			close(p.stopReset)
			p.stopReset = nil
		}
	}
	p.backend = ""
}

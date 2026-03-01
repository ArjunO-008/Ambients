//go:build windows

package power

import (
	"sync"

	"golang.org/x/sys/windows"
)

const (
	esContinuous      = 0x80000000
	esSystemRequired  = 0x00000001
	esDisplayRequired = 0x00000002
)

type PowerService struct {
	mu     sync.Mutex
	proc   *windows.LazyProc
	active bool
}

func NewPowerService() *PowerService {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	return &PowerService{
		proc: kernel32.NewProc("SetThreadExecutionState"),
	}
}

func (p *PowerService) Prevent() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.active {
		return
	}

	if p.proc != nil {
		p.proc.Call(
			uintptr(esContinuous | esSystemRequired | esDisplayRequired),
		)
		p.active = true
	}
}

func (p *PowerService) Restore() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.active {
		return
	}

	if p.proc != nil {
		p.proc.Call(
			uintptr(esContinuous),
		)
		p.active = false
	}
}

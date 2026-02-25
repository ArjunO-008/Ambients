package bridge

import (
	"Ambients/core/clock"
	"context"
)

type Bridge struct {
	ctx          context.Context
	clockService *clock.ClockService
}

func NewBridge() *Bridge {
	return &Bridge{}
}

func (b *Bridge) SetContext(ctx context.Context) {
	b.ctx = ctx

	b.clockService = clock.NewClockService(ctx)
}

func (b *Bridge) StartClock(user24Hour bool) {
	b.clockService.Start(user24Hour)
}

func (b *Bridge) StopClock() {
	b.clockService.Stop()
}

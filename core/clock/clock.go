package clock

import (
	"context"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ClockService struct {
	ctx  context.Context
	stop chan struct{}
}

func NewClockService(ctx context.Context) *ClockService {
	return &ClockService{
		ctx:  ctx,
		stop: make(chan struct{}),
	}
}

type TimePayload struct {
	Hours   string `json:"hours"`
	Minutes string `json:"minutes"`
	Seconds string `json:"seconds"`
	Date    string `json:"date"`
	AmPm    string `json:"ampm"`
}

func (c *ClockService) Start(user24 bool) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				now := time.Now()
				payload := buildPayload(now, user24)

				runtime.EventsEmit(c.ctx, "clock:tick", payload)
			case <-c.stop:
				return
			}
		}
	}()
}

func (c *ClockService) Stop() {
	c.stop <- struct{}{}
}

func buildPayload(t time.Time, use24Hour bool) TimePayload {
	if use24Hour {
		return TimePayload{
			Hours:   fmt.Sprintf("%02d", t.Hour()),
			Minutes: fmt.Sprintf("%02d", t.Minute()),
			Seconds: fmt.Sprintf("%02d", t.Second()),
			Date:    t.Format("Monday, 02 January 2006"),
			AmPm:    "",
		}
	}

	hour := t.Hour() % 12
	if hour == 0 {
		hour = 12
	}

	ampm := "AM"

	if t.Hour() >= 12 {
		ampm = "PM"
	}

	return TimePayload{
		Hours:   fmt.Sprintf("%02d", hour),
		Minutes: fmt.Sprintf("%02d", t.Minute()),
		Seconds: fmt.Sprintf("%02d", t.Second()),
		Date:    t.Format("Monday, 02 January 2006"),
		AmPm:    ampm,
	}
}

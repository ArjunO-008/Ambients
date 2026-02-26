// This file only compiles on Linux.
//go:build linux

package audio

import (
	"errors"
	"strings"

	"github.com/gordonklaus/portaudio"
)

// findLoopbackDevice on Linux looks for a PulseAudio monitor device.
// PulseAudio automatically creates a "monitor" source for every output sink —
// this is the loopback. It typically appears as "Monitor of <device name>".
//
// If using PipeWire (modern Ubuntu/Fedora), it also exposes monitor sources
// the same way for compatibility.
func findLoopbackDevice() (*portaudio.DeviceInfo, error) {
	devices, err := portaudio.Devices()
	if err != nil {
		return nil, err
	}

	for _, d := range devices {
		if d.MaxInputChannels > 0 {
			name := strings.ToLower(d.Name)
			// PulseAudio monitor sources follow this naming pattern
			if strings.Contains(name, "monitor") {
				return d, nil
			}
		}
	}

	return nil, errors.New(
		"no monitor source found — ensure PulseAudio or PipeWire is running. " +
			"Run `pactl list sources` to check available sources",
	)
}

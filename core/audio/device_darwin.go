// This file only compiles on macOS.
//go:build darwin

package audio

import (
	"errors"
	"strings"

	"github.com/gordonklaus/portaudio"
)

// findLoopbackDevice on macOS looks for a virtual audio loopback device.
// macOS does NOT natively expose system audio as an input — the user must
// install a virtual audio driver. BlackHole (free) is the recommended option:
// https://existential.audio/blackhole/
//
// Once BlackHole is installed, it appears as an audio input device.
// The user should set it as their output device OR use a Multi-Output Device
// in Audio MIDI Setup to route both speakers and BlackHole simultaneously.
func findLoopbackDevice() (*portaudio.DeviceInfo, error) {
	devices, err := portaudio.Devices()
	if err != nil {
		return nil, err
	}

	for _, d := range devices {
		if d.MaxInputChannels > 0 {
			name := strings.ToLower(d.Name)
			// BlackHole is the most common free loopback driver for macOS
			if strings.Contains(name, "blackhole") ||
				strings.Contains(name, "loopback") ||
				strings.Contains(name, "soundflower") { // older alternative
				return d, nil
			}
		}
	}

	return nil, errors.New(
		"no loopback device found on macOS — install BlackHole: " +
			"https://existential.audio/blackhole/",
	)
}

//go:build windows

package audio

import (
	"errors"
	"strings"

	"github.com/gordonklaus/portaudio"
)

// findLoopbackDevice targets the WASAPI Stereo Mix device on Windows.
// From device enumeration we know this appears as:
// "Stereo Mix (Realtek(R) Audio)" with api "Windows WASAPI" and in:2
// We prefer WASAPI over MME/DirectSound for lower latency and better accuracy.
func findLoopbackDevice() (*portaudio.DeviceInfo, error) {
	devices, err := portaudio.Devices()
	if err != nil {
		return nil, err
	}

	// First pass: ideal — WASAPI Stereo Mix
	// This is device [11] on this machine
	for _, d := range devices {
		if d.MaxInputChannels > 0 && d.HostApi != nil {
			isWASAPI := strings.Contains(d.HostApi.Name, "WASAPI")
			isStereoMix := strings.Contains(strings.ToLower(d.Name), "stereo mix")
			if isWASAPI && isStereoMix {
				return d, nil
			}
		}
	}

	// Second pass: any Stereo Mix regardless of API
	// Catches MME/DirectSound versions as fallback
	for _, d := range devices {
		if d.MaxInputChannels > 0 {
			if strings.Contains(strings.ToLower(d.Name), "stereo mix") {
				return d, nil
			}
		}
	}

	// Third pass: any WASAPI input device
	// Last resort before giving up
	for _, d := range devices {
		if d.MaxInputChannels > 0 && d.HostApi != nil {
			if strings.Contains(d.HostApi.Name, "WASAPI") {
				return d, nil
			}
		}
	}

	return nil, errors.New(
		"no loopback device found.\n" +
			"Fix: Right-click speaker → Sounds → Recording tab → " +
			"right-click empty area → Show Disabled Devices → " +
			"enable Stereo Mix",
	)
}

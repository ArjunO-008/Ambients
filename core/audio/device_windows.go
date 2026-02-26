// device_windows.go
//go:build windows

package audio

import (
	"errors"
	"strings"

	"github.com/gordonklaus/portaudio"
)

// findLoopbackDevice prefers WASAPI loopback devices.
// Falls back to Stereo Mix if available.
func findLoopbackDevice() (*portaudio.DeviceInfo, error) {
	devices, err := portaudio.Devices()
	if err != nil {
		return nil, err
	}

	// 1️⃣ Prefer WASAPI loopback (output devices)
	for _, d := range devices {
		if d.HostApi != nil &&
			strings.Contains(strings.ToLower(d.HostApi.Name), "wasapi") &&
			d.MaxOutputChannels > 0 {
			return d, nil
		}
	}

	// 2️⃣ Fallback: Stereo Mix / virtual cables
	for _, d := range devices {
		if d.MaxInputChannels > 0 {
			name := strings.ToLower(d.Name)
			if strings.Contains(name, "stereo mix") ||
				strings.Contains(name, "loopback") ||
				strings.Contains(name, "virtual") {
				return d, nil
			}
		}
	}

	return nil, errors.New(
		"no loopback device found — enable Stereo Mix or use a virtual audio cable",
	)
}

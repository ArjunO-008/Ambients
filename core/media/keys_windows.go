//go:build windows

package media

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	vkMediaPlayPause = 0xB3
	vkMediaNextTrack = 0xB0
	vkMediaPrevTrack = 0xB1
	vkMediaStop      = 0xB2
)

type keyboardInput struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

type input struct {
	inputType uint32
	ki        keyboardInput
	padding   [8]byte
}

var (
	user32    = windows.NewLazySystemDLL("user32.dll")
	sendInput = user32.NewProc("SendInput")
)

func sendKey(vk uint16) {
	down := input{
		inputType: 1,
		ki: keyboardInput{
			wVk: vk,
		},
	}
	up := input{
		inputType: 1,
		ki: keyboardInput{
			wVk:     vk,
			dwFlags: 0x0002,
		},
	}

	inputs := []input{down, up}
	sendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		uintptr(unsafe.Sizeof(inputs[0])),
	)

}

func (m *MediaService) PlayPause() { sendKey(vkMediaPlayPause) }
func (m *MediaService) Next()      { sendKey(vkMediaNextTrack) }
func (m *MediaService) Previous()  { sendKey(vkMediaPrevTrack) }
func (m *MediaService) Stop()      { sendKey(vkMediaStop) }

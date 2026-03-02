package shortcut

import (
	"strings"

	"golang.design/x/hotkey"
)

type ShortcutService struct {
	hk      *hotkey.Hotkey
	stop    chan struct{}
	running bool
}

func NewShortcutService() *ShortcutService {
	return &ShortcutService{
		stop: make(chan struct{}, 1),
	}
}

func (s *ShortcutService) Listen(shortcutStr string, onTrigger func()) error {
	s.Stop()

	mods, key, err := parseShortcut(shortcutStr)
	if err != nil {
		return err
	}

	hk := hotkey.New(mods, key)
	if err := hk.Register(); err != nil {
		return err
	}

	s.hk = hk
	s.running = true

	go func() {
		for {
			select {
			case <-hk.Keydown():
				onTrigger()
			case <-s.stop:
				hk.Unregister()
				s.running = false
				return
			}
		}

	}()

	return nil
}

func (s *ShortcutService) Stop() {
	if s.running {
		s.stop <- struct{}{}
	}
}

func parseShortcut(s string) ([]hotkey.Modifier, hotkey.Key, error) {
	parts := strings.Split(strings.ToLower(s), "+")

	var mods []hotkey.Modifier
	var key hotkey.Key

	for _, p := range parts {
		p = strings.TrimSpace(p)
		switch p {
		case "ctrl", "control":
			mods = append(mods, hotkey.ModCtrl)
		case "shift":
			mods = append(mods, hotkey.ModShift)
		case "alt":
			mods = append(mods, modAlt())
		case "super", "win", "cmd":
			mods = append(mods, modSuper())
		default:
			k, err := stringToKey(p)
			if err != nil {
				return nil, 0, err
			}
			key = k
		}
	}

	return mods, key, nil
}

func stringToKey(s string) (hotkey.Key, error) {
	keys := map[string]hotkey.Key{
		"a": hotkey.KeyA, "b": hotkey.KeyB, "c": hotkey.KeyC,
		"d": hotkey.KeyD, "e": hotkey.KeyE, "f": hotkey.KeyF,
		"g": hotkey.KeyG, "h": hotkey.KeyH, "i": hotkey.KeyI,
		"j": hotkey.KeyJ, "k": hotkey.KeyK, "l": hotkey.KeyL,
		"m": hotkey.KeyM, "n": hotkey.KeyN, "o": hotkey.KeyO,
		"p": hotkey.KeyP, "q": hotkey.KeyQ, "r": hotkey.KeyR,
		"s": hotkey.KeyS, "t": hotkey.KeyT, "u": hotkey.KeyU,
		"v": hotkey.KeyV, "w": hotkey.KeyW, "x": hotkey.KeyX,
		"y": hotkey.KeyY, "z": hotkey.KeyZ,
		"space":  hotkey.KeySpace,
		"return": hotkey.KeyReturn,
		"f1":     hotkey.KeyF1, "f2": hotkey.KeyF2,
		"f3": hotkey.KeyF3, "f4": hotkey.KeyF4,
		"f5": hotkey.KeyF5, "f6": hotkey.KeyF6,
		"f7": hotkey.KeyF7, "f8": hotkey.KeyF8,
		"f9": hotkey.KeyF9, "f10": hotkey.KeyF10,
		"f11": hotkey.KeyF11, "f12": hotkey.KeyF12,
	}

	if k, ok := keys[s]; ok {
		return k, nil
	}
	return 0, nil
}

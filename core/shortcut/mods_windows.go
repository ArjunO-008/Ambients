//go:build windows

package shortcut

import "golang.design/x/hotkey"

func modAlt() hotkey.Modifier   { return hotkey.ModAlt }
func modSuper() hotkey.Modifier { return hotkey.ModWin }

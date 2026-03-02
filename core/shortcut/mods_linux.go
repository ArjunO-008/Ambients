//go:build linux

package shortcut

import "golang.design/x/hotkey"

func modAlt() hotkey.Modifier   { return hotkey.Mod1 }
func modSuper() hotkey.Modifier { return hotkey.Mod4 }

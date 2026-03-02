//go:build darwin

package shortcut

import "golang.design/x/hotkey"

func modAlt() hotkey.Modifier   { return hotkey.ModOption }
func modSuper() hotkey.Modifier { return hotkey.ModCmd }

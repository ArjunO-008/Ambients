//go:build linux

package bridge

import "github.com/sqweek/dialog"

func pickExeFile() (string, error) {
	return dialog.File().
		Title("Select Music Player").
		Load()
}

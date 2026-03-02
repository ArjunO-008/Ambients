//go:build windows

package bridge

import "github.com/sqweek/dialog"

func pickExeFile() (string, error) {
	return dialog.File().
		Title("Select Music Player").
		Filter("Executable", "exe").
		Load()
}

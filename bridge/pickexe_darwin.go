//go:build darwin

package bridge

import "github.com/sqweek/dialog"

func pickExeFile() (string, error) {
	return dialog.File().
		Title("Select Music Player").
		Filter("Application", "app").
		Load()
}

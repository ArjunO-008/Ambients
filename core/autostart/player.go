package autostart

import (
	"Ambients/core/media"
	"os/exec"
	"time"

	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func LaunchAndPlay(playerPath string, ctx context.Context) error {
	if playerPath == "" {
		return nil
	}

	cmd := exec.Command(playerPath)
	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {

		time.Sleep(5 * time.Second)

		m := media.NewMediaService()
		m.PlayPause()

		time.Sleep(300 * time.Millisecond)
		runtime.WindowSetAlwaysOnTop(ctx, true)
		time.Sleep(100 * time.Millisecond)
		runtime.WindowSetAlwaysOnTop(ctx, false)
	}()

	return nil
}

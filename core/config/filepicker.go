package config

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sqweek/dialog"
)

var SupportedImageExts = []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}
var SupportedVideoExts = []string{".mp4", ".webm", ".mov", ".mkv"}

func PickBackgroundFile() (path string, mediaType string, errMsg string) {
	// sqweek/dialog on Windows needs filters added one by one
	// passing a slice directly doesn't work correctly
	dlg := dialog.File().Title("Choose Background Image or Video")

	// add image filter
	dlg = dlg.Filter("Image Files (jpg, jpeg, png, webp, gif)", "jpg", "jpeg", "png", "webp", "gif")

	// add video filter
	dlg = dlg.Filter("Video Files (mp4, webm, mov, mkv)", "mp4", "webm", "mov", "mkv")

	// add all files fallback so user isn't stuck
	dlg = dlg.Filter("All Files", "*")

	filePath, err := dlg.Load()

	if err != nil {
		if err == dialog.ErrCancelled {
			return "", "", ""
		}
		return "", "", "file picker error: " + err.Error()
	}

	ext := strings.ToLower(filepath.Ext(filePath))

	for _, e := range SupportedImageExts {
		if ext == e {
			return filePath, "image", ""
		}
	}

	for _, e := range SupportedVideoExts {
		if ext == e {
			if err := validateVideoDuration(filePath); err != nil {
				return "", "", err.Error()
			}
			return filePath, "video", ""
		}
	}

	return "", "", "unsupported file type: " + ext
}

func validateVideoDuration(path string) error {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	var duration float64
	fmt.Sscanf(strings.TrimSpace(string(out)), "%f", &duration)

	if duration > 60 {
		return fmt.Errorf(
			"video is %.0f seconds — maximum allowed is 60 seconds",
			duration,
		)
	}

	_ = runtime.GOOS
	_ = time.Second
	return nil
}

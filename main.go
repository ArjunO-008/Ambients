package main

import (
	"Ambients/bridge"
	"Ambients/core/tray"
	"context"
	"embed"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS
var appCtx context.Context

func main() {
	app := NewApp()
	b := bridge.NewBridge()

	err := wails.Run(&options.App{
		Title:             "Ambients",
		Width:             1024,
		Height:            768,
		StartHidden:       true,
		HideWindowOnClose: true,
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Middleware: mediaMiddleware, // ← middleware runs BEFORE asset lookup
		},
		BackgroundColour: &options.RGBA{R: 8, G: 8, B: 8, A: 1},
		OnStartup: func(ctx context.Context) {
			appCtx = ctx
			app.startup(ctx)
			b.SetContext(ctx)
			b.StartShortcutListener(ctx)

			go systray.Run(onTrayReady, onTrayExit)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			runtime.WindowHide(ctx)
			return true
		},
		Bind: []interface{}{
			app,
			b,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// mediaMiddleware intercepts /media?path=... requests and serves local files.
// Middleware runs before Wails checks the embedded assets — so /media is
// caught here before it ever reaches the frontend bundle.
func mediaMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// only intercept /media requests
		if r.URL.Path != "/media" {
			next.ServeHTTP(w, r) // pass everything else through normally
			return
		}
		println("MEDIA REQUEST:", r.URL.String())

		filePath := r.URL.Query().Get("path")
		if filePath == "" {
			http.Error(w, "missing path", 400)
			return
		}

		// security: only allow known media extensions
		allowed := map[string]string{
			".jpg":  "image/jpeg",
			".jpeg": "image/jpeg",
			".png":  "image/png",
			".webp": "image/webp",
			".gif":  "image/gif",
			".mp4":  "video/mp4",
			".webm": "video/webm",
			".mov":  "video/quicktime",
			".mkv":  "video/x-matroska",
		}

		ext := strings.ToLower(filepath.Ext(filePath))
		mime, ok := allowed[ext]
		if !ok {
			http.Error(w, "file type not allowed", 403)
			return
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "file not found: "+err.Error(), 404)
			return
		}

		w.Header().Set("Content-Type", mime)
		w.Header().Set("Cache-Control", "private, max-age=3600") // cache for 1hr
		w.Write(data)
	})
}

func onTrayReady() {
	// set a simple icon — 32x32 ICO bytes
	// for now use a minimal placeholder
	systray.SetTitle("Ambients")
	systray.SetTooltip("AmbientSpace")

	// set icon — we'll use a bundled icon
	systray.SetIcon(getTrayIcon())

	// menu items
	mOverlay := systray.AddMenuItem("Show Overlay", "Open the ambient overlay")
	mSettings := systray.AddMenuItem("Settings", "Open settings")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Exit AmbientSpace")

	// handle menu clicks in goroutine
	go func() {
		for {
			select {
			case <-mOverlay.ClickedCh:
				if appCtx != nil {
					runtime.WindowShow(appCtx)
					runtime.WindowFullscreen(appCtx)
					runtime.EventsEmit(appCtx, "overlay:show")
				}

			case <-mSettings.ClickedCh:
				if appCtx != nil {
					runtime.WindowShow(appCtx)
					runtime.WindowUnfullscreen(appCtx)
					runtime.EventsEmit(appCtx, "settings:open")
				}

			case <-mQuit.ClickedCh:
				systray.Quit()
				if appCtx != nil {
					runtime.Quit(appCtx)
				}
			}
		}
	}()
}
func onTrayExit() {
	// cleanup when tray exits
}
func getTrayIcon() []byte {
	return tray.Icon
}

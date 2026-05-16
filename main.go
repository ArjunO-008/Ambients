package main

import (
	"Ambients/bridge"
	"context"
	"embed"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	b := bridge.NewBridge()
	frameless := runtime.GOOS == "darwin"

	opts := &options.App{
		Title:            "Ambients",
		Width:            1024,
		Height:           768,
		Frameless:        frameless,
		BackgroundColour: &options.RGBA{R: 8, G: 8, B: 8, A: 1},
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Middleware: mediaMiddleware,
		},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			b.SetContext(ctx)
		},
		Bind: []interface{}{
			app,
			b,
		},
	}

	if runtime.GOOS == "darwin" {
		opts.Mac = &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "Ambients",
				Message: "A minimal ambient overlay\nMIT License — Arjun.O",
			},
		}
	}

	if err := wails.Run(opts); err != nil {
		println("Error:", err.Error())
	}
}

func mediaMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/media" {
			next.ServeHTTP(w, r)
			return
		}

		filePath := r.URL.Query().Get("path")
		if filePath == "" {
			http.Error(w, "missing path", 400)
			return
		}

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
		w.Header().Set("Cache-Control", "private, max-age=3600")
		w.Write(data)
	})
}

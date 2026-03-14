package main

import (
	"Ambients/bridge"
	"context"
	"embed"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	b := bridge.NewBridge()

	err := wails.Run(&options.App{
		Title:  "Ambients",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Middleware: mediaMiddleware,
		},
		BackgroundColour: &options.RGBA{R: 8, G: 8, B: 8, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			b.SetContext(ctx)
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

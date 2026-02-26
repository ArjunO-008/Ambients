package main

import (
	"Ambients/bridge"
	"context"
	"embed"

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
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},

		// Context is guaranteed ready here
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			b.SetContext(ctx)
		},

		// Safe because Bridge now guards initialization
		Bind: []interface{}{
			app,
			b,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

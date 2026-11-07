package main

import (
	"embed"
	"log/slog"

	"github.com/ColinEge/8bd-ultimate-appgithub.com/ColinEge/8bitdo-ultimate-app/internal/app"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/build
var assets embed.FS

func main() {
	log := slog.Default()

	// Create an instance of the app structure
	app := app.New(log)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "8BitDo Mapper Ultimatest Software V1",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		OnShutdown:       app.Shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Error(err.Error())
	}
}
